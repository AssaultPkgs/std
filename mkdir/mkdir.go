package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	help := flag.Bool("help", false, "Show help information")
	name := flag.String("name", "", "Name of the directory to create")
	parent := flag.String("parent", ".", "Parent directory path (default is current directory)")

	flag.Parse()

	if *help {
		fmt.Println("Usage: mkdir [options]")
		fmt.Println("Options:")
		fmt.Println("  --help          Show help information")
		fmt.Println("  --name <dir>    Name of the directory to create")
		fmt.Println("  --parent <dir>  Parent directory path (default is current directory)")
		return
	}

	if *name == "" {
		//fmt.Println("Error: Directory name is required")
		if *&os.Args[1] != "" {
			name = &os.Args[1]
		} else {
			fmt.Println("Error: Directory name is required")
		}
	}

	dirPath := fmt.Sprintf("%s/%s", *parent, *name)
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Directory '%s' created successfully.\n", dirPath)
}
