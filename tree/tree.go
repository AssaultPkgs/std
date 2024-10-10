package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[94m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
)

var (
	indent   = "    "
	maxDepth = 10
)

func main() {
	pathFlag := flag.String("p", ".", "Path to the directory")
	allFlag := flag.Bool("a", false, "Include hidden files")
	depthFlag := flag.Int("d", maxDepth, "Maximum display depth of the directory tree")
	flag.Parse()

	printTree(*pathFlag, *allFlag, *depthFlag)
}

func printTree(path string, showHidden bool, maxDepth int) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	if !showHidden {
		entries = filterHidden(entries)
	}

	fmt.Println(path)

	printDirectoryTree(entries, path, 0, maxDepth, "")
}

func filterHidden(entries []os.DirEntry) []os.DirEntry {
	var filtered []os.DirEntry
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func printDirectoryTree(entries []os.DirEntry, path string, level int, maxDepth int, prefix string) {
	for i, entry := range entries {
		isLast := i == len(entries)-1
		newPrefix := getPrefix(prefix, level, isLast)

		if entry.IsDir() {
			fmt.Printf("%s%s%s%s -> \n", newPrefix, getColorForFile(entry), entry.Name(), colorReset)

			if level < maxDepth {
				subEntries, err := os.ReadDir(filepath.Join(path, entry.Name()))
				if err != nil {
					fmt.Println("Error reading directory:", err)
					continue
				}

				printDirectoryTree(subEntries, filepath.Join(path, entry.Name()), level+1, maxDepth, newPrefix)
			}
		} else {
			fmt.Printf("%s%s%s%s\n", newPrefix, getColorForFile(entry), entry.Name(), colorReset)
		}
	}
}

func getPrefix(prefix string, level int, isLast bool) string {
	if level == 0 {
		return prefix
	}

	if isLast {
		return prefix + "│ "
	}

	return prefix + "│ "
}

func getColorForFile(entry os.DirEntry) string {
	if entry.IsDir() {
		return colorBlue
	}

	ext := strings.ToLower(filepath.Ext(entry.Name()))
	switch ext {
	case ".zip", ".tar", ".gz", ".rar", ".7z":
		return colorRed
	case ".txt", ".md", ".log":
		return colorYellow
	case ".jpg", ".png", ".gif", ".bmp":
		return colorCyan
	case ".exe", ".dll":
		return colorCyan
	default:
		return colorGreen
	}
}

func getLinePrefix(level int, isLast bool) string {
	if level == 0 {
		return ""
	}

	if isLast {
		return "└─"
	}

	return "├─"
}
