package build

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/utils"
)

// SetupWizard provides interactive build configuration setup
type SetupWizard struct {
	ProjectPath string
	scanner     *bufio.Scanner
}

// NewSetupWizard creates a new setup wizard
func NewSetupWizard(projectPath string) *SetupWizard {
	return &SetupWizard{
		ProjectPath: projectPath,
		scanner:     bufio.NewScanner(os.Stdin),
	}
}

// Run executes the interactive setup wizard
func (sw *SetupWizard) Run() (*BuildConfig, error) {
	fmt.Println("\n" + utils.Separator("=", 60))
	utils.Success("🚀 FDAWG Build Configuration Setup")
	fmt.Println(utils.Separator("=", 60))
	fmt.Println("Let's configure your build pipeline...")

	// Start with default configuration
	config := DefaultBuildConfig()

	// Detect existing setup
	detectedSetup := sw.detectExistingSetup()
	sw.displayDetectedSetup(detectedSetup)

	// Configure metadata
	if err := sw.configureMetadata(&config.Metadata); err != nil {
		return nil, err
	}

	// Configure pre-build steps
	if err := sw.configurePreBuild(&config.PreBuild, detectedSetup); err != nil {
		return nil, err
	}

	// Configure platforms
	if err := sw.configurePlatforms(&config.Platforms); err != nil {
		return nil, err
	}

	// Configure artifacts
	if err := sw.configureArtifacts(&config.Artifacts); err != nil {
		return nil, err
	}

	// Configure execution
	if err := sw.configureExecution(&config.Execution); err != nil {
		return nil, err
	}

	// Display summary
	sw.displayConfigSummary(config)

	// Confirm configuration
	if !sw.promptYesNo("Save this configuration?", true) {
		return nil, fmt.Errorf("configuration cancelled by user")
	}

	return config, nil
}

// DetectedSetup contains information about detected project setup
type DetectedSetup struct {
	HasBuildRunner    bool
	HasLauncherIcons  bool
	HasNativeAssets   bool
	HasNodeJS         bool
	CustomScripts     []string
	AvailablePlatforms []Platform
}

// detectExistingSetup detects existing project setup
func (sw *SetupWizard) detectExistingSetup() *DetectedSetup {
	setup := &DetectedSetup{}

	// Check for common build tools
	if sw.fileExists("build.yaml") {
		setup.HasBuildRunner = true
	}

	if sw.fileExists("flutter_launcher_icons.yaml") {
		setup.HasLauncherIcons = true
	}

	if sw.fileExists("flutter_native_splash.yaml") {
		setup.HasNativeAssets = true
	}

	if sw.fileExists("package.json") {
		setup.HasNodeJS = true
	}

	// Check for custom scripts
	scriptsDir := filepath.Join(sw.ProjectPath, "scripts")
	if sw.dirExists(scriptsDir) {
		setup.CustomScripts = sw.findScripts(scriptsDir)
	}

	// Detect available platforms
	setup.AvailablePlatforms = sw.detectAvailablePlatforms()

	return setup
}

// displayDetectedSetup displays what was detected in the project
func (sw *SetupWizard) displayDetectedSetup(setup *DetectedSetup) {
	fmt.Println("\n📋 Detected Project Setup:")
	
	if setup.HasBuildRunner {
		fmt.Println("  ✅ build_runner detected")
	}
	if setup.HasLauncherIcons {
		fmt.Println("  ✅ flutter_launcher_icons detected")
	}
	if setup.HasNativeAssets {
		fmt.Println("  ✅ flutter_native_splash detected")
	}
	if setup.HasNodeJS {
		fmt.Println("  ✅ Node.js project detected")
	}
	if len(setup.CustomScripts) > 0 {
		fmt.Printf("  ✅ Custom scripts found: %s\n", strings.Join(setup.CustomScripts, ", "))
	}

	fmt.Printf("  📱 Available platforms: %s\n", sw.formatPlatforms(setup.AvailablePlatforms))
}

