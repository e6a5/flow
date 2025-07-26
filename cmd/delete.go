package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a session",
	Long: `Deletes a session from the log.

This command interactively lists your recent sessions and allows you to select one to delete.
You will be asked to confirm before the session is permanently removed.

Example:
  flow delete`,
	Run: func(cmd *cobra.Command, args []string) {
		sessions, err := core.GetRecentSessions(10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting recent sessions: %v\n", err)
			os.Exit(1)
		}

		if len(sessions) == 0 {
			fmt.Println("No sessions to delete.")
			return
		}

		fmt.Println("Select a session to delete:")
		for i, session := range sessions {
			fmt.Printf("%d: %s - %s (%s)\n", i+1, session.StartTime.Format("2006-01-02 15:04"), session.Tag, session.Duration)
		}

		fmt.Print("Enter the number of the session to delete (or 0 to cancel): ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 0 || choice > len(sessions) {
			fmt.Println("Invalid selection.")
			return
		}

		if choice == 0 {
			fmt.Println("Operation cancelled.")
			return
		}

		sessionToDelete := sessions[choice-1]

		fmt.Printf("\nYou have selected to delete the following session:\n")
		fmt.Printf("%s - %s (%s)\n", sessionToDelete.StartTime.Format("2006-01-02 15:04"), sessionToDelete.Tag, sessionToDelete.Duration)
		fmt.Print("Are you sure you want to delete this session? (y/N) ")

		scanner.Scan()
		confirmation := scanner.Text()

		if strings.ToLower(confirmation) == "y" {
			if err := core.DeleteLogEntry(sessionToDelete); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting session: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Session deleted.")
		} else {
			fmt.Println("Operation cancelled.")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
