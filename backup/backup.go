package backup

import (
	"fmt"
	"minecraft-server-manager/log"
	"minecraft-server-manager/server"
	"minecraft-server-manager/system"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Start() {
	if server.IsServerRunning() {
		err := broadcastMessage("Server will restart in 30 seconds...")
		if err != nil {
			log.Error(err)
			return
		}
		time.Sleep(25 * time.Second)
		for i := 5; i > 0; i-- {
			err := broadcastMessage(fmt.Sprintf("Server will restart in %d seconds...", i))
			if err != nil {
				log.Error(err)
				return
			}
			time.Sleep(time.Second)
		}
		err = server.Stop()
		if err != nil {
			log.Error(err)
			return
		}
	}

	server, err := server.GetSelected()
	if err != nil {
		log.Error(err)
		return
	}
	dirPath := filepath.Dir(server.JarPath)

	uploadServerData(dirPath)
}

// uploadServerData commits and pushes the latest changes to a Git repository
func uploadServerData(dir string) {

	if err := os.Chdir(dir); err != nil {
		log.Error(fmt.Sprintf("Failed to change directory: %v", err))
		return
	}

	// Capture the current timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("Uploading server data: %s\n", timestamp)

	// Run Git commands
	commands := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", fmt.Sprintf("update %s", timestamp)},
		{"git", "push"},
	}

	for _, cmd := range commands {
		fmt.Println("Running:", cmd)
		if err := system.RunCommand(cmd...); err != nil {
			log.Error(fmt.Sprintf("Command failed: %v", err))
			break
		}
	}

	fmt.Println("Upload done")
	server.StartInBackground()
}

func broadcastMessage(message string) error {
	if !server.IsServerRunning() {
		return fmt.Errorf("Server is not running!")
	}

	sessionName := server.GetSelectedServerSessionName()
	cmd := exec.Command("screen", "-S", sessionName, "-X", "stuff", fmt.Sprintf("say %s\n", message))

	// Run the command
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
