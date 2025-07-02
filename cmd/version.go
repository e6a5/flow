package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Flow",
	Long:  `All software has versions. This is Flow's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Flow %s\n", version)
		if commit != "none" {
			fmt.Printf("Commit: %s\n", commit)
		}
		if date != "unknown" {
			fmt.Printf("Built: %s\n", date)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
