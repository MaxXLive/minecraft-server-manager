package main

import (
	"fmt"
	"minecraft-server-manager/log"
	"os/exec"
)

func main() {
	fmt.Println("--------- [ MINECRAFT SERVER MANAGER ] ---------")
	checkPrerequisites()

}

func checkPrerequisites() {
	cmd := exec.Command("screen", "--version")
	err := cmd.Run()
	if err != nil {
		log.Error("Screen is not installed! Please use apt install screen")
	}
}
