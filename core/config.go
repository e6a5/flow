package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/goccy/go-yaml"
)

// Config holds all application configuration.
type Config struct {
	Watch WatchConfig `yaml:"watch"`
}

// WatchConfig holds configuration specific to the 'watch' command.
type WatchConfig struct {
	Interval          time.Duration `yaml:"interval"`
	RemindAfterIdle   time.Duration `yaml:"remind_after_idle"`
	RemindAfterPause  time.Duration `yaml:"remind_after_pause"`
	RemindAfterActive time.Duration `yaml:"remind_after_active"`
}

var defaultConfig = Config{
	Watch: WatchConfig{
		Interval:          5 * time.Minute,
		RemindAfterIdle:   15 * time.Minute,
		RemindAfterPause:  5 * time.Minute,
		RemindAfterActive: 2 * time.Hour,
	},
}

// LoadConfig loads the configuration from the YAML file, applying defaults.
func LoadConfig() (Config, error) {
	cfg := defaultConfig

	configPath, err := getConfigPath()
	if err != nil {
		return cfg, fmt.Errorf("could not determine config path: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No config file, return default config. This is not an error.
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("could not read config file: %w", err)
	}

	// Temporary struct to read user-provided duration strings
	var userCfg struct {
		Watch struct {
			Interval          string `yaml:"interval"`
			RemindAfterIdle   string `yaml:"remind_after_idle"`
			RemindAfterPause  string `yaml:"remind_after_pause"`
			RemindAfterActive string `yaml:"remind_after_active"`
		} `yaml:"watch"`
	}

	if err := yaml.Unmarshal(data, &userCfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Parse user strings and apply them over defaults
	if userCfg.Watch.Interval != "" {
		if d, err := time.ParseDuration(userCfg.Watch.Interval); err == nil {
			cfg.Watch.Interval = d
		}
	}
	if userCfg.Watch.RemindAfterIdle != "" {
		if d, err := time.ParseDuration(userCfg.Watch.RemindAfterIdle); err == nil {
			cfg.Watch.RemindAfterIdle = d
		}
	}
	if userCfg.Watch.RemindAfterPause != "" {
		if d, err := time.ParseDuration(userCfg.Watch.RemindAfterPause); err == nil {
			cfg.Watch.RemindAfterPause = d
		}
	}
	if userCfg.Watch.RemindAfterActive != "" {
		if d, err := time.ParseDuration(userCfg.Watch.RemindAfterActive); err == nil {
			cfg.Watch.RemindAfterActive = d
		}
	}

	return cfg, nil
}

// getConfigPath determines the expected path for the configuration file.
func getConfigPath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "flow", "config.yml"), nil
}
