package main

import (
	"fmt"
	"os"

	"github.com/roma-glushko/bp/cmd"
	"github.com/roma-glushko/bp/internal/version"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		if err := cmd.ServeCommand(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "--version", "-v", "version":
		fmt.Println(version.FullVersion)
	case "--help", "-h", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`bp - A blood pressure journal (%s)

Usage:
  bp serve [flags]    Start the local web UI
  bp version          Print version information
  bp help             Show this help

Serve flags:
  --port <port>       Port to listen on (default: 7391)
  --data <path>       Data directory path
  --no-open           Do not open browser automatically
  --debug             Enable debug logging
`, version.Version)
}
