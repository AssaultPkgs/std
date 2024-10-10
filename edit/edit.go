package main

import (
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	// Define command-line flags for editor
	var editor string
	flag.StringVar(&editor, "e", "vscode", "Specify the editor to use: vscode, notepad, or notepad++")
	flag.Parse()

	// Get the file names from command-line arguments
	files := flag.Args()
	if len(files) == 0 {
		fmt.Println("Usage: go run edit.go -e [editor] <filename1> <filename2> ...")
		return
	}

	for _, file := range files {
		err := openFileInEditor(file, editor)
		if err != nil {
			fmt.Printf("Error opening file %s in %s: %v\n", file, editor, err)
		}
	}
}

func openFileInEditor(file, editor string) error {
	var cmd *exec.Cmd

	// Determine the command based on the editor specified
	switch editor {
	case "vscode":
		cmd = exec.Command("code", file)
	case "notepad":
		cmd = exec.Command("notepad", file)
	case "notepad++":
		cmd = exec.Command("notepad++", file)
	default:
		return fmt.Errorf("unknown editor: %s", editor)
	}

	// Check if the editor command is available
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start editor: %v", err)
	}

	return nil
}
