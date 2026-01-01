// Package logger provides a simple structured logging utility.
package logger

import (
	"log"
	"os"
)

// Logger provides structured logging capabilities.
type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	debug *log.Logger
}

// New creates a new Logger instance.
func New() *Logger {
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC

	return &Logger{
		info:  log.New(os.Stdout, "[INFO]  ", flags),
		warn:  log.New(os.Stdout, "[WARN]  ", flags),
		error: log.New(os.Stderr, "[ERROR] ", flags),
		debug: log.New(os.Stdout, "[DEBUG] ", flags),
	}
}

// Info logs an informational message.
func (l *Logger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.info.Printf(msg, args...)
	} else {
		l.info.Println(msg)
	}
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.warn.Printf(msg, args...)
	} else {
		l.warn.Println(msg)
	}
}

// Error logs an error message.
func (l *Logger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.error.Printf(msg, args...)
	} else {
		l.error.Println(msg)
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.debug.Printf(msg, args...)
	} else {
		l.debug.Println(msg)
	}
}

// Fatal logs an error message and exits the program.
func (l *Logger) Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.error.Fatalf(msg, args...)
	} else {
		l.error.Fatal(msg)
	}
}
