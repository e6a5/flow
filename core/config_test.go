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
	// Temporarily unset env vars to ensure we are testing defaults
	t.Setenv("XDG_CONFIG_HOME", "/tmp/non-existent-dir")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg.Watch.Interval != 5*time.Minute {
		t.Errorf("expected Interval to be %v, got %v", 5*time.Minute, cfg.Watch.Interval)
	}
	if cfg.Watch.RemindAfterIdle != 15*time.Minute {
		t.Errorf("expected RemindAfterIdle to be %v, got %v", 15*time.Minute, cfg.Watch.RemindAfterIdle)
	}
	if cfg.Watch.RemindAfterPause != 5*time.Minute {
		t.Errorf("expected RemindAfterPause to be %v, got %v", 5*time.Minute, cfg.Watch.RemindAfterPause)
	}
	if cfg.Watch.RemindAfterActive != 2*time.Hour {
		t.Errorf("expected RemindAfterActive to be %v, got %v", 2*time.Hour, cfg.Watch.RemindAfterActive)
	}
}

func TestLoadConfig_UserOverrides(t *testing.T) {
	content := `
watch:
  interval: "1m"
  remind_after_idle: "30m"
  remind_after_pause: "10m"
  remind_after_active: "1h30m"
`
	path, cleanup := createTestConfigFile(t, content)
	defer cleanup()

	// Temporarily set the config path to our test file
	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(filepath.Dir(path)))

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg.Watch.Interval != 1*time.Minute {
		t.Errorf("expected Interval to be %v, got %v", 1*time.Minute, cfg.Watch.Interval)
	}
	if cfg.Watch.RemindAfterIdle != 30*time.Minute {
		t.Errorf("expected RemindAfterIdle to be %v, got %v", 30*time.Minute, cfg.Watch.RemindAfterIdle)
	}
	if cfg.Watch.RemindAfterPause != 10*time.Minute {
		t.Errorf("expected RemindAfterPause to be %v, got %v", 10*time.Minute, cfg.Watch.RemindAfterPause)
	}
	if cfg.Watch.RemindAfterActive != 90*time.Minute {
		t.Errorf("expected RemindAfterActive to be %v, got %v", 90*time.Minute, cfg.Watch.RemindAfterActive)
	}
}

func TestLoadConfig_Partial(t *testing.T) {
	content := `
watch:
  remind_after_pause: "1m"
`
	path, cleanup := createTestConfigFile(t, content)
	defer cleanup()

	// Temporarily set the config path to our test file
	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(filepath.Dir(path)))

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}
	// Check that the overridden value is set
	if cfg.Watch.RemindAfterPause != 1*time.Minute {
		t.Errorf("expected RemindAfterPause to be %v, got %v", 1*time.Minute, cfg.Watch.RemindAfterPause)
	}
	// Check that other values are still the default
	if cfg.Watch.Interval != 5*time.Minute {
		t.Errorf("expected Interval to be %v, got %v", 5*time.Minute, cfg.Watch.Interval)
	}
	if cfg.Watch.RemindAfterIdle != 15*time.Minute {
		t.Errorf("expected RemindAfterIdle to be %v, got %v", 15*time.Minute, cfg.Watch.RemindAfterIdle)
	}
}

func TestLoadConfig_Malformed(t *testing.T) {
	content := `
watch:
  remind_after_idle: "invalid-duration"
`
	path, cleanup := createTestConfigFile(t, content)
	defer cleanup()

	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(filepath.Dir(path)))

	cfg, err := LoadConfig()
	if err != nil {
		// We don't expect an error from LoadConfig itself, as it should
		// gracefully handle a parsing error for a single field.
		t.Fatalf("LoadConfig() returned an unexpected error: %v", err)
	}

	if cfg.Watch.RemindAfterIdle == 0 {
		t.Errorf("expected RemindAfterIdle to fall back to default, but it was zero")
	}

	if cfg.Watch.RemindAfterIdle != defaultConfig.Watch.RemindAfterIdle {
		t.Errorf("expected RemindAfterIdle to be default %v, got %v", defaultConfig.Watch.RemindAfterIdle, cfg.Watch.RemindAfterIdle)
	}
}

func TestLoadConfig_MalformedYAML(t *testing.T) {
	content := `not: valid: yaml`
	path, cleanup := createTestConfigFile(t, content)
	defer cleanup()

	t.Setenv("XDG_CONFIG_HOME", filepath.Dir(filepath.Dir(path)))

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("LoadConfig() should have failed for malformed YAML, but didn't")
	}
}

func TestStaleSessionThresholdConfig(t *testing.T) {
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

	// Test default value
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	expected := 8 * time.Hour
	if config.ParsedStaleSessionThreshold() != expected {
		t.Errorf("Expected default stale session threshold to be %v, got %v", expected, config.ParsedStaleSessionThreshold())
	}

	// Test custom value
	tempDir := t.TempDir()
	flowConfigDir := filepath.Join(tempDir, "flow")
	if err := os.MkdirAll(flowConfigDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	configPath := filepath.Join(flowConfigDir, "config.yml")
	configData := `stale_session_threshold: "6h"`
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Temporarily override XDG_CONFIG_HOME
	if err := os.Setenv("XDG_CONFIG_HOME", tempDir); err != nil {
		t.Fatalf("Failed to set XDG_CONFIG_HOME: %v", err)
	}

	config, err = LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load custom config: %v", err)
	}

	expected = 6 * time.Hour
	if config.ParsedStaleSessionThreshold() != expected {
		t.Errorf("Expected custom stale session threshold to be %v, got %v", expected, config.ParsedStaleSessionThreshold())
	}
}
