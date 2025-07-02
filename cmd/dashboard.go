package cmd

import (
	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Show a yearly contribution graph of your focus sessions",
	Long:  `Visualizes your deep work history over the last year, similar to a GitHub contribution graph.`,
	Run: func(cmd *cobra.Command, args []string) {
		core.HandleDashboard()
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}
