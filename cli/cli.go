package cli

import (
	"fmt"
	"github.com/chzyer/readline"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"minecraft-server-manager/config"
	"minecraft-server-manager/log"
	"minecraft-server-manager/server"
	"minecraft-server-manager/system"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func AddServer() {
	// Ask for server name
	name := promptInput("Enter server name")
	if name == "" {
		fmt.Println("\033[31mError: Server name cannot be empty\033[0m")
		return
	}

	// Ask for path to server .jar file
	jarPath := promptInput("Enter full path to server .jar file")
	absPath, err := filepath.Abs(jarPath)
	if err != nil {
		fmt.Println("\033[31mError: Invalid file path\033[0m")
		return
	}
	// Ensure the .jar file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Error("File does not exist")
		return
	}

	// Ask for path to java executable
	version, err := system.GetJavaVersion()
	javaStr := "No default found"
	if err == nil {
		javaStr = fmt.Sprintf("default: %s", version)
	}
	javaPath := promptInput(fmt.Sprintf("Enter full path to java executable (%s)", javaStr))
	if len(javaPath) == 0 {
		javaPath = "java"
	} else {
		javaPath, err = filepath.Abs(javaPath)
		if err != nil {
			fmt.Println("\033[31mError: Invalid file path\033[0m")
			return
		}
	}

	// Ask for max RAM in MB
	ramStr := promptInput("Enter max RAM (MB)")
	maxRam, err := strconv.Atoi(ramStr)
	if err != nil || maxRam <= 0 {
		fmt.Println("\033[31mError: Invalid RAM value\033[0m")
		return
	}

	serverID := uuid.New().String()

	// Create new server entry
	newServer := config.Server{
		ID:       serverID,
		Name:     name,
		JarPath:  jarPath,
		JavaPath: javaPath,
		MaxRAM:   maxRam,
	}

	server.Add(newServer)
}

func RemoveServer() {
	servers, err := server.List()
	if err != nil {
		return
	}
	PrintServerList()

	str := promptInput("Enter index of server to remove")
	selectedIndex, err := strconv.Atoi(str)
	if err != nil || selectedIndex <= 0 || selectedIndex > len(servers) {
		fmt.Println("\033[31mError: Invalid value\033[0m")
		return
	}

	str = promptInput("Are you sure you want to delete this server? (yes/no)")
	switch str {
	case "yes", "y":
		err = server.Remove(servers[selectedIndex-1].ID)
		if err == nil {
			println("Successfully removed server: " + servers[selectedIndex-1].Name)
		}
		return
	default:
		log.Info("Deleting aborted")
	}
}

func PrintServerList() {
	servers, err := server.List()
	if err != nil {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("Server List")
	t.AppendHeader(table.Row{"#", "Name", "Selected"})

	for i, server := range servers {
		// Check if server is selected and format accordingly
		formatSelected := " "
		if server.IsSelected {
			formatSelected = "X"
		}

		t.AppendRows([]table.Row{{i + 1, server.Name, formatSelected}})
	}
	t.Render()
}

func PrintHelp(appName string, version string) {
	fmt.Println("")
	fmt.Println(fmt.Sprintf("Usage: %s [command]", appName))
	fmt.Println("")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Command", "Description"})
	t.AppendRows([]table.Row{{"start", "Start the minecraft server"}})
	t.AppendRows([]table.Row{{"stop", "Stop the minecraft server"}})
	t.AppendRows([]table.Row{{"console", "Attach to the minecraft server's console"}})
	t.AppendRows([]table.Row{{"status", "Show the status of the minecraft server"}})
	t.AppendRows([]table.Row{{"restart", "Restart the minecraft server"}})
	t.AppendSeparator()
	t.AppendRows([]table.Row{{"list", "List saved servers in config"}})
	t.AppendRows([]table.Row{{"select", "Select from servers in config"}})
	t.AppendRows([]table.Row{{"add", "Add new server to config"}})
	t.AppendRows([]table.Row{{"remove", "Remove a server from config"}})
	t.AppendSeparator()
	t.AppendRows([]table.Row{{"help", "Show this help message"}})
	t.AppendFooter(table.Row{"Version", version})
	t.Render()

}

func SelectServer() {
	servers, err := server.List()
	if err != nil {
		return
	}
	PrintServerList()

	str := promptInput("Enter index of server to select")
	selectedIndex, err := strconv.Atoi(str)
	if err != nil || selectedIndex <= 0 || selectedIndex > len(servers) {
		fmt.Println("\033[31mError: Invalid value\033[0m")
		return
	}

	err = server.Select(servers[selectedIndex-1].ID)
	if err != nil {
		return
	}
	println("Successfully selected server: " + servers[selectedIndex-1].Name)
}

func promptInput(prompt string) string {
	// Setup readline for interactive input
	rl, err := readline.New(prompt + ": ")
	if err != nil {
		fmt.Println("\033[31mError creating readline instance:\033[0m", err)
		return ""
	}
	defer rl.Close()

	// Read input from the user
	line, err := rl.Readline()
	if err != nil {
		fmt.Println("\033[31mError reading input:\033[0m", err)
		return ""
	}

	return strings.TrimSpace(line)
}
