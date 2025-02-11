package main

import (
	"fmt"
	"minecraft-server-manager/cli"
	"minecraft-server-manager/log"
	"minecraft-server-manager/server"
	"os"
	"os/exec"
)

var version = "1.1"

func main() {
	fmt.Println("--------- [ MINECRAFT SERVER MANAGER ] ---------")
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
	case "stop":
		server.Stop()
		return
	case "console":
		server.Attach()
		return
	case "status":
		server.Status()
		return
	case "select":
		cli.SelectServer()
		return
	case "version":
		log.Info("Version: " + version)
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