// configureMetadata configures build metadata
func (sw *SetupWizard) configureMetadata(metadata *MetadataConfig) error {
	fmt.Println("\n" + utils.Separator("-", 40))
	utils.Info("📝 Build Metadata Configuration")

	// App name source
	fmt.Println("\nHow should we get the app name for artifacts?")
	fmt.Println("1. From namer configuration (recommended)")
	fmt.Println("2. From pubspec.yaml")
	fmt.Println("3. Custom name")

	choice := sw.promptChoice("Choose option", []string{"1", "2", "3"}, "1")
	switch choice {
	case "1":
		metadata.AppNameSource = "namer"
	case "2":
		metadata.AppNameSource = "pubspec"
	case "3":
		metadata.AppNameSource = "custom"
		metadata.CustomAppName = sw.promptString("Enter custom app name", "MyApp")
	}

	// Version source
	fmt.Println("\nHow should we get the version for artifacts?")
	fmt.Println("1. From pubspec.yaml (recommended)")
	fmt.Println("2. Custom version")

	choice = sw.promptChoice("Choose option", []string{"1", "2"}, "1")
	switch choice {
	case "1":
		metadata.VersionSource = "pubspec"
	case "2":
		metadata.VersionSource = "custom"
		metadata.CustomVersion = sw.promptString("Enter custom version", "1.0.0")
	}

	return nil
}

// configurePreBuild configures pre-build steps
func (sw *SetupWizard) configurePreBuild(preBuild *PreBuildConfig, setup *DetectedSetup) error {
	fmt.Println("\n" + utils.Separator("-", 40))
	utils.Info("⚙️  Pre-build Steps Configuration")

	// Global pre-build steps
	fmt.Println("\nGlobal pre-build steps (run before all platform builds):")

	// build_runner
	if setup.HasBuildRunner {
		if sw.promptYesNo("Include code generation with build_runner?", true) {
			preBuild.Global = append(preBuild.Global, BuildStep{
				Name:     "Generate code",
				Command:  "dart run build_runner build --delete-conflicting-outputs",
				Required: false,
				Timeout:  600,
				Condition: "file_exists:build.yaml",
			})
		}
	}

	// flutter_launcher_icons
	if setup.HasLauncherIcons {
		if sw.promptYesNo("Include launcher icon generation?", true) {
			preBuild.Global = append(preBuild.Global, BuildStep{
				Name:     "Generate launcher icons",
				Command:  "dart run flutter_launcher_icons:main",
				Required: false,
				Timeout:  300,
				Condition: "file_exists:flutter_launcher_icons.yaml",
			})
		}
	}

	// flutter_native_splash
	if setup.HasNativeAssets {
		if sw.promptYesNo("Include native splash screen generation?", true) {
			preBuild.Global = append(preBuild.Global, BuildStep{
				Name:     "Generate splash screens",
				Command:  "dart run flutter_native_splash:create",
				Required: false,
				Timeout:  300,
				Condition: "file_exists:flutter_native_splash.yaml",
			})
		}
	}

	// Custom scripts
	if len(setup.CustomScripts) > 0 {
		fmt.Printf("\nFound custom scripts: %s\n", strings.Join(setup.CustomScripts, ", "))
		if sw.promptYesNo("Include custom scripts in pre-build?", false) {
			for _, script := range setup.CustomScripts {
				scriptPath := filepath.Join("scripts", script)
				preBuild.Global = append(preBuild.Global, BuildStep{
					Name:     fmt.Sprintf("Run %s", script),
					Command:  fmt.Sprintf("./%s", scriptPath),
					Required: false,
					Timeout:  300,
					Condition: fmt.Sprintf("file_exists:%s", scriptPath),
				})
			}
		}
	}

	// Platform-specific pre-build steps
	if sw.promptYesNo("Configure platform-specific pre-build steps?", false) {
		sw.configurePlatformPreBuild(preBuild, setup)
	}

	return nil
}

// configurePlatformPreBuild configures platform-specific pre-build steps
func (sw *SetupWizard) configurePlatformPreBuild(preBuild *PreBuildConfig, setup *DetectedSetup) {
	for _, platform := range setup.AvailablePlatforms {
		if !sw.promptYesNo(fmt.Sprintf("Configure pre-build steps for %s?", platform), false) {
			continue
		}

		switch platform {
		case PlatformAndroid:
			if sw.promptYesNo("Clean Android build before building?", false) {
				preBuild.Android = append(preBuild.Android, BuildStep{
					Name:     "Clean Android build",
					Command:  "flutter clean",
					Required: false,
					Timeout:  120,
				})
			}

		case PlatformIOS:
			if sw.promptYesNo("Run pod install before iOS builds?", true) {
				preBuild.IOS = append(preBuild.IOS, BuildStep{
					Name:     "Pod install",
					Command:  "cd ios && pod install",
					Required: true,
					Timeout:  600,
					Condition: "platform_available:ios",
				})
			}

		case PlatformWeb:
			if setup.HasNodeJS && sw.promptYesNo("Build web assets with npm?", false) {
				preBuild.Web = append(preBuild.Web, BuildStep{
					Name:     "Build web assets",
					Command:  "npm run build:web",
					Required: false,
					Timeout:  300,
					Condition: "file_exists:package.json",
				})
			}
		}
	}
}

