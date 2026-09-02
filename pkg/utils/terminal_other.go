//go:build !windows

package utils

import (
	"os"
)

func isColorSupported() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}

	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}
