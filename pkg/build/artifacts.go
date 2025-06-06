package build

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/namer"
)

// ArtifactManager manages build artifacts
type ArtifactManager struct {
	ProjectPath string
	Config      *ArtifactsConfig
}

// NewArtifactManager creates a new artifact manager
func NewArtifactManager(projectPath string, config *ArtifactsConfig) *ArtifactManager {
	if config == nil {
		// Use default config
		defaultConfig := DefaultBuildConfig()
		config = &defaultConfig.Artifacts
	}

	return &ArtifactManager{
		ProjectPath: projectPath,
		Config:      config,
	}
}

// OrganizeArtifact organizes a build artifact into the proper directory structure
func (am *ArtifactManager) OrganizeArtifact(artifact *BuildArtifact) error {
	// Create date-based directory
	dateDir := artifact.BuildTime.Format(am.Config.Organization.DateFormat)
	outputPath := filepath.Join(am.ProjectPath, am.Config.BaseOutputDir, dateDir)

	// Add platform directory if enabled
	if am.Config.Organization.ByPlatform {
		outputPath = filepath.Join(outputPath, string(artifact.Platform))
	}

	// Add build type directory if enabled and build type exists
	if am.Config.Organization.ByBuildType && artifact.BuildType != "" {
		outputPath = filepath.Join(outputPath, artifact.BuildType)
	}

	// Ensure directory exists
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate final filename
	finalName := am.generateArtifactName(artifact)
	finalPath := filepath.Join(outputPath, finalName)

	// Update artifact with new information
	artifact.FileName = finalName
	artifact.AppName = am.getAppName()
	artifact.Version = am.getAppVersion()

	// Move/copy artifact to organized location
	if err := am.moveArtifact(artifact.FilePath, finalPath); err != nil {
		return fmt.Errorf("failed to move artifact: %w", err)
	}

	// Update artifact file path
	artifact.FilePath = finalPath

	// Get file size
	if stat, err := os.Stat(finalPath); err == nil {
		artifact.Size = stat.Size()
	}

	return nil
}

// generateArtifactName generates the final artifact name based on configuration
func (am *ArtifactManager) generateArtifactName(artifact *BuildArtifact) string {
	// Get app name
	appName := am.getAppName()
	if appName == "" {
		appName = am.Config.Naming.FallbackAppName
	}

	// Get version
	version := am.getAppVersion()
	if version == "" {
		version = "1.0.0"
	}

	// Get architecture
	arch := artifact.Architecture
	if arch == "" {
		arch = "universal"
	}

	// Get file extension from original
	ext := filepath.Ext(artifact.FileName)

	// Apply naming pattern
	name := am.Config.Naming.Pattern
	name = strings.ReplaceAll(name, "{app_name}", appName)
	name = strings.ReplaceAll(name, "{version}", version)
	name = strings.ReplaceAll(name, "{arch}", arch)
	name = strings.ReplaceAll(name, "{platform}", string(artifact.Platform))
	name = strings.ReplaceAll(name, "{build_type}", artifact.BuildType)

	// Add extension
	name += ext

	// Sanitize filename
	return am.sanitizeFileName(name)
}

// getAppName gets the app name based on configuration
func (am *ArtifactManager) getAppName() string {
	// Try to get from namer first (most reliable for cross-platform names)
	if appNamesResult, err := namer.GetAppNames(am.ProjectPath, nil); err == nil {
		// Look for a universal name or use the first available
		for _, nameInfo := range appNamesResult.AppNames {
			if nameInfo.DisplayName != "" {
				// Clean the name for filename use
				return am.sanitizeAppName(nameInfo.DisplayName)
			}
		}
	}

	// Fallback to pubspec.yaml
	if result, err := flutter.ValidateProject(am.ProjectPath); err == nil {
		if result.PubspecInfo != nil && result.PubspecInfo.Name != "" {
			return am.sanitizeAppName(result.PubspecInfo.Name)
		}
	}

	return am.Config.Naming.FallbackAppName
}

