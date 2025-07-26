package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDeleteLogEntry(t *testing.T) {
	// Create temporary directory for test log files
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Set environment variable to use our test directory
	originalLogPath := os.Getenv("FLOW_LOG_PATH")
	t.Setenv("FLOW_LOG_PATH", filepath.Join(tempDir, "test.log"))
	defer t.Setenv("FLOW_LOG_PATH", originalLogPath)

	// Create test entries
	now := time.Now()
	entry1 := LogEntry{
		Tag:       "test session 1",
		StartTime: now.Add(-1 * time.Hour),
		EndTime:   now,
		Duration:  1 * time.Hour,
	}
	entry2 := LogEntry{
		Tag:       "test session 2",
		StartTime: now.Add(-2 * time.Hour),
		EndTime:   now.Add(-1 * time.Hour),
		Duration:  1 * time.Hour,
	}
	entry3 := LogEntry{
		Tag:       "test session 3",
		StartTime: now.Add(-3 * time.Hour),
		EndTime:   now.Add(-2 * time.Hour),
		Duration:  1 * time.Hour,
	}

	// Create a log file with test entries
	logPath, err := GetLogPath(now)
	if err != nil {
		t.Fatalf("Failed to get log path: %v", err)
	}

	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Write test entries to log file
	file, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	entries := []LogEntry{entry1, entry2, entry3}
	for _, entry := range entries {
		data, err := json.Marshal(entry)
		if err != nil {
			t.Fatalf("Failed to marshal entry: %v", err)
		}
		if _, err := file.WriteString(string(data) + "\n"); err != nil {
			t.Fatalf("Failed to write entry: %v", err)
		}
	}

	tests := []struct {
		name           string
		entryToDelete  LogEntry
		wantErr        bool
		expectedError  string
		remainingCount int
	}{
		{
			name:           "delete middle entry",
			entryToDelete:  entry2,
			wantErr:        false,
			remainingCount: 2,
		},
		{
			name:           "delete first entry",
			entryToDelete:  entry1,
			wantErr:        false,
			remainingCount: 1,
		},
		{
			name:           "delete last entry",
			entryToDelete:  entry3,
			wantErr:        false,
			remainingCount: 0,
		},
		{
			name: "delete non-existent entry",
			entryToDelete: LogEntry{
				Tag:       "non-existent",
				StartTime: now.Add(-10 * time.Hour),
				EndTime:   now.Add(-9 * time.Hour),
				Duration:  1 * time.Hour,
			},
			wantErr:       true,
			expectedError: "log entry not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Delete the entry
			err := DeleteLogEntry(tt.entryToDelete)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteLogEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err.Error() != tt.expectedError {
					t.Errorf("DeleteLogEntry() error = %v, want %v", err, tt.expectedError)
				}
				return
			}

			// Verify the entry was actually deleted
			remainingEntries, err := GetRecentSessions(10)
			if err != nil {
				t.Errorf("Failed to get remaining entries: %v", err)
				return
			}

			if len(remainingEntries) != tt.remainingCount {
				t.Errorf("Expected %d remaining entries, got %d", tt.remainingCount, len(remainingEntries))
			}

			// Verify the deleted entry is not in the remaining entries
			for _, entry := range remainingEntries {
				if entry.StartTime.Equal(tt.entryToDelete.StartTime) && entry.Tag == tt.entryToDelete.Tag {
					t.Errorf("Deleted entry still found in log: %+v", entry)
				}
			}
		})
	}
}

func TestDeleteLogEntry_EmptyFile(t *testing.T) {
	// Create temporary directory for test log files
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Set environment variable to use our test directory
	originalLogPath := os.Getenv("FLOW_LOG_PATH")
	t.Setenv("FLOW_LOG_PATH", filepath.Join(tempDir, "test.log"))
	defer t.Setenv("FLOW_LOG_PATH", originalLogPath)

	// Create an empty log file
	logPath, err := GetLogPath(time.Now())
	if err != nil {
		t.Fatalf("Failed to get log path: %v", err)
	}

	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Create empty file
	file, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	// Try to delete a non-existent entry
	entryToDelete := LogEntry{
		Tag:       "non-existent",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		Duration:  1 * time.Hour,
	}

	err = DeleteLogEntry(entryToDelete)
	if err == nil {
		t.Error("Expected error when deleting from empty file, got nil")
		return
	}

	if err.Error() != "log entry not found" {
		t.Errorf("Expected 'log entry not found' error, got: %v", err)
	}
}

