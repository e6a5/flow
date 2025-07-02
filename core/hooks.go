package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunHook executes a custom script for a given event
func RunHook(event string, args ...string) {
	hookPath, err := getHookScriptPath(event)
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
	cmd := exec.Command(hookPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // We run hooks on a best-effort basis. Ignore errors.
}

// getHookScriptPath finds the path for a given hook event
func getHookScriptPath(event string) (string, error) {
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

	return filepath.Join(configDir, "flow", "hooks", event), nil
}
