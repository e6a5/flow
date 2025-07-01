package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

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
	fmt.Printf("\nDeep work session initiated.\n")
	fmt.Printf("%sUse 'flow status' to check, 'flow end' to complete.%s\n\n", Gray, Reset)

	runHook("on_start", session.Tag)
}

func handleStatus() {
	// Support --raw flag for script-friendly output
	if len(os.Args) > 2 && os.Args[2] == "--raw" {
		if sessionExists() {
			session, err := loadSession()
			if err == nil {
				fmt.Print(session.Tag)
			}
		}
		// Exit cleanly after raw output
		return
	}

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
	runHook("on_pause", session.Tag)
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
	runHook("on_resume", session.Tag)
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

	endTime := time.Now()
	totalDuration := time.Since(session.StartTime) - session.TotalPaused
	if session.IsPaused {
		totalDuration = session.PausedAt.Sub(session.StartTime) - session.TotalPaused
		endTime = session.PausedAt
	}

	// Log the completed session before removing the session file
	logEntry := LogEntry{
		Tag:         session.Tag,
		StartTime:   session.StartTime,
		EndTime:     endTime,
		Duration:    totalDuration,
		TotalPaused: session.TotalPaused,
	}

	if err := logSession(logEntry); err != nil {
		// Don't fail the session end if logging fails, just warn
		fmt.Fprintf(os.Stderr, "Warning: failed to log session: %v\n", err)
	}

	// Remove session file
	sessionPath, err := getSessionPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error determining session path: %v\n", err)
		os.Exit(1)
	}
	if err := os.Remove(sessionPath); err != nil {
		// This is not a critical error, so we'll just warn the user.
		fmt.Fprintf(os.Stderr, "Warning: could not remove session file: %v\n", err)
	}

	fmt.Printf("✨ Session complete: %s\n", session.Tag)
	fmt.Printf("Total focus time: %s\n", formatDuration(totalDuration))
	fmt.Printf("\n%sCarry this focus forward.%s\n", Dim, Reset)
	runHook("on_end", session.Tag)
}