// configurePlatforms configures platform build settings
func (sw *SetupWizard) configurePlatforms(platforms *PlatformsConfig) error {
	fmt.Println("\n" + utils.Separator("-", 40))
	utils.Info("📱 Platform Build Configuration")

	availablePlatforms := sw.detectAvailablePlatforms()

	for _, platform := range availablePlatforms {
		if !sw.promptYesNo(fmt.Sprintf("Configure builds for %s?", platform), true) {
			sw.disablePlatform(platforms, platform)
			continue
		}

		switch platform {
		case PlatformAndroid:
			sw.configureAndroidPlatform(&platforms.Android)
		case PlatformIOS:
			sw.configureIOSPlatform(&platforms.IOS)
		case PlatformWeb:
			sw.configureWebPlatform(&platforms.Web)
		case PlatformMacOS:
			sw.configureMacOSPlatform(&platforms.MacOS)
		case PlatformLinux:
			sw.configureLinuxPlatform(&platforms.Linux)
		case PlatformWindows:
			sw.configureWindowsPlatform(&platforms.Windows)
		}
	}

	return nil
}

// configureAndroidPlatform configures Android build settings
func (sw *SetupWizard) configureAndroidPlatform(config *AndroidBuildConfig) {
	fmt.Printf("\n🤖 Android Configuration:\n")

	config.Enabled = true
	config.BuildTypes = []AndroidBuildType{}

	if sw.promptYesNo("Build APK files?", true) {
		buildType := AndroidBuildType{
			Name:      "release_apk",
			Type:      "apk",
			BuildMode: "release",
		}

		if sw.promptYesNo("Split APK by architecture?", true) {
			buildType.SplitPerABI = true
		}

		if sw.promptYesNo("Enable code obfuscation?", true) {
			buildType.CustomArgs = append(buildType.CustomArgs, "--obfuscate")
		}

		config.BuildTypes = append(config.BuildTypes, buildType)
	}

	if sw.promptYesNo("Build App Bundle (AAB)?", true) {
		buildType := AndroidBuildType{
			Name:      "release_bundle",
			Type:      "appbundle",
			BuildMode: "release",
		}

		if sw.promptYesNo("Enable code obfuscation for AAB?", true) {
			buildType.CustomArgs = append(buildType.CustomArgs, "--obfuscate")
		}

		config.BuildTypes = append(config.BuildTypes, buildType)
	}
}

// configureIOSPlatform configures iOS build settings
func (sw *SetupWizard) configureIOSPlatform(config *IOSBuildConfig) {
	fmt.Printf("\n🍎 iOS Configuration:\n")

	config.Enabled = true
	config.BuildTypes = []IOSBuildType{}

	if sw.promptYesNo("Build iOS Archive (.xcarchive)?", true) {
		buildType := IOSBuildType{
			Name:       "archive",
			Type:       "archive",
			BuildMode:  "release",
			CustomArgs: []string{"--no-codesign"},
		}
		config.BuildTypes = append(config.BuildTypes, buildType)
	}

	if sw.promptYesNo("Build IPA file?", false) {
		buildType := IOSBuildType{
			Name:         "ipa",
			Type:         "ipa",
			BuildMode:    "release",
			ExportMethod: "development",
		}
		config.BuildTypes = append(config.BuildTypes, buildType)
	}
}

// configureWebPlatform configures Web build settings
func (sw *SetupWizard) configureWebPlatform(config *WebBuildConfig) {
	fmt.Printf("\n🌐 Web Configuration:\n")

	config.Enabled = true
	config.BuildTypes = []WebBuildType{
		{
			Name:       "release",
			Type:       "web",
			BuildMode:  "release",
			PWA:        sw.promptYesNo("Enable PWA features?", true),
			CustomArgs: []string{"--web-renderer", "canvaskit"},
		},
	}
}

// configureMacOSPlatform configures macOS build settings
func (sw *SetupWizard) configureMacOSPlatform(config *MacOSBuildConfig) {
	fmt.Printf("\n🖥️  macOS Configuration:\n")

	config.Enabled = true
	config.BuildTypes = []MacOSBuildType{
		{
			Name:      "release",
			Type:      "macos",
			BuildMode: "release",
		},
	}
}

