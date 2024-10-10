package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return
	}

	// Define the default path to the user's home .vin_env/bin directory
	defaultPath := filepath.Join(homeDir, ".vin_env", "bin")

	// Define the command-line flag for the directory path with the default value set to the home .vin_env/bin directory
	pathFlag := flag.String("path", defaultPath, "Path to the directory")
	flag.Parse()

	// Open the directory
	entries, err := os.ReadDir(*pathFlag)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Get the terminal width to calculate the number of columns
	termWidth := getTerminalWidth()

	// Find the maximum length of file/directory names
	maxLength := getMaxLength(entries)

	// Calculate how many items can fit per row based on terminal width
	// Adding 2 to maxLength for spacing between columns
	itemsPerRow := termWidth / (maxLength + 2)

	// Track how many items have been printed in the current row
	count := 0

	// Iterate over directory entries
	for _, entry := range entries {
		// Get the name and apply color based on type
		var coloredName string
		if entry.IsDir() {
			coloredName = fmt.Sprintf("\033[94m%s\033[0m", entry.Name()) // Light Blue for directories
		} else if strings.HasSuffix(entry.Name(), ".zip") {
			coloredName = fmt.Sprintf("\033[91m%s\033[0m", entry.Name()) // Light Red for zip files
		} else {
			coloredName = fmt.Sprintf("\033[92m%s\033[0m", entry.Name()) // Light Green for regular files
		}

		// Print the name with adjusted width
		fmt.Printf("%-*s ", maxLength+2, coloredName)

		// Increment the count and check if we need to start a new row
		count++
		if count%itemsPerRow == 0 {
			fmt.Println() // Print a newline after every row
		}
	}

	// If the last row didn't finish, print a newline to tidy up the output
	if count%itemsPerRow != 0 {
		fmt.Println()
	}
}

// getMaxLength returns the maximum length of file or directory names
func getMaxLength(entries []os.DirEntry) int {
	maxLength := 0
	for _, entry := range entries {
		if len(entry.Name()) > maxLength {
			maxLength = len(entry.Name())
		}
	}
	return maxLength
}

// getTerminalWidth retrieves the width of the terminal
func getTerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 80 // Default to 80 if unable to determine
	}

	dimensions := strings.Fields(string(out))
	width, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return 80 // Default to 80 if conversion fails
	}
	return width
}
