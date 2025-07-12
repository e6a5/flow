// main_e2e_test.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var binaryPath string

// TestMain runs before all other tests in this file.
// It builds the binary for E2E testing.
func TestMain(m *testing.M) {
	var err error

	// Create a temporary directory for the binary
	tempDir, err := os.MkdirTemp("", "flow-e2e-tests")
	if err != nil {
		panic("failed to create temp dir for binary: " + err.Error())
	}
	defer func() {
		_ = os.RemoveAll(tempDir) // Ignore error in test cleanup
	}()

	binaryPath = filepath.Join(tempDir, "flow")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		panic("failed to build binary for E2E tests: " + err.Error() + "\nOutput:\n" + string(output))
	}

	// Run the tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func runFlowCommand(t *testing.T, args ...string) (string, string, error) {
	cmd := exec.Command(binaryPath, args...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func TestE2EWorkflow(t *testing.T) {
	// Create a temporary directory for the session file
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir) // Override home dir to control where session file is created

	// 1. Start a session
	stdout, stderr, err := runFlowCommand(t, "start", "--tag", "e2e test")
	if err != nil {
		t.Fatalf("Expected 'start' to succeed, but got error: %v\nStderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Starting deep work: e2e test") {
		t.Errorf("Expected start output to contain 'Starting deep work: e2e test', got:\n%s", stdout)
	}

	// 2. Check status
	stdout, stderr, err = runFlowCommand(t, "status")
	if err != nil {
		t.Fatalf("Expected 'status' to succeed, but got error: %v\nStderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Deep work: e2e test") {
		t.Errorf("Expected status output to contain 'Deep work: e2e test', got:\n%s", stdout)
	}

	// 3. Pause the session
	stdout, stderr, err = runFlowCommand(t, "pause")
	if err != nil {
		t.Fatalf("Expected 'pause' to succeed, but got error: %v\nStderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Paused session: e2e test") {
		t.Errorf("Expected pause output to contain 'Paused session: e2e test', got:\n%s", stdout)
	}

	// 4. Resume the session
	stdout, stderr, err = runFlowCommand(t, "resume")
	if err != nil {
		t.Fatalf("Expected 'resume' to succeed, but got error: %v\nStderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Resumed: e2e test") {
		t.Errorf("Expected resume output to contain 'Resumed: e2e test', got:\n%s", stdout)
	}

	// 5. End the session
	stdout, stderr, err = runFlowCommand(t, "end")
	if err != nil {
		t.Fatalf("Expected 'end' to succeed, but got error: %v\nStderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Session complete: e2e test") {
		t.Errorf("Expected end output to contain 'Session complete: e2e test', got:\n%s", stdout)
	}

	// 6. Check status again, should be no active session
	stdout, stderr, err = runFlowCommand(t, "status")
	if err != nil {
		t.Fatalf("Expected 'status' to succeed, but got error: %v\nStderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "No active session") {
		t.Errorf("Expected status output to contain 'No active session', got:\n%s", stdout)
	}
}

func TestE2ERawOutput(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	// 1. Check raw status with no session active
	stdout, _, err := runFlowCommand(t, "status", "--raw")
	if err != nil {
		t.Fatalf("Expected 'status --raw' to succeed, but got error: %v", err)
	}
	if stdout != "" {
		t.Errorf("Expected empty output for raw status with no session, got: %q", stdout)
	}

	// 2. Start a session
	_, _, err = runFlowCommand(t, "start", "--tag", "raw test")
	if err != nil {
		t.Fatalf("Failed to start session for raw test: %v", err)
	}

	// 3. Check raw status with an active session
	stdout, _, err = runFlowCommand(t, "status", "--raw")
	if err != nil {
		t.Fatalf("Expected 'status --raw' to succeed, but got error: %v", err)
	}
	if stdout != "raw test" {
		t.Errorf("Expected raw output to be 'raw test', got: %q", stdout)
	}

	// 4. End the session
	_, _, err = runFlowCommand(t, "end")
	if err != nil {
		t.Fatalf("Failed to end session for raw test: %v", err)
	}
}

func TestE2EHooks(t *testing.T) {
	// Create a temporary directory and force all configuration paths to use it.
	// This makes the test reliable regardless of the user's actual home directory
	// or how the CI environment is configured.
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)
	t.Setenv("XDG_CONFIG_HOME", tempDir)
	defer func() {
		t.Setenv("HOME", "")
		t.Setenv("XDG_CONFIG_HOME", "")
	}()

	hooksDir := filepath.Join(tempDir, "flow", "hooks")
	err := os.MkdirAll(hooksDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create hooks directory: %v", err)
	}

	// Create a hook script for the 'on_start' event
	hookScriptPath := filepath.Join(hooksDir, "on_start")
	hookOutputPath := filepath.Join(tempDir, "hook_output.txt")

	scriptContent := fmt.Sprintf("#!/bin/sh\necho \"Hook triggered for tag: $1\" > %s\n", hookOutputPath)
	err = os.WriteFile(hookScriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to write hook script: %v", err)
	}

	// Start a flow session, which should trigger the hook
	_, _, err = runFlowCommand(t, "start", "--tag", "hook test")
	if err != nil {
		t.Fatalf("Failed to start session for hook test: %v", err)
	}

	// Check if the hook was executed by reading the output file
	outputBytes, err := os.ReadFile(hookOutputPath)
	if err != nil {
		t.Fatalf("Failed to read hook output file: %v", err)
	}

	expectedOutput := "Hook triggered for tag: hook test\n"
	if string(outputBytes) != expectedOutput {
		t.Errorf("Expected hook output to be %q, got %q", expectedOutput, string(outputBytes))
	}

	// Clean up by ending the session
	_, _, err = runFlowCommand(t, "end")
	if err != nil {
		t.Logf("Failed to end session during hook test cleanup (non-critical): %v", err)
	}
}

func TestE2EDataCommands(t *testing.T) {
	// --- Setup ---
	tempDir := t.TempDir()
	// Use XDG_DATA_HOME to control where logs are stored, which is more robust
	// than relying on the HOME-based fallback.
	logDir := filepath.Join(tempDir, "flow", "logs")
	t.Setenv("XDG_DATA_HOME", tempDir)
	defer t.Setenv("XDG_DATA_HOME", "")

	// Create a mock log file with a known entry
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create log directory: %v", err)
	}

	now := time.Now()
	logFileName := fmt.Sprintf("%s_sessions.jsonl", now.Format("200601"))
	logFilePath := filepath.Join(logDir, logFileName)
	// Note: duration is in nanoseconds for JSON
	// Use today's date for the mock entry so --today filter will work
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	logContent := fmt.Sprintf(`{"tag":"e2e data test","start_time":"%s","end_time":"%s","duration":3600000000000,"total_paused":0}`, todayStart.Format(time.RFC3339), todayEnd.Format(time.RFC3339)) + "\n"
	err = os.WriteFile(logFilePath, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock log file: %v", err)
	}

	// --- Dashboard Test ---
	t.Run("dashboard", func(t *testing.T) {
		stdout, stderr, err := runFlowCommand(t, "dashboard")
		if err != nil {
			t.Fatalf("Expected 'dashboard' to succeed, but got error: %v\nStderr: %s", err, stderr)
		}

		if !strings.Contains(stdout, "Your Deep Work History") {
			t.Error("Dashboard output should contain title 'Your Deep Work History'")
		}
		if !strings.Contains(stdout, "Yearly Stats") {
			t.Error("Dashboard output should contain 'Yearly Stats'")
		}
	})

	// --- Export Test ---
	t.Run("export", func(t *testing.T) {
		// Test CSV export to stdout
		stdout, stderr, err := runFlowCommand(t, "export", "--format", "csv")
		if err != nil {
			t.Fatalf("Expected 'export' to succeed, but got error: %v\nStderr: %s", err, stderr)
		}

		// A simple check to ensure the header and a data row are present
		if !strings.Contains(stdout, "tag,start_time,end_time,duration_seconds") {
			t.Errorf("Expected CSV header in export output, got:\n%s", stdout)
		}
		if !strings.Contains(stdout, "e2e data test") {
			t.Errorf("Expected log entry in export output, got:\n%s", stdout)
		}
	})

	// --- Log Test ---
	t.Run("log", func(t *testing.T) {
		// Create a second log entry for a different day to test filtering
		yesterday := now.AddDate(0, 0, -1)
		yesterdayLogFileName := fmt.Sprintf("%s_sessions.jsonl", yesterday.Format("200601"))
		yesterdayLogFilePath := filepath.Join(logDir, yesterdayLogFileName)
		yesterdayLogContent := `{"tag":"yesterday's test","start_time":"` + yesterday.Format(time.RFC3339) + `","end_time":"` + yesterday.Format(time.RFC3339) + `","duration":1800000000000}` + "\n"

		// This check is necessary because if it's the first day of the month,
		// yesterday belongs to a different file. Otherwise, we append.
		if now.Day() == 1 {
			err = os.WriteFile(yesterdayLogFilePath, []byte(yesterdayLogContent), 0644)
		} else {
			f, err_append := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err_append != nil {
				t.Fatalf("Failed to open log file for appending: %v", err_append)
			}
			_, err = f.WriteString(yesterdayLogContent)
			_ = f.Close()
		}
		if err != nil {
			t.Fatalf("Failed to write second log entry: %v", err)
		}

		// Test `log --today`
		stdout, stderr, err := runFlowCommand(t, "log", "--today")
		if err != nil {
			t.Fatalf("Expected 'log --today' to succeed, but got error: %v\nStderr: %s", err, stderr)
		}

		if !strings.Contains(stdout, "e2e data test") {
			t.Errorf("Log output should contain today's entry ('e2e data test'), but it didn't.\nOutput: %s", stdout)
		}
		if strings.Contains(stdout, "yesterday's test") {
			t.Errorf("Log output should not contain yesterday's entry, but it did.\nOutput: %s", stdout)
		}
	})
}