// configureLinuxPlatform configures Linux build settings
func (sw *SetupWizard) configureLinuxPlatform(config *LinuxBuildConfig) {
	fmt.Printf("\n🐧 Linux Configuration:\n")

	config.Enabled = true
	config.BuildTypes = []LinuxBuildType{
		{
			Name:      "release",
			Type:      "linux",
			BuildMode: "release",
		},
	}
}

// configureWindowsPlatform configures Windows build settings
func (sw *SetupWizard) configureWindowsPlatform(config *WindowsBuildConfig) {
	fmt.Printf("\n🪟 Windows Configuration:\n")

	config.Enabled = true
	config.BuildTypes = []WindowsBuildType{
		{
			Name:      "release",
			Type:      "windows",
			BuildMode: "release",
		},
	}
}

// configureArtifacts configures artifact management
func (sw *SetupWizard) configureArtifacts(artifacts *ArtifactsConfig) error {
	fmt.Println("\n" + utils.Separator("-", 40))
	utils.Info("📦 Artifact Management Configuration")

	// Output directory
	defaultDir := artifacts.BaseOutputDir
	artifacts.BaseOutputDir = sw.promptString("Output directory", defaultDir)

	// Organization
	artifacts.Organization.ByDate = sw.promptYesNo("Organize artifacts by date?", true)
	artifacts.Organization.ByPlatform = sw.promptYesNo("Organize artifacts by platform?", true)
	artifacts.Organization.ByBuildType = sw.promptYesNo("Organize artifacts by build type?", true)

	// Naming
	fmt.Printf("\nCurrent naming pattern: %s\n", artifacts.Naming.Pattern)
	fmt.Println("Available variables: {app_name}, {version}, {arch}, {platform}, {build_type}")
	
	if sw.promptYesNo("Customize naming pattern?", false) {
		artifacts.Naming.Pattern = sw.promptString("Naming pattern", artifacts.Naming.Pattern)
	}

	// Cleanup
	artifacts.Cleanup.Enabled = sw.promptYesNo("Enable automatic cleanup?", true)
	if artifacts.Cleanup.Enabled {
		keepBuilds := sw.promptInt("Keep last N builds", artifacts.Cleanup.KeepLastBuilds)
		artifacts.Cleanup.KeepLastBuilds = keepBuilds

		maxAge := sw.promptInt("Maximum age in days", artifacts.Cleanup.MaxAgeDays)
		artifacts.Cleanup.MaxAgeDays = maxAge
	}

	return nil
}

// configureExecution configures build execution settings
func (sw *SetupWizard) configureExecution(execution *ExecutionConfig) error {
	fmt.Println("\n" + utils.Separator("-", 40))
	utils.Info("⚡ Build Execution Configuration")

	execution.SaveLogs = sw.promptYesNo("Save build logs?", true)
	execution.ContinueOnError = sw.promptYesNo("Continue building other platforms if one fails?", false)
	execution.ParallelBuilds = sw.promptYesNo("Enable parallel builds? (experimental)", false)

	if execution.ParallelBuilds {
		maxParallel := sw.promptInt("Maximum parallel builds", execution.MaxParallel)
		execution.MaxParallel = maxParallel
	}

	return nil
}

// Helper methods

func (sw *SetupWizard) fileExists(filename string) bool {
	path := filepath.Join(sw.ProjectPath, filename)
	_, err := os.Stat(path)
	return err == nil
}

func (sw *SetupWizard) dirExists(dirname string) bool {
	path := filepath.Join(sw.ProjectPath, dirname)
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func (sw *SetupWizard) findScripts(scriptsDir string) []string {
	var scripts []string
	entries, err := os.ReadDir(scriptsDir)
	if err != nil {
		return scripts
	}

	for _, entry := range entries {
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".sh") || strings.HasSuffix(entry.Name(), ".bat")) {
			scripts = append(scripts, entry.Name())
		}
	}

	return scripts
}

func (sw *SetupWizard) detectAvailablePlatforms() []Platform {
	var platforms []Platform

	platformDirs := map[Platform]string{
		PlatformAndroid: "android",
		PlatformIOS:     "ios",
		PlatformWeb:     "web",
		PlatformMacOS:   "macos",
		PlatformLinux:   "linux",
		PlatformWindows: "windows",
	}

	for platform, dir := range platformDirs {
		if sw.dirExists(dir) {
			platforms = append(platforms, platform)
		}
	}

	return platforms
}

