package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Session represents a Flow work session
type Session struct {
	Tag         string        `json:"tag"`
	StartTime   time.Time     `json:"start_time"`
	PausedAt    time.Time     `json:"paused_at,omitempty"`
	IsPaused    bool          `json:"is_paused"`
	TotalPaused time.Duration `json:"total_paused"`
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
func getSessionPath() (string, error) {
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

// getLogPath returns the path to the session log file for a specific month
func getLogPath(date time.Time) (string, error) {
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

// getLogDir returns the directory containing all log files
func getLogDir() (string, error) {
	// Use same logic as getLogPath but return the logs directory
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

func sessionExists() bool {
	path, err := getSessionPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

func loadSession() (Session, error) {
	var session Session
	path, err := getSessionPath()
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

func saveSession(session Session) error {
	path, err := getSessionPath()
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

// logSession appends a completed session to the appropriate monthly log file
func logSession(entry LogEntry) error {
	logPath, err := getLogPath(entry.EndTime)
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
