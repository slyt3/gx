package main

import (
	"fmt"
	"os"

	"github.com/slyt3/gx/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gx <command> [args]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: no script file specified")
			os.Exit(1)
		}
		exitCode := cmd.Run(os.Args[2], os.Args[3:])
		os.Exit(exitCode)
	case "version":
		fmt.Println("gx version 0.1.0")

	case "clean":
		if err := cmd.Clean(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "watch":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: no script file specified")
			os.Exit(1)
		}
		if err := cmd.Watch(os.Args[2], os.Args[3:]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}
