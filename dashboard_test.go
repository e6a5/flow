package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDisplayDashboardStats(t *testing.T) {
	now := time.Date(2023, 10, 27, 10, 0, 0, 0, time.UTC)

	testCases := []struct {
		name           string
		dailyTotals    map[time.Time]time.Duration
		expectedOutput []string
	}{
		{
			name:        "No data",
			dailyTotals: map[time.Time]time.Duration{},
			expectedOutput: []string{
				"Total Focus Time: 0s",
				"Daily Average:    0s",
				"Current Streak:   0 days",
			},
		},
		{
			name: "Single entry today",
			dailyTotals: map[time.Time]time.Duration{
				time.Date(2023, 10, 27, 0, 0, 0, 0, time.UTC): 2 * time.Hour,
			},
			expectedOutput: []string{
				"Total Focus Time: 2h 0m",
				"Daily Average:    19s", // 2h / 365 days = 19.7s, formatted to 19s
				"Current Streak:   1 days",
			},
		},
		{
			name: "3-day streak ending yesterday",
			dailyTotals: map[time.Time]time.Duration{
				time.Date(2023, 10, 26, 0, 0, 0, 0, time.UTC): 1 * time.Hour,
				time.Date(2023, 10, 25, 0, 0, 0, 0, time.UTC): 1 * time.Hour,
				time.Date(2023, 10, 24, 0, 0, 0, 0, time.UTC): 1 * time.Hour,
			},
			expectedOutput: []string{
				"Total Focus Time: 3h 0m",
				"Current Streak:   0 days",
			},
		},
		{
			name: "3-day streak ending today",
			dailyTotals: map[time.Time]time.Duration{
				now.Truncate(24 * time.Hour):                   30 * time.Minute,
				now.AddDate(0, 0, -1).Truncate(24 * time.Hour): 1 * time.Hour,
				now.AddDate(0, 0, -2).Truncate(24 * time.Hour): 1 * time.Hour,
			},
			expectedOutput: []string{
				"Total Focus Time: 2h 30m",
				"Current Streak:   3 days",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			displayDashboardStats(tc.dailyTotals, now)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, r); err != nil {
				t.Fatalf("Failed to read captured output: %v", err)
			}
			output := buf.String()

			for _, expected := range tc.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', but it didn't.\nOutput:\n%s", expected, output)
				}
			}
		})
	}
}
