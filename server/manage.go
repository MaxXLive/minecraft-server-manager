package server

import (
	"fmt"
	"github.com/google/uuid"
	"minecraft-server-manager/config"
	"minecraft-server-manager/log"
)

func List() ([]config.Server, error) {
	managerConfig, err := config.LoadConfig()
	if err != nil {
		log.Error("Could not open config " + err.Error())
		return nil, err
	}
	if len(managerConfig.Servers) == 0 {
		err := fmt.Errorf("No servers available. Use add command")
		log.Error(err)
		return managerConfig.Servers, err
	}

	return managerConfig.Servers, nil
}

func Add(server config.Server) {
	managerConfig, err := config.LoadConfig()
	if err != nil {
		log.Error("Could not get config: " + err.Error())
		return
	}

	server.ID = uuid.New().String()
	// Add to list and save
	managerConfig.Servers = append(managerConfig.Servers, server)

	err = config.SaveConfig(managerConfig)
	if err != nil {
		return
	}

	fmt.Println("\033[32mServer added successfully!\033[0m")
	Select(server.ID)
}

func Remove(id string) error {
	managerConfig, err := config.LoadConfig()
	if err != nil {
		return err
	}

	for i, server := range managerConfig.Servers {
		if server.ID == id {
			// Remove the server from the list
			managerConfig.Servers = append(managerConfig.Servers[:i], managerConfig.Servers[i+1:]...)
			return config.SaveConfig(managerConfig)
		}
	}
	return fmt.Errorf("server not found")
}

func Select(id string) error {
	managerConfig, err := config.LoadConfig()
	if err != nil {
		log.Error("Could not get config: " + err.Error())
		return err
	}

	for i, server := range managerConfig.Servers {
		managerConfig.Servers[i].IsSelected = server.ID == id
	}

	err = config.SaveConfig(managerConfig)
	if err != nil {
		return err
	}
	return nil
}

func GetSelected() (config.Server, error) {
	managerConfig, err := config.LoadConfig()
	if err != nil {
		log.Error("Could not get config: " + err.Error())
		return config.Server{}, err
	}

	for _, server := range managerConfig.Servers {
		if server.IsSelected {
			return server, nil
		}
	}
	return config.Server{}, fmt.Errorf("No server selected")
}

func GetSelectedServerSessionName() string {
	server, err := GetSelected()
	if err != nil {
		return ""
	}
	prefix, err := config.GetServerPrefix()
	if err != nil {
		return ""
	}
	return prefix + server.ID
}
