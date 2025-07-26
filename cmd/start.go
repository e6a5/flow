package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Begin a deep work session",
	Long: `Starts a new deep work session.

A session is a single, uninterrupted period of focus.
You can add a descriptive tag to your session to remember what you worked on.
If a session is already active, 'start' will show you the status instead.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		config, err := core.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Check if session already exists
		if core.SessionExists() {
			session, err := core.LoadSession()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading existing session: %v\n", err)
				os.Exit(1)
			}

			// Check if the session is stale (running for too long)
			if core.IsSessionStale(session, config.ParsedStaleSessionThreshold()) {
				duration := time.Since(session.StartTime) - session.TotalPaused
				if session.IsPaused {
					duration = session.PausedAt.Sub(session.StartTime) - session.TotalPaused
				}

				// Automatically clean up the stale session
				if err := core.CleanupStaleSession(session, true); err != nil {
					fmt.Fprintf(os.Stderr, "Error cleaning up stale session: %v\n", err)
					os.Exit(1)
				}

				thresholdStr := core.FormatDuration(config.ParsedStaleSessionThreshold())
				fmt.Printf("⚠️  Found and cleaned up a stale session: %s\n", session.Tag)
				fmt.Printf("   Duration: %s (logged as abandoned)\n", core.FormatDuration(duration))
				fmt.Printf("   Threshold: %s\n", thresholdStr)
				fmt.Printf("   Starting fresh session...\n\n")
			} else {
				// Normal existing session (not stale)
				if session.IsPaused {
					fmt.Printf("🌊 You have a paused session: %s\n", session.Tag)
					fmt.Printf("Use 'flow resume' to continue or 'flow end' to finish.\n")
				} else {
					duration := time.Since(session.StartTime) - session.TotalPaused
					fmt.Printf("🌊 Already in deep work: %s\n", session.Tag)
					fmt.Printf("Working for %s. Use 'flow end' to complete.\n", core.FormatDuration(duration))
				}
				fmt.Printf("\n%sOne thing at a time.%s\n", core.Dim, core.Reset)
				return
			}
		}

		tag, _ := cmd.Flags().GetString("tag")
		targetStr, _ := cmd.Flags().GetString("target")
		var targetDuration time.Duration
		if targetStr != "" {
			var err error
			targetDuration, err = time.ParseDuration(targetStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Invalid duration format for --target: %v\n", err)
				os.Exit(1)
			}
		}

		// Create new session
		session := core.Session{
			Tag:            tag,
			StartTime:      time.Now(),
			IsPaused:       false,
			TargetDuration: targetDuration,
		}

		if err := core.SaveSession(session); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting session: %v\n", err)
			os.Exit(1)
		}

		// Show mindful start
		fmt.Printf("\n🌊 Starting deep work: %s\n", tag)
		fmt.Printf("\n%s   Clear your mind%s\n", core.Dim, core.Reset)
		fmt.Printf("%s   Focus on what matters%s\n", core.Dim, core.Reset)
		fmt.Printf("%s   Let distractions pass%s\n", core.Dim, core.Reset)
		fmt.Printf("\nDeep work session initiated.\n")
		fmt.Printf("%sUse 'flow status' to check, 'flow end' to complete.%s\n\n", core.Gray, core.Reset)

		core.RunHook("on_start", session.Tag)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("tag", "t", "Deep Work", "A description of the work session")
	startCmd.Flags().String("target", "", "Set a target duration for the session (e.g., '1h30m', '2h')")
}
