package build

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Jerinji2016/fdawg/pkg/utils"
)

// Platform represents a supported build platform
type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformWeb     Platform = "web"
	PlatformMacOS   Platform = "macos"
	PlatformLinux   Platform = "linux"
	PlatformWindows Platform = "windows"
)

// BuildManager manages the build process
type BuildManager struct {
	ProjectPath     string
	Config          *BuildConfig
	ArtifactManager *ArtifactManager
	Logger          *utils.Logger
}

// BuildOptions contains options for the build process
type BuildOptions struct {
	SkipPreBuild    bool
	ContinueOnError bool
	DryRun          bool
	Parallel        bool
	Environment     string
}

// BuildResult contains the results of a build process
type BuildResult struct {
	Success         bool
	PlatformResults map[Platform]*PlatformBuildResult
	Artifacts       []*BuildArtifact
	BuildTime       time.Time
	Duration        time.Duration
	LogFile         string
}

// PlatformBuildResult contains the result of building for a specific platform
type PlatformBuildResult struct {
	Platform  Platform
	Success   bool
	Artifacts []*BuildArtifact
	Error     error
	Duration  time.Duration
}

// BuildArtifact represents a build output file
type BuildArtifact struct {
	Platform     Platform  `json:"platform"`
	BuildType    string    `json:"build_type"`
	Architecture string    `json:"architecture"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	Size         int64     `json:"size"`
	BuildTime    time.Time `json:"build_time"`
	AppName      string    `json:"app_name"`
	Version      string    `json:"version"`
}

// BuildStatus represents the current build status
type BuildStatus struct {
	TotalArtifacts int
	LastBuildTime  time.Time
	RecentBuilds   []RecentBuild
}

// RecentBuild represents information about a recent build
type RecentBuild struct {
	Date          string
	ArtifactCount int
}

// ArtifactFilters contains filters for listing artifacts
type ArtifactFilters struct {
	Date     string
	Platform string
}

// NewBuildManager creates a new build manager
func NewBuildManager(projectPath string, config *BuildConfig) (*BuildManager, error) {
	// Validate project path
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project path does not exist: %s", projectPath)
	}

	// Create artifact manager
	artifactManager := NewArtifactManager(projectPath, &config.Artifacts)

	// Create logger
	logger := utils.NewLogger("BUILD")

	return &BuildManager{
		ProjectPath:     projectPath,
		Config:          config,
		ArtifactManager: artifactManager,
		Logger:          logger,
	}, nil
}

// ExecuteBuild executes the build process for specified platforms
func (bm *BuildManager) ExecuteBuild(platforms []Platform, options BuildOptions) (*BuildResult, error) {
	startTime := time.Now()
	result := &BuildResult{
		PlatformResults: make(map[Platform]*PlatformBuildResult),
		BuildTime:       startTime,
	}

	// Setup build logging
	logFile := bm.setupBuildLogging(startTime)
	result.LogFile = logFile

	bm.Logger.Info("Starting build process for platforms: %v", platforms)

	// Execute pre-build steps
	if !options.SkipPreBuild {
		bm.Logger.Info("Executing pre-build steps...")
		if err := bm.executePreBuildSteps(); err != nil {
			bm.Logger.Error("Pre-build failed: %v", err)
			return result, fmt.Errorf("pre-build failed: %w", err)
		}
		bm.Logger.Success("Pre-build steps completed")
	} else {
		bm.Logger.Info("Skipping pre-build steps")
	}

	// Execute platform builds
	var allArtifacts []*BuildArtifact

	for _, platform := range platforms {
		bm.Logger.Info("Building for platform: %s", platform)

		platformResult, err := bm.buildPlatformWithOptions(platform, options)
		result.PlatformResults[platform] = platformResult

		if err != nil {
			bm.Logger.Error("Platform %s build failed: %v", platform, err)
			if !options.ContinueOnError {
				result.Duration = time.Since(startTime)
				return result, fmt.Errorf("platform %s build failed: %w", platform, err)
			}
			continue
		}

		bm.Logger.Success("Platform %s build completed with %d artifacts", platform, len(platformResult.Artifacts))

		// Organize artifacts
		for _, artifact := range platformResult.Artifacts {
			if err := bm.ArtifactManager.OrganizeArtifact(artifact); err != nil {
				bm.Logger.Warning("Failed to organize artifact %s: %v", artifact.FileName, err)
			} else {
				allArtifacts = append(allArtifacts, artifact)
			}
		}
	}

	result.Artifacts = allArtifacts
	result.Duration = time.Since(startTime)
	result.Success = len(allArtifacts) > 0

	// Generate build summary
	bm.generateBuildSummary(result)

	bm.Logger.Success("Build process completed in %v with %d artifacts", result.Duration, len(result.Artifacts))

	return result, nil
}

// ShowBuildPlan shows what would be executed in a dry run
func (bm *BuildManager) ShowBuildPlan(platforms []Platform, options BuildOptions) error {
	fmt.Println("\n" + utils.Separator("=", 60))
	utils.Info("Build Plan (Dry Run)")
	fmt.Println(utils.Separator("=", 60))

	// Show pre-build steps
	if !options.SkipPreBuild {
		fmt.Println("\n📋 Pre-build Steps:")
		if err := bm.showPreBuildPlan(); err != nil {
			return err
		}
	}

	// Show platform builds
	fmt.Println("\n🏗️  Platform Builds:")
	for _, platform := range platforms {
		fmt.Printf("  • %s\n", platform)
		if err := bm.showPlatformBuildPlan(platform); err != nil {
			utils.Warning("Failed to show plan for %s: %v", platform, err)
		}
	}

	// Show artifact organization
	fmt.Println("\n📦 Artifact Organization:")
	fmt.Printf("  • Output directory: %s\n", bm.ArtifactManager.GetOutputDir())
	fmt.Printf("  • Naming pattern: %s\n", bm.Config.Artifacts.Naming.Pattern)

	return nil
}

// buildPlatformWithOptions builds for a specific platform with build options
func (bm *BuildManager) buildPlatformWithOptions(platform Platform, options BuildOptions) (*PlatformBuildResult, error) {
	startTime := time.Now()
	result := &PlatformBuildResult{
		Platform: platform,
	}

	// Check if platform is available
	if !bm.isPlatformAvailable(platform) {
		result.Error = fmt.Errorf("platform %s not available in project", platform)
		return result, result.Error
	}

	// Get platform configuration
	platformConfig := bm.getPlatformConfig(platform)
	if platformConfig == nil {
		result.Error = fmt.Errorf("no configuration found for platform %s", platform)
		return result, result.Error
	}

	// Execute platform-specific pre-build steps
	if err := bm.executePlatformPreBuildSteps(platform); err != nil {
		result.Error = fmt.Errorf("platform pre-build failed: %w", err)
		return result, result.Error
	}

	// Build platform with options
	artifacts, err := bm.executePlatformBuildWithOptions(platform, platformConfig, options)
	if err != nil {
		result.Error = err
		return result, err
	}

	result.Artifacts = artifacts
	result.Success = len(artifacts) > 0
	result.Duration = time.Since(startTime)

	return result, nil
}

// isPlatformAvailable checks if a platform is available in the project
func (bm *BuildManager) isPlatformAvailable(platform Platform) bool {
	platformDirs := map[Platform]string{
		PlatformAndroid: "android",
		PlatformIOS:     "ios",
		PlatformWeb:     "web",
		PlatformMacOS:   "macos",
		PlatformLinux:   "linux",
		PlatformWindows: "windows",
	}

	dir, exists := platformDirs[platform]
	if !exists {
		return false
	}

	platformPath := filepath.Join(bm.ProjectPath, dir)
	if _, err := os.Stat(platformPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// getPlatformConfig gets the configuration for a specific platform
func (bm *BuildManager) getPlatformConfig(platform Platform) interface{} {
	switch platform {
	case PlatformAndroid:
		return &bm.Config.Platforms.Android
	case PlatformIOS:
		return &bm.Config.Platforms.IOS
	case PlatformWeb:
		return &bm.Config.Platforms.Web
	case PlatformMacOS:
		return &bm.Config.Platforms.MacOS
	case PlatformLinux:
		return &bm.Config.Platforms.Linux
	case PlatformWindows:
		return &bm.Config.Platforms.Windows
	default:
		return nil
	}
}

// setupBuildLogging sets up logging for the build process
func (bm *BuildManager) setupBuildLogging(startTime time.Time) string {
	if !bm.Config.Execution.SaveLogs {
		return ""
	}

	logDir := filepath.Join(bm.ProjectPath, bm.Config.Artifacts.BaseOutputDir, "build-logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		bm.Logger.Warning("Failed to create log directory: %v", err)
		return ""
	}

	logFileName := fmt.Sprintf("%s.log", startTime.Format("January-2_15-04-05"))
	logFile := filepath.Join(logDir, logFileName)

	// TODO: Setup file logging
	return logFile
}

// generateBuildSummary generates a summary of the build process
func (bm *BuildManager) generateBuildSummary(_ *BuildResult) {
	// TODO: Generate detailed build summary
	bm.Logger.Info("Build summary generated")
}

// executePreBuildSteps executes global pre-build steps
func (bm *BuildManager) executePreBuildSteps() error {
	executor := NewCommandExecutor(bm.ProjectPath, bm.Logger)

	for _, step := range bm.Config.PreBuild.Global {
		if err := executor.ExecuteStep(step); err != nil {
			if step.Required {
				return fmt.Errorf("required pre-build step '%s' failed: %w", step.Name, err)
			}
			bm.Logger.Warning("Optional pre-build step '%s' failed: %v", step.Name, err)
		}
	}

	return nil
}

// executePlatformPreBuildSteps executes platform-specific pre-build steps
func (bm *BuildManager) executePlatformPreBuildSteps(platform Platform) error {
	executor := NewCommandExecutor(bm.ProjectPath, bm.Logger)

	var steps []BuildStep
	switch platform {
	case PlatformAndroid:
		steps = bm.Config.PreBuild.Android
	case PlatformIOS:
		steps = bm.Config.PreBuild.IOS
	case PlatformWeb:
		steps = bm.Config.PreBuild.Web
	}

	for _, step := range steps {
		if err := executor.ExecuteStep(step); err != nil {
			if step.Required {
				return fmt.Errorf("required platform pre-build step '%s' failed: %w", step.Name, err)
			}
			bm.Logger.Warning("Optional platform pre-build step '%s' failed: %v", step.Name, err)
		}
	}

	return nil
}

// executePlatformBuildWithOptions executes the build for a specific platform with options
func (bm *BuildManager) executePlatformBuildWithOptions(platform Platform, config interface{}, options BuildOptions) ([]*BuildArtifact, error) {
	switch platform {
	case PlatformAndroid:
		return bm.buildAndroidWithOptions(config.(*AndroidBuildConfig), options)
	case PlatformIOS:
		return bm.buildIOSWithOptions(config.(*IOSBuildConfig), options)
	case PlatformWeb:
		return bm.buildWebWithOptions(config.(*WebBuildConfig), options)
	case PlatformMacOS:
		return bm.buildMacOSWithOptions(config.(*MacOSBuildConfig), options)
	case PlatformLinux:
		return bm.buildLinuxWithOptions(config.(*LinuxBuildConfig), options)
	case PlatformWindows:
		return bm.buildWindowsWithOptions(config.(*WindowsBuildConfig), options)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

// showPreBuildPlan shows the pre-build plan
func (bm *BuildManager) showPreBuildPlan() error {
	for _, step := range bm.Config.PreBuild.Global {
		status := "optional"
		if step.Required {
			status = "required"
		}
		fmt.Printf("    %s (%s): %s\n", step.Name, status, step.Command)
	}
	return nil
}

// showPlatformBuildPlan shows the build plan for a platform
func (bm *BuildManager) showPlatformBuildPlan(platform Platform) error {
	config := bm.getPlatformConfig(platform)
	if config == nil {
		return fmt.Errorf("no configuration for platform %s", platform)
	}

	// Show platform-specific build types
	switch platform {
	case PlatformAndroid:
		androidConfig := config.(*AndroidBuildConfig)
		for _, buildType := range androidConfig.BuildTypes {
			fmt.Printf("    • %s (%s)\n", buildType.Name, buildType.Type)
		}
	case PlatformIOS:
		iosConfig := config.(*IOSBuildConfig)
		for _, buildType := range iosConfig.BuildTypes {
			fmt.Printf("    • %s (%s)\n", buildType.Name, buildType.Type)
		}
		// Add other platforms as needed
	}

	return nil
}

// getEnvironmentFilePath returns the path to the environment file
func (bm *BuildManager) getEnvironmentFilePath(envName string) string {
	return filepath.Join(bm.ProjectPath, ".environment", envName+".json")
}
