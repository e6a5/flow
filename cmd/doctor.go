package cmd

import (
	"fmt"
	"os"

	"github.com/e6a5/flow/core"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Run a diagnostic check on your Flow setup",
	Long:  `Checks for common problems with your configuration, session files, and log data.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🩺 Running diagnostics...")
		allGood := true

		// Check 1: Config file
		cfgPath, err := core.GetConfigPath()
		if err != nil {
			fmt.Println("❌ Config Path: Could not determine config path.")
			allGood = false
		} else {
			_, err := os.Stat(cfgPath)
			if os.IsNotExist(err) {
				fmt.Printf("✅ Config File: OK (No config file found, using defaults).\n")
			} else if err != nil {
				fmt.Printf("❌ Config File: Error checking config at %s: %v\n", cfgPath, err)
				allGood = false
			} else {
				// Try to load it
				_, err := core.LoadConfig()
				if err != nil {
					fmt.Printf("❌ Config File: Found at %s, but could not parse: %v\n", cfgPath, err)
					allGood = false
				} else {
					fmt.Printf("✅ Config File: OK (Loaded successfully from %s).\n", cfgPath)
				}
			}
		}

		// Check 2: Session file
		sessionPath, err := core.GetSessionPath()
		if err != nil {
			fmt.Println("❌ Session Path: Could not determine session path.")
			allGood = false
		} else {
			if core.SessionExists() {
				_, err := core.LoadSession()
				if err != nil {
					fmt.Printf("❌ Session File: Corrupted or unreadable at %s: %v\n", sessionPath, err)
					allGood = false
				} else {
					fmt.Printf("✅ Session File: OK (Readable at %s).\n", sessionPath)
				}
			} else {
				fmt.Printf("✅ Session File: OK (No active session).\n")
			}
		}

		// Check 3: Log directory
		logDir, err := core.GetLogDir()
		if err != nil {
			fmt.Println("❌ Log Directory: Could not determine log directory.")
			allGood = false
		} else {
			info, err := os.Stat(logDir)
			if os.IsNotExist(err) {
				fmt.Printf("✅ Log Directory: OK (Will be created at %s).\n", logDir)
			} else if err != nil || !info.IsDir() {
				fmt.Printf("❌ Log Directory: Path at %s is not a valid directory.\n", logDir)
				allGood = false
			} else {
				fmt.Printf("✅ Log Directory: OK (Exists at %s).\n", logDir)
			}
		}

		fmt.Println()
		if allGood {
			fmt.Println("✨ Your Flow setup looks healthy! ✨")
		} else {
			fmt.Println("⚠️  Found issues with your setup. Please review the messages above.")
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
