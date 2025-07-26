package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Session represents a Flow work session
type Session struct {
	Tag            string        `json:"tag"`
	StartTime      time.Time     `json:"start_time"`
	TargetDuration time.Duration `json:"target_duration,omitempty"`
	PausedAt       time.Time     `json:"paused_at,omitempty"`
	IsPaused       bool          `json:"is_paused"`
	TotalPaused    time.Duration `json:"total_paused"`
}

// LogEntry represents a completed session for logging
type LogEntry struct {
	Tag         string        `json:"tag"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	TotalPaused time.Duration `json:"total_paused,omitempty"`
}

// Session file management
func GetSessionPath() (string, error) {
	// 1. Check for FLOW_SESSION_PATH environment variable
	if path := os.Getenv("FLOW_SESSION_PATH"); path != "" {
		return path, nil
	}

	// 2. Check for XDG_DATA_HOME environment variable
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "flow", "session"), nil
	}

	// 3. Fallback to ~/.local/share/flow/session (default XDG)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %w", err)
	}
	xdgDefaultPath := filepath.Join(homeDir, ".local", "share", "flow", "session")

	// 4. For backward compatibility, check if the old ~/.flow-session file exists
	legacyPath := filepath.Join(homeDir, ".flow-session")
	if _, err := os.Stat(legacyPath); err == nil {
		return legacyPath, nil
	}

	return xdgDefaultPath, nil
}

// GetLogPath returns the path to the session log file for a specific month
func GetLogPath(date time.Time) (string, error) {
	// 1. Check for FLOW_LOG_PATH environment variable (base directory)
	baseDir := ""
	if path := os.Getenv("FLOW_LOG_PATH"); path != "" {
		baseDir = filepath.Dir(path)
	} else if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		baseDir = filepath.Join(xdgDataHome, "flow")
	} else {
		// 3. Fallback to ~/.local/share/flow (default XDG)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home directory: %w", err)
		}
		baseDir = filepath.Join(homeDir, ".local", "share", "flow")
	}

	// Generate filename with YYYYMM format
	monthStr := date.Format("200601") // YYYYMM format
	filename := fmt.Sprintf("%s_sessions.jsonl", monthStr)

	return filepath.Join(baseDir, "logs", filename), nil
}

// GetLogDir returns the directory containing all log files
func GetLogDir() (string, error) {
	// Use same logic as GetLogPath but return the logs directory
	if path := os.Getenv("FLOW_LOG_PATH"); path != "" {
		return filepath.Join(filepath.Dir(path), "logs"), nil
	} else if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return filepath.Join(xdgDataHome, "flow", "logs"), nil
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home directory: %w", err)
		}
		return filepath.Join(homeDir, ".local", "share", "flow", "logs"), nil
	}
}

func SessionExists() bool {
	path, err := GetSessionPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

func LoadSession() (Session, error) {
	var session Session
	path, err := GetSessionPath()
	if err != nil {
		return session, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return session, err
	}
	err = json.Unmarshal(data, &session)
	return session, err
}

func SaveSession(session Session) error {
	path, err := GetSessionPath()
	if err != nil {
		return err
	}
	// Ensure the directory exists
	if err := ensureDir(path); err != nil {
		return err
	}
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LogSession appends a completed session to the appropriate monthly log file
func LogSession(entry LogEntry) error {
	logPath, err := GetLogPath(entry.EndTime)
	if err != nil {
		return err
	}

	// Ensure the directory exists
	if err := ensureDir(logPath); err != nil {
		return err
	}

	// Serialize as JSON Lines format
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	// Append to log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			// We can't return the error, but we can log it.
			// This is a common pattern for deferred close on write operations.
			fmt.Fprintf(os.Stderr, "Error closing log file: %v\n", err)
		}
	}()

	_, err = file.WriteString(string(data) + "\n")
	return err
}

// ensureDir creates the directory for the given path if it doesn't already exist.
func ensureDir(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// IsSessionStale checks if a session has been running for an unreasonable amount of time
func IsSessionStale(session Session, threshold time.Duration) bool {
	if session.IsPaused {
		// For paused sessions, check if they've been paused for too long
		return time.Since(session.PausedAt) > threshold
	}
	// For active sessions, check total running time
	return time.Since(session.StartTime) > threshold
}

// CleanupStaleSession removes a stale session file and optionally logs it as abandoned
func CleanupStaleSession(session Session, logAsAbandoned bool) error {
	if logAsAbandoned {
		// Log the session as abandoned with a special tag
		endTime := time.Now()
		if session.IsPaused {
			endTime = session.PausedAt
		}

		totalDuration := endTime.Sub(session.StartTime) - session.TotalPaused
		if totalDuration < 0 {
			totalDuration = 0 // Ensure non-negative duration
		}

		logEntry := LogEntry{
			Tag:         session.Tag + " [ABANDONED]",
			StartTime:   session.StartTime,
			EndTime:     endTime,
			Duration:    totalDuration,
			TotalPaused: session.TotalPaused,
		}

		if err := LogSession(logEntry); err != nil {
			return fmt.Errorf("failed to log abandoned session: %w", err)
		}
	}

	// Remove the session file
	sessionPath, err := GetSessionPath()
	if err != nil {
		return fmt.Errorf("failed to get session path: %w", err)
	}

	return os.Remove(sessionPath)
}
