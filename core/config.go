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
	Watch                 WatchConfig `yaml:"watch"`
	DailyGoal             string      `yaml:"daily_goal"`
	StaleSessionThreshold string      `yaml:"stale_session_threshold"`
	parsedGoal            time.Duration
	parsedStaleThreshold  time.Duration
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
	StaleSessionThreshold: "8h", // Default to 8 hours
	parsedStaleThreshold:  8 * time.Hour,
}

// ParsedDailyGoal returns the parsed daily goal duration.
func (c *Config) ParsedDailyGoal() time.Duration {
	return c.parsedGoal
}

// ParsedStaleSessionThreshold returns the parsed stale session threshold duration.
func (c *Config) ParsedStaleSessionThreshold() time.Duration {
	return c.parsedStaleThreshold
}

// LoadConfig loads the configuration from the YAML file, applying defaults.
func LoadConfig() (Config, error) {
	cfg := defaultConfig

	configPath, err := GetConfigPath()
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

	// A temporary struct for all user settings to avoid direct manipulation
	var tempCfg struct {
		Watch struct {
			Interval          string `yaml:"interval"`
			RemindAfterIdle   string `yaml:"remind_after_idle"`
			RemindAfterPause  string `yaml:"remind_after_pause"`
			RemindAfterActive string `yaml:"remind_after_active"`
		} `yaml:"watch"`
		DailyGoal             string `yaml:"daily_goal"`
		StaleSessionThreshold string `yaml:"stale_session_threshold"`
	}

	if err := yaml.Unmarshal(data, &tempCfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Parse user strings and apply them over defaults
	if tempCfg.Watch.Interval != "" {
		if d, err := time.ParseDuration(tempCfg.Watch.Interval); err == nil {
			cfg.Watch.Interval = d
		}
	}
	if tempCfg.Watch.RemindAfterIdle != "" {
		if d, err := time.ParseDuration(tempCfg.Watch.RemindAfterIdle); err == nil {
			cfg.Watch.RemindAfterIdle = d
		}
	}
	if tempCfg.Watch.RemindAfterPause != "" {
		if d, err := time.ParseDuration(tempCfg.Watch.RemindAfterPause); err == nil {
			cfg.Watch.RemindAfterPause = d
		}
	}
	if tempCfg.Watch.RemindAfterActive != "" {
		if d, err := time.ParseDuration(tempCfg.Watch.RemindAfterActive); err == nil {
			cfg.Watch.RemindAfterActive = d
		}
	}
	if tempCfg.DailyGoal != "" {
		cfg.DailyGoal = tempCfg.DailyGoal
		if d, err := time.ParseDuration(tempCfg.DailyGoal); err == nil {
			cfg.parsedGoal = d
		}
	}
	if tempCfg.StaleSessionThreshold != "" {
		cfg.StaleSessionThreshold = tempCfg.StaleSessionThreshold
		if d, err := time.ParseDuration(tempCfg.StaleSessionThreshold); err == nil {
			cfg.parsedStaleThreshold = d
		}
	}

	return cfg, nil
}

// GetConfigPath determines the expected path for the configuration file.
func GetConfigPath() (string, error) {
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
