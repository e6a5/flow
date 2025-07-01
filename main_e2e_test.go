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
	if !strings.Contains(stdout, "Paused: e2e test") {
		t.Errorf("Expected pause output to contain 'Paused: e2e test', got:\n%s", stdout)
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
