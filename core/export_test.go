package core

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestExportCSV(t *testing.T) {
	entry1 := LogEntry{
		Tag:       "Test 1",
		StartTime: time.Now(),
		Duration:  time.Hour,
	}
	entry2 := LogEntry{
		Tag:       "Test 2",
		StartTime: time.Now().Add(-24 * time.Hour),
		Duration:  30 * time.Minute,
	}
	entries := []LogEntry{entry1, entry2}

	var buf bytes.Buffer
	exportCSV(&buf, entries)

	output := buf.String()
	expectedHeader := "tag,start_time,end_time,duration_seconds,total_paused_seconds,duration_formatted,total_paused_formatted"
	if !strings.HasPrefix(output, expectedHeader) {
		t.Errorf("Expected CSV header '%s', got '%s'", expectedHeader, output)
	}
	if strings.Count(output, "\n") != 3 { // Header + 2 rows
		t.Errorf("Expected 3 lines in CSV output, got %d", strings.Count(output, "\n"))
	}
}

func TestExportJSON(t *testing.T) {
	entry1 := LogEntry{
		Tag:       "Test 1",
		StartTime: time.Now(),
		Duration:  time.Hour,
	}
	entries := []LogEntry{entry1}

	var buf bytes.Buffer
	exportJSON(&buf, entries)

	var decoded []LogEntry
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(decoded) != 1 || decoded[0].Tag != "Test 1" {
		t.Errorf("JSON output did not match expected output")
	}
}

func TestExportJSONEmpty(t *testing.T) {
	var entries []LogEntry // Empty slice
	var buf bytes.Buffer
	exportJSON(&buf, entries)

	output := buf.String()
	if strings.TrimSpace(output) != "[]" {
		t.Errorf("Expected empty JSON array, got '%s'", output)
	}
}
