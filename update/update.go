package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"minecraft-server-manager/log"

	"github.com/hashicorp/go-version"
)

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
	currentVersionValue, err := version.NewVersion(currentVersion)
	if err != nil {
		log.Error("Could not compare versions")
		return false
	}
	latestVersionValue, err := version.NewVersion(latestVersion)
	if err != nil {
		log.Error("Could not compare versions")
		return false
	}
	return latestVersionValue.GreaterThan(currentVersionValue)
}

func getAppFullPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return exePath, nil
}

func RunUpdate(currentVersion string, force bool) {
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

	if !force && !updateAvailable(currentVersion, latestVersion) {
		log.Info("You are up to date! No need to update")
		return
	}

	log.Info("Downloading to path: " + appPath)

	err = updateScript(appPath, latestVersion)
	if err != nil {
		log.Error(err)
		return
	}
}

func updateScript(appPath string, latestVersion string) error {
	fmt.Println("Starting self-update process...")

	// Create the updater script
	scriptPath := "/tmp/msm_updater.sh"

	script := "#!/bin/bash\n"
	script += "echo 'Starting self-updater...'\n"
	script += "sleep 2\n"
	script += fmt.Sprintf("curl --progress-bar -L -o %s https://github.com/%s/releases/download/v%s/%s_%s_%s\n", appPath, githubRepo, latestVersion, "minecraft-server-manager", runtime.GOOS, runtime.GOARCH)
	script += fmt.Sprintf("chmod +x %s\n", appPath)
	script += "echo 'Update done'\n"
	script += fmt.Sprintf("rm %s", scriptPath)

	// Write the script to a file
	err := os.WriteFile(scriptPath, []byte(script), 0755)
	if err != nil {
		return err
	}

	// Run the update script in a detached screen session
	cmd := exec.Command("screen", "-dmS", "msm_updater", "bash", scriptPath)
	err = cmd.Run()
	if err != nil {
		return err
	}

	log.Info("Update started in background. Check progress with: screen -r msm_updater")
	return nil
}
