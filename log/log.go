package log

import (
	"fmt"
	"minecraft-server-manager/config"
	"os"
	"time"
)

const (
	red   = "\033[31m"
	reset = "\033[0m"
)

var (
	fileEnabled bool
	filePath    string
)

func Init() {
	fileEnabled = config.IsLogFileEnabled()
	if fileEnabled {
		path, err := config.GetLogFilePath()
		if err == nil {
			filePath = path
			writeSessionSeparator()
		} else {
			fileEnabled = false
		}
	}
}

func writeSessionSeparator() {
	if !fileEnabled || filePath == "" {
		return
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	f.WriteString(fmt.Sprintf("\n%s\n", "════════════════════════════════════════════════════════════"))
	f.WriteString(fmt.Sprintf("[%s] Session started\n", timestamp))
	f.WriteString(fmt.Sprintf("%s\n", "════════════════════════════════════════════════════════════"))
}

func writeToFile(level string, message string) {
	if !fileEnabled || filePath == "" {
		return
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	f.WriteString(fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message))
}

func Info(str interface{}) {
	fmt.Println(str)
	writeToFile("INFO", fmt.Sprintf("%v", str))
}

func Error(err interface{}) {
	fmt.Fprintf(os.Stderr, "%sError: %s%s\n", red, err, reset)
	writeToFile("ERROR", fmt.Sprintf("%v", err))
}
