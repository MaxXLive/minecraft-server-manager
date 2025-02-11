package system

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// getJavaVersion runs `java -version` and returns the Java version string

func GetJavaVersion() (string, error) {
	// Run the command `java -version`
	cmd := exec.Command("java", "-version")

	// Capture the standard error output, as java prints version information to stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to get Java version: %v", err)
	}

	// Parse the stderr to extract the Java version information
	versionInfo := stderr.String()

	// Extract the Java version from the output
	// The output usually starts with something like:
	// java version "1.8.0_241"
	// or
	// openjdk version "11.0.8" 2020-07-14
	// or
	// openjdk version "17.0.1" 2021-10-19
	// So we need to extract the first version number found.

	// Split the version info by lines
	lines := strings.Split(versionInfo, "\n")

	// Extract the version from the first line (which typically contains the version information)
	for _, line := range lines {
		// Check for a line that contains the version info
		if strings.HasPrefix(line, "java version") || strings.HasPrefix(line, "openjdk version") {
			// Extract the version number (between quotes)
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("could not find Java version in output")
}
