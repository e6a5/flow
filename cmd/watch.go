package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch the current session and provide gentle reminders",
	Long: `Runs in the foreground and periodically checks the session status.
It provides gentle, timestamped nudges to help you remember to start,
pause, resume, or end a session. Designed to be run in a separate,
dedicated terminal tab.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := core.LoadConfig()
		if err != nil {
			// If config fails to load, print a warning but continue with defaults.
			fmt.Fprintf(os.Stderr, "Warning: could not load config file: %v\n", err)
		}

		fmt.Printf("[%s] 🌊 Flow Watcher started. Checking every %s.\n", time.Now().Format("03:04 PM"), cfg.Watch.Interval)

		runOnce, _ := cmd.Flags().GetBool("_test_run_once")

		watcher := core.NewWatcher()
		for {
			watcher.CheckSessionAndNudge(cfg)

			if runOnce {
				break
			}
			time.Sleep(cfg.Watch.Interval)
		}
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().Bool("_test_run_once", false, "Run the watch loop only once for testing.")
	watchCmd.Flags().MarkHidden("_test_run_once")
}
