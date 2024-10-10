package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	help := flag.Bool("help", false, "Show help information")
	path := flag.String("path", "", "File or directory to remove")
	recursive := flag.Bool("r", false, "Remove directories and their contents recursively")
	force := flag.Bool("f", false, "Ignore nonexistent files and arguments")

	flag.Parse()

	if *help {

		fmt.Println("Usage: rm [options]")
		fmt.Println("Options:")
		fmt.Println("  --help          Show help information")
		fmt.Println("  --path <path>   File or directory to remove")
		fmt.Println("  --r             Remove directories and their contents recursively")
		fmt.Println("  --f             Ignore nonexistent files and arguments")
		return
	}

	if *path == "" {
		//fmt.Println("Error: Path is required")
		//return
		if *&os.Args[1] != "" {
			path = &os.Args[1]
		} else {
			fmt.Println("Error: Path is required")
		}
	}

	if *recursive {
		err := os.RemoveAll(*path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	} else {
		err := os.Remove(*path)
		if err != nil {
			if !os.IsNotExist(err) || !*force {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}
	}

	fmt.Printf("thing '%s' removed successfully.\n", *path)
}
