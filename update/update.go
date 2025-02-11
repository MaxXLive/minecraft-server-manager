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
	currentVersionInt, err := strconv.ParseInt(strings.ReplaceAll(currentVersion, ".", ""), 10, 32)
	if err != nil {
		log.Error("Could not compare versions")
		return false
	}
	latestVersionInt, err := strconv.ParseInt(strings.ReplaceAll(latestVersion, ".", ""), 10, 32)
	if err != nil {
		log.Error("Could not compare versions")
		return false
	}
	return latestVersionInt > currentVersionInt
}

func getAppFullPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return exePath, nil
}

func RunUpdate(currentVersion string) {
	appPath, err := getAppFullPath()
	if err != nil {
		log.Error(err)
		return
	}

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

	log.Info("Downloading to path: " + appPath)

	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s_%s_%s", githubRepo, latestVersion, "minecraft-server-manager", runtime.GOOS, runtime.GOARCH)

	// Use curl to download the latest release (you can also use Go's http client here)
	cmd := exec.Command("curl", "--progress-bar", "-L", "-o", appPath, downloadURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Error(err)
	}

	// Make the new executable executable
	err = os.Chmod(appPath, 0755)
	if err != nil {
		log.Error(fmt.Errorf("failed to make new executable: %v", err))
	}
}