func (sw *SetupWizard) formatPlatforms(platforms []Platform) string {
	var strs []string
	for _, p := range platforms {
		strs = append(strs, string(p))
	}
	return strings.Join(strs, ", ")
}

func (sw *SetupWizard) disablePlatform(platforms *PlatformsConfig, platform Platform) {
	switch platform {
	case PlatformAndroid:
		platforms.Android.Enabled = false
	case PlatformIOS:
		platforms.IOS.Enabled = false
	case PlatformWeb:
		platforms.Web.Enabled = false
	case PlatformMacOS:
		platforms.MacOS.Enabled = false
	case PlatformLinux:
		platforms.Linux.Enabled = false
	case PlatformWindows:
		platforms.Windows.Enabled = false
	}
}

func (sw *SetupWizard) promptString(prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	sw.scanner.Scan()
	input := strings.TrimSpace(sw.scanner.Text())
	if input == "" {
		return defaultValue
	}
	return input
}

func (sw *SetupWizard) promptYesNo(prompt string, defaultValue bool) bool {
	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}

	fmt.Printf("%s [%s]: ", prompt, defaultStr)
	sw.scanner.Scan()
	input := strings.ToLower(strings.TrimSpace(sw.scanner.Text()))

	if input == "" {
		return defaultValue
	}

	return input == "y" || input == "yes"
}

func (sw *SetupWizard) promptChoice(prompt string, choices []string, defaultChoice string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultChoice)
	sw.scanner.Scan()
	input := strings.TrimSpace(sw.scanner.Text())

	if input == "" {
		return defaultChoice
	}

	for _, choice := range choices {
		if input == choice {
			return input
		}
	}

	return defaultChoice
}

func (sw *SetupWizard) promptInt(prompt string, defaultValue int) int {
	fmt.Printf("%s [%d]: ", prompt, defaultValue)
	sw.scanner.Scan()
	input := strings.TrimSpace(sw.scanner.Text())

	if input == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(input); err == nil {
		return value
	}

	return defaultValue
}

func (sw *SetupWizard) displayConfigSummary(config *BuildConfig) {
	fmt.Println("\n" + utils.Separator("=", 60))
	utils.Success("📋 Configuration Summary")
	fmt.Println(utils.Separator("=", 60))

	fmt.Printf("App Name Source: %s\n", config.Metadata.AppNameSource)
	fmt.Printf("Version Source: %s\n", config.Metadata.VersionSource)
	fmt.Printf("Output Directory: %s\n", config.Artifacts.BaseOutputDir)
	fmt.Printf("Naming Pattern: %s\n", config.Artifacts.Naming.Pattern)

	fmt.Println("\nEnabled Platforms:")
	if config.Platforms.Android.Enabled {
		fmt.Printf("  ✅ Android (%d build types)\n", len(config.Platforms.Android.BuildTypes))
	}
	if config.Platforms.IOS.Enabled {
		fmt.Printf("  ✅ iOS (%d build types)\n", len(config.Platforms.IOS.BuildTypes))
	}
	if config.Platforms.Web.Enabled {
		fmt.Printf("  ✅ Web (%d build types)\n", len(config.Platforms.Web.BuildTypes))
	}
	if config.Platforms.MacOS.Enabled {
		fmt.Printf("  ✅ macOS (%d build types)\n", len(config.Platforms.MacOS.BuildTypes))
	}
	if config.Platforms.Linux.Enabled {
		fmt.Printf("  ✅ Linux (%d build types)\n", len(config.Platforms.Linux.BuildTypes))
	}
	if config.Platforms.Windows.Enabled {
		fmt.Printf("  ✅ Windows (%d build types)\n", len(config.Platforms.Windows.BuildTypes))
	}

	fmt.Printf("\nPre-build Steps: %d global", len(config.PreBuild.Global))
	if len(config.PreBuild.Android) > 0 {
		fmt.Printf(", %d Android", len(config.PreBuild.Android))
	}
	if len(config.PreBuild.IOS) > 0 {
		fmt.Printf(", %d iOS", len(config.PreBuild.IOS))
	}
	if len(config.PreBuild.Web) > 0 {
		fmt.Printf(", %d Web", len(config.PreBuild.Web))
	}
	fmt.Println()
}
