package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/build"
	"github.com/Jerinji2016/fdawg/pkg/environment"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/utils"
	"github.com/urfave/cli/v2"
)

// BuildCommand returns the CLI command for building Flutter applications
func BuildCommand() *cli.Command {
	return &cli.Command{
		Name:        "build",
		Usage:       "Build Flutter applications for multiple platforms",
		Description: "Comprehensive build management with pre-build setup and artifact organization",
		Subcommands: []*cli.Command{
			{
				Name:        "run",
				Usage:       "Execute builds for specified platforms",
				Description: "Build Flutter applications for one or more platforms with customizable pre-build steps",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "platforms",
						Aliases: []string{"p"},
						Usage:   "Platforms to build (android, ios, web, macos, linux, windows, all)",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "Path to build configuration file",
						Value:   ".fdawg/build.yaml",
					},
					&cli.BoolFlag{
						Name:  "skip-pre-build",
						Usage: "Skip pre-build setup steps",
					},
					&cli.BoolFlag{
						Name:  "continue-on-error",
						Usage: "Continue building other platforms if one fails",
					},
					&cli.BoolFlag{
						Name:  "dry-run",
						Usage: "Show what would be executed without running builds",
					},
					&cli.BoolFlag{
						Name:  "parallel",
						Usage: "Run platform builds in parallel (experimental)",
					},
					&cli.StringFlag{
						Name:    "env",
						Aliases: []string{"e"},
						Usage:   "Environment to use for build (uses --dart-define-from-file)",
					},
				},
				Action: runBuild,
			},
			{
				Name:        "setup",
				Usage:       "Interactive build configuration wizard",
				Description: "Set up build configuration with guided prompts",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "default",
						Usage: "Use default configuration without interactive prompts",
					},
					&cli.BoolFlag{
						Name:  "force",
						Usage: "Overwrite existing configuration",
					},
				},
				Action: setupBuild,
			},
			{
				Name:        "status",
				Usage:       "Show build status and available artifacts",
				Description: "Display information about recent builds and available artifacts",
				Action:      showBuildStatus,
			},
			{
				Name:        "clean",
				Usage:       "Clean build artifacts and cache",
				Description: "Remove build artifacts and clean build cache",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "older-than",
						Usage: "Remove artifacts older than specified duration (e.g., 7d, 2w, 1m)",
					},
					&cli.BoolFlag{
						Name:  "all",
						Usage: "Remove all build artifacts",
					},
				},
				Action: cleanBuild,
			},
			{
				Name:        "list",
				Usage:       "List available build artifacts",
				Description: "Show available build artifacts organized by date and platform",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "date",
						Usage: "Filter artifacts by date (e.g., June-6)",
					},
					&cli.StringFlag{
						Name:  "platform",
						Usage: "Filter artifacts by platform",
					},
				},
				Action: listBuildArtifacts,
			},
		},
	}
}

// runBuild executes the build process
func runBuild(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForBuild()
	if err != nil {
		return err
	}

	// Check if build configuration exists
	configPath := c.String("config")
	if !buildConfigExists(project.ProjectPath, configPath) {
		utils.Error("Build configuration not found!")
		utils.Info("Please run 'fdawg build setup' first to configure your build settings.")
		utils.Info("Or use 'fdawg build setup --default' for quick setup with default settings.")
		return fmt.Errorf("build configuration required")
	}

	// Parse platforms
	platforms, err := parseBuildPlatforms(c.StringSlice("platforms"))
	if err != nil {
		return err
	}

	if len(platforms) == 0 {
		utils.Error("No platforms specified")
		utils.Info("Usage: fdawg build run --platforms android,ios")
		utils.Info("Available platforms: android, ios, web, macos, linux, windows, all")
		return fmt.Errorf("no platforms specified")
	}

	// Load build configuration
	buildConfig, err := build.LoadBuildConfig(project.ProjectPath, configPath)
	if err != nil {
		utils.Error("Failed to load build config: %v", err)
		return err
	}

	// Validate environment if specified
	envName := c.String("env")
	if envName != "" {
		if err := validateEnvironment(project.ProjectPath, envName); err != nil {
			utils.Error("Environment validation failed: %v", err)
			return err
		}
		utils.Info("Using environment: %s", envName)
	}

	// Create build manager
	buildManager, err := build.NewBuildManager(project.ProjectPath, buildConfig)
	if err != nil {
		utils.Error("Failed to create build manager: %v", err)
		return err
	}

	// Build options
	options := build.BuildOptions{
		SkipPreBuild:    c.Bool("skip-pre-build"),
		ContinueOnError: c.Bool("continue-on-error"),
		DryRun:          c.Bool("dry-run"),
		Parallel:        c.Bool("parallel"),
		Environment:     c.String("env"),
	}

	if options.DryRun {
		utils.Info("Dry run mode - showing what would be executed:")
		return buildManager.ShowBuildPlan(platforms, options)
	}

	// Execute build
	utils.Info("Starting build process...")
	result, err := buildManager.ExecuteBuild(platforms, options)
	if err != nil {
		utils.Error("Build failed: %v", err)
		return err
	}

	// Display results
	displayBuildResults(result)
	return nil
}

