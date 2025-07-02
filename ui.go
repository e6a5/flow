package main

import (
	"fmt"
	"time"
)

// ANSI color codes for calm, minimal output
const (
	Reset = "\033[0m"
	Bold  = "\033[1m"
	Dim   = "\033[2m"
	Gray  = "\033[90m"

	// High-contrast, 5-tier blue color scale for the dashboard
	// Using 256-color ANSI codes for better terminal compatibility.
	Color0 = "\033[38;5;250m" // Light Gray (for empty days)
	Blue1  = "\033[38;5;117m" // Lightest Blue
	Blue2  = "\033[38;5;75m"  // Light Blue
	Blue3  = "\033[38;5;33m"  // Medium Blue
	Blue4  = "\033[38;5;21m"  // Darkest Blue
)

func showUsage() {
	fmt.Print(`🌊 Flow: A Terminal-Based Tool for Deep Work

Protect your attention. Enter deep work.

USAGE:
  flow start [--tag "description"]   Start a deep work session
  flow status                        Check current session
  flow status [--raw]                Show only the session tag (for scripting)
  flow pause                         Pause current session
  flow resume                        Resume paused session
  flow end                          End current session
  	flow log [--today|--week|--month|--stats|--all]  Show session history
  flow dashboard                     Show a dashboard of your activity
  flow export [--format csv|json]    Export session history
  flow completion [bash|zsh]         Output shell completion script

EXAMPLES:
  flow start --tag "writing docs"
  flow start                    # Start without tag
  flow status                   # Check what you're working on
  flow pause                    # Take a break
  flow resume                   # Continue working
  flow end                      # Complete session
  flow log                      # Show recent sessions
  flow log --today              # Today's sessions only
  flow log --month              # This month's sessions
  flow log 2025-07              # Specific month (YYYY-MM)
  flow log --stats              # Show summary statistics
  flow dashboard                # See your progress visually
  flow export --format csv      # Export all sessions to CSV

One session at a time. No tracking. Pure focus.

`)
}

func showVersion() {
	fmt.Printf("Flow %s\n", version)
	if commit != "none" {
		fmt.Printf("Commit: %s\n", commit)
	}
	if date != "unknown" {
		fmt.Printf("Built: %s\n", date)
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	return fmt.Sprintf("%ds", int(d.Seconds()))
}
