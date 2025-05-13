package utils

import "fmt"

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