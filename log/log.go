package log

import (
	"fmt"
	"os"
)

const (
	red   = "\033[31m"
	reset = "\033[0m"
)

func Info(str interface{}) {
	fmt.Println(str)
}

func Error(err interface{}) {
	fmt.Fprintf(os.Stderr, "%sError: %s%s\n", red, err, reset)
}
