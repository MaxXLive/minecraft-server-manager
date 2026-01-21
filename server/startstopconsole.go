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

	cmd := exec.Command("screen", "-S", sessionName, "bash", "-c", fmt.Sprintf("cd %s && %s -Xms%dM -Xmx%dM -jar %s --nogui", dirPath, server.JavaPath, server.MaxRAM/2, server.MaxRAM, server.JarPath))

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

func StartInBackground() {
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

	for attempt := 1; attempt <= maxStartRetries; attempt++ {
		if attempt > 1 {
			log.Info(fmt.Sprintf("Retry attempt %d/%d...", attempt, maxStartRetries))
		}

		log.Info("Starting in background...")

		if err := startServerProcess(server); err != nil {
			log.Error(fmt.Sprintf("Error starting server: %v", err))
			continue
		}

		// If health check is not enabled, we're done
		if !server.HealthCheckEnabled {
			return
		}

		log.Info(fmt.Sprintf("Waiting for server to become healthy (timeout: %v)...", healthCheckTimeout))

		if WaitForHealthy() {
			log.Info("Server is healthy and ready!")
			return
		}

		log.Error("Health check failed - server did not respond in time")

		if IsServerRunning() {
			log.Info("Killing screen session...")
			if err := Kill(); err != nil {
				log.Error(fmt.Sprintf("Failed to kill session: %v", err))
			}
		}
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	log.Error(fmt.Sprintf("Failed to start server after %d attempts", maxStartRetries))
}

func startServerProcess(server config.Server) error {
	dirPath := filepath.Dir(server.JarPath)
	sessionName := GetSelectedServerSessionName()

	cmd := exec.Command("screen", "-dmS", sessionName, "bash", "-c",
		fmt.Sprintf("cd %s && %s -Xms%dM -Xmx%dM -jar %s --nogui",
			dirPath, server.JavaPath, server.MaxRAM/2, server.MaxRAM, server.JarPath))

	return cmd.Run()
}

func Stop() error {
	if !IsServerRunning() {
		return fmt.Errorf("Server is not running!")
	}

	sessionName := GetSelectedServerSessionName()
	cmd := exec.Command("screen", "-S", sessionName, "-X", "stuff", "\nstop\n")

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to send command:")
	}

	// Start the spinner in a separate goroutine
	stopChan := make(chan bool)
	go showSpinner(stopChan)

	// Wait until the screen session is no longer listed
	i := 0
	maxTries := 50

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
		return nil
	}

	return fmt.Errorf("Could not stop server. Try connecting and stopping manually")
}

func Restart() {
	log.Info("Restarting server...")

	if IsServerRunning() {
		err := Stop()
		if err != nil {
			log.Error(err)
			return
		}
		time.Sleep(2 * time.Second) // Wait a moment before starting again
	}

	StartInBackground()
}

func Kill() error {
	if !IsServerRunning() {
		return fmt.Errorf("server is not running")
	}

	sessionName := GetSelectedServerSessionName()
	cmd := exec.Command("screen", "-S", sessionName, "-X", "quit")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill session: %v", err)
	}

	// Wait for session to fully terminate
	time.Sleep(2 * time.Second)
	return nil
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
	server, err := GetSelected()
	if err != nil {
		log.Error(err)
		return
	}

	if IsServerRunning() {
		log.Info("Server is running...")

		if server.HealthCheckEnabled {
			if checkHealthEndpoint() {
				log.Info("Health check: healthy")
			} else {
				log.Error("Health check: unhealthy")
			}
		}
	} else {
		log.Info("Server is stopped!")
	}
}

func KillEntities() error {
	if !IsServerRunning() {
		return fmt.Errorf("server is not running")
	}

	sessionName := GetSelectedServerSessionName()
	cmd := exec.Command("screen", "-S", sessionName, "-X", "stuff", "\nkill @e[type=item]\n")
	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to send command")
	}
	return nil
}