// setupBuild runs the interactive build configuration wizard
func setupBuild(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForBuild()
	if err != nil {
		return err
	}

	configPath := ".fdawg/build.yaml"

	// Check if configuration already exists
	if buildConfigExists(project.ProjectPath, configPath) && !c.Bool("force") {
		utils.Warning("Build configuration already exists at %s", configPath)
		utils.Info("Use --force to overwrite existing configuration")
		utils.Info("Or run 'fdawg build run --platforms <platforms>' to use existing configuration")
		return fmt.Errorf("configuration already exists")
	}

	var config *build.BuildConfig

	if c.Bool("default") {
		// Use default configuration
		utils.Info("Setting up build configuration with default settings for: %s", project.PubspecInfo.Name)
		config = createDefaultConfigForProject(project.ProjectPath)
		utils.Success("Default configuration created")
	} else {
		// Run interactive setup
		utils.Info("Setting up build configuration for: %s", project.PubspecInfo.Name)

		wizard := build.NewSetupWizard(project.ProjectPath)
		config, err = wizard.Run()
		if err != nil {
			utils.Error("Setup failed: %v", err)
			return err
		}
	}

	// Save configuration
	if err := build.SaveBuildConfig(project.ProjectPath, configPath, config); err != nil {
		utils.Error("Failed to save configuration: %v", err)
		return err
	}

	utils.Success("Build configuration saved to %s", configPath)
	utils.Info("You can now run: fdawg build run --platforms all")

	return nil
}

// showBuildStatus displays build status and artifacts
func showBuildStatus(c *cli.Context) error {
	project, err := validateFlutterProjectForBuild()
	if err != nil {
		return err
	}

	// Check if build configuration exists
	configPath := ".fdawg/build.yaml"
	if !buildConfigExists(project.ProjectPath, configPath) {
		utils.Error("Build configuration not found!")
		utils.Info("Please run 'fdawg build setup' first to configure your build settings.")
		return fmt.Errorf("build configuration required")
	}

	artifactManager := build.NewArtifactManager(project.ProjectPath, nil)
	status, err := artifactManager.GetBuildStatus()
	if err != nil {
		utils.Error("Failed to get build status: %v", err)
		return err
	}

	displayBuildStatus(status)
	return nil
}

// cleanBuild cleans build artifacts
func cleanBuild(c *cli.Context) error {
	project, err := validateFlutterProjectForBuild()
	if err != nil {
		return err
	}

	artifactManager := build.NewArtifactManager(project.ProjectPath, nil)

	if c.Bool("all") {
		return artifactManager.CleanAll()
	}

	olderThan := c.String("older-than")
	if olderThan != "" {
		return artifactManager.CleanOlderThan(olderThan)
	}

	utils.Info("Use --all to clean all artifacts or --older-than to specify age")
	return nil
}

// listBuildArtifacts lists available build artifacts
func listBuildArtifacts(c *cli.Context) error {
	project, err := validateFlutterProjectForBuild()
	if err != nil {
		return err
	}

	// Check if build configuration exists
	configPath := ".fdawg/build.yaml"
	if !buildConfigExists(project.ProjectPath, configPath) {
		utils.Error("Build configuration not found!")
		utils.Info("Please run 'fdawg build setup' first to configure your build settings.")
		return fmt.Errorf("build configuration required")
	}

	artifactManager := build.NewArtifactManager(project.ProjectPath, nil)

	filters := build.ArtifactFilters{
		Date:     c.String("date"),
		Platform: c.String("platform"),
	}

	artifacts, err := artifactManager.ListArtifacts(filters)
	if err != nil {
		utils.Error("Failed to list artifacts: %v", err)
		return err
	}

	displayArtifactList(artifacts)
	return nil
}

// Helper functions

func validateFlutterProjectForBuild() (*flutter.ValidationResult, error) {
	result, err := flutter.ValidateProject(".")
	if err != nil {
		utils.Error("Not a valid Flutter project: %v", err)
		return nil, err
	}
	return result, nil
}