// getAppVersion gets the app version from pubspec.yaml
func (am *ArtifactManager) getAppVersion() string {
	if result, err := flutter.ValidateProject(am.ProjectPath); err == nil {
		if result.PubspecInfo != nil && result.PubspecInfo.Version != "" {
			// Remove build number if present (e.g., "1.0.0+1" -> "1.0.0")
			version := result.PubspecInfo.Version
			if plusIndex := strings.Index(version, "+"); plusIndex != -1 {
				version = version[:plusIndex]
			}
			return version
		}
	}

	return "1.0.0"
}

// sanitizeAppName sanitizes app name for filename use
func (am *ArtifactManager) sanitizeAppName(name string) string {
	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")

	// Remove or replace invalid filename characters
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	name = reg.ReplaceAllString(name, "")

	// Remove multiple consecutive underscores
	reg = regexp.MustCompile(`_+`)
	name = reg.ReplaceAllString(name, "_")

	// Trim underscores from start and end
	name = strings.Trim(name, "_")

	return name
}

// sanitizeFileName sanitizes a filename
func (am *ArtifactManager) sanitizeFileName(filename string) string {
	// Remove or replace invalid filename characters
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	filename = reg.ReplaceAllString(filename, "_")

	// Remove multiple consecutive underscores
	reg = regexp.MustCompile(`_+`)
	filename = reg.ReplaceAllString(filename, "_")

	return filename
}

// moveArtifact moves an artifact from source to destination
func (am *ArtifactManager) moveArtifact(sourcePath, destPath string) error {
	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("source artifact not found: %s", sourcePath)
	}

	// If destination already exists, remove it
	if _, err := os.Stat(destPath); err == nil {
		if err := os.Remove(destPath); err != nil {
			return fmt.Errorf("failed to remove existing artifact: %w", err)
		}
	}

	// For directories (like .xcarchive), use recursive copy
	if info, err := os.Stat(sourcePath); err == nil && info.IsDir() {
		return am.copyDirectory(sourcePath, destPath)
	}

	// For files, use rename (move)
	if err := os.Rename(sourcePath, destPath); err != nil {
		// If rename fails (cross-device), try copy and delete
		return am.copyFile(sourcePath, destPath)
	}

	return nil
}

// copyFile copies a file from source to destination
func (am *ArtifactManager) copyFile(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := destFile.ReadFrom(sourceFile); err != nil {
		return err
	}

	// Remove source file after successful copy
	return os.Remove(sourcePath)
}

// copyDirectory recursively copies a directory
func (am *ArtifactManager) copyDirectory(sourcePath, destPath string) error {
	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		sourceEntryPath := filepath.Join(sourcePath, entry.Name())
		destEntryPath := filepath.Join(destPath, entry.Name())

		if entry.IsDir() {
			if err := am.copyDirectory(sourceEntryPath, destEntryPath); err != nil {
				return err
			}
		} else {
			if err := am.copyFile(sourceEntryPath, destEntryPath); err != nil {
				return err
			}
		}
	}

	// Remove source directory after successful copy
	return os.RemoveAll(sourcePath)
}

// GetOutputDir returns the base output directory
func (am *ArtifactManager) GetOutputDir() string {
	return filepath.Join(am.ProjectPath, am.Config.BaseOutputDir)
}

// ListArtifacts lists artifacts based on filters
func (am *ArtifactManager) ListArtifacts(filters ArtifactFilters) ([]*BuildArtifact, error) {
	outputDir := am.GetOutputDir()

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return []*BuildArtifact{}, nil
	}

	var artifacts []*BuildArtifact

	// Walk through the output directory
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and log files
		if info.IsDir() || strings.HasSuffix(path, ".log") {
			return nil
		}

		// Parse artifact information from path
		artifact := am.parseArtifactFromPath(path, info)
		if artifact == nil {
			return nil
		}

		// Apply filters
		if filters.Date != "" && artifact.BuildTime.Format(am.Config.Organization.DateFormat) != filters.Date {
			return nil
		}

		if filters.Platform != "" && string(artifact.Platform) != filters.Platform {
			return nil
		}

		artifacts = append(artifacts, artifact)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}

	// Sort artifacts by build time (newest first)
	sort.Slice(artifacts, func(i, j int) bool {
		return artifacts[i].BuildTime.After(artifacts[j].BuildTime)
	})

	return artifacts, nil
}

