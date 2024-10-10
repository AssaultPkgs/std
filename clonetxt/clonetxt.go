package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	help := flag.Bool("help", false, "Show help information")
	srcFile := flag.String("src", "", "Source file to copy")
	destFile := flag.String("dest", "", "Destination file path")
	overwrite := flag.Bool("overwrite", false, "Overwrite the destination file if it exists")

	flag.Parse()

	if *help {
		fmt.Println("Usage: cp [options]")
		fmt.Println("Options:")
		fmt.Println("  --help          Show help information")
		fmt.Println("  --src <file>    Source file to copy")
		fmt.Println("  --dest <file>   Destination file path")
		fmt.Println("  --overwrite     Overwrite the destination file if it exists")
		return
	}

	if *srcFile == "" || *destFile == "" {
		if *&os.Args[1] != "" && *&os.Args[2] != "" {
			srcFile = &os.Args[1]
			destFile = &os.Args[2]
		} else {
			fmt.Println("Error: Both source and destination files are required")
			return
		}
	}

	// Open the source file
	input, err := os.Open(*srcFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer input.Close()

	// Check if the destination file exists and handle accordingly
	var output *os.File
	if _, err := os.Stat(*destFile); err == nil && !*overwrite {
		fmt.Println("Error: Destination file already exists. Use -overwrite to replace it.")
		return
	} else {
		output, err = os.Create(*destFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		defer output.Close()
	}

	// Copy the content from source to destination
	_, err = io.Copy(output, input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("File '%s' copied to '%s' successfully.\n", *srcFile, *destFile)
}
