package server

import (
	"fmt"
	"minecraft-server-manager/config"
	"minecraft-server-manager/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func Start() {
	server, err := GetSelected()
	if err != nil {
		log.Error(err)
		return
	}
	if IsServerRunning() {
		log.Error("Server is already running! Stop server first")
		return
	}

	log.Info("Selected Server: " + server.Name)
	log.Info("Starting...")

	dirPath := filepath.Dir(server.JarPath)

	sessionName := GetSelectedServerSessionName()

	cmd := exec.Command("screen", "-S", sessionName, "bash", "-c", fmt.Sprintf("cd %s && %s -Xms%dM -Xmx%dM -jar %s --nogui", dirPath, server.JavaPath, server.MaxRAM, server.MaxRAM, server.JavaPath))

	// Set the standard input, output, and error to the current process
	cmd.Stdin = os.Stdin   // Enable user input
	cmd.Stdout = os.Stdout // Show command output
	cmd.Stderr = os.Stderr // Show error messages

	// Start the command and attach to it
	err = cmd.Run()
	if err != nil {
		fmt.Println("\033[31mError executing command:\033[0m", err)
		return
	}
}

func Stop() {
	if !IsServerRunning() {
		log.Error("Server is not running!")
		return
	}

	sessionName := GetSelectedServerSessionName()
	cmd := exec.Command("screen", "-S", sessionName, "-X", "stuff", "stop\n")

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Error("failed to send command:")
	}

	// Start the spinner in a separate goroutine
	stopChan := make(chan bool)
	go showSpinner(stopChan)

	// Wait until the screen session is no longer listed
	i := 0
	maxTries := 20

	for i = 0; i < maxTries; i++ {
		if !IsServerRunning() {
			break
		}
		time.Sleep(500 * time.Millisecond) // Check every 500ms
	}

	// Stop the spinner and print "Done"
	stopChan <- true
	if i < maxTries {
		fmt.Println("Done!")
	} else {
		log.Error("Could not stop server. Try connecting and stopping manually")
	}
}

func IsServerRunning() bool {
	server, err := GetSelected()
	if err != nil {
		log.Error(err)
		return false
	}
	prefix, err := config.GetServerPrefix()
	if err != nil {
		return false
	}
	sessionName := prefix + server.ID

	output, _ := exec.Command("screen", "-list").Output()
	return strings.Contains(string(output), sessionName)
}

func showSpinner(stopChan chan bool) {
	spinnerChars := []rune{'|', '/', '-', '\\'}
	i := 0
	for {
		select {
		case <-stopChan:
			fmt.Print("\r\033[K") // Clear spinner line
			return
		default:
			fmt.Printf("\rStopping server... %c", spinnerChars[i%len(spinnerChars)])
			i++
			time.Sleep(150 * time.Millisecond) // Adjust speed of the spinner
		}
	}
}

func Attach() {
	if !IsServerRunning() {
		log.Error("Server is not running!")
		return
	}

	sessionName := GetSelectedServerSessionName()
	cmd := exec.Command("screen", "-r", sessionName)

	// Set the standard input, output, and error to the current process
	cmd.Stdin = os.Stdin   // Enable user input
	cmd.Stdout = os.Stdout // Show command output
	cmd.Stderr = os.Stderr // Show error messages

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Error("Failed to attach to server" + err.Error())
	}
}

func Status() {
	if IsServerRunning() {
		println("Server is running...")
	} else {
		println("Server is stopped!")
	}
}
