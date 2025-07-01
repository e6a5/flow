package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "seconds only",
			duration: 45 * time.Second,
			expected: "45s",
		},
		{
			name:     "minutes only",
			duration: 5 * time.Minute,
			expected: "5m",
		},
		{
			name:     "minutes and seconds",
			duration: 5*time.Minute + 30*time.Second,
			expected: "5m",
		},
		{
			name:     "hours and minutes",
			duration: 2*time.Hour + 30*time.Minute,
			expected: "2h 30m",
		},
		{
			name:     "hours only",
			duration: 3 * time.Hour,
			expected: "3h 0m",
		},
		{
			name:     "complex duration",
			duration: 1*time.Hour + 23*time.Minute + 45*time.Second,
			expected: "1h 23m",
		},
		{
			name:     "zero duration",
			duration: 0,
			expected: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
		})
	}
}

func TestSessionPersistence(t *testing.T) {
	// Create temporary directory for test session files
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", originalHomeDir)
	}()
	os.Setenv("HOME", tempDir)

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
			os.Remove(getSessionPath())
		})
	}
}

func TestSessionStateTransitions(t *testing.T) {
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", originalHomeDir)
	}()
	os.Setenv("HOME", tempDir)

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
			os.Remove(getSessionPath())

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
				saveSession(session)
			case "pause":
				session, _ := loadSession()
				session.IsPaused = true
				session.PausedAt = time.Now()
				saveSession(session)
			case "resume":
				session, _ := loadSession()
				session.TotalPaused += time.Since(session.PausedAt)
				session.IsPaused = false
				session.PausedAt = time.Time{}
				saveSession(session)
			case "end":
				os.Remove(getSessionPath())
			}

			// Check resulting state
			actualState := getCurrentSessionState()
			if actualState != tt.expectedState {
				t.Errorf("After %s operation, got state %v, want %v", tt.operation, actualState, tt.expectedState)
			}

			// Clean up
			os.Remove(getSessionPath())
		})
	}
}

func TestTagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "no tag provided",
			args:     []string{"flow", "start"},
			expected: "Deep Work",
		},
		{
			name:     "tag with --tag flag",
			args:     []string{"flow", "start", "--tag", "writing docs"},
			expected: "writing docs",
		},
		{
			name:     "tag with --tag= syntax",
			args:     []string{"flow", "start", "--tag=refactoring"},
			expected: "refactoring",
		},
		{
			name:     "tag with quotes",
			args:     []string{"flow", "start", "--tag", "code review session"},
			expected: "code review session",
		},
		{
			name:     "tag with special characters",
			args:     []string{"flow", "start", "--tag", "work 🌊 session"},
			expected: "work 🌊 session",
		},
		{
			name:     "empty tag with --tag flag",
			args:     []string{"flow", "start", "--tag", ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTagFromArgs(tt.args[2:]) // Skip "flow start"
			if result != tt.expected {
				t.Errorf("parseTagFromArgs(%v) = %q, want %q", tt.args[2:], result, tt.expected)
			}
		})
	}
}

func TestSessionEnforcement(t *testing.T) {
	tempDir := t.TempDir()
	originalHomeDir := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", originalHomeDir)
	}()
	os.Setenv("HOME", tempDir)

	tests := []struct {
		name            string
		existingSession *Session
		newTag          string
		shouldAllow     bool
		expectedMessage string
	}{
		{
			name:            "no existing session",
			existingSession: nil,
			newTag:          "new work",
			shouldAllow:     true,
			expectedMessage: "",
		},
		{
			name: "existing active session",
			existingSession: &Session{
				Tag:       "existing work",
				StartTime: time.Now().Add(-30 * time.Minute),
				IsPaused:  false,
			},
			newTag:          "new work",
			shouldAllow:     false,
			expectedMessage: "Already in deep work",
		},
		{
			name: "existing paused session",
			existingSession: &Session{
				Tag:         "existing work",
				StartTime:   time.Now().Add(-30 * time.Minute),
				PausedAt:    time.Now().Add(-10 * time.Minute),
				IsPaused:    true,
				TotalPaused: 5 * time.Minute,
			},
			newTag:          "new work",
			shouldAllow:     false,
			expectedMessage: "You have a paused session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before each test
			os.Remove(getSessionPath())

			// Set up existing session if provided
			if tt.existingSession != nil {
				err := saveSession(*tt.existingSession)
				if err != nil {
					t.Fatalf("Failed to save existing session: %v", err)
				}
			}

			// Test session enforcement
			canStart, message := canStartNewSession()
			if canStart != tt.shouldAllow {
				t.Errorf("canStartNewSession() = %v, want %v", canStart, tt.shouldAllow)
			}

			if tt.expectedMessage != "" && message != tt.expectedMessage {
				// For this test, we just check that the message contains expected keywords
				if tt.expectedMessage == "Already in deep work" && !containsKeywords(message, []string{"Already", "deep work"}) {
					t.Errorf("Expected message to contain 'Already in deep work', got: %s", message)
				}
				if tt.expectedMessage == "You have a paused session" && !containsKeywords(message, []string{"paused session"}) {
					t.Errorf("Expected message to contain 'paused session', got: %s", message)
				}
			}

			// Clean up
			os.Remove(getSessionPath())
		})
	}
}

// Helper types and functions for testing

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

func parseTagFromArgs(args []string) string {
	if len(args) == 0 {
		return "Deep Work"
	}

	for i, arg := range args {
		if arg == "--tag" && i+1 < len(args) {
			return args[i+1]
		}
		if len(arg) > 6 && arg[:6] == "--tag=" {
			return arg[6:]
		}
	}

	return "Deep Work"
}

func canStartNewSession() (bool, string) {
	if !sessionExists() {
		return true, ""
	}

	session, err := loadSession()
	if err != nil {
		return true, ""
	}

	if session.IsPaused {
		return false, "You have a paused session"
	}

	return false, "Already in deep work"
}

func containsKeywords(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(text, keyword) {
			return false
		}
	}
	return true
}

// Test basic functions don't panic
func TestBasicFunctions(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "showUsage",
			fn:   showUsage,
		},
		{
			name: "showVersion",
			fn:   showVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s panicked: %v", tt.name, r)
				}
			}()
			tt.fn()
		})
	}
}
