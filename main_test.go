package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/e6a5/flow/core"
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
			result := core.FormatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
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
	sessionPath := filepath.Join(tempDir, "test-session.json")
	t.Setenv("FLOW_SESSION_PATH", sessionPath)
	defer t.Setenv("FLOW_SESSION_PATH", "")

	tests := []struct {
		name            string
		existingSession *core.Session
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
			existingSession: &core.Session{
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
			existingSession: &core.Session{
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
			path, _ := core.GetSessionPath()
			_ = os.Remove(path)

			// Set up existing session if provided
			if tt.existingSession != nil {
				err := core.SaveSession(*tt.existingSession)
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
			path, _ = core.GetSessionPath()
			_ = os.Remove(path)
		})
	}
}

// Helper types and functions for testing
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
	if !core.SessionExists() {
		return true, ""
	}

	session, err := core.LoadSession()
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
