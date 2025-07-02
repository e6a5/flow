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
	Long: `Provides details about the currently active or paused deep work session.
Includes the session tag and the elapsed time.
The --raw flag can be used to output only the session tag for scripting purposes.`,
	Run: func(cmd *cobra.Command, args []string) {
		raw, _ := cmd.Flags().GetBool("raw")

		if raw {
			if core.SessionExists() {
				session, err := core.LoadSession()
				if err == nil {
					fmt.Print(session.Tag)
				}
			}
			return
		}

		if !core.SessionExists() {
			fmt.Printf("🌊 No active session.\n")
			fmt.Printf("Use 'flow start' to begin deep work.\n")
			return
		}

		session, err := core.LoadSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading session: %v\n", err)
			os.Exit(1)
		}

		if session.IsPaused {
			pausedDuration := time.Since(session.PausedAt)
			fmt.Printf("⏸️  Session paused: %s\n", session.Tag)
			fmt.Printf("Paused for %s. Use 'flow resume' to continue.\n", core.FormatDuration(pausedDuration))
		} else {
			activeDuration := time.Since(session.StartTime) - session.TotalPaused
			fmt.Printf("🌊 Deep work: %s\n", session.Tag)
			fmt.Printf("Active for %s.\n", core.FormatDuration(activeDuration))
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().Bool("raw", false, "Output only the session tag for scripting")
}
