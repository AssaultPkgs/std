package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	help := flag.Bool("help", false, "Show help information")
	src := flag.String("src", "", "Source file or directory to copy")
	dest := flag.String("dest", "", "Destination file or directory path")
	overwrite := flag.Bool("o", false, "Overwrite files in destination if they exist")
	copyDirs := flag.Bool("d", false, "Copy directories recursively")
	threshold := flag.Int("th", 200, "File count threshold for user confirmation")

	flag.Parse()

	if *help {
		fmt.Println("Usage: cp [options]")
		fmt.Println("Options:")
		fmt.Println("  --help          Show help information")
		fmt.Println("  --src <path>    Source file or directory to copy")
		fmt.Println("  --dest <path>   Destination file or directory path")
		fmt.Println("  --o             Overwrite files in destination if they exist")
		fmt.Println("  --d             Copy directories recursively")
		fmt.Println("  --th            File count threshold for user confirmation")
		return
	}

	if *src == "" || *dest == "" {
		fmt.Println("Error: Both source and destination paths are required")
		return
	}

	srcInfo, err := os.Stat(*src)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if srcInfo.IsDir() && !*copyDirs {
		fmt.Println("Error: Source is a directory but --dirs flag is not set")
		return
	}

	if srcInfo.IsDir() {
		if err := copyDir(*src, *dest, *overwrite, *threshold); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		if err := copyFile(*src, *dest, *overwrite); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func copyDir(src string, dest string, overwrite bool, threshold int) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source directory does not exist: %w", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source path is not a directory")
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("could not create destination directory: %w", err)
	}

	var fileCount int
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)
		if info.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("could not create directory %s: %w", destPath, err)
			}
			return nil
		}

		fileCount++
		if fileCount > threshold {
			fmt.Printf("Warning: Number of files (%d) exceeds %d. Do you want to continue? (y/n): ", fileCount, threshold)
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) != "y" {
				return fmt.Errorf("operation cancelled by user")
			}
		}

		if err := copyFile(path, destPath, overwrite); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("Directory '%s' copied to '%s' successfully.\n", src, dest)
	return nil
}

func copyFile(src string, dest string, overwrite bool) error {
	if _, err := os.Stat(dest); err == nil && !overwrite {
		return fmt.Errorf("destination file '%s' already exists. Use --overwrite to replace it", dest)
	}

	input, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}
	defer input.Close()

	output, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}

	return nil
}
