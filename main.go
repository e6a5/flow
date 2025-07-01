package main

import (
	"fmt"
	"os"
)

// Version information (set by build flags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "start":
		handleStart()
	case "status":
		handleStatus()
	case "pause":
		handlePause()
	case "resume":
		handleResume()
	case "end":
		handleEnd()
	case "log":
		handleLog()
	case "completion":
		handleCompletion()
	case "--help", "-h":
		showUsage()
	case "--version", "-v":
		showVersion()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		showUsage()
		os.Exit(1)
	}
}
