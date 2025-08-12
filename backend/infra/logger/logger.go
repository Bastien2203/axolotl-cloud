package logger

import (
	"fmt"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "INFO"
	LevelError LogLevel = "ERROR"
	LevelDebug LogLevel = "DEBUG"
)

type Logger struct {
	output func(level LogLevel, msg string, args ...any)
}

func NewLogger(output func(level LogLevel, msg string, args ...any)) *Logger {
	return &Logger{
		output: output,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.output(LevelInfo, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.output(LevelError, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.output(LevelDebug, msg, args...)
}

// Global log functions -------------

var globalLogger = NewLogger(func(level LogLevel, msg string, args ...any) {
	switch level {
	case LevelInfo:
		fmt.Printf("[INFO] "+msg+"\n", args...)
	case LevelError:
		fmt.Printf("[ERROR] "+msg+"\n", args...)
	case LevelDebug:
		fmt.Printf("[DEBUG] "+msg+"\n", args...)
	}
})

func Info(msg string) {
	globalLogger.Info(msg)
}
func Error(msg string, err error) {
	globalLogger.Error(msg, err)
}
func Debug(msg string) {
	globalLogger.Debug(msg)
}
