package config

import (
	"encoding/json"
	"fmt"
	"minecraft-server-manager/log"
	"os"
	"path/filepath"
)

func getConfigFilePath() (string, error) {
	// Get the path to the currently running executable
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	// Extract the directory from the executable path
	executableDir := filepath.Dir(executablePath)

	// Construct the full path to the config file (config.json)
	configFilePath := filepath.Join(executableDir, "config.json")

	return configFilePath, nil
}

func LoadConfig() (ManagerConfig, error) {

	var config ManagerConfig
	config.ScreenName = "minecraft_server_"

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return ManagerConfig{}, err
	}
	// Check if config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return config, nil
	}

	// Read the file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Error("Could not load config: " + err.Error())
		return config, err
	}

	// Parse JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Error("Could not load config: " + err.Error())
		return config, err
	}

	return config, nil
}

// Save the server list to config.json

func SaveConfig(config ManagerConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Error("Could not save config: " + err.Error())
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, data, 0644)
}

func GetServerPrefix() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}
	return config.ScreenName, nil
}
