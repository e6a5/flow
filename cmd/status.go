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

		session, err := core.LoadSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading session: %v\n", err)
			os.Exit(1)
		}

		if raw {
			fmt.Print(session.Tag)
			return
		}

		if session.IsPaused {
			pausedDuration := time.Since(session.PausedAt)
			fmt.Printf("⏸️  Session paused: %s\n", session.Tag)
			fmt.Printf("Paused for %s. Use 'flow resume' to continue.\n", core.FormatDuration(pausedDuration))
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
