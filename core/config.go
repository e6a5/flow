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
	StaleSessionThreshold string `yaml:"stale_session_threshold"`
	parsedStaleThreshold  time.Duration
}

var defaultConfig = Config{
	StaleSessionThreshold: "8h", // Default to 8 hours
	parsedStaleThreshold:  8 * time.Hour,
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
		StaleSessionThreshold string `yaml:"stale_session_threshold"`
	}

	if err := yaml.Unmarshal(data, &tempCfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Parse user strings and apply them over defaults
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
