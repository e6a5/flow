package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "Complete the session and log it",
	Long:  `Completes the current deep work session, logs the total focus time, and cleans up the session file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !core.SessionExists() {
			fmt.Printf("🌊 No active session to end.\n")
			return
		}

		session, err := core.LoadSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading session: %v\n", err)
			os.Exit(1)
		}

		endTime := time.Now()
		totalDuration := time.Since(session.StartTime) - session.TotalPaused
		if session.IsPaused {
			totalDuration = session.PausedAt.Sub(session.StartTime) - session.TotalPaused
			endTime = session.PausedAt
		}

		// Log the completed session before removing the session file
		logEntry := core.LogEntry{
			Tag:         session.Tag,
			StartTime:   session.StartTime,
			EndTime:     endTime,
			Duration:    totalDuration,
			TotalPaused: session.TotalPaused,
		}

		if err := core.LogSession(logEntry); err != nil {
			// Don't fail the session end if logging fails, just warn
			fmt.Fprintf(os.Stderr, "Warning: failed to log session: %v\n", err)
		}

		// Remove session file
		sessionPath, err := core.GetSessionPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining session path: %v\n", err)
			os.Exit(1)
		}
		if err := os.Remove(sessionPath); err != nil {
			// This is not a critical error, so we'll just warn the user.
			fmt.Fprintf(os.Stderr, "Warning: could not remove session file: %v\n", err)
		}

		fmt.Printf("✨ Session complete: %s\n", session.Tag)
		fmt.Printf("Total focus time: %s\n", core.FormatDuration(totalDuration))
		fmt.Printf("\n%sCarry this focus forward.%s\n", core.Dim, core.Reset)
		core.RunHook("on_end", session.Tag)
	},
}

func init() {
	rootCmd.AddCommand(endCmd)
}
