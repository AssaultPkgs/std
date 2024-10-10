package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[94m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
)

func main() {
	// Get the current working directory
	defaultPath, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}

	pathFlag := flag.String("path", defaultPath, "Path to the directory")
	longFlag := flag.Bool("l", false, "Long format (detailed listing)")
	allFlag := flag.Bool("a", false, "Include hidden files")
	flag.Parse()

	entries, err := os.ReadDir(*pathFlag)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	if !*allFlag {
		entries = filterHidden(entries)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	//fmt.Fprintf(w, "PS C:\\Users\\Hamza\\.vin_env\\bin\\ls> go run ls.go --path ./test_dir/dir3 -l --path ../../../../../../../../../\n")

	for _, entry := range entries {
		if *longFlag {
			printLongFormat(w, entry, *pathFlag)
		} else {
			printShortFormat(w, entry)
		}
	}

	//fmt.Fprintf(w, "PS C:\\Users\\Hamza\\.vin_env\\bin\\ls> ")
	w.Flush()
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

func printLongFormat(w *tabwriter.Writer, entry os.DirEntry, path string) {
	fullPath := filepath.Join(path, entry.Name())
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	fileType := "-"
	if fileInfo.IsDir() {
		fileType = "d"
	}

	mode := fileInfo.Mode()
	perms := fmt.Sprintf("%s", mode.Perm())
	modTime := fileInfo.ModTime().Format("2006-01-02T15:04:05-07:00")
	size := fileInfo.Size()

	fmt.Fprintf(w, "%s%s\t%s\t%s\t%s\t%d\t%s\t%s%s%s\n",
		fileType, perms, "1", "owner group", "owner group", size, modTime,
		getColorForFile(entry), entry.Name(), colorReset)
}

func printShortFormat(w *tabwriter.Writer, entry os.DirEntry) {
	fmt.Fprintf(w, "%s%s%s\t", getColorForFile(entry), entry.Name(), colorReset)
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
		return colorGreen
	default:
		return colorGreen
	}
}