func parseBuildPlatforms(platformStrings []string) ([]build.Platform, error) {
	if len(platformStrings) == 0 {
		return nil, nil
	}

	var platforms []build.Platform
	allPlatforms := []build.Platform{
		build.PlatformAndroid,
		build.PlatformIOS,
		build.PlatformWeb,
		build.PlatformMacOS,
		build.PlatformLinux,
		build.PlatformWindows,
	}

	for _, platformStr := range platformStrings {
		// Handle comma-separated platforms in a single string
		parts := strings.Split(platformStr, ",")
		for _, part := range parts {
			part = strings.TrimSpace(strings.ToLower(part))

			if part == "all" {
				return allPlatforms, nil
			}

			platform := build.Platform(part)
			if !isValidPlatform(platform) {
				return nil, fmt.Errorf("invalid platform: %s", part)
			}

			platforms = append(platforms, platform)
		}
	}

	return platforms, nil
}

func isValidPlatform(platform build.Platform) bool {
	validPlatforms := []build.Platform{
		build.PlatformAndroid,
		build.PlatformIOS,
		build.PlatformWeb,
		build.PlatformMacOS,
		build.PlatformLinux,
		build.PlatformWindows,
	}

	for _, valid := range validPlatforms {
		if platform == valid {
			return true
		}
	}
	return false
}

func displayBuildResults(result *build.BuildResult) {
	fmt.Println("\n" + utils.Separator("=", 60))
	utils.Success("Build Completed")
	fmt.Println(utils.Separator("=", 60))

	fmt.Printf("Duration: %v\n", result.Duration)
	fmt.Printf("Artifacts: %d\n", len(result.Artifacts))

	if result.LogFile != "" {
		fmt.Printf("Log file: %s\n", result.LogFile)
	}

	// Display platform results
	fmt.Println(utils.Separator("-", 60))
	utils.Info("Platform Results")

	for platform, platformResult := range result.PlatformResults {
		status := "✅ Success"
		if !platformResult.Success {
			status = "❌ Failed"
		}
		fmt.Printf("%-10s: %s (%d artifacts)\n", platform, status, len(platformResult.Artifacts))
	}

	// Display artifacts
	if len(result.Artifacts) > 0 {
		fmt.Println(utils.Separator("-", 60))
		utils.Info("Generated Artifacts")

		for _, artifact := range result.Artifacts {
			fmt.Printf("📦 %s\n", artifact.FileName)
			fmt.Printf("   Platform: %s | Size: %s\n", artifact.Platform, utils.FormatFileSize(artifact.Size))
		}
	}
}

func displayBuildStatus(status *build.BuildStatus) {
	fmt.Println("\n" + utils.Separator("=", 50))
	utils.Success("Build Status")
	fmt.Println(utils.Separator("=", 50))

	fmt.Printf("Total Artifacts: %d\n", status.TotalArtifacts)
	fmt.Printf("Last Build: %s\n", status.LastBuildTime.Format("2006-01-02 15:04:05"))

	if len(status.RecentBuilds) > 0 {
		fmt.Println(utils.Separator("-", 50))
		utils.Info("Recent Builds")

		for _, build := range status.RecentBuilds {
			fmt.Printf("%s - %d artifacts\n", build.Date, build.ArtifactCount)
		}
	}
}

func displayArtifactList(artifacts []*build.BuildArtifact) {
	if len(artifacts) == 0 {
		utils.Info("No artifacts found")
		return
	}

	fmt.Println("\n" + utils.Separator("=", 60))
	utils.Success("Build Artifacts")
	fmt.Println(utils.Separator("=", 60))

	currentDate := ""
	for _, artifact := range artifacts {
		date := artifact.BuildTime.Format("January-2")
		if date != currentDate {
			currentDate = date
			fmt.Printf("\n📅 %s\n", date)
			fmt.Println(utils.Separator("-", 40))
		}

		fmt.Printf("📦 %s\n", artifact.FileName)
		fmt.Printf("   Platform: %s | Arch: %s | Size: %s\n",
			artifact.Platform, artifact.Architecture, utils.FormatFileSize(artifact.Size))
	}
}

// buildConfigExists checks if build configuration file exists
func buildConfigExists(projectPath, configPath string) bool {
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(projectPath, configPath)
	}

	_, err := os.Stat(configPath)
	return err == nil
}

