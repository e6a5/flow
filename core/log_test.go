package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogReaderPerformanceLimits(t *testing.T) {
	tempDir := t.TempDir()

	// Set environment variable for testing
	t.Setenv("XDG_DATA_HOME", tempDir)

	reader, err := NewLogReader()
	if err != nil {
		t.Fatalf("Failed to create log reader: %v", err)
	}

	// Test limit enforcement
	entries, err := reader.ReadRecentEntries(2000, false, false) // Above maxEntriesLimit
	if err != nil {
		t.Fatalf("Failed to read entries: %v", err)
	}

	// Should be capped at maxEntriesLimit even though we requested more
	if len(entries) > maxEntriesLimit {
		t.Errorf("Expected at most %d entries, got %d", maxEntriesLimit, len(entries))
	}
}

func TestLogSessionWithPartitioning(t *testing.T) {
	tempDir := t.TempDir()

	// Set environment variable for testing
	t.Setenv("XDG_DATA_HOME", tempDir)

	// Test logging sessions from different months
	jan2024 := LogEntry{
		Tag:       "January Work",
		StartTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC),
		Duration:  time.Hour,
	}

	feb2024 := LogEntry{
		Tag:       "February Work",
		StartTime: time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2024, 2, 15, 11, 0, 0, 0, time.UTC),
		Duration:  time.Hour,
	}

	// Log the sessions
	if err := LogSession(jan2024); err != nil {
		t.Fatalf("Failed to log January session: %v", err)
	}
	if err := LogSession(feb2024); err != nil {
		t.Fatalf("Failed to log February session: %v", err)
	}

	// Verify files were created with correct names
	logDir := filepath.Join(tempDir, "flow", "logs")
	janFile := filepath.Join(logDir, "202401_sessions.jsonl")
	febFile := filepath.Join(logDir, "202402_sessions.jsonl")

	if _, err := os.Stat(janFile); os.IsNotExist(err) {
		t.Errorf("Expected January file to be created: %s", janFile)
	}

	if _, err := os.Stat(febFile); os.IsNotExist(err) {
		t.Errorf("Expected February file to be created: %s", febFile)
	}

	// Verify content
	janData, err := os.ReadFile(janFile)
	if err != nil {
		t.Fatalf("Failed to read January file: %v", err)
	}

	if !contains(string(janData), "January Work") {
		t.Errorf("January file doesn't contain expected data")
	}

	febData, err := os.ReadFile(febFile)
	if err != nil {
		t.Fatalf("Failed to read February file: %v", err)
	}

	if !contains(string(febData), "February Work") {
		t.Errorf("February file doesn't contain expected data")
	}
}

func TestLogReaderMonthFiltering(t *testing.T) {
	tempDir := t.TempDir()
	logDir := filepath.Join(tempDir, "flow", "logs")

	// Set environment variable for testing
	t.Setenv("XDG_DATA_HOME", tempDir)

	// Create test files with sessions
	createTestMonthFile(t, logDir, "202401_sessions.jsonl", []LogEntry{
		{Tag: "Jan1", EndTime: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC), Duration: time.Hour},
		{Tag: "Jan2", EndTime: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), Duration: time.Hour},
	})

	createTestMonthFile(t, logDir, "202402_sessions.jsonl", []LogEntry{
		{Tag: "Feb1", EndTime: time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC), Duration: time.Hour},
		{Tag: "Feb2", EndTime: time.Date(2024, 2, 15, 10, 0, 0, 0, time.UTC), Duration: time.Hour},
	})

	reader, err := NewLogReader()
	if err != nil {
		t.Fatalf("Failed to create reader: %v", err)
	}

	// Test reading specific month
	jan2024 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	entries, err := reader.ReadMonthEntries(jan2024, 100)
	if err != nil {
		t.Fatalf("Failed to read January entries: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 January entries, got %d", len(entries))
	}

	for _, entry := range entries {
		if entry.Tag != "Jan1" && entry.Tag != "Jan2" {
			t.Errorf("Unexpected entry in January: %s", entry.Tag)
		}
	}

	// Test reading all entries
	allEntries, err := reader.ReadAllEntries()
	if err != nil {
		t.Fatalf("Failed to read all entries: %v", err)
	}

	if len(allEntries) != 4 {
		t.Errorf("Expected 4 total entries, got %d", len(allEntries))
	}
}

