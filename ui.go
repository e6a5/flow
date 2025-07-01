package main

import (
	"fmt"
	"time"
)

// ANSI color codes for calm, minimal output
const (
	Reset = "\033[0m"
	Dim   = "\033[2m"
	Blue  = "\033[34m"
	Gray  = "\033[90m"
)

func showUsage() {
	fmt.Print(`🌊 Flow

Protect your attention. Enter deep work.

USAGE:
  flow start [--tag "description"]   Start a deep work session
  flow status                        Check current session
  flow status [--raw]                Show only the session tag (for scripting)
  flow pause                         Pause current session
  flow resume                        Resume paused session
  flow end                          End current session
  flow completion [bash|zsh]         Output shell completion script

EXAMPLES:
  flow start --tag "writing docs"
  flow start                    # Start without tag
  flow status                   # Check what you're working on
  flow pause                    # Take a break
  flow resume                   # Continue working
  flow end                      # Complete session

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
