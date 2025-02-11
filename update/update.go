package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minecraft-server-manager/log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Define the GitHub repository (replace with your repository details)
const githubRepo = "maxxlive/minecraft-server-manager"

// Get the current version of the program (you can customize this)
func getCurrentVersion() string {
	// Replace this with your actual method to get the current version,
	// e.g., by reading from a file or using a version flag.
	// For example:
	return "1.0.0"
}

// Fetch the latest release version from GitHub API
func getLatestReleaseVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubRepo)

	// Send a GET request to the GitHub API
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var release Release

	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal body: %v", err)
	}

	parts := strings.Split(release.HTMLURL, "/")

	return strings.ReplaceAll(parts[len(parts)-1], "v", ""), nil
}

// Download the latest release from GitHub

// Replace the old executable with the new one
func replaceExecutable() error {
	// We assume that the current executable is the one we are replacing
	err := os.Remove("minecraft-server-manager")
	if err != nil {
		return fmt.Errorf("failed to remove old executable: %v", err)
	}

	// Rename the downloaded file (assuming it's downloaded in the current directory)
	err = os.Rename("minecraft-server-manager", "/usr/local/bin/minecraft-server-manager")
	if err != nil {
		return fmt.Errorf("failed to replace executable: %v", err)
	}

	// Make the new executable executable
	err = os.Chmod("/usr/local/bin/minecraft-server-manager", 0755)
	if err != nil {
		return fmt.Errorf("failed to make new executable: %v", err)
	}

	return nil
}

//func Update() {
//	// Get the current version
//	currentVersion := getCurrentVersion()
//
//	// Get the latest release version from GitHub
//	latestVersion, err := getLatestReleaseVersion()
//	if err != nil {
//		fmt.Println("Error fetching latest release:", err)
//		return
//	}
//
//	// Compare versions using the semver package
//	currentSemver, err := semver.Parse(currentVersion)
//	if err != nil {
//		fmt.Println("Error parsing current version:", err)
//		return
//	}
//
//	latestSemver, err := semver.Parse(latestVersion)
//	if err != nil {
//		fmt.Println("Error parsing latest version:", err)
//		return
//	}
//
//	// Compare versions
//	if latestSemver.GT(currentSemver) {
//		// New version is available
//		fmt.Printf("New version available: %s (current: %s)\n", latestVersion, currentVersion)
//
//		// Download and replace the executable
//		err := downloadLatestRelease(latestVersion)
//		if err != nil {
//			fmt.Println("Error downloading latest release:", err)
//			return
//		}
//
//		// Replace the old executable with the new one
//		err = replaceExecutable()
//		if err != nil {
//			fmt.Println("Error replacing executable:", err)
//			return
//		}
//
//		// Done
//		fmt.Println("Update successful. Please restart the program.")
//	} else {
//		// Already up-to-date
//		fmt.Println("You are already running the latest version.")
//	}
//}

func CheckForUpdate(currentVersion string) {
	// Get the current version
	log.Info("Current version: v" + currentVersion)

	// Get the latest release version from GitHub
	latestVersion, err := getLatestReleaseVersion()
	if err != nil {
		fmt.Println("Error fetching latest release:", err)
		return
	}
	log.Info("Latest version: v" + latestVersion)
	if updateAvailable(currentVersion, latestVersion) {
		log.Info("Update available! Use \"update\" command to update")
	} else {
		log.Info("You are up to date!")
	}
}

func updateAvailable(currentVersion string, latestVersion string) bool {
	currentVersionFloat, err := strconv.ParseFloat(currentVersion, 32)
	if err != nil {
		log.Error("Could not compare versions")
		return false
	}
	latestVersionFloat, err := strconv.ParseFloat(latestVersion, 32)
	if err != nil {
		log.Error("Could not compare versions")
		return false
	}
	return latestVersionFloat > currentVersionFloat
}

func RunUpdate(appName string, currentVersion string) {
	// Get the current version
	log.Info("Current version: v" + currentVersion)

	// Get the latest release version from GitHub
	latestVersion, err := getLatestReleaseVersion()
	if err != nil {
		fmt.Println("Error fetching latest release:", err)
		return
	}
	log.Info("Latest version: v" + latestVersion)

	if !updateAvailable(currentVersion, latestVersion) {
		log.Info("You are up to date! No need to update")
		return
	}

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s_%s_%s", githubRepo, latestVersion, "minecraft-server-manager", runtime.GOOS, runtime.GOARCH)

	// Use curl to download the latest release (you can also use Go's http client here)
	cmd := exec.Command("curl", "--progress-bar", "-L", "-o", appName, downloadURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Error(err)
	}
}
