package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "flow",
	Short: "Flow is a terminal-based tool for deep work.",
	Long: `A minimalist command-line tool for focused, single-tasking work sessions.
It protects your attention, helps you build a deep work habit, and provides
powerful insights into your focus patterns—all without leaving your terminal.`,
	Version: version, // This will be handled by a version flag
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{printf "Flow %s\n" .Version}}`)
}
