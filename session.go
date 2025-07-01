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

// ensureDir creates the directory for the given path if it doesn't already exist.
func ensureDir(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
