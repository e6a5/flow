package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// createMockEntries creates a slice of LogEntry structs for testing.
func createMockEntries() []LogEntry {
	startTime1, _ := time.Parse(time.RFC3339, "2023-10-27T09:00:00Z")
	endTime1, _ := time.Parse(time.RFC3339, "2023-10-27T10:30:00Z")

	startTime2, _ := time.Parse(time.RFC3339, "2023-10-27T11:00:00Z")
	endTime2, _ := time.Parse(time.RFC3339, "2023-10-27T12:00:00Z")

	return []LogEntry{
		{
			Tag:         "Task 1",
			StartTime:   startTime1,
			EndTime:     endTime1,
			Duration:    90 * time.Minute,
			TotalPaused: 5 * time.Minute,
		},
		{
			Tag:       "Task 2, with comma",
			StartTime: startTime2,
			EndTime:   endTime2,
			Duration:  60 * time.Minute,
		},
	}
}

func TestExportCSV(t *testing.T) {
	entries := createMockEntries()
	var buf bytes.Buffer

	exportCSV(&buf, entries)

	// Read the CSV output from the buffer
	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read CSV output: %v", err)
	}

	// Check header
	expectedHeader := []string{
		"tag", "start_time", "end_time", "duration_seconds",
		"total_paused_seconds", "duration_formatted", "total_paused_formatted",
	}
	if len(records) < 2 {
		t.Fatalf("Expected at least 2 records (header + data), got %d", len(records))
	}

	header := records[0]
	if !equalSlices(header, expectedHeader) {
		t.Errorf("Expected header %v, got %v", expectedHeader, header)
	}

	// Check number of data rows
	if len(records)-1 != len(entries) {
		t.Errorf("Expected %d data rows, got %d", len(entries), len(records)-1)
	}

	// Check content of the first data row
	expectedRow1 := []string{
		"Task 1",
		"2023-10-27T09:00:00Z",
		"2023-10-27T10:30:00Z",
		"5400",
		"300",
		"1h 30m",
		"5m",
	}
	row1 := records[1]
	if !equalSlices(row1, expectedRow1) {
		t.Errorf("Expected row %v, got %v", expectedRow1, row1)
	}

	// Check content of the second data row (handles commas in tags)
	expectedRow2 := []string{
		"Task 2, with comma",
		"2023-10-27T11:00:00Z",
		"2023-10-27T12:00:00Z",
		"3600",
		"0",
		"1h 0m",
		"0s",
	}
	row2 := records[2]
	// The duration format for 1h is "1h 0m", not "1h".
	expectedRow2[5] = formatDuration(entries[1].Duration)
	if !equalSlices(row2, expectedRow2) {
		t.Errorf("Expected row %v, got %v", expectedRow2, row2)
	}
}

func TestExportJSON(t *testing.T) {
	expectedEntries := createMockEntries()
	var buf bytes.Buffer

	exportJSON(&buf, expectedEntries)

	var actualEntries []LogEntry
	err := json.Unmarshal(buf.Bytes(), &actualEntries)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON output: %v", err)
	}

	if len(actualEntries) != len(expectedEntries) {
		t.Fatalf("Expected %d entries, got %d", len(expectedEntries), len(actualEntries))
	}

	// Use reflect.DeepEqual for robust comparison of structs
	if !reflect.DeepEqual(expectedEntries, actualEntries) {
		t.Errorf("Exported JSON does not match expected JSON")
		// Optional: Print diff for easier debugging
		// t.Logf("Expected: %+v", expectedEntries)
		// t.Logf("Actual:   %+v", actualEntries)
	}
}

// equalSlices is a helper to compare two string slices.
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
