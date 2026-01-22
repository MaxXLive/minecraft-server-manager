package server

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"minecraft-server-manager/log"
)

// GetServerPropertiesPath returns the path to server.properties for the given jar path
func GetServerPropertiesPath(jarPath string) string {
	serverDir := filepath.Dir(jarPath)
	return filepath.Join(serverDir, "server.properties")
}

// ReadServerProperties reads the server.properties file and returns a map of key-value pairs
func ReadServerProperties(jarPath string) (map[string]string, error) {
	propsPath := GetServerPropertiesPath(jarPath)
	props := make(map[string]string)

	file, err := os.Open(propsPath)
	if err != nil {
		return props, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			props[parts[0]] = parts[1]
		}
	}

	return props, scanner.Err()
}

// SetServerProperty updates a single property in server.properties
func SetServerProperty(jarPath string, key string, value string) error {
	propsPath := GetServerPropertiesPath(jarPath)

	// Read the file
	file, err := os.Open(propsPath)
	if err != nil {
		return err
	}

	var lines []string
	found := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+"=") {
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
			found = true
		} else {
			lines = append(lines, line)
		}
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		return err
	}

	// If key wasn't found, append it
	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	// Write the file back
	return os.WriteFile(propsPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

// SetServerPort updates the server-port in server.properties
func SetServerPort(jarPath string, port int) error {
	log.Info(fmt.Sprintf("Setting server port to %d in server.properties", port))
	return SetServerProperty(jarPath, "server-port", fmt.Sprintf("%d", port))
}