func TestCalculateStats(t *testing.T) {
	now := time.Now()
	entries := []LogEntry{
		{
			Tag:       "coding",
			StartTime: now.Add(-2 * time.Hour),
			EndTime:   now.Add(-90 * time.Minute),
			Duration:  30 * time.Minute,
		},
		{
			Tag:       "coding",
			StartTime: now.Add(-1 * time.Hour),
			EndTime:   now.Add(-30 * time.Minute),
			Duration:  30 * time.Minute,
		},
		{
			Tag:       "writing",
			StartTime: now.Add(-30 * time.Minute),
			EndTime:   now,
			Duration:  30 * time.Minute,
		},
	}

	stats := CalculateStats(entries)

	// Test basic statistics
	if stats.TotalSessions != 3 {
		t.Errorf("Expected 3 sessions, got %d", stats.TotalSessions)
	}

	expectedTotal := 90 * time.Minute
	if stats.TotalTime != expectedTotal {
		t.Errorf("Expected total time %v, got %v", expectedTotal, stats.TotalTime)
	}

	expectedAvg := 30 * time.Minute
	if stats.AverageTime != expectedAvg {
		t.Errorf("Expected average time %v, got %v", expectedAvg, stats.AverageTime)
	}

	// Test top activities
	if len(stats.TopActivities) != 2 {
		t.Errorf("Expected 2 top activities, got %d", len(stats.TopActivities))
	}

	// "coding" should be first (60 minutes total)
	if stats.TopActivities[0].Tag != "coding" {
		t.Errorf("Expected 'coding' to be top activity, got '%s'", stats.TopActivities[0].Tag)
	}

	if stats.TopActivities[0].Duration != 60*time.Minute {
		t.Errorf("Expected coding duration 60m, got %v", stats.TopActivities[0].Duration)
	}

	if stats.TopActivities[0].Count != 2 {
		t.Errorf("Expected coding count 2, got %d", stats.TopActivities[0].Count)
	}

	if stats.TotalTime != 4*time.Hour {
		t.Errorf("Expected total time %v, got %v", 4*time.Hour, stats.TotalTime)
	}
}

func TestEmptyStatsCalculation(t *testing.T) {
	stats := CalculateStats([]LogEntry{})

	if stats.TotalSessions != 0 {
		t.Errorf("Expected 0 sessions for empty input, got %d", stats.TotalSessions)
	}

	if stats.TotalTime != 0 {
		t.Errorf("Expected 0 total time for empty input, got %v", stats.TotalTime)
	}

	if len(stats.TopActivities) != 0 {
		t.Errorf("Expected 0 top activities for empty input, got %d", len(stats.TopActivities))
	}
}

func TestDateFilteringFunctions(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	lastWeek := today.AddDate(0, 0, -8)

	// Test isToday
	if !isToday(today, now) {
		t.Error("Expected today to be identified as today")
	}

	if isToday(yesterday, now) {
		t.Error("Expected yesterday not to be identified as today")
	}

	// Test isThisWeek
	if !isThisWeek(today, now) {
		t.Error("Expected today to be in this week")
	}

	if !isThisWeek(yesterday, now) {
		t.Error("Expected yesterday to be in this week")
	}

	if isThisWeek(lastWeek, now) {
		t.Error("Expected last week not to be in this week")
	}
}

func TestLogReaderWithMalformedData(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmpDir)
	logDir := filepath.Join(tmpDir, "flow", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}
	logFile := filepath.Join(logDir, "202301_sessions.jsonl")

	// Setup: create a log file with one valid and one malformed line
	malformedData := "this is not json\n"
	validEntry := LogEntry{Tag: "valid", Duration: time.Hour}
	validData, _ := json.Marshal(validEntry)
	content := string(validData) + "\n" + malformedData
	if err := os.WriteFile(logFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test log file: %v", err)
	}

	reader, err := NewLogReader()
	if err != nil {
		t.Fatalf("Failed to create new log reader: %v", err)
	}

	entries, err := reader.ReadAllEntries()
	if err != nil {
		t.Fatalf("ReadAllEntries returned an error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].Tag != "valid" {
		t.Errorf("Expected entry tag 'valid', got '%s'", entries[0].Tag)
	}
}

// Helper functions
func createTestMonthFile(t *testing.T, logDir, filename string, entries []LogEntry) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	file, err := os.Create(filepath.Join(logDir, filename))
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Errorf("Failed to close test file: %v", err)
		}
	}()

	for _, entry := range entries {
		data, err := json.Marshal(entry)
		if err != nil {
			t.Fatalf("Failed to marshal entry: %v", err)
		}
		if _, err := fmt.Fprintf(file, "%s\n", data); err != nil {
			t.Fatalf("Failed to write to test file: %v", err)
		}
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) &&
		(haystack == needle ||
			(len(haystack) > len(needle) &&
				(haystack[:len(needle)] == needle ||
					haystack[len(haystack)-len(needle):] == needle ||
					findInString(haystack, needle))))
}

func findInString(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
