package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the current session status",
	Long: `Shows the status of the current deep work session.
This includes the session tag, how long it has been active, and progress if a target was set.
The --raw flag can be used to output only the session tag for scripting purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		raw, _ := cmd.Flags().GetBool("raw")

		if !core.SessionExists() {
			if raw {
				// Print nothing if no session exists and raw is requested
				return
			}
			fmt.Printf("🌊 No active session.\n")
			fmt.Printf("Use 'flow start' to begin deep work.\n")
			return
		}

		// Load configuration
		config, err := core.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		session, err := core.LoadSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading session: %v\n", err)
			os.Exit(1)
		}

		if raw {
			fmt.Print(session.Tag)
			return
		}

		// Check if session is stale and warn the user
		if core.IsSessionStale(session, config.ParsedStaleSessionThreshold()) {
			duration := time.Since(session.StartTime) - session.TotalPaused
			if session.IsPaused {
				duration = session.PausedAt.Sub(session.StartTime) - session.TotalPaused
			}

			thresholdStr := core.FormatDuration(config.ParsedStaleSessionThreshold())
			fmt.Printf("⚠️  WARNING: This session has been running for over %s!\n", thresholdStr)
			fmt.Printf("   Duration: %s\n", core.FormatDuration(duration))
			fmt.Printf("   You likely forgot to end the previous session.\n")
			fmt.Printf("   Run 'flow start' to automatically clean up and start fresh.\n\n")
		}

		if session.IsPaused {
			pausedDuration := time.Since(session.PausedAt)
			fmt.Printf("⏸️  Session paused: %s\n", session.Tag)

			// Calculate working time up to when the session was paused
			workingTime := session.PausedAt.Sub(session.StartTime) - session.TotalPaused
			if workingTime < 0 {
				workingTime = 0
			}

			fmt.Printf("Worked for %s • Paused for %s\n",
				core.FormatDuration(workingTime),
				core.FormatDuration(pausedDuration))
			fmt.Printf("Use 'flow resume' to continue or 'flow end' to finish.\n")
		} else {
			duration := time.Since(session.StartTime) - session.TotalPaused
			baseMsg := fmt.Sprintf("🌊 Deep work: %s (Active for %s)", session.Tag, core.FormatDuration(duration))
			if session.TargetDuration > 0 {
				// Adjust for pauses to get accurate end time
				effectiveEndTime := session.StartTime.Add(session.TargetDuration).Add(session.TotalPaused)
				remaining := time.Until(effectiveEndTime)

				// Don't show negative remaining time
				if remaining < 0 {
					remaining = 0
				}
				fmt.Printf("%s / %s (%s remaining)\n", baseMsg, core.FormatDuration(session.TargetDuration), core.FormatDuration(remaining))
			} else {
				fmt.Printf("%s\n", baseMsg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().Bool("raw", false, "Output only the session tag for scripting")
}
