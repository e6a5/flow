package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused session",
	Long:  `Resumes a previously paused deep work session, restarting the timer.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !core.SessionExists() {
			fmt.Printf("🌊 No session to resume.\n")
			return
		}

		session, err := core.LoadSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading session: %v\n", err)
			os.Exit(1)
		}

		if !session.IsPaused {
			fmt.Printf("🌊 Session already active: %s\n", session.Tag)
			return
		}

		// Calculate total paused time
		session.TotalPaused += time.Since(session.PausedAt)
		session.IsPaused = false
		session.PausedAt = time.Time{}

		if err := core.SaveSession(session); err != nil {
			fmt.Fprintf(os.Stderr, "Error resuming session: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🌊 Resumed: %s\n", session.Tag)
		fmt.Printf("Continue your deep work.\n")
		core.RunHook("on_resume", session.Tag)
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
