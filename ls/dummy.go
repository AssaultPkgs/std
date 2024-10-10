package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Define the base directory for the dummy structure
	baseDir := "./test_dir"

	// Create the base directory
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating base directory: %v\n", err)
		return
	}

	// Define files and directories to create
	entries := []struct {
		Name     string
		IsDir    bool
		SubFiles []string
	}{
		{"dir1", true, []string{"file1.txt", "file2.md", "file3.zip"}},
		{"dir2", true, []string{"file4.gz", "file5.html", "file6.css"}},
		{"dir3", true, []string{"file7.js", "file8.go", "file9.pdf"}},
		{"file1.txt", false, nil},
		{"file2.md", false, nil},
		{"file3.zip", false, nil},
		{"file4.gz", false, nil},
		{"file5.html", false, nil},
		{"file6.css", false, nil},
		{"file7.js", false, nil},
		{"file8.go", false, nil},
		{"file9.pdf", false, nil},
	}

	// Create the files and directories
	for _, entry := range entries {
		path := filepath.Join(baseDir, entry.Name)
		if entry.IsDir {
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				fmt.Printf("Error creating directory %s: %v\n", path, err)
				return
			}
			for _, subFile := range entry.SubFiles {
				subFilePath := filepath.Join(path, subFile)
				file, err := os.Create(subFilePath)
				if err != nil {
					fmt.Printf("Error creating file %s: %v\n", subFilePath, err)
					return
				}
				file.Close()
			}
		} else {
			_, err := os.Create(path)
			if err != nil {
				fmt.Printf("Error creating file %s: %v\n", path, err)
				return
			}
		}
	}

	fmt.Println("Dummy files and directories created successfully!")
}
