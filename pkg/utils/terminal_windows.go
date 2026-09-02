//go:build windows

package utils

import (
	"os"
	"syscall"
)

const (
	enableVirtualTerminalProcessing = 0x0004
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

func isColorSupported() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}

	handle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil || handle == syscall.InvalidHandle {
		return false
	}

	var mode uint32
	err = syscall.GetConsoleMode(handle, &mode)
	if err != nil {
		// Not a console handle (e.g. redirected to a file or pipe)
		return false
	}

	if mode&enableVirtualTerminalProcessing != 0 {
		return true
	}

	r, _, _ := procSetConsoleMode.Call(uintptr(handle), uintptr(mode|enableVirtualTerminalProcessing))
	return r != 0
}
