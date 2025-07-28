package logger

import "fmt"

type Logger struct{}

func Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}
func Error(msg string, err error) {
	fmt.Printf("[ERROR] %s: %v\n", msg, err)
}
func Debug(msg string) {
	fmt.Printf("[DEBUG] %s\n", msg)
}
