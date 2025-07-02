package core

import (
	"fmt"
	"os"
	"time"
)

func HandleDashboard() {
	reader, err := NewLogReader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log reader: %v\n", err)
		return
	}

	// Read all entries. We'll filter them by date later.
	entries, err := reader.ReadAllEntries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading log entries: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No sessions logged. Use 'flow start' to begin.")
		return
	}

	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)
	dailyTotals := make(map[time.Time]time.Duration)

	for _, entry := range entries {
		if !entry.EndTime.IsZero() && entry.EndTime.After(oneYearAgo) {
			day := time.Date(entry.EndTime.Year(), entry.EndTime.Month(), entry.EndTime.Day(), 0, 0, 0, 0, time.UTC)
			dailyTotals[day] += entry.Duration
		}
	}

	renderContributionGraph(dailyTotals)
	displayDashboardStats(dailyTotals, now)
}

func renderContributionGraph(dailyTotals map[time.Time]time.Duration) {
	now := time.Now()

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	weekday := int(today.Weekday())
	lastSunday := today.AddDate(0, 0, -weekday)
	graphStartDate := lastSunday.AddDate(0, 0, -(51 * 7))

	fmt.Printf("\n%sYour Deep Work History (Last Year)%s\n", Bold, Reset)

	// --- Header Row ---
	// Create a character buffer for the header to ensure perfect alignment.
	// 52 weeks * 2 chars/week ("■ ") = 104 chars wide.
	header := make([]rune, 104)
	for i := range header {
		header[i] = ' '
	}

	var lastMonth time.Month
	for week := 0; week < 52; week++ {
		// Use a representative day to find the month for this column.
		representativeDay := graphStartDate.AddDate(0, 0, week*7+3)
		month := representativeDay.Month()
		if month != lastMonth {
			monthLabel := month.String()[:3]
			// Place the label at the calculated position.
			position := week * 2
			for i, char := range monthLabel {
				if position+i < len(header) {
					header[position+i] = char
				}
			}
			lastMonth = month
		}
	}
	// Print the fully constructed header with padding for day labels.
	fmt.Printf("     %s\n", string(header))

	// --- Grid ---
	dayLabels := [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
		if dayOfWeek%2 != 0 {
			fmt.Printf("%-3s  ", dayLabels[dayOfWeek])
		} else {
			fmt.Printf("%-3s  ", " ")
		}

		for week := 0; week < 52; week++ {
			currentDay := graphStartDate.AddDate(0, 0, week*7+dayOfWeek)
			total := dailyTotals[currentDay]

			var color string
			switch {
			case total >= 6*time.Hour:
				color = Blue4 // Darkest Blue
			case total >= 4*time.Hour:
				color = Blue3 // Medium Blue
			case total >= 2*time.Hour:
				color = Blue2 // Light Blue
			case total > 0:
				color = Blue1 // Lightest Blue
			default:
				color = Color0 // Use high-contrast gray for empty
			}
			fmt.Printf("%s■ %s", color, Reset)
		}
		fmt.Println()
	}

	// --- Legend ---
	// A complete legend showing all 5 tiers from no activity to high activity.
	fmt.Printf("\n  Less %s■%s %s■%s %s■%s %s■%s %s■%s More\n",
		Color0, Reset,
		Blue1, Reset,
		Blue2, Reset,
		Blue3, Reset,
		Blue4, Reset,
	)
	fmt.Println()
}

func displayDashboardStats(dailyTotals map[time.Time]time.Duration, now time.Time) {
	oneYearAgo := now.AddDate(-1, 0, 0)

	var totalTime time.Duration
	for _, duration := range dailyTotals {
		totalTime += duration
	}

	var avgDailyTime time.Duration
	// Use 365 days for a true yearly average
	if totalTime > 0 {
		avgDailyTime = totalTime / 365
	}

	// Calculate streak
	var currentStreak int
	for i := 0; ; i++ {
		day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -i)
		if day.Before(oneYearAgo) {
			break
		}
		if dailyTotals[day] > 0 {
			currentStreak++
		} else {
			// A streak is broken by any day with no activity.
			break
		}
	}

	fmt.Printf("%sYearly Stats%s\n", Bold, Reset)
	fmt.Printf("  Total Focus Time: %s\n", FormatDuration(totalTime))
	fmt.Printf("  Daily Average:    %s\n", FormatDuration(avgDailyTime))
	fmt.Printf("  Current Streak:   %d days\n", currentStreak)
	fmt.Println()
}
