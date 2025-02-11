package config

import (
	"encoding/json"
	"minecraft-server-manager/log"
	"os"
)

const configFile = "config.json"

// Load existing config or create a new one

func LoadConfig() (ManagerConfig, error) {

	var config ManagerConfig
	config.ScreenName = "minecraft_server_"

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return config, nil
	}

	// Read the file
	data, err := os.ReadFile(configFile)
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

	return os.WriteFile(configFile, data, 0644)
}

func GetServerPrefix() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}
	return config.ScreenName, nil
}