// createDefaultConfigForProject creates a default configuration tailored to the project
func createDefaultConfigForProject(projectPath string) *build.BuildConfig {
	config := build.DefaultBuildConfig()

	// Detect available platforms and enable only those
	availablePlatforms := detectAvailablePlatforms(projectPath)

	// Disable all platforms first
	config.Platforms.Android.Enabled = false
	config.Platforms.IOS.Enabled = false
	config.Platforms.Web.Enabled = false
	config.Platforms.MacOS.Enabled = false
	config.Platforms.Linux.Enabled = false
	config.Platforms.Windows.Enabled = false

	// Enable only available platforms
	for _, platform := range availablePlatforms {
		switch platform {
		case build.PlatformAndroid:
			config.Platforms.Android.Enabled = true
		case build.PlatformIOS:
			config.Platforms.IOS.Enabled = true
		case build.PlatformWeb:
			config.Platforms.Web.Enabled = true
		case build.PlatformMacOS:
			config.Platforms.MacOS.Enabled = true
		case build.PlatformLinux:
			config.Platforms.Linux.Enabled = true
		case build.PlatformWindows:
			config.Platforms.Windows.Enabled = true
		}
	}

	// Detect and configure common build tools
	if fileExists(projectPath, "build.yaml") {
		config.PreBuild.Global = append(config.PreBuild.Global, build.BuildStep{
			Name:      "Generate code",
			Command:   "dart run build_runner build --delete-conflicting-outputs",
			Required:  false,
			Timeout:   600,
			Condition: "file_exists:build.yaml",
		})
	}

	if fileExists(projectPath, "flutter_launcher_icons.yaml") {
		config.PreBuild.Global = append(config.PreBuild.Global, build.BuildStep{
			Name:      "Generate launcher icons",
			Command:   "dart run flutter_launcher_icons:main",
			Required:  false,
			Timeout:   300,
			Condition: "file_exists:flutter_launcher_icons.yaml",
		})
	}

	if fileExists(projectPath, "flutter_native_splash.yaml") {
		config.PreBuild.Global = append(config.PreBuild.Global, build.BuildStep{
			Name:      "Generate splash screens",
			Command:   "dart run flutter_native_splash:create",
			Required:  false,
			Timeout:   300,
			Condition: "file_exists:flutter_native_splash.yaml",
		})
	}

	// Add iOS pod install if iOS is available
	if contains(availablePlatforms, build.PlatformIOS) {
		config.PreBuild.IOS = append(config.PreBuild.IOS, build.BuildStep{
			Name:      "Pod install",
			Command:   "cd ios && pod install",
			Required:  true,
			Timeout:   600,
			Condition: "platform_available:ios",
		})
	}

	return config
}

// detectAvailablePlatforms detects which platforms are available in the project
func detectAvailablePlatforms(projectPath string) []build.Platform {
	var platforms []build.Platform

	platformDirs := map[build.Platform]string{
		build.PlatformAndroid: "android",
		build.PlatformIOS:     "ios",
		build.PlatformWeb:     "web",
		build.PlatformMacOS:   "macos",
		build.PlatformLinux:   "linux",
		build.PlatformWindows: "windows",
	}

	for platform, dir := range platformDirs {
		if dirExists(projectPath, dir) {
			platforms = append(platforms, platform)
		}
	}

	return platforms
}

// fileExists checks if a file exists in the project
func fileExists(projectPath, filename string) bool {
	path := filepath.Join(projectPath, filename)
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// dirExists checks if a directory exists in the project
func dirExists(projectPath, dirname string) bool {
	path := filepath.Join(projectPath, dirname)
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// contains checks if a slice contains a platform
func contains(platforms []build.Platform, platform build.Platform) bool {
	for _, p := range platforms {
		if p == platform {
			return true
		}
	}
	return false
}

// validateEnvironment validates that the specified environment exists
func validateEnvironment(projectPath, envName string) error {
	// Check if environment file exists
	_, err := environment.GetEnvFile(projectPath, envName)
	if err != nil {
		// List available environments for helpful error message
		envFiles, listErr := environment.ListEnvFiles(projectPath)
		if listErr != nil {
			return fmt.Errorf("environment '%s' not found", envName)
		}

		if len(envFiles) == 0 {
			return fmt.Errorf("environment '%s' not found. No environments exist. Create one with: fdawg env create %s", envName, envName)
		}

		var availableEnvs []string
		for _, envFile := range envFiles {
			availableEnvs = append(availableEnvs, envFile.Name)
		}

		return fmt.Errorf("environment '%s' not found. Available environments: %s", envName, strings.Join(availableEnvs, ", "))
	}

	return nil
}
