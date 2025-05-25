package helpers

import (
	"os"
	"path/filepath"
)

// IsPlatformDirectoryAvailable checks if a platform directory exists
func IsPlatformDirectoryAvailable(projectPath, platformID string) bool {
	// This is a simplified check - the actual namer package has more sophisticated logic
	switch platformID {
	case "android":
		return DirExists(filepath.Join(projectPath, "android"))
	case "ios":
		return DirExists(filepath.Join(projectPath, "ios"))
	case "macos":
		return DirExists(filepath.Join(projectPath, "macos"))
	case "linux":
		return DirExists(filepath.Join(projectPath, "linux"))
	case "windows":
		return DirExists(filepath.Join(projectPath, "windows"))
	case "web":
		return DirExists(filepath.Join(projectPath, "web"))
	default:
		return false
	}
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
