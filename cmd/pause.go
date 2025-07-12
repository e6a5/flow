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
			fmt.Println("No active session to pause. Use 'flow start' to begin.")
			return
		}

		session, err := core.LoadSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading session: %v\n", err)
			os.Exit(1)
		}

		if session.IsPaused {
			fmt.Printf("Session '%s' is already paused.\n", session.Tag)
			return
		}

		session.IsPaused = true
		session.PausedAt = time.Now()

		if err := core.SaveSession(session); err != nil {
			fmt.Fprintf(os.Stderr, "Error pausing session: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("⏸️  Paused session: %s\n", session.Tag)
		core.RunHook("on_pause", session.Tag)
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
