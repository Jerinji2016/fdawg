package helpers

import "os"

// IsPlatformDirectoryAvailable checks if a platform directory exists
func IsPlatformDirectoryAvailable(projectPath, platformID string) bool {
	// This is a simplified check - the actual namer package has more sophisticated logic
	switch platformID {
	case "android":
		return DirExists(projectPath + "/android")
	case "ios":
		return DirExists(projectPath + "/ios")
	case "macos":
		return DirExists(projectPath + "/macos")
	case "linux":
		return DirExists(projectPath + "/linux")
	case "windows":
		return DirExists(projectPath + "/windows")
	case "web":
		return DirExists(projectPath + "/web")
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
