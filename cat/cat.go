package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Define command-line flags
	var showLineNumbers bool
	flag.BoolVar(&showLineNumbers, "n", false, "Show line numbers")
	flag.Parse()

	// Get the file name from command-line arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: go run cat.go [options] <filename>")
		return
	}

	fileName := args[0]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	// Read and print the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		if showLineNumbers {
			fmt.Printf("%d\t%s\n", lineNumber, line)
		} else {
			fmt.Println(line)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}
