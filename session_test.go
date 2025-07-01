package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSessionPersistence(t *testing.T) {
	// Create temporary directory for test session files
	tempDir := t.TempDir()
	// Set a specific file path for the session to avoid interacting
	// with the new XDG-based logic during this test.
	sessionPath := filepath.Join(tempDir, "test-session.json")
	t.Setenv("FLOW_SESSION_PATH", sessionPath)
	defer t.Setenv("FLOW_SESSION_PATH", "")

	tests := []struct {
		name    string
		session Session
		wantErr bool
	}{
		{
			name: "basic session",
			session: Session{
				Tag:       "test work",
				StartTime: time.Now(),
				IsPaused:  false,
			},
			wantErr: false,
		},
		{
			name: "paused session",
			session: Session{
				Tag:         "paused work",
				StartTime:   time.Now().Add(-1 * time.Hour),
				PausedAt:    time.Now().Add(-30 * time.Minute),
				IsPaused:    true,
				TotalPaused: 15 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "session with special characters",
			session: Session{
				Tag:       "work with émojis 🌊 and unicode",
				StartTime: time.Now(),
				IsPaused:  false,
			},
			wantErr: false,
		},
		{
			name: "empty tag session",
			session: Session{
				Tag:       "",
				StartTime: time.Now(),
				IsPaused:  false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test save
			err := saveSession(tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("saveSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Test session exists
			if !sessionExists() {
				t.Error("sessionExists() = false, want true after saving")
			}

			// Test load
			loaded, err := loadSession()
			if err != nil {
				t.Errorf("loadSession() error = %v", err)
				return
			}

			// Compare loaded session (ignoring time precision differences)
			if loaded.Tag != tt.session.Tag {
				t.Errorf("loaded.Tag = %q, want %q", loaded.Tag, tt.session.Tag)
			}
			if loaded.IsPaused != tt.session.IsPaused {
				t.Errorf("loaded.IsPaused = %v, want %v", loaded.IsPaused, tt.session.IsPaused)
			}

			// Clean up for next test
			path, _ := getSessionPath()
			_ = os.Remove(path)
		})
	}
}

type sessionState int

const (
	sessionStateNone sessionState = iota
	sessionStateActive
	sessionStatePaused
)

func getCurrentSessionState() sessionState {
	if !sessionExists() {
		return sessionStateNone
	}

	session, err := loadSession()
	if err != nil {
		return sessionStateNone
	}

	if session.IsPaused {
		return sessionStatePaused
	}

	return sessionStateActive
}

func TestSessionStateTransitions(t *testing.T) {
	tempDir := t.TempDir()
	sessionPath := filepath.Join(tempDir, "test-session.json")
	t.Setenv("FLOW_SESSION_PATH", sessionPath)
	defer t.Setenv("FLOW_SESSION_PATH", "")

	tests := []struct {
		name           string
		initialSession *Session
		operation      string
		expectedState  sessionState
	}{
		{
			name:           "no session to active",
			initialSession: nil,
			operation:      "start",
			expectedState:  sessionStateActive,
		},
		{
			name: "active to paused",
			initialSession: &Session{
				Tag:       "test work",
				StartTime: time.Now().Add(-30 * time.Minute),
				IsPaused:  false,
			},
			operation:     "pause",
			expectedState: sessionStatePaused,
		},
		{
			name: "paused to active",
			initialSession: &Session{
				Tag:         "test work",
				StartTime:   time.Now().Add(-30 * time.Minute),
				PausedAt:    time.Now().Add(-10 * time.Minute),
				IsPaused:    true,
				TotalPaused: 5 * time.Minute,
			},
			operation:     "resume",
			expectedState: sessionStateActive,
		},
		{
			name: "active to ended",
			initialSession: &Session{
				Tag:       "test work",
				StartTime: time.Now().Add(-30 * time.Minute),
				IsPaused:  false,
			},
			operation:     "end",
			expectedState: sessionStateNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before each test
			path, _ := getSessionPath()
			_ = os.Remove(path)

			// Set up initial session if provided
			if tt.initialSession != nil {
				err := saveSession(*tt.initialSession)
				if err != nil {
					t.Fatalf("Failed to save initial session: %v", err)
				}
			}

			// Perform operation and check resulting state
			switch tt.operation {
			case "start":
				// Simulate start operation
				session := Session{
					Tag:       "new work",
					StartTime: time.Now(),
					IsPaused:  false,
				}
				if err := saveSession(session); err != nil {
					t.Fatalf("Failed to save session for start operation: %v", err)
				}
			case "pause":
				session, _ := loadSession()
				session.IsPaused = true
				session.PausedAt = time.Now()
				if err := saveSession(session); err != nil {
					t.Fatalf("Failed to save session for pause operation: %v", err)
				}
			case "resume":
				session, _ := loadSession()
				session.TotalPaused += time.Since(session.PausedAt)
				session.IsPaused = false
				session.PausedAt = time.Time{}
				if err := saveSession(session); err != nil {
					t.Fatalf("Failed to save session for resume operation: %v", err)
				}
			case "end":
				path, _ := getSessionPath()
				_ = os.Remove(path)
			}

			// Check resulting state
			actualState := getCurrentSessionState()
			if actualState != tt.expectedState {
				t.Errorf("After %s operation, got state %v, want %v", tt.operation, actualState, tt.expectedState)
			}

			// Clean up
			path, _ = getSessionPath()
			_ = os.Remove(path)
		})
	}
}
