package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	// Performance limits to prevent memory issues
	defaultMaxEntries = 10
	maxEntriesLimit   = 1000
	warningThreshold  = 10000
)

// LogReader provides efficient reading of log entries from partitioned files
type LogReader struct {
	logDir string
}

// NewLogReader creates a new log reader
func NewLogReader() (*LogReader, error) {
	logDir, err := GetLogDir()
	if err != nil {
		return nil, err
	}
	return &LogReader{logDir: logDir}, nil
}

// getRelevantLogFiles returns the list of log files to read based on filters
func (lr *LogReader) getRelevantLogFiles(filterToday, filterWeek, filterMonth bool, targetMonth ...time.Time) ([]string, error) {
	// Ensure logs directory exists
	if _, err := os.Stat(lr.logDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	files, err := filepath.Glob(filepath.Join(lr.logDir, "*_sessions.jsonl"))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return []string{}, nil
	}

	// If no filters, return all files
	if !filterToday && !filterWeek && !filterMonth && len(targetMonth) == 0 {
		return files, nil
	}

	now := time.Now()
	var relevantFiles []string

	for _, file := range files {
		basename := filepath.Base(file)
		// Extract YYYYMM from filename like "202507_sessions.jsonl"
		if len(basename) < 6 {
			continue
		}
		monthStr := basename[:6] // Extract YYYYMM

		// Parse the month
		fileDate, err := time.Parse("200601", monthStr)
		if err != nil {
			continue // Skip malformed filenames
		}

		// Check if this file is relevant based on filters
		if len(targetMonth) > 0 {
			// Specific month requested
			if fileDate.Year() == targetMonth[0].Year() && fileDate.Month() == targetMonth[0].Month() {
				relevantFiles = append(relevantFiles, file)
			}
		} else if filterToday {
			// For today filter, include files from current month (filtering by actual date happens later)
			if fileDate.Year() == now.Year() && fileDate.Month() == now.Month() {
				relevantFiles = append(relevantFiles, file)
			}
		} else if filterWeek {
			// For week filter, include files from current and potentially previous month
			currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			previousMonth := currentMonth.AddDate(0, -1, 0)

			if (fileDate.Year() == currentMonth.Year() && fileDate.Month() == currentMonth.Month()) ||
				(fileDate.Year() == previousMonth.Year() && fileDate.Month() == previousMonth.Month()) {
				relevantFiles = append(relevantFiles, file)
			}
		} else if filterMonth {
			// Current month
			if fileDate.Year() == now.Year() && fileDate.Month() == now.Month() {
				relevantFiles = append(relevantFiles, file)
			}
		} else if filterWeek || filterToday {
			// For week/today filters, we might need current month and potentially previous month
			currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			previousMonth := currentMonth.AddDate(0, -1, 0)

			if (fileDate.Year() == currentMonth.Year() && fileDate.Month() == currentMonth.Month()) ||
				(fileDate.Year() == previousMonth.Year() && fileDate.Month() == previousMonth.Month()) {
				relevantFiles = append(relevantFiles, file)
			}
		}
	}

	return relevantFiles, nil
}

// ReadRecentEntries reads the most recent entries efficiently
func (lr *LogReader) ReadRecentEntries(limit int, filterToday, filterWeek bool) ([]LogEntry, error) {
	entries, err := lr.readEntries(limit, filterToday, filterWeek, false, nil)
	if err != nil {
		return nil, err
	}

	if filterToday || filterWeek {
		now := time.Now()
		filteredEntries := []LogEntry{}
		for _, entry := range entries {
			if filterToday && !isToday(entry.EndTime, now) {
				continue
			}
			if filterWeek && !isThisWeek(entry.EndTime, now) {
				continue
			}
			filteredEntries = append(filteredEntries, entry)
		}
		// Respect the original limit after filtering
		if limit > 0 && len(filteredEntries) > limit {
			return filteredEntries[:limit], nil
		}
		return filteredEntries, nil
	}

	return entries, nil
}

// ReadMonthEntries reads entries from a specific month
func (lr *LogReader) ReadMonthEntries(month time.Time, limit int) ([]LogEntry, error) {
	return lr.readEntries(limit, false, false, false, &month)
}

// ReadAllEntries reads all entries (use with caution for large datasets)
func (lr *LogReader) ReadAllEntries() ([]LogEntry, error) {
	return lr.readEntries(0, false, false, true, nil)
}

