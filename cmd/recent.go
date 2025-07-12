package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var recentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Show today's completed sessions",
	Long:  `Displays a summary of all deep work sessions completed today.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := core.NewLogReader()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating log reader: %v\n", err)
			return
		}

		// Read entries for today, with a reasonable limit for performance
		entries, err := reader.ReadRecentEntries(100, true, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading log entries: %v\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No sessions completed today. Keep up the focus!")
			return
		}

		fmt.Printf("✨ Today's Completed Sessions ✨\n\n")
		var totalTime time.Duration
		for _, entry := range entries {
			fmt.Printf("  - %s (%s)\n", entry.Tag, core.FormatDuration(entry.Duration))
			totalTime += entry.Duration
		}
		fmt.Printf("\nTotal focus time today: %s\n", core.FormatDuration(totalTime))
	},
}

func init() {
	rootCmd.AddCommand(recentCmd)
}
