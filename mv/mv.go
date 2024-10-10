package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	help := flag.Bool("help", false, "Show help information")
	srcFile := flag.String("src", "", "Source file to move or rename")
	destFile := flag.String("dest", "", "Destination file path")
	force := flag.Bool("f", false, "Overwrite the destination file if it exists")

	flag.Parse()

	// Show help information if requested
	if *help {
		fmt.Println("Usage: mv [options] <source> <destination>")
		fmt.Println("Options:")
		fmt.Println("  --help          Show help information")
		fmt.Println("  --src <file>    Source file to move or rename")
		fmt.Println("  --dest <file>   Destination file path")
		fmt.Println("  --f             Overwrite the destination file if it exists")
		return
	}

	// Handle positional arguments if flags are not used
	args := flag.Args()
	if len(args) > 0 {
		if *srcFile == "" {
			*srcFile = args[0]
		}
		if len(args) > 1 && *destFile == "" {
			*destFile = args[1]
		}
	}

	// Ensure both source and destination files are specified
	if *srcFile == "" || *destFile == "" {
		fmt.Println("Error: Both source and destination files are required")
		return
	}

	// Check if destination file exists and handle overwrite option
	if _, err := os.Stat(*destFile); err == nil && !*force {
		fmt.Println("Error: Destination file already exists. Use -f to replace it.")
		return
	}

	// Attempt to move or rename the file
	err := os.Rename(*srcFile, *destFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("File '%s' moved/renamed to '%s' successfully.\n", *srcFile, *destFile)
}