// parseArtifactFromPath parses artifact information from file path
func (am *ArtifactManager) parseArtifactFromPath(path string, info os.FileInfo) *BuildArtifact {
	relPath, err := filepath.Rel(am.GetOutputDir(), path)
	if err != nil {
		return nil
	}

	parts := strings.Split(relPath, string(filepath.Separator))
	if len(parts) < 2 {
		return nil
	}

	// Parse date from first part
	dateStr := parts[0]
	buildTime, err := time.Parse(am.Config.Organization.DateFormat, dateStr)
	if err != nil {
		// Try to parse as a date
		buildTime = info.ModTime()
	}

	// Determine platform and build type from path structure
	var platform Platform
	var buildType string

	if len(parts) >= 3 {
		platform = Platform(parts[1])
		buildType = parts[2]
	} else if len(parts) >= 2 {
		platform = Platform(parts[1])
	}

	// Parse filename for additional information
	filename := info.Name()
	architecture := am.parseArchitectureFromFilename(filename)

	return &BuildArtifact{
		Platform:     platform,
		BuildType:    buildType,
		Architecture: architecture,
		FileName:     filename,
		FilePath:     path,
		Size:         info.Size(),
		BuildTime:    buildTime,
	}
}

// parseArchitectureFromFilename attempts to parse architecture from filename
func (am *ArtifactManager) parseArchitectureFromFilename(filename string) string {
	// Common architecture patterns
	archPatterns := []string{
		"arm64-v8a", "armeabi-v7a", "x86_64", "x86",
		"arm64", "amd64", "universal",
	}

	lowerFilename := strings.ToLower(filename)
	for _, arch := range archPatterns {
		if strings.Contains(lowerFilename, strings.ToLower(arch)) {
			return arch
		}
	}

	return "universal"
}

// GetBuildStatus returns the current build status
func (am *ArtifactManager) GetBuildStatus() (*BuildStatus, error) {
	artifacts, err := am.ListArtifacts(ArtifactFilters{})
	if err != nil {
		return nil, err
	}

	status := &BuildStatus{
		TotalArtifacts: len(artifacts),
	}

	if len(artifacts) > 0 {
		status.LastBuildTime = artifacts[0].BuildTime

		// Group by date for recent builds
		dateGroups := make(map[string]int)
		for _, artifact := range artifacts {
			date := artifact.BuildTime.Format(am.Config.Organization.DateFormat)
			dateGroups[date]++
		}

		// Convert to recent builds slice
		for date, count := range dateGroups {
			status.RecentBuilds = append(status.RecentBuilds, RecentBuild{
				Date:          date,
				ArtifactCount: count,
			})
		}

		// Sort recent builds by date (newest first)
		sort.Slice(status.RecentBuilds, func(i, j int) bool {
			dateI, _ := time.Parse(am.Config.Organization.DateFormat, status.RecentBuilds[i].Date)
			dateJ, _ := time.Parse(am.Config.Organization.DateFormat, status.RecentBuilds[j].Date)
			return dateI.After(dateJ)
		})

		// Keep only recent builds (last 10)
		if len(status.RecentBuilds) > 10 {
			status.RecentBuilds = status.RecentBuilds[:10]
		}
	}

	return status, nil
}

// CleanAll removes all build artifacts
func (am *ArtifactManager) CleanAll() error {
	outputDir := am.GetOutputDir()

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return nil // Nothing to clean
	}

	return os.RemoveAll(outputDir)
}

// CleanOlderThan removes artifacts older than the specified duration
func (am *ArtifactManager) CleanOlderThan(duration string) error {
	// Parse duration (e.g., "7d", "2w", "1m")
	// This is a simplified implementation
	// TODO: Implement proper duration parsing
	return fmt.Errorf("CleanOlderThan not yet implemented")
}
