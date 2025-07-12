package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var insightsCmd = &cobra.Command{
	Use:   "insights",
	Short: "Show insights about your work patterns",
	Long:  `Analyzes your session history to provide insights, such as your most productive days and average session duration.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := core.NewLogReader()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating log reader: %v\n", err)
			return
		}

		// Read all entries for analysis
		entries, err := reader.ReadAllEntries()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading log entries: %v\n", err)
			return
		}

		if len(entries) < 10 { // Require a minimum amount of data for meaningful insights
			fmt.Printf("You have logged %d sessions. At least 10 are needed for meaningful insights. Keep up the great work!\n", len(entries))
			return
		}

		// Calculate insights
		report := calculateInsights(entries)

		// Display insights
		fmt.Printf("📊 Your Focus Insights (based on %d sessions)\n", report.TotalSessions)
		fmt.Println("----------------------------------------------------")
		fmt.Printf("Total Time Focused:     %s\n", core.FormatDuration(report.TotalTime))
		fmt.Printf("Average Session Length: %s\n\n", core.FormatDuration(report.AvgSessionLength))
		fmt.Printf("Busiest Day:            %s\n", report.BusiestDay)
		fmt.Printf("  - You focus an average of %s on %ss.\n", core.FormatDuration(report.BusiestDayAvg), report.BusiestDay)
		fmt.Printf("  - Your average on other days is %s.\n\n", core.FormatDuration(report.OtherDaysAvg))

		if len(report.TopActivities) > 0 {
			fmt.Println("Top Activities (by time):")
			for _, activity := range report.TopActivities {
				fmt.Printf("  - %-20s %-10s (%d%%)\n", activity.Tag, core.FormatDuration(activity.Duration), activity.Percent)
			}
		}
		fmt.Println("----------------------------------------------------")
	},
}

type InsightReport struct {
	TotalSessions    int
	TotalTime        time.Duration
	AvgSessionLength time.Duration
	BusiestDay       time.Weekday
	BusiestDayAvg    time.Duration
	OtherDaysAvg     time.Duration
	TopActivities    []ActivityStat
}

type ActivityStat struct {
	Tag      string
	Duration time.Duration
	Percent  int
}

func calculateInsights(entries []core.LogEntry) InsightReport {
	report := InsightReport{TotalSessions: len(entries)}
	if len(entries) == 0 {
		return report
	}

	dailyTotals := make(map[time.Weekday]time.Duration)
	dailyCounts := make(map[time.Weekday]int)
	tagTotals := make(map[string]time.Duration)

	for _, entry := range entries {
		report.TotalTime += entry.Duration
		dailyTotals[entry.EndTime.Weekday()] += entry.Duration
		dailyCounts[entry.EndTime.Weekday()]++
		tagTotals[entry.Tag] += entry.Duration
	}

	report.AvgSessionLength = report.TotalTime / time.Duration(len(entries))

	var maxDuration time.Duration
	for day, duration := range dailyTotals {
		if duration > maxDuration {
			maxDuration = duration
			report.BusiestDay = day
		}
	}

	busiestDayTotalTime := dailyTotals[report.BusiestDay]
	busiestDaySessionCount := dailyCounts[report.BusiestDay]
	if busiestDaySessionCount > 0 {
		report.BusiestDayAvg = busiestDayTotalTime / time.Duration(busiestDaySessionCount)
	}

	otherDaysTotalTime := report.TotalTime - busiestDayTotalTime
	otherDaysSessionCount := len(entries) - busiestDaySessionCount
	if otherDaysSessionCount > 0 {
		report.OtherDaysAvg = otherDaysTotalTime / time.Duration(otherDaysSessionCount)
	}

	// Calculate top activities
	type tagStatPair struct {
		tag      string
		duration time.Duration
	}
	var sortedTags []tagStatPair
	for tag, duration := range tagTotals {
		sortedTags = append(sortedTags, tagStatPair{tag, duration})
	}
	// Sort tags by duration descending
	sort.Slice(sortedTags, func(i, j int) bool {
		return sortedTags[i].duration > sortedTags[j].duration
	})

	// Get top 3 activities
	for i, pair := range sortedTags {
		if i >= 3 {
			break
		}
		percent := 0
		if report.TotalTime > 0 {
			percent = int((float64(pair.duration) / float64(report.TotalTime)) * 100)
		}
		report.TopActivities = append(report.TopActivities, ActivityStat{
			Tag:      pair.tag,
			Duration: pair.duration,
			Percent:  percent,
		})
	}

	return report
}

func init() {
	rootCmd.AddCommand(insightsCmd)
}
