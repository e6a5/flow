package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the active session",
	Long:  `Pauses the currently active deep work session, freezing the timer.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !core.SessionExists() {
			fmt.Printf("🌊 No active session to pause.\n")
			return
		}

		session, err := core.LoadSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading session: %v\n", err)
			os.Exit(1)
		}

		if session.IsPaused {
			fmt.Printf("⏸️  Session already paused: %s\n", session.Tag)
			return
		}

		session.IsPaused = true
		session.PausedAt = time.Now()

		if err := core.SaveSession(session); err != nil {
			fmt.Fprintf(os.Stderr, "Error pausing session: %v\n", err)
			os.Exit(1)
		}

		duration := time.Since(session.StartTime) - session.TotalPaused
		fmt.Printf("⏸️  Paused: %s\n", session.Tag)
		fmt.Printf("Worked for %s. Use 'flow resume' when ready.\n", core.FormatDuration(duration))
		core.RunHook("on_pause", session.Tag)
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
