package core

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

func ShowVersion(version, commit, date string) {
	fmt.Printf("Flow %s\n", version)
	if commit != "none" {
		fmt.Printf("Commit: %s\n", commit)
	}
	if date != "unknown" {
		fmt.Printf("Built: %s\n", date)
	}
}

func FormatDuration(d time.Duration) string {
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
