package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/e6a5/flow/core"
)

func TestRecentCmd(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tempDir)

	// Helper to redirect stdout
	captureOutput := func(f func()) string {
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		f()

		w.Close()
		os.Stdout = old
		var buf bytes.Buffer
		io.Copy(&buf, r)
		return buf.String()
	}

	// 1. Test with no log files
	t.Run("no logs", func(t *testing.T) {
		output := captureOutput(func() {
			recentCmd.Run(recentCmd, []string{})
		})
		if !strings.Contains(output, "No sessions completed today") {
			t.Errorf("Expected 'No sessions' message, got: %s", output)
		}
	})

	// 2. Setup some log entries
	now := time.Now()
	logEntries := []core.LogEntry{
		// Today's entry
		{Tag: "Work on feature A", Duration: 1 * time.Hour, EndTime: now},
		// Another entry for today
		{Tag: "Review PRs", Duration: 30 * time.Minute, EndTime: now.Add(-2 * time.Hour)},
		// Yesterday's entry (should not be shown)
		{Tag: "Plan week", Duration: 45 * time.Minute, EndTime: now.AddDate(0, 0, -1)},
	}

	// Manually create log file and write entries
	for _, entry := range logEntries {
		core.LogSession(entry)
	}

	// 3. Test with logs
	t.Run("with logs", func(t *testing.T) {
		output := captureOutput(func() {
			recentCmd.Run(recentCmd, []string{})
		})

		if !strings.Contains(output, "Work on feature A (1h 0m)") {
			t.Errorf("Expected to find 'Work on feature A', got: %s", output)
		}
		if !strings.Contains(output, "Review PRs (30m)") {
			t.Errorf("Expected to find 'Review PRs', got: %s", output)
		}
		if strings.Contains(output, "Plan week") {
			t.Errorf("Did not expect to find yesterday's entry 'Plan week', but got: %s", output)
		}
		if !strings.Contains(output, "Total focus time today: 1h 30m") {
			t.Errorf("Expected total time to be 1h 30m, got: %s", output)
		}
	})
}
