package core

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func HandleExport() {
	// --- Argument Parsing ---
	var filterToday, filterWeek, filterMonth, showAll bool
	var targetMonth *time.Time
	format := "csv"  // Default format
	outputFile := "" // Default to stdout

	// Manually parse flags to handle values like --format=csv and --output=file.csv
	args := os.Args[2:] // Skip "flow" and "export"
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--today":
			filterToday = true
		case arg == "--week":
			filterWeek = true
		case arg == "--month":
			filterMonth = true
		case arg == "--all":
			showAll = true
		case strings.HasPrefix(arg, "--format="):
			format = strings.TrimPrefix(arg, "--format=")
		case arg == "--format":
			if i+1 < len(args) {
				format = args[i+1]
				i++
			}
		case strings.HasPrefix(arg, "--output="):
			outputFile = strings.TrimPrefix(arg, "--output=")
		case arg == "--output":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		default:
			if t, err := time.Parse("2006-01", arg); err == nil {
				targetMonth = &t
			}
		}
	}

	// --- Data Fetching ---
	reader, err := NewLogReader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log reader: %v\n", err)
		return
	}

	var entries []LogEntry
	if showAll {
		entries, err = reader.ReadAllEntries()
	} else if targetMonth != nil {
		entries, err = reader.ReadMonthEntries(*targetMonth, 0)
	} else if filterToday || filterWeek || filterMonth {
		// ReadRecentEntries can handle these filters. Limit 0 means no limit.
		entries, err = reader.ReadRecentEntries(0, filterToday, filterWeek)
		if filterMonth && !filterToday && !filterWeek { // handle --month separately
			now := time.Now()
			entries, err = reader.ReadMonthEntries(now, 0)
		}
	} else {
		// Default: read recent entries (same as `flow log`)
		entries, err = reader.ReadRecentEntries(defaultMaxEntries, false, false)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading log entries: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stderr, "No log entries found for the selected period.")
		return
	}

	// --- Output Handling ---
	var writer io.Writer = os.Stdout
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing output file: %v\n", err)
			}
		}()
		writer = file
		fmt.Fprintf(os.Stderr, "Exporting %d entries to %s...\n", len(entries), outputFile)
	}

	switch format {
	case "csv":
		exportCSV(writer, entries)
	case "json":
		exportJSON(writer, entries)
	default:
		fmt.Fprintf(os.Stderr, "Unknown format: %s. Supported formats are csv, json.\n", format)
	}
}

func exportCSV(writer io.Writer, entries []LogEntry) {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	headers := []string{
		"tag", "start_time", "end_time", "duration_seconds",
		"total_paused_seconds", "duration_formatted", "total_paused_formatted",
	}
	if err := csvWriter.Write(headers); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing CSV header: %v\n", err)
		return
	}

	// Write rows
	for _, entry := range entries {
		row := []string{
			entry.Tag,
			entry.StartTime.Format(time.RFC3339),
			entry.EndTime.Format(time.RFC3339),
			strconv.FormatInt(int64(entry.Duration.Seconds()), 10),
			strconv.FormatInt(int64(entry.TotalPaused.Seconds()), 10),
			FormatDuration(entry.Duration),
			FormatDuration(entry.TotalPaused),
		}
		if err := csvWriter.Write(row); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV row: %v\n", err)
			continue
		}
	}
}

func exportJSON(writer io.Writer, entries []LogEntry) {
	if entries == nil {
		if _, err := writer.Write([]byte("[]\n")); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing empty JSON: %v\n", err)
		}
		return
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	if err := encoder.Encode(entries); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
	}
}
