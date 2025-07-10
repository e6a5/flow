package core

import (
	"fmt"
	"os"
	"time"
)

// Watcher holds the state for the session watcher.
type Watcher struct {
	noSessionSince      time.Time
	lastActiveNudgeTime time.Time
	lastPausedNudgeTime time.Time
}

// NewWatcher creates a new Watcher instance.
func NewWatcher() *Watcher {
	return &Watcher{}
}

// CheckSessionAndNudge evaluates the current session state and provides a reminder if necessary.
func (w *Watcher) CheckSessionAndNudge(cfg Config) {
	if SessionExists() {
		w.noSessionSince = time.Time{} // Reset timer when a session is active.
		session, err := LoadSession()
		if err != nil {
			return
		}

		if session.IsPaused {
			w.handlePausedSession(session, cfg)
		} else {
			w.handleActiveSession(session, cfg)
		}
	} else {
		// No session exists, reset the other timers.
		w.lastActiveNudgeTime = time.Time{}
		w.lastPausedNudgeTime = time.Time{}
		w.handleNoSession(cfg)
	}
}

func (w *Watcher) handleActiveSession(s Session, cfg Config) {
	if time.Since(s.StartTime) > cfg.Watch.RemindAfterActive {
		// Only nudge if we haven't nudged before, or if enough time has passed since the last nudge.
		if w.lastActiveNudgeTime.IsZero() || time.Since(w.lastActiveNudgeTime) > cfg.Watch.RemindAfterActive {
			printNudge(fmt.Sprintf("🏃 Session active for over %s. Time for a break?", FormatDuration(cfg.Watch.RemindAfterActive)))
			w.lastActiveNudgeTime = time.Now()
		}
	}
}

func (w *Watcher) handlePausedSession(s Session, cfg Config) {
	if time.Since(s.PausedAt) > cfg.Watch.RemindAfterPause {
		if w.lastPausedNudgeTime.IsZero() || time.Since(w.lastPausedNudgeTime) > cfg.Watch.RemindAfterPause {
			printNudge(fmt.Sprintf("🤔 Session paused for over %s. Ready to resume?", FormatDuration(cfg.Watch.RemindAfterPause)))
			w.lastPausedNudgeTime = time.Now()
		}
	}
}

func (w *Watcher) handleNoSession(cfg Config) {
	if w.noSessionSince.IsZero() {
		w.noSessionSince = time.Now()
		return
	}
	if time.Since(w.noSessionSince) > cfg.Watch.RemindAfterIdle {
		printNudge(fmt.Sprintf("💡 No active session for over %s. Ready to start one?", FormatDuration(cfg.Watch.RemindAfterIdle)))
		w.noSessionSince = time.Now() // Reset timer after nudging.
	}
}

func printNudge(message string) {
	fmt.Fprintf(os.Stderr, "[%s] %s\n", time.Now().Format("03:04 PM"), message)
}
