package cmd

import (
	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export session data to CSV or JSON",
	Long:  `Exports your session history to a structured format like CSV or JSON for analysis or invoicing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Like log, this delegates to the old handler.
		// A full refactor would move the flag parsing from core.HandleExport here.
		core.HandleExport()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().String("format", "csv", "Export format (csv or json)")
	exportCmd.Flags().String("output", "", "Output file path (default is stdout)")
	exportCmd.Flags().Bool("today", false, "Export sessions from today")
	exportCmd.Flags().Bool("week", false, "Export sessions from this week")
	exportCmd.Flags().Bool("month", false, "Export sessions from this month")
	exportCmd.Flags().Bool("all", false, "Export all session history")
}
