package utils

import (
	"fmt"
	"time"
)

// Color codes for console output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

// Error prints a message in red
func Error(format string, a ...interface{}) {
	fmt.Printf(ColorRed+"ERROR: "+format+ColorReset+"\n", a...)
}

// Success prints a message in green
func Success(format string, a ...interface{}) {
	fmt.Printf(ColorGreen+format+ColorReset+"\n", a...)
}

// Warning prints a message in yellow
func Warning(format string, a ...interface{}) {
	fmt.Printf(ColorYellow+"WARNING: "+format+ColorReset+"\n", a...)
}

// Info prints a message in blue
func Info(format string, a ...interface{}) {
	fmt.Printf(ColorBlue+"INFO: "+format+ColorReset+"\n", a...)
}

// Log prints a regular message
func Log(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

// Logger represents a structured logger
type Logger struct {
	prefix string
}

// NewLogger creates a new logger with a prefix
func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

// Error logs an error message
func (l *Logger) Error(format string, a ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s%s ERROR: %s%s\n", timestamp, ColorRed, l.prefix, fmt.Sprintf(format, a...), ColorReset)
}

// Success logs a success message
func (l *Logger) Success(format string, a ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s%s SUCCESS: %s%s\n", timestamp, ColorGreen, l.prefix, fmt.Sprintf(format, a...), ColorReset)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, a ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s%s WARNING: %s%s\n", timestamp, ColorYellow, l.prefix, fmt.Sprintf(format, a...), ColorReset)
}

// Info logs an info message
func (l *Logger) Info(format string, a ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s%s INFO: %s%s\n", timestamp, ColorBlue, l.prefix, fmt.Sprintf(format, a...), ColorReset)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, a ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s DEBUG: %s%s\n", timestamp, l.prefix, fmt.Sprintf(format, a...), ColorReset)
}

// FormatFileSize formats a file size in bytes to human readable format
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
