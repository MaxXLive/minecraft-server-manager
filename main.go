package main

import (
	"fmt"
	"minecraft-server-manager/backup"
	"minecraft-server-manager/cli"
	"minecraft-server-manager/log"
	"minecraft-server-manager/server"
	"minecraft-server-manager/update"
	"os"
	"os/exec"
)

var version = "1.6.7"

func main() {
	fmt.Println("--------- [ MINECRAFT SERVER MANAGER ] ---------")
	log.Init()
	checkPrerequisites()

	if len(os.Args) < 2 {
		cli.PrintHelp(os.Args[0], version)
		return
	}

	switch os.Args[1] {
	case "help":
		cli.PrintHelp(os.Args[0], version)
		return
	case "list":
		cli.PrintServerList()
		return
	case "add":
		cli.AddServer()
		return
	case "remove":
		cli.RemoveServer()
		return
	case "start":
		server.Start()
		return
	case "start-bg":
		server.StartInBackground()
		return
	case "stop":
		_ = server.Stop()
		return
	case "restart", "r":
		server.Restart()
		return
	case "kill":
		_ = server.Kill()
		return
	case "console", "c":
		server.Attach()
		return
	case "status", "s":
		server.Status()
		return
	case "select":
		cli.SelectServer()
		return
	case "version":
		log.Info("Version: " + version)
		return
	case "check":
		update.CheckForUpdate(version)
		return
	case "update":
		update.RunUpdate(version, includes(os.Args, "--force"))
		return
	case "backup", "b":
		backup.Start()
		return
	case "logfile":
		cli.LogFile(os.Args[2:])
		return
	default:
		cli.PrintHelp(os.Args[0], version)
		return
	}
}

func checkPrerequisites() {
	err := exec.Command("which", "screen").Run()
	if err != nil {
		log.Error("Screen is not installed! Please use apt install screen")
	}
}

func includes(args []string, target string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}
	return false
}
