package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createTestConfigFile(t *testing.T, content string) (string, func()) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "flow")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create temp config dir: %v", err)
	}

	path := filepath.Join(configDir, "config.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp config file: %v", err)
	}
	return path, func() {
		if removeErr := os.Remove(path); removeErr != nil {
			t.Errorf("Failed to remove test config file: %v", removeErr)
		}
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Temporarily move the real config file if it exists
	realConfigPath := ""
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		realConfigPath = filepath.Join(xdgConfigHome, "flow", "config.yml")
	} else {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			realConfigPath = filepath.Join(homeDir, ".config", "flow", "config.yml")
		}
	}

	// Move the real config file temporarily if it exists
	if realConfigPath != "" {
		if _, err := os.Stat(realConfigPath); err == nil {
			tempBackup := realConfigPath + ".testbackup"
			if err := os.Rename(realConfigPath, tempBackup); err == nil {
				defer func() {
					if restoreErr := os.Rename(tempBackup, realConfigPath); restoreErr != nil {
						t.Logf("Failed to restore config file: %v", restoreErr)
					}
				}()
			}
		}
	}

	// Clear any existing XDG_CONFIG_HOME to ensure we get defaults
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	if err := os.Unsetenv("XDG_CONFIG_HOME"); err != nil {
		t.Logf("Failed to unset XDG_CONFIG_HOME: %v", err)
	}
	defer func() {
		if err := os.Setenv("XDG_CONFIG_HOME", originalXDG); err != nil {
			t.Logf("Failed to restore XDG_CONFIG_HOME: %v", err)
		}
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg.ParsedStaleSessionThreshold() != 8*time.Hour {
		t.Errorf("expected stale session threshold to be %v, got %v", 8*time.Hour, cfg.ParsedStaleSessionThreshold())
	}
}

func TestLoadConfig_UserOverrides(t *testing.T) {
	content := `stale_session_threshold: "6h"`
	path, cleanup := createTestConfigFile(t, content)
	defer cleanup()

	// Temporarily set the config path to our test file
	if err := os.Setenv("XDG_CONFIG_HOME", filepath.Dir(filepath.Dir(path))); err != nil {
		t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg.ParsedStaleSessionThreshold() != 6*time.Hour {
		t.Errorf("expected stale session threshold to be %v, got %v", 6*time.Hour, cfg.ParsedStaleSessionThreshold())
	}
}

func TestLoadConfig_Partial(t *testing.T) {
	content := `stale_session_threshold: "4h"`
	path, cleanup := createTestConfigFile(t, content)
	defer cleanup()

	// Temporarily set the config path to our test file
	if err := os.Setenv("XDG_CONFIG_HOME", filepath.Dir(filepath.Dir(path))); err != nil {
		t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg.ParsedStaleSessionThreshold() != 4*time.Hour {
		t.Errorf("expected stale session threshold to be %v, got %v", 4*time.Hour, cfg.ParsedStaleSessionThreshold())
	}
}

func TestLoadConfig_Malformed(t *testing.T) {
	content := `stale_session_threshold: "invalid-duration"`
	path, cleanup := createTestConfigFile(t, content)
	defer cleanup()

	if err := os.Setenv("XDG_CONFIG_HOME", filepath.Dir(filepath.Dir(path))); err != nil {
		t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		// We don't expect an error from LoadConfig itself, as it should
		// gracefully handle a parsing error for a single field.
		t.Fatalf("LoadConfig() returned an unexpected error: %v", err)
	}

	// Should fall back to default when parsing fails
	if cfg.ParsedStaleSessionThreshold() != 8*time.Hour {
		t.Errorf("expected stale session threshold to fall back to default %v, got %v", 8*time.Hour, cfg.ParsedStaleSessionThreshold())
	}
}

func TestLoadConfig_MalformedYAML(t *testing.T) {
	content := `not: valid: yaml`
	path, cleanup := createTestConfigFile(t, content)
	defer cleanup()

	if err := os.Setenv("XDG_CONFIG_HOME", filepath.Dir(filepath.Dir(path))); err != nil {
		t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
	}

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("LoadConfig() should have failed for malformed YAML, but didn't")
	}
}
