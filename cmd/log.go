package cmd

import (
	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log [YYYY-MM]",
	Short: "View completed session history",
	Long: `Displays a log of completed deep work sessions.
You can filter the log by time periods (today, week, month) or view statistics.`,
	Run: func(cmd *cobra.Command, args []string) {
		showStats, _ := cmd.Flags().GetBool("stats")
		filterToday, _ := cmd.Flags().GetBool("today")
		filterWeek, _ := cmd.Flags().GetBool("week")
		filterMonth, _ := cmd.Flags().GetBool("month")
		showAll, _ := cmd.Flags().GetBool("all")

		monthStr := ""
		if len(args) > 0 {
			monthStr = args[0]
		}

		core.HandleLog(showStats, filterToday, filterWeek, filterMonth, showAll, monthStr)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.Flags().Bool("today", false, "Show sessions from today")
	logCmd.Flags().Bool("week", false, "Show sessions from this week")
	logCmd.Flags().Bool("month", false, "Show sessions from this month")
	logCmd.Flags().Bool("stats", false, "Show summary statistics")
	logCmd.Flags().Bool("all", false, "Show all session history")
}
