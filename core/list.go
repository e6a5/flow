package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// GetRecentSessions retrieves the most recent log entries.
func GetRecentSessions(limit int) ([]LogEntry, error) {
	logDir, err := GetLogDir()
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(logDir, "*_sessions.jsonl"))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return []LogEntry{}, nil
	}

	var allEntries []LogEntry
	totalLines := 0

	// Sort files in reverse order (newest first) for better performance with limits
	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j] // Lexicographically, newer YYYYMM comes after older
	})

	for _, file := range files {
		fileEntries, lines, err := readSingleFile(file)
		if err != nil {
			// Log error but continue with other files
			fmt.Fprintf(os.Stderr, "Warning: error reading %s: %v\n", file, err)
			continue
		}

		allEntries = append(allEntries, fileEntries...)
		totalLines += lines

		// If we have enough entries and not reading all, break early
		if limit > 0 && len(allEntries) >= limit {
			break
		}
	}

	// Sort all entries by end time (most recent first)
	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].EndTime.After(allEntries[j].EndTime)
	})

	// Apply limit after sorting
	if limit > 0 && len(allEntries) > limit {
		allEntries = allEntries[:limit]
	}

	return allEntries, nil
}

// readSingleFile reads entries from a single log file
func readSingleFile(filePath string) (entries []LogEntry, lineCount int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't return it as it's in a defer
			fmt.Fprintf(os.Stderr, "Warning: failed to close file %s: %v\n", filePath, closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Skip malformed lines
			continue
		}

		entries = append(entries, entry)
	}

	return entries, lineCount, scanner.Err()
}
