package bundler

import (
	"fmt"
	"os"
	"path/filepath"
)

// Platform represents a supported platform
type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformMacOS   Platform = "macos"
	PlatformLinux   Platform = "linux"
	PlatformWindows Platform = "windows"
	PlatformWeb     Platform = "web"
)

// BundleIDInfo represents bundle ID information for a platform
type BundleIDInfo struct {
	Platform  Platform `json:"platform"`
	BundleID  string   `json:"bundle_id"`
	Namespace string   `json:"namespace,omitempty"` // For Android
	Available bool     `json:"available"`
	Error     string   `json:"error,omitempty"`
}

// BundleIDResult represents the result of getting bundle IDs
type BundleIDResult struct {
	ProjectPath string         `json:"project_path"`
	BundleIDs   []BundleIDInfo `json:"bundle_ids"`
}

// BundleIDRequest represents a request to set bundle IDs
type BundleIDRequest struct {
	ProjectPath string              `json:"project_path"`
	Universal   string              `json:"universal,omitempty"`
	Platforms   map[Platform]string `json:"platforms,omitempty"`
}

// GetBundleIDs retrieves bundle IDs for specified platforms
func GetBundleIDs(projectPath string, platforms []Platform) (*BundleIDResult, error) {
	// Validate project path
	if !isFlutterProject(projectPath) {
		return nil, fmt.Errorf("not a valid Flutter project: %s", projectPath)
	}

	result := &BundleIDResult{
		ProjectPath: projectPath,
		BundleIDs:   make([]BundleIDInfo, 0),
	}

	// If no platforms specified, get all available platforms
	if len(platforms) == 0 {
		platforms = getAvailablePlatforms(projectPath)
	}

	for _, platform := range platforms {
		bundleIDInfo := getBundleIDForPlatform(projectPath, platform)
		result.BundleIDs = append(result.BundleIDs, bundleIDInfo)
	}

	return result, nil
}

// SetBundleIDs sets bundle IDs for specified platforms
func SetBundleIDs(request *BundleIDRequest) error {
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
		return fmt.Errorf("either universal bundle ID or platform-specific bundle IDs must be provided")
	}

	// Create backups before making changes
	backupDir := filepath.Join(request.ProjectPath, ".fdawg-backups", "bundler")
	if err := createBackups(request.ProjectPath, platformsToUpdate, backupDir); err != nil {
		return fmt.Errorf("failed to create backups: %v", err)
	}

	// Apply changes to each platform
	var errors []string
	for platform, bundleID := range platformsToUpdate {
		if err := setBundleIDForPlatform(request.ProjectPath, platform, bundleID); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", platform, err))
		}
	}

	if len(errors) > 0 {
		// If there were errors, attempt to restore from backups
		restoreFromBackups(request.ProjectPath, platformsToUpdate, backupDir)
		return fmt.Errorf("failed to set bundle IDs: %v", errors)
	}

	return nil
}

// isFlutterProject checks if the given path is a Flutter project
func isFlutterProject(projectPath string) bool {
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if _, err := os.Stat(pubspecPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// getAvailablePlatforms returns all available platforms in the project
func getAvailablePlatforms(projectPath string) []Platform {
	var platforms []Platform

	platformDirs := map[Platform]string{
		PlatformAndroid: "android",
		PlatformIOS:     "ios",
		PlatformMacOS:   "macos",
		PlatformLinux:   "linux",
		PlatformWindows: "windows",
		PlatformWeb:     "web",
	}

	for platform := range platformDirs {
		if isPlatformAvailable(projectPath, platform) {
			platforms = append(platforms, platform)
		}
	}

	return platforms
}

// isPlatformAvailable checks if a platform is available in the project
func isPlatformAvailable(projectPath string, platform Platform) bool {
	platformDirs := map[Platform]string{
		PlatformAndroid: "android",
		PlatformIOS:     "ios",
		PlatformMacOS:   "macos",
		PlatformLinux:   "linux",
		PlatformWindows: "windows",
		PlatformWeb:     "web",
	}

	dir, exists := platformDirs[platform]
	if !exists {
		return false
	}

	platformPath := filepath.Join(projectPath, dir)
	if _, err := os.Stat(platformPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// getBundleIDForPlatform retrieves bundle ID for a specific platform
func getBundleIDForPlatform(projectPath string, platform Platform) BundleIDInfo {
	info := BundleIDInfo{
		Platform:  platform,
		Available: isPlatformAvailable(projectPath, platform),
	}

	if !info.Available {
		info.Error = "Platform not available in project"
		return info
	}

	switch platform {
	case PlatformAndroid:
		return getAndroidBundleID(projectPath)
	case PlatformIOS:
		return getIOSBundleID(projectPath)
	case PlatformMacOS:
		return getMacOSBundleID(projectPath)
	case PlatformLinux:
		return getLinuxBundleID(projectPath)
	case PlatformWindows:
		return getWindowsBundleID(projectPath)
	case PlatformWeb:
		return getWebBundleID(projectPath)
	default:
		info.Error = "Unsupported platform"
		return info
	}
}

// setBundleIDForPlatform sets bundle ID for a specific platform
func setBundleIDForPlatform(projectPath string, platform Platform, bundleID string) error {
	if !isPlatformAvailable(projectPath, platform) {
		return fmt.Errorf("platform %s not available in project", platform)
	}

	switch platform {
	case PlatformAndroid:
		return setAndroidBundleID(projectPath, bundleID)
	case PlatformIOS:
		return setIOSBundleID(projectPath, bundleID)
	case PlatformMacOS:
		return setMacOSBundleID(projectPath, bundleID)
	case PlatformLinux:
		return setLinuxBundleID(projectPath, bundleID)
	case PlatformWindows:
		return setWindowsBundleID(projectPath, bundleID)
	case PlatformWeb:
		return setWebBundleID(projectPath, bundleID)
	default:
		return fmt.Errorf("unsupported platform: %s", platform)
	}
}

// createBackups creates backups of platform-specific files before making changes
func createBackups(projectPath string, platformsToUpdate map[Platform]string, backupDir string) error {
	// Create backup directory
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	// Create backups for each platform
	for platform := range platformsToUpdate {
		if err := createPlatformBackup(projectPath, platform, backupDir); err != nil {
			return fmt.Errorf("failed to backup %s: %v", platform, err)
		}
	}

	return nil
}

// restoreFromBackups restores files from backups in case of errors
func restoreFromBackups(projectPath string, platformsToUpdate map[Platform]string, backupDir string) {
	for platform := range platformsToUpdate {
		restorePlatformBackup(projectPath, platform, backupDir)
	}
}