func TestDeleteLogEntry_MalformedLines(t *testing.T) {
	// Create temporary directory for test log files
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Set environment variable to use our test directory
	originalLogPath := os.Getenv("FLOW_LOG_PATH")
	t.Setenv("FLOW_LOG_PATH", filepath.Join(tempDir, "test.log"))
	defer t.Setenv("FLOW_LOG_PATH", originalLogPath)

	// Create a log file with malformed lines
	logPath, err := GetLogPath(time.Now())
	if err != nil {
		t.Fatalf("Failed to get log path: %v", err)
	}

	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Write test entries with malformed lines
	file, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	now := time.Now()
	entry := LogEntry{
		Tag:       "test session",
		StartTime: now.Add(-1 * time.Hour),
		EndTime:   now,
		Duration:  1 * time.Hour,
	}

	// Write malformed line, valid entry, another malformed line
	if _, err := file.WriteString("invalid json line\n"); err != nil {
		t.Fatalf("Failed to write invalid line: %v", err)
	}
	data, _ := json.Marshal(entry)
	if _, err := file.WriteString(string(data) + "\n"); err != nil {
		t.Fatalf("Failed to write valid entry: %v", err)
	}
	if _, err := file.WriteString("another invalid line\n"); err != nil {
		t.Fatalf("Failed to write invalid line: %v", err)
	}

	// Try to delete the valid entry
	err = DeleteLogEntry(entry)
	if err != nil {
		t.Errorf("Failed to delete entry from file with malformed lines: %v", err)
	}

	// Verify the entry was deleted
	remainingEntries, err := GetRecentSessions(10)
	if err != nil {
		t.Errorf("Failed to get remaining entries: %v", err)
		return
	}

	if len(remainingEntries) != 0 {
		t.Errorf("Expected 0 remaining entries, got %d", len(remainingEntries))
	}
}

func TestDeleteLogEntry_FileNotFound(t *testing.T) {
	// Create temporary directory for test log files
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Set environment variable to use our test directory
	originalLogPath := os.Getenv("FLOW_LOG_PATH")
	t.Setenv("FLOW_LOG_PATH", filepath.Join(tempDir, "test.log"))
	defer t.Setenv("FLOW_LOG_PATH", originalLogPath)

	// Try to delete from a non-existent file
	entryToDelete := LogEntry{
		Tag:       "test session",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(1 * time.Hour),
		Duration:  1 * time.Hour,
	}

	err := DeleteLogEntry(entryToDelete)
	if err == nil {
		t.Error("Expected error when deleting from non-existent file, got nil")
	}
}

func TestDeleteLogEntry_DuplicateEntries(t *testing.T) {
	// Create temporary directory for test log files
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Set environment variable to use our test directory
	originalLogPath := os.Getenv("FLOW_LOG_PATH")
	t.Setenv("FLOW_LOG_PATH", filepath.Join(tempDir, "test.log"))
	defer t.Setenv("FLOW_LOG_PATH", originalLogPath)

	// Create a log file with duplicate entries
	logPath, err := GetLogPath(time.Now())
	if err != nil {
		t.Fatalf("Failed to get log path: %v", err)
	}

	// Ensure log directory exists
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	// Write test entries with duplicates
	file, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Errorf("Failed to close file: %v", closeErr)
		}
	}()

	now := time.Now()
	entry := LogEntry{
		Tag:       "duplicate session",
		StartTime: now.Add(-1 * time.Hour),
		EndTime:   now,
		Duration:  1 * time.Hour,
	}

	// Write the same entry twice
	data, _ := json.Marshal(entry)
	if _, err := file.WriteString(string(data) + "\n"); err != nil {
		t.Fatalf("Failed to write first entry: %v", err)
	}
	if _, err := file.WriteString(string(data) + "\n"); err != nil {
		t.Fatalf("Failed to write second entry: %v", err)
	}

	// Delete the entry (should delete both instances)
	err = DeleteLogEntry(entry)
	if err != nil {
		t.Errorf("Failed to delete duplicate entries: %v", err)
	}

	// Verify both entries were deleted
	remainingEntries, err := GetRecentSessions(10)
	if err != nil {
		t.Errorf("Failed to get remaining entries: %v", err)
		return
	}

	if len(remainingEntries) != 0 {
		t.Errorf("Expected 0 remaining entries, got %d", len(remainingEntries))
	}
}