// readEntries is the internal method that handles all reading scenarios
func (lr *LogReader) readEntries(limit int, filterToday, filterWeek, readAll bool, targetMonth *time.Time) ([]LogEntry, error) {
	if limit > maxEntriesLimit && !readAll {
		limit = maxEntriesLimit
	}

	// Determine which files to read
	var files []string
	var err error

	if targetMonth != nil {
		files, err = lr.getRelevantLogFiles(false, false, false, *targetMonth)
	} else {
		// Pass the filters to getRelevantLogFiles so it can select the right month files
		files, err = lr.getRelevantLogFiles(filterToday, filterWeek, false)
	}

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
		fileEntries, lines, err := lr.readSingleFile(file)
		if err != nil {
			// Log error but continue with other files
			fmt.Fprintf(os.Stderr, "Warning: error reading %s: %v\n", file, err)
			continue
		}

		allEntries = append(allEntries, fileEntries...)
		totalLines += lines

		// If we have enough entries and not reading all, break early
		if !readAll && limit > 0 && len(allEntries) >= limit {
			break
		}
	}

	// Warn if we have a very large number of entries
	if totalLines > warningThreshold {
		fmt.Fprintf(os.Stderr, "⚠️  Large dataset detected (%d entries across %d files). Consider using more specific filters.\n", totalLines, len(files))
	}

	// Sort all entries by end time (most recent first)
	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].EndTime.After(allEntries[j].EndTime)
	})

	// Apply limit after sorting
	if !readAll && limit > 0 && len(allEntries) > limit {
		allEntries = allEntries[:limit]
	}

	return allEntries, nil
}

// readSingleFile reads entries from a single log file
func (lr *LogReader) readSingleFile(filePath string) (entries []LogEntry, lineCount int, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if closeErr := file.Close(); err == nil {
			err = closeErr
		}
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineCount++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var entry LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Skip malformed lines but continue processing
			continue
		}

		entries = append(entries, entry)
	}

	return entries, lineCount, scanner.Err()
}

// LogStats contains aggregated statistics
type LogStats struct {
	TotalTime     time.Duration
	TotalSessions int
	AverageTime   time.Duration
	TopActivities []ActivityStat
	DateRange     string
}

// ActivityStat represents statistics for a specific activity
type ActivityStat struct {
	Tag      string
	Duration time.Duration
	Count    int
}

// CalculateStats computes statistics from log entries
func CalculateStats(entries []LogEntry) LogStats {
	if len(entries) == 0 {
		return LogStats{}
	}

	stats := LogStats{
		TotalSessions: len(entries),
	}

	tagCounts := make(map[string]int)
	tagTimes := make(map[string]time.Duration)
	var earliest, latest time.Time

	for i, entry := range entries {
		stats.TotalTime += entry.Duration
		tagCounts[entry.Tag]++
		tagTimes[entry.Tag] += entry.Duration

		// Track date range
		if i == 0 {
			earliest = entry.EndTime
			latest = entry.EndTime
		} else {
			if entry.EndTime.Before(earliest) {
				earliest = entry.EndTime
			}
			if entry.EndTime.After(latest) {
				latest = entry.EndTime
			}
		}
	}

	stats.AverageTime = stats.TotalTime / time.Duration(stats.TotalSessions)

	// Set date range
	if earliest.Format("2006-01-02") == latest.Format("2006-01-02") {
		stats.DateRange = earliest.Format("Jan 2, 2006")
	} else {
		stats.DateRange = fmt.Sprintf("%s - %s",
			earliest.Format("Jan 2"), latest.Format("Jan 2, 2006"))
	}

	// Calculate top activities
	for tag, duration := range tagTimes {
		stats.TopActivities = append(stats.TopActivities, ActivityStat{
			Tag:      tag,
			Duration: duration,
			Count:    tagCounts[tag],
		})
	}

	// Sort by total time
	sort.Slice(stats.TopActivities, func(i, j int) bool {
		return stats.TopActivities[i].Duration > stats.TopActivities[j].Duration
	})

	// Limit to top 10
	if len(stats.TopActivities) > 10 {
		stats.TopActivities = stats.TopActivities[:10]
	}

	return stats
}

