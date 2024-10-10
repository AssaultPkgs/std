package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const (
	orgURL     = "https://api.github.com/orgs/AssaultPkgs/repos"
	baseDir    = "C:\\Users\\%s\\vin_env\\bin"
	deletedDir = "DeletedThings"
	lightBlue  = "\033[1;36m" // Light blue color
	limeGreen  = "\033[1;32m" // Lime green color
	resetColor = "\033[0m"    // Reset color
)

type Repo struct {
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	CloneURL  string `json:"clone_url"`
	UpdatedAt string `json:"updated_at"`
	HTMLURL   string `json:"html_url"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run apt.go [command]")
		return
	}

	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return
	}

	// Set the base path correctly
	basePath := fmt.Sprintf(baseDir, currentUser.Username)
	args := os.Args[1:]

	switch args[0] {
	case "--list":
		listRepos()
	case "--list-installed":
		listInstalled(basePath)
	case "--delete":
		if len(args) < 2 {
			fmt.Println("Please specify a package name to delete.")
			return
		}
		deletePackage(basePath, args[1])
	case "--update":
		if len(args) < 2 {
			fmt.Println("Please specify a package name to update.")
			return
		}
		updatePackage(basePath, args[1])
	case "--install":
		if len(args) < 2 {
			fmt.Println("Please specify a package name to install.")
			return
		}
		installPackage(basePath, args[1])
	case "--install-all":
		installAllPackages(basePath)
	case "--list-upgradeable":
		listUpgradeable()
	case "--list-new":
		listNewRepos()
	default:
		fmt.Println("Unknown command:", args[0])
	}
}

func listRepos() {
	resp, err := http.Get(orgURL)
	if err != nil {
		fmt.Println("Error fetching repositories:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get repositories:", resp.Status)
		return
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println("Available Repositories:")
	for _, repo := range repos {
		fmt.Printf("%s- %s (%s)%s\n", lightBlue, repo.Name, repo.HTMLURL, resetColor)
	}
}

func listInstalled(basePath string) {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println("Error reading installed packages:", err)
		return
	}

	fmt.Println("Installed Packages:")
	for _, file := range files {
		if file.IsDir() {
			fmt.Println("- " + file.Name())
		}
	}
}

func deletePackage(basePath, name string) {
	src := filepath.Join(basePath, name)
	dest := filepath.Join(basePath, deletedDir, name)

	if err := os.MkdirAll(filepath.Join(basePath, deletedDir), os.ModePerm); err != nil {
		fmt.Println("Error creating DeletedThings directory:", err)
		return
	}

	if err := os.Rename(src, dest); err != nil {
		fmt.Println("Error deleting package:", err)
		return
	}

	fmt.Printf("Package '%s' moved to DeletedThings.\n", name)
}

func updatePackage(basePath, name string) {
	installPackage(basePath, name) // Reuse install logic for update
}

func installPackage(basePath, name string) {
	repoURL := fmt.Sprintf("https://github.com/AssaultPkgs/%s", name)
	zipURL := fmt.Sprintf("https://github.com/AssaultPkgs/%s/archive/refs/heads/main.zip", name)

	// Create the directory for the package
	destDir := filepath.Join(basePath, name, "PKG")
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		fmt.Println("Error creating package directory:", err)
		return
	}

	// Download the zip file
	resp, err := http.Get(zipURL)
	if err != nil {
		fmt.Println("Error downloading package:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to download package:", resp.Status)
		return
	}

	// Save the zip file
	zipFile := filepath.Join(destDir, name+".zip")
	out, err := os.Create(zipFile)
	if err != nil {
		fmt.Println("Error creating zip file:", err)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		fmt.Println("Error saving zip file:", err)
		return
	}

	fmt.Printf("Package '%s' installed successfully from %s.\n", name, repoURL)
}

func installAllPackages(basePath string) {
	resp, err := http.Get(orgURL)
	if err != nil {
		fmt.Println("Error fetching repositories:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get repositories:", resp.Status)
		return
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	for _, repo := range repos {
		installPackage(basePath, repo.Name)
	}
}

func listUpgradeable() {
	resp, err := http.Get(orgURL)
	if err != nil {
		fmt.Println("Error fetching repositories:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get repositories:", resp.Status)
		return
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Get the current time
	now := time.Now()
	// Time range for the last 7 days
	sevenDaysAgo := now.Add(-7 * 24 * time.Hour)

	fmt.Println("Upgradeable Packages (Last Updated in Last 7 Days):")
	count := 0
	for _, repo := range repos {
		updatedAt, err := time.Parse(time.RFC3339, repo.UpdatedAt)
		if err != nil {
			fmt.Println("Error parsing updated time:", err)
			continue
		}
		if updatedAt.After(sevenDaysAgo) {
			fmt.Printf("%s- %s (Updated at: %s, URL: %s)%s\n", lightBlue, repo.Name, updatedAt.Format(time.RFC1123), repo.HTMLURL, resetColor)
			count++
			if count >= 20 {
				break // Limit to last 20 updated repos
			}
		}
	}
}

func listNewRepos() {
	resp, err := http.Get(orgURL)
	if err != nil {
		fmt.Println("Error fetching repositories:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get repositories:", resp.Status)
		return
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println("New Repositories:")
	for _, repo := range repos {
		fmt.Printf("%s- %s (URL: %s)%s\n", lightBlue, repo.Name, repo.HTMLURL, resetColor)
	}
}
