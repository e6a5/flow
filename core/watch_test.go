package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// captureStderr captures everything written to stderr during the execution of a function.
func captureStderr(f func()) string {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		// In a test context, if we can't copy the buffer, something is
		// very wrong, and we should fail the test.
		panic(fmt.Sprintf("failed to copy stderr buffer: %v", err))
	}
	os.Stderr = oldStderr

	return buf.String()
}

func TestHandleNoSession_NudgeLogic(t *testing.T) {
	// 1. First call, should set the timer but not nudge
	cfg := defaultConfig
	watcher := NewWatcher()

	watcher.handleNoSession(cfg)
	if watcher.noSessionSince.IsZero() {
		t.Fatal("expected noSessionSince to be set, but it was zero")
	}

	// 2. Second call, before idle time, should not nudge
	output := captureStderr(func() {
		watcher.handleNoSession(cfg)
	})
	if output != "" {
		t.Errorf("expected no output, but got %q", output)
	}

	// 3. Third call, after idle time, should nudge
	// Advance the timer manually
	watcher.noSessionSince = time.Now().Add(-(cfg.Watch.RemindAfterIdle + time.Second))
	output = captureStderr(func() {
		watcher.handleNoSession(cfg)
	})
	if !strings.Contains(output, "No active session") {
		t.Errorf("expected nudge for no session, but got %q", output)
	}

	// 4. Immediately after a nudge, the timer should be reset.
	// We check if it's recent (within 1 second)
	if time.Since(watcher.noSessionSince) > time.Second {
		t.Errorf("expected noSessionSince to be reset, but it was not")
	}
}

func TestHandleActiveSession_BreakReminder(t *testing.T) {
	s := Session{StartTime: time.Now().Add(-3 * time.Hour)}
	cfg := Config{
		Watch: WatchConfig{
			RemindAfterActive: 2 * time.Hour,
		},
	}
	watcher := NewWatcher()
	// First call should produce a nudge
	output := captureStderr(func() {
		watcher.handleActiveSession(s, cfg)
	})
	if !strings.Contains(output, "Session active for over 2h") {
		t.Errorf("Expected output to contain break reminder, got %q", output)
	}

	// Immediate second call should not produce a nudge
	output = captureStderr(func() {
		watcher.handleActiveSession(s, cfg)
	})
	if output != "" {
		t.Errorf("Expected no output on second call, got %q", output)
	}
}

func TestHandlePausedSession_ResumeReminder(t *testing.T) {
	s := Session{IsPaused: true, PausedAt: time.Now().Add(-45 * time.Minute)}
	cfg := Config{
		Watch: WatchConfig{
			RemindAfterPause: 30 * time.Minute,
		},
	}
	watcher := NewWatcher()
	// First call should produce a nudge
	output := captureStderr(func() {
		watcher.handlePausedSession(s, cfg)
	})
	if !strings.Contains(output, "Session paused for over 30m") {
		t.Errorf("Expected output to contain 'Session paused for over 30m', got %q", output)
	}

	// Immediate second call should not produce a nudge
	output = captureStderr(func() {
		watcher.handlePausedSession(s, cfg)
	})
	if output != "" {
		t.Errorf("Expected no output on second call, got %q", output)
	}
}