// HandleLog handles the log command with improved performance
func HandleLog(showStats, filterToday, filterWeek, filterMonth, showAll bool, monthStr string) {
	reader, err := NewLogReader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log reader: %v\n", err)
		os.Exit(1)
	}

	var targetMonth *time.Time
	if monthStr != "" {
		if month, err := time.Parse("2006-01", monthStr); err == nil {
			targetMonth = &month
		} else {
			// Handle case where monthStr is not a valid date, e.g. from args
			if t, err := time.Parse("2006-01-02", monthStr); err == nil {
				targetMonth = &t
			} else {
				fmt.Fprintf(os.Stderr, "Error: Invalid month format '%s'. Please use YYYY-MM.\n", monthStr)
				os.Exit(1)
			}
		}
	}

	var entries []LogEntry
	maxEntries := defaultMaxEntries // Use default unless --all is specified

	if targetMonth != nil {
		entries, err = reader.ReadMonthEntries(*targetMonth, maxEntries)
	} else if showAll {
		var allEntries []LogEntry
		allEntries, err = reader.ReadAllEntries()
		if err == nil && (filterToday || filterWeek || filterMonth) {
			// If --all is combined with filters, apply them after loading
			now := time.Now()
			entries = []LogEntry{} // Reset entries to fill with filtered results
			for _, entry := range allEntries {
				if filterToday && !isToday(entry.EndTime, now) {
					continue
				}
				if filterWeek && !isThisWeek(entry.EndTime, now) {
					continue
				}
				if filterMonth && (entry.EndTime.Year() != now.Year() || entry.EndTime.Month() != now.Month()) {
					continue
				}
				entries = append(entries, entry)
			}
		} else {
			entries = allEntries
		}
	} else {
		// Default path for recent entries, passing filters
		entries, err = reader.ReadRecentEntries(maxEntries, filterToday, filterWeek)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading log entries: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("No sessions logged for the selected period. Use 'flow start' to begin.")
		return
	}

	if showStats {
		displayStats(entries, filterToday, filterWeek, filterMonth, targetMonth)
	} else {
		displayEntries(entries, filterToday, filterWeek, filterMonth, targetMonth, showAll)
	}
}

// displayEntries shows session entries in a user-friendly format
func displayEntries(entries []LogEntry, filterToday, filterWeek, filterMonth bool, targetMonth *time.Time, showAll bool) {
	// Determine header
	period := "Recent sessions"
	if targetMonth != nil {
		period = fmt.Sprintf("%s sessions", targetMonth.Format("January 2006"))
	} else if filterToday {
		period = "Today's sessions"
	} else if filterWeek {
		period = "This week's sessions"
	} else if filterMonth {
		period = "This month's sessions"
	} else if showAll {
		period = "All sessions"
	}

	fmt.Printf("🌊 %s:\n\n", period)

	// Display entries
	for _, entry := range entries {
		date := entry.EndTime.Format("Jan 2")
		timeRange := fmt.Sprintf("%s-%s",
			entry.StartTime.Format("15:04"),
			entry.EndTime.Format("15:04"))

		fmt.Printf("%s %s %s %s\n",
			date,
			timeRange,
			FormatDuration(entry.Duration),
			entry.Tag)
	}

	// Show summary
	if len(entries) > 0 {
		totalTime := time.Duration(0)
		for _, entry := range entries {
			totalTime += entry.Duration
		}
		fmt.Printf("\n%sTotal: %s across %d sessions%s\n",
			Dim, FormatDuration(totalTime), len(entries), Reset)
	}
}

// displayStats shows statistical analysis
func displayStats(entries []LogEntry, filterToday, filterWeek, filterMonth bool, targetMonth *time.Time) {
	stats := CalculateStats(entries)

	// Header based on filter
	period := "All Time"
	if targetMonth != nil {
		period = targetMonth.Format("January 2006")
	} else if filterToday {
		period = "Today"
	} else if filterWeek {
		period = "This Week"
	} else if filterMonth {
		period = "This Month"
	}

	fmt.Printf("🌊 Deep Work Statistics (%s):\n\n", period)
	fmt.Printf("Total time:     %s\n", FormatDuration(stats.TotalTime))
	fmt.Printf("Sessions:       %d\n", stats.TotalSessions)
	fmt.Printf("Average:        %s per session\n", FormatDuration(stats.AverageTime))

	if stats.DateRange != "" {
		fmt.Printf("Date range:     %s\n", stats.DateRange)
	}

	// Show top activities
	if len(stats.TopActivities) > 1 {
		fmt.Printf("\nTop activities:\n")
		for i, activity := range stats.TopActivities {
			if i >= 5 {
				break
			}
			percentage := (float64(activity.Duration) / float64(stats.TotalTime)) * 100
			fmt.Printf("  %d. %s (%d sessions, %s, %.1f%%)\n", i+1, activity.Tag, activity.Count, FormatDuration(activity.Duration), percentage)
		}
	}
}

// Date filtering helper functions
func isToday(t, now time.Time) bool {
	ty, tm, td := t.Date()
	ny, nm, nd := now.Date()
	return ty == ny && tm == nm && td == nd
}

func isThisWeek(t, now time.Time) bool {
	// Get start of current week (Sunday)
	weekday := int(now.Weekday()) // Sunday = 0, Monday = 1, etc.
	weekStart := now.AddDate(0, 0, -weekday)
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	// Get end of current week (Saturday)
	weekEnd := weekStart.AddDate(0, 0, 6)
	weekEnd = time.Date(weekEnd.Year(), weekEnd.Month(), weekEnd.Day(), 23, 59, 59, 999999999, weekEnd.Location())

	return (t.After(weekStart) || t.Equal(weekStart)) && (t.Before(weekEnd) || t.Equal(weekEnd))
}
