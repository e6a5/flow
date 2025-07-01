package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ANSI color codes for calm, minimal output
const (
	Reset = "\033[0m"
	Dim   = "\033[2m"
	Blue  = "\033[34m"
	Gray  = "\033[90m"
)

// Version information (set by build flags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Session represents a Flow work session
type Session struct {
	Tag         string        `json:"tag"`
	StartTime   time.Time     `json:"start_time"`
	PausedAt    time.Time     `json:"paused_at,omitempty"`
	IsPaused    bool          `json:"is_paused"`
	TotalPaused time.Duration `json:"total_paused"`
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "start":
		handleStart()
	case "status":
		handleStatus()
	case "pause":
		handlePause()
	case "resume":
		handleResume()
	case "end":
		handleEnd()
	case "--help", "-h":
		showUsage()
	case "--version", "-v":
		showVersion()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		showUsage()
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Print(`🌊 Flow

Protect your attention. Enter deep work.

USAGE:
  flow start [--tag "description"]   Start a deep work session
  flow status                        Check current session
  flow pause                         Pause current session
  flow resume                        Resume paused session
  flow end                          End current session

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

func handleStart() {
	// Check if session already exists
	if sessionExists() {
		session, err := loadSession()
		if err == nil {
			if session.IsPaused {
				fmt.Printf("🌊 You have a paused session: %s\n", session.Tag)
				fmt.Printf("Use 'flow resume' to continue or 'flow end' to finish.\n")
			} else {
				duration := time.Since(session.StartTime) - session.TotalPaused
				fmt.Printf("🌊 Already in deep work: %s\n", session.Tag)
				fmt.Printf("Working for %s. Use 'flow end' to complete.\n", formatDuration(duration))
			}
			fmt.Printf("\n%sOne thing at a time.%s\n", Dim, Reset)
			return
		}
	}

	// Parse tag if provided
	tag := "Deep Work"
	if len(os.Args) >= 3 {
		if os.Args[2] == "--tag" && len(os.Args) >= 4 {
			tag = os.Args[3]
		} else if strings.HasPrefix(os.Args[2], "--tag=") {
			tag = strings.TrimPrefix(os.Args[2], "--tag=")
		}
	}

	// Create new session
	session := Session{
		Tag:       tag,
		StartTime: time.Now(),
		IsPaused:  false,
	}

	if err := saveSession(session); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting session: %v\n", err)
		os.Exit(1)
	}

	// Show mindful start
	fmt.Printf("\n🌊 Starting deep work: %s\n", tag)
	fmt.Printf("\n%s   Clear your mind%s\n", Dim, Reset)
	fmt.Printf("%s   Focus on what matters%s\n", Dim, Reset)
	fmt.Printf("%s   Let distractions pass%s\n", Dim, Reset)
	fmt.Printf("\nSession active in background.\n")
	fmt.Printf("%sUse 'flow status' to check, 'flow end' to complete.%s\n\n", Gray, Reset)
}

func handleStatus() {
	if !sessionExists() {
		fmt.Printf("🌊 No active session.\n")
		fmt.Printf("Use 'flow start' to begin deep work.\n")
		return
	}

	session, err := loadSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading session: %v\n", err)
		os.Exit(1)
	}

	if session.IsPaused {
		pausedDuration := time.Since(session.PausedAt)
		fmt.Printf("⏸️  Session paused: %s\n", session.Tag)
		fmt.Printf("Paused for %s. Use 'flow resume' to continue.\n", formatDuration(pausedDuration))
	} else {
		activeDuration := time.Since(session.StartTime) - session.TotalPaused
		fmt.Printf("🌊 Deep work: %s\n", session.Tag)
		fmt.Printf("Active for %s.\n", formatDuration(activeDuration))
	}
}

func handlePause() {
	if !sessionExists() {
		fmt.Printf("🌊 No active session to pause.\n")
		return
	}

	session, err := loadSession()
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

	if err := saveSession(session); err != nil {
		fmt.Fprintf(os.Stderr, "Error pausing session: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(session.StartTime) - session.TotalPaused
	fmt.Printf("⏸️  Paused: %s\n", session.Tag)
	fmt.Printf("Worked for %s. Use 'flow resume' when ready.\n", formatDuration(duration))
}

func handleResume() {
	if !sessionExists() {
		fmt.Printf("🌊 No session to resume.\n")
		return
	}

	session, err := loadSession()
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

	if err := saveSession(session); err != nil {
		fmt.Fprintf(os.Stderr, "Error resuming session: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("🌊 Resumed: %s\n", session.Tag)
	fmt.Printf("Continue your deep work.\n")
}

func handleEnd() {
	if !sessionExists() {
		fmt.Printf("🌊 No active session to end.\n")
		return
	}

	session, err := loadSession()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading session: %v\n", err)
		os.Exit(1)
	}

	totalDuration := time.Since(session.StartTime) - session.TotalPaused
	if session.IsPaused {
		totalDuration = session.PausedAt.Sub(session.StartTime) - session.TotalPaused
	}

	// Remove session file
	sessionPath := getSessionPath()
	os.Remove(sessionPath)

	fmt.Printf("✨ Session complete: %s\n", session.Tag)
	fmt.Printf("Total focus time: %s\n", formatDuration(totalDuration))
	fmt.Printf("\n%sCarry this focus forward.%s\n", Dim, Reset)
}

// Session file management
func getSessionPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".flow-session")
}

func sessionExists() bool {
	_, err := os.Stat(getSessionPath())
	return err == nil
}

func loadSession() (Session, error) {
	var session Session
	data, err := os.ReadFile(getSessionPath())
	if err != nil {
		return session, err
	}
	err = json.Unmarshal(data, &session)
	return session, err
}

func saveSession(session Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return os.WriteFile(getSessionPath(), data, 0644)
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
