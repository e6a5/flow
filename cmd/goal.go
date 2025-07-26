package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/e6a5/flow/core"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var goalCmd = &cobra.Command{
	Use:   "goal",
	Short: "Set or view your daily focus goal",
	Long:  `Manages your daily focus goal. Use --set to define a new goal (e.g., '4h', '3h30m'). Run without flags to view your current goal and today's progress.`,
	Run: func(cmd *cobra.Command, args []string) {
		set, _ := cmd.Flags().GetString("set")

		if set != "" {
			handleSetGoal(set)
		} else {
			handleViewGoal()
		}
	},
}

func handleSetGoal(goalStr string) {
	// Validate duration format first
	_, err := time.ParseDuration(goalStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid duration format for goal: %v\n", err)
		os.Exit(1)
	}

	cfgPath, err := core.GetConfigPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config path: %v\n", err)
		os.Exit(1)
	}

	// Read existing config or create new one
	var configData map[string]interface{}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}
		configData = make(map[string]interface{})
	} else {
		if err := yaml.Unmarshal(data, &configData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing existing config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Set or update the daily_goal
	configData["daily_goal"] = goalStr

	// Write back to file
	updatedData, err := yaml.Marshal(configData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling config data: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(cfgPath, updatedData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Daily focus goal set to: %s\n", goalStr)
}

func handleViewGoal() {
	cfg, err := core.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return
	}

	goal := cfg.ParsedDailyGoal()
	if goal == 0 {
		fmt.Println("No daily goal set. Use 'flow goal --set <duration>' to set one.")
		return
	}

	// Get today's progress
	reader, err := core.NewLogReader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating log reader: %v\n", err)
		return
	}
	entries, err := reader.ReadRecentEntries(1000, true, false) // High limit for today
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading entries: %v\n", err)
		return
	}
	var totalTime time.Duration
	for _, entry := range entries {
		totalTime += entry.Duration
	}

	// Display progress
	percentage := 0.0
	if goal > 0 {
		percentage = (float64(totalTime) / float64(goal)) * 100
	}
	fmt.Printf("🎯 Daily Goal: %s / %s (%d%%)\n", core.FormatDuration(totalTime), core.FormatDuration(goal), int(percentage))
}

func init() {
	rootCmd.AddCommand(goalCmd)
	goalCmd.Flags().String("set", "", "Set your daily focus goal (e.g., '4h', '3h30m')")
}
