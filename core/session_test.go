package core

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
			err := SaveSession(tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Test session exists
			if !SessionExists() {
				t.Error("SessionExists() = false, want true after saving")
			}

			// Test load
			loaded, err := LoadSession()
			if err != nil {
				t.Errorf("LoadSession() error = %v", err)
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
			path, _ := GetSessionPath()
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
	if !SessionExists() {
		return sessionStateNone
	}

	session, err := LoadSession()
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
			path, _ := GetSessionPath()
			_ = os.Remove(path)

			// Set up initial session if provided
			if tt.initialSession != nil {
				err := SaveSession(*tt.initialSession)
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
				if err := SaveSession(session); err != nil {
					t.Fatalf("Failed to save session for start operation: %v", err)
				}
			case "pause":
				session, _ := LoadSession()
				session.IsPaused = true
				session.PausedAt = time.Now()
				if err := SaveSession(session); err != nil {
					t.Fatalf("Failed to save session for pause operation: %v", err)
				}
			case "resume":
				session, _ := LoadSession()
				session.TotalPaused += time.Since(session.PausedAt)
				session.IsPaused = false
				session.PausedAt = time.Time{}
				if err := SaveSession(session); err != nil {
					t.Fatalf("Failed to save session for resume operation: %v", err)
				}
			case "end":
				path, _ := GetSessionPath()
				_ = os.Remove(path)
			}

			// Check resulting state
			actualState := getCurrentSessionState()
			if actualState != tt.expectedState {
				t.Errorf("After %s operation, got state %v, want %v", tt.operation, actualState, tt.expectedState)
			}

			// Clean up
			path, _ = GetSessionPath()
			_ = os.Remove(path)
		})
	}
}

// Test logging functionality
func TestSessionLogging(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tempDir)

	// Create test log entries
	entries := []LogEntry{
		{
			Tag:       "test work",
			StartTime: time.Now().Add(-2 * time.Hour),
			EndTime:   time.Now().Add(-90 * time.Minute),
			Duration:  30 * time.Minute,
		},
		{
			Tag:       "another task",
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   time.Now().Add(-30 * time.Minute),
			Duration:  30 * time.Minute,
		},
	}

	// Test logging entries
	for _, entry := range entries {
		err := LogSession(entry)
		if err != nil {
			t.Fatalf("Failed to log session: %v", err)
		}
	}

	// Test loading entries using new LogReader
	reader, err := NewLogReader()
	if err != nil {
		t.Fatalf("Failed to create log reader: %v", err)
	}

	loadedEntries, err := reader.ReadAllEntries()
	if err != nil {
		t.Fatalf("Failed to load log entries: %v", err)
	}

	if len(loadedEntries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(loadedEntries))
	}

	// Verify entry content (entries are sorted newest first)
	if loadedEntries[0].Tag != "another task" {
		t.Errorf("Expected tag 'another task', got '%s'", loadedEntries[0].Tag)
	}

	if loadedEntries[0].Duration != 30*time.Minute {
		t.Errorf("Expected duration 30m, got %v", loadedEntries[0].Duration)
	}

	// Verify second entry
	if loadedEntries[1].Tag != "test work" {
		t.Errorf("Expected tag 'test work', got '%s'", loadedEntries[1].Tag)
	}
}

func TestEmptyLogHandling(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tempDir)

	// Test loading from non-existent log using new LogReader
	reader, err := NewLogReader()
	if err != nil {
		t.Fatalf("Failed to create log reader: %v", err)
	}

	entries, err := reader.ReadAllEntries()
	if err != nil {
		t.Fatalf("Expected no error for non-existent log, got: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected empty slice, got %d entries", len(entries))
	}
}

func TestLogReaderFiltering(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "filter-test.jsonl")
	t.Setenv("FLOW_LOG_PATH", logPath)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	lastWeek := today.AddDate(0, 0, -8)

	entries := []LogEntry{
		{Tag: "today1", StartTime: today, EndTime: today.Add(30 * time.Minute), Duration: 30 * time.Minute},
		{Tag: "today2", StartTime: today.Add(2 * time.Hour), EndTime: today.Add(2*time.Hour + 30*time.Minute), Duration: 30 * time.Minute},
		{Tag: "yesterday", StartTime: yesterday, EndTime: yesterday.Add(30 * time.Minute), Duration: 30 * time.Minute},
		{Tag: "lastweek", StartTime: lastWeek, EndTime: lastWeek.Add(30 * time.Minute), Duration: 30 * time.Minute},
	}

	// Log all entries
	for _, entry := range entries {
		err := LogSession(entry)
		if err != nil {
			t.Fatalf("Failed to log session: %v", err)
		}
	}

	reader, err := NewLogReader()
	if err != nil {
		t.Fatalf("Failed to create log reader: %v", err)
	}

	// Test today filter
	todayEntries, err := reader.ReadRecentEntries(100, true, false)
	if err != nil {
		t.Fatalf("Failed to read today entries: %v", err)
	}
	if len(todayEntries) != 2 {
		t.Errorf("Expected 2 today entries, got %d", len(todayEntries))
	}

	// Test week filter
	weekEntries, err := reader.ReadRecentEntries(100, false, true)
	if err != nil {
		t.Fatalf("Failed to read week entries: %v", err)
	}
	if len(weekEntries) != 3 { // today1, today2, yesterday
		t.Errorf("Expected 3 week entries, got %d", len(weekEntries))
	}

	// Test no filter
	allEntries, err := reader.ReadAllEntries()
	if err != nil {
		t.Fatalf("Failed to read all entries: %v", err)
	}
	if len(allEntries) != 4 {
		t.Errorf("Expected 4 total entries, got %d", len(allEntries))
	}
}

