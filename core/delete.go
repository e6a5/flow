package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// DeleteLogEntry removes a specific log entry from the log files.
func DeleteLogEntry(entryToDelete LogEntry) error {
	logPath, err := GetLogPath(entryToDelete.EndTime)
	if err != nil {
		return err
	}

	file, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log the error but don't return it as it's in a defer
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
		}
	}()

	tempFile, err := os.CreateTemp(filepath.Dir(logPath), "temp_log_")
	if err != nil {
		return err
	}
	defer func() {
		if removeErr := os.Remove(tempFile.Name()); removeErr != nil {
			// Log the error but don't return it as it's in a defer
			fmt.Fprintf(os.Stderr, "Warning: failed to remove temp file: %v\n", removeErr)
		}
	}()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tempFile)
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Skip malformed lines
			continue
		}

		if entry.StartTime.Equal(entryToDelete.StartTime) && entry.Tag == entryToDelete.Tag {
			found = true
		} else {
			if _, writeErr := fmt.Fprintln(writer, line); writeErr != nil {
				if closeErr := tempFile.Close(); closeErr != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to close temp file: %v\n", closeErr)
				}
				return writeErr
			}
		}
	}

	if err := scanner.Err(); err != nil {
		if closeErr := tempFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close temp file: %v\n", closeErr)
		}
		return err
	}

	if err := writer.Flush(); err != nil {
		if closeErr := tempFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close temp file: %v\n", closeErr)
		}
		return err
	}

	if closeErr := tempFile.Close(); closeErr != nil {
		return closeErr
	}

	if found {
		return os.Rename(tempFile.Name(), logPath)
	}

	return fmt.Errorf("log entry not found")
}
