package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// runHook executes a script for a given event if it exists and is executable.
// The session tag is passed as the first argument to the script.
func runHook(eventName string, sessionTag string) {
	hookPath, err := getHookPath(eventName)
	if err != nil {
		// Silently fail if we can't even determine the hook path.
		// This is a power-user feature and should not crash the main app.
		return
	}

	// Check if the hook script exists and is executable.
	info, err := os.Stat(hookPath)
	if err != nil || info.IsDir() || info.Mode()&0111 == 0 {
		// The hook doesn't exist, isn't a file, or isn't executable.
		return
	}

	// Execute the hook script.
	cmd := exec.Command(hookPath, sessionTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // We run hooks on a best-effort basis. Ignore errors.
}

// getHookPath determines the path for a hook script based on the event name.
// Hooks are expected to be in ~/.config/flow/hooks/.
func getHookPath(eventName string) (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")

	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home directory: %w", err)
		}
		// In a test environment, HOME might be set, but os.UserHomeDir() might not respect it.
		// To make this testable, we'll allow HOME to override.
		if home := os.Getenv("HOME"); home != "" {
			homeDir = home
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, "flow", "hooks", eventName), nil
}