func TestGetSessionPath_Default(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	expectedPath := filepath.Join(homeDir, ".local", "share", "flow", "session")
	path, err := GetSessionPath()
	if err != nil {
		t.Fatalf("getSessionPath() error = %v", err)
	}

	if path != expectedPath {
		t.Errorf("Expected path %q, got %q", expectedPath, path)
	}
}

func TestGetSessionPath_Legacy(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}
	legacyPath := filepath.Join(homeDir, ".flow-session")
	if _, err := os.Create(legacyPath); err != nil {
		t.Fatalf("Failed to create legacy session file: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Remove(legacyPath)
	})

	path, err := GetSessionPath()
	if err != nil {
		t.Fatalf("GetSessionPath() error = %v", err)
	}

	if path != legacyPath {
		t.Errorf("Expected path %q, got %q", legacyPath, path)
	}
}

func TestGetLogPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}
	expectedDir := filepath.Join(homeDir, ".local", "share", "flow", "logs")
	testDate := time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC)
	logPath, err := GetLogPath(testDate)
	if err != nil {
		t.Fatalf("GetLogPath() error = %v", err)
	}

	expectedLogFile := "202507_sessions.jsonl"
	expectedFullPath := filepath.Join(expectedDir, expectedLogFile)

	if logPath != expectedFullPath {
		t.Errorf("Expected log path %q, got %q", expectedFullPath, logPath)
	}
}

func TestSaveLoadSession(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("FLOW_SESSION_PATH", filepath.Join(tmpDir, "test-session"))

	session := Session{
		Tag:         "test work",
		StartTime:   time.Now(),
		IsPaused:    false,
		TotalPaused: 30 * time.Minute,
	}

	err := SaveSession(session)
	if err != nil {
		t.Fatalf("saveSession() error = %v", err)
	}

	loaded, err := LoadSession()
	if err != nil {
		t.Fatalf("loadSession() error = %v", err)
	}

	if loaded.Tag != session.Tag {
		t.Errorf("loaded.Tag = %q, want %q", loaded.Tag, session.Tag)
	}
	if loaded.IsPaused != session.IsPaused {
		t.Errorf("loaded.IsPaused = %v, want %v", loaded.IsPaused, session.IsPaused)
	}
	if loaded.TotalPaused != session.TotalPaused {
		t.Errorf("loaded.TotalPaused = %v, want %v", loaded.TotalPaused, session.TotalPaused)
	}
}

func TestLogSession(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)

	entry := LogEntry{
		Tag:       "test work",
		StartTime: time.Now().Add(-2 * time.Hour),
		EndTime:   time.Now().Add(-90 * time.Minute),
		Duration:  time.Hour,
	}

	err := LogSession(entry)
	if err != nil {
		t.Fatalf("logSession() error = %v", err)
	}

	logPath, _ := GetLogPath(entry.EndTime)
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty log file, got empty")
	}
}

func TestIsSessionStale(t *testing.T) {
	now := time.Now()
	threshold := 8 * time.Hour

	tests := []struct {
		name     string
		session  Session
		expected bool
	}{
		{
			name: "fresh active session",
			session: Session{
				StartTime: now.Add(-1 * time.Hour),
				IsPaused:  false,
			},
			expected: false,
		},
		{
			name: "stale active session",
			session: Session{
				StartTime: now.Add(-9 * time.Hour),
				IsPaused:  false,
			},
			expected: true,
		},
		{
			name: "fresh paused session",
			session: Session{
				StartTime: now.Add(-2 * time.Hour),
				PausedAt:  now.Add(-1 * time.Hour),
				IsPaused:  true,
			},
			expected: false,
		},
		{
			name: "stale paused session",
			session: Session{
				StartTime: now.Add(-2 * time.Hour),
				PausedAt:  now.Add(-9 * time.Hour),
				IsPaused:  true,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSessionStale(tt.session, threshold)
			if result != tt.expected {
				t.Errorf("IsSessionStale() = %v, want %v", result, tt.expected)
			}
		})
	}
}
