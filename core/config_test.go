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
	return path, func() { os.Remove(path) }
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
		// This is not a fatal error, as we expect parsing to be lenient.
		// The config should fall back to the default value.
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
