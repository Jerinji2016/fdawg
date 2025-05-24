package namer

import (
	"fmt"
	"os"
	"path/filepath"
)

// Platform represents a Flutter platform
type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformMacOS   Platform = "macos"
	PlatformLinux   Platform = "linux"
	PlatformWindows Platform = "windows"
	PlatformWeb     Platform = "web"
)

// AllPlatforms returns all supported platforms
func AllPlatforms() []Platform {
	return []Platform{
		PlatformAndroid,
		PlatformIOS,
		PlatformMacOS,
		PlatformLinux,
		PlatformWindows,
		PlatformWeb,
	}
}

// AppNameInfo holds app name information for a platform
type AppNameInfo struct {
	Platform    Platform `json:"platform"`
	DisplayName string   `json:"display_name"`
	InternalName string  `json:"internal_name,omitempty"`
	Available   bool     `json:"available"`
	Error       string   `json:"error,omitempty"`
}

// AppNameResult holds the result of getting app names
type AppNameResult struct {
	ProjectPath string        `json:"project_path"`
	AppNames    []AppNameInfo `json:"app_names"`
}

// SetAppNameRequest represents a request to set app names
type SetAppNameRequest struct {
	ProjectPath string            `json:"project_path"`
	Universal   string            `json:"universal,omitempty"`
	Platforms   map[Platform]string `json:"platforms,omitempty"`
}

// GetAppNames retrieves app names for specified platforms
func GetAppNames(projectPath string, platforms []Platform) (*AppNameResult, error) {
	// Validate project path
	if !isFlutterProject(projectPath) {
		return nil, fmt.Errorf("not a valid Flutter project: %s", projectPath)
	}

	result := &AppNameResult{
		ProjectPath: projectPath,
		AppNames:    make([]AppNameInfo, 0),
	}

	// If no platforms specified, get all available platforms
	if len(platforms) == 0 {
		platforms = getAvailablePlatforms(projectPath)
	}

	for _, platform := range platforms {
		appNameInfo := getAppNameForPlatform(projectPath, platform)
		result.AppNames = append(result.AppNames, appNameInfo)
	}

	return result, nil
}

// SetAppNames sets app names for specified platforms
func SetAppNames(request *SetAppNameRequest) error {
	// Validate project path
	if !isFlutterProject(request.ProjectPath) {
		return fmt.Errorf("not a valid Flutter project: %s", request.ProjectPath)
	}

	// Determine which platforms to update
	var platformsToUpdate map[Platform]string

	if request.Universal != "" {
		// Universal update - apply to all available platforms
		platformsToUpdate = make(map[Platform]string)
		availablePlatforms := getAvailablePlatforms(request.ProjectPath)
		for _, platform := range availablePlatforms {
			platformsToUpdate[platform] = request.Universal
		}
	} else if len(request.Platforms) > 0 {
		// Platform-specific updates
		platformsToUpdate = request.Platforms
	} else {
		return fmt.Errorf("either universal name or platform-specific names must be provided")
	}

	// Create backups before making changes
	backupDir := filepath.Join(request.ProjectPath, ".fdawg-backups", "namer")
	if err := createBackups(request.ProjectPath, platformsToUpdate, backupDir); err != nil {
		return fmt.Errorf("failed to create backups: %v", err)
	}

	// Apply changes to each platform
	var errors []string
	for platform, appName := range platformsToUpdate {
		if err := setAppNameForPlatform(request.ProjectPath, platform, appName); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", platform, err))
		}
	}

	if len(errors) > 0 {
		// If there were errors, attempt to restore from backups
		restoreFromBackups(request.ProjectPath, platformsToUpdate, backupDir)
		return fmt.Errorf("failed to set app names: %v", errors)
	}

	return nil
}

// isFlutterProject checks if the given path is a Flutter project
func isFlutterProject(projectPath string) bool {
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if _, err := os.Stat(pubspecPath); os.IsNotExist(err) {
		return false
	}

	// Additional check for Flutter-specific content in pubspec.yaml
	// This is a basic check - could be enhanced
	return true
}

// getAvailablePlatforms returns platforms that are available in the project
func getAvailablePlatforms(projectPath string) []Platform {
	var available []Platform

	for _, platform := range AllPlatforms() {
		if isPlatformAvailable(projectPath, platform) {
			available = append(available, platform)
		}
	}

	return available
}

// isPlatformAvailable checks if a platform directory exists
func isPlatformAvailable(projectPath string, platform Platform) bool {
	platformPath := filepath.Join(projectPath, string(platform))
	if _, err := os.Stat(platformPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// getAppNameForPlatform retrieves app name for a specific platform
func getAppNameForPlatform(projectPath string, platform Platform) AppNameInfo {
	info := AppNameInfo{
		Platform:  platform,
		Available: isPlatformAvailable(projectPath, platform),
	}

	if !info.Available {
		info.Error = "Platform not available in project"
		return info
	}

	switch platform {
	case PlatformAndroid:
		return getAndroidAppName(projectPath)
	case PlatformIOS:
		return getIOSAppName(projectPath)
	case PlatformMacOS:
		return getMacOSAppName(projectPath)
	case PlatformLinux:
		return getLinuxAppName(projectPath)
	case PlatformWindows:
		return getWindowsAppName(projectPath)
	case PlatformWeb:
		return getWebAppName(projectPath)
	default:
		info.Error = "Unsupported platform"
		return info
	}
}

// setAppNameForPlatform sets app name for a specific platform
func setAppNameForPlatform(projectPath string, platform Platform, appName string) error {
	if !isPlatformAvailable(projectPath, platform) {
		return fmt.Errorf("platform %s not available in project", platform)
	}

	switch platform {
	case PlatformAndroid:
		return setAndroidAppName(projectPath, appName)
	case PlatformIOS:
		return setIOSAppName(projectPath, appName)
	case PlatformMacOS:
		return setMacOSAppName(projectPath, appName)
	case PlatformLinux:
		return setLinuxAppName(projectPath, appName)
	case PlatformWindows:
		return setWindowsAppName(projectPath, appName)
	case PlatformWeb:
		return setWebAppName(projectPath, appName)
	default:
		return fmt.Errorf("unsupported platform: %s", platform)
	}
}

// createBackups creates backup files before making changes
func createBackups(projectPath string, platforms map[Platform]string, backupDir string) error {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	for platform := range platforms {
		if err := createPlatformBackup(projectPath, platform, backupDir); err != nil {
			return fmt.Errorf("failed to backup %s: %v", platform, err)
		}
	}

	return nil
}

// restoreFromBackups restores files from backups in case of errors
func restoreFromBackups(projectPath string, platforms map[Platform]string, backupDir string) {
	for platform := range platforms {
		restorePlatformBackup(projectPath, platform, backupDir)
	}
}
