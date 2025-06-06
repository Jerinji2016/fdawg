package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Platform-specific build implementations

// buildAndroid builds for Android platform
func (bm *BuildManager) buildAndroid(config *AndroidBuildConfig) ([]*BuildArtifact, error) {
	var allArtifacts []*BuildArtifact

	executor := NewCommandExecutor(bm.ProjectPath, bm.Logger)
	
	// Set Android environment variables
	if len(config.Environment) > 0 {
		executor.SetEnvironment(config.Environment)
	}

	for _, buildType := range config.BuildTypes {
		bm.Logger.Info("Building Android %s (%s)", buildType.Type, buildType.Name)

		// Construct Flutter build command
		args := []string{"build", buildType.Type}

		if buildType.BuildMode != "" {
			args = append(args, "--"+buildType.BuildMode)
		}

		// Handle split APKs
		if buildType.Type == "apk" && buildType.SplitPerABI {
			args = append(args, "--split-per-abi")
		}

		// Add custom arguments
		args = append(args, buildType.CustomArgs...)

		// Execute build
		if err := executor.ExecuteFlutterBuild(args, PlatformAndroid); err != nil {
			return allArtifacts, fmt.Errorf("Android %s build failed: %w", buildType.Type, err)
		}

		// Collect artifacts based on build type
		var artifacts []*BuildArtifact
		var err error

		if buildType.Type == "apk" && buildType.SplitPerABI {
			artifacts, err = bm.collectSplitAPKs()
		} else {
			artifacts, err = bm.collectAndroidArtifact(buildType.Type)
		}

		if err != nil {
			bm.Logger.Warning("Failed to collect Android artifacts: %v", err)
			continue
		}

		// Set build type for all artifacts
		for _, artifact := range artifacts {
			artifact.BuildType = buildType.Type
			artifact.BuildTime = time.Now()
		}

		allArtifacts = append(allArtifacts, artifacts...)
	}

	return allArtifacts, nil
}

// collectSplitAPKs collects split APK artifacts
func (bm *BuildManager) collectSplitAPKs() ([]*BuildArtifact, error) {
	var artifacts []*BuildArtifact

	apkDir := filepath.Join(bm.ProjectPath, "build", "app", "outputs", "flutter-apk")

	// Common APK architectures
	architectures := []string{"arm64-v8a", "armeabi-v7a", "x86_64"}

	for _, arch := range architectures {
		apkPattern := fmt.Sprintf("app-%s-release.apk", arch)
		apkPath := filepath.Join(apkDir, apkPattern)

		if _, err := os.Stat(apkPath); err == nil {
			artifact := &BuildArtifact{
				Platform:     PlatformAndroid,
				Architecture: arch,
				FileName:     apkPattern,
				FilePath:     apkPath,
			}

			artifacts = append(artifacts, artifact)
		}
	}

	if len(artifacts) == 0 {
		return nil, fmt.Errorf("no split APKs found in %s", apkDir)
	}

	return artifacts, nil
}

// collectAndroidArtifact collects a single Android artifact
func (bm *BuildManager) collectAndroidArtifact(buildType string) ([]*BuildArtifact, error) {
	var artifactPath string
	var fileName string

	switch buildType {
	case "apk":
		artifactPath = filepath.Join(bm.ProjectPath, "build", "app", "outputs", "flutter-apk", "app-release.apk")
		fileName = "app-release.apk"
	case "appbundle":
		artifactPath = filepath.Join(bm.ProjectPath, "build", "app", "outputs", "bundle", "release", "app-release.aab")
		fileName = "app-release.aab"
	default:
		return nil, fmt.Errorf("unknown Android build type: %s", buildType)
	}

	if _, err := os.Stat(artifactPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Android artifact not found: %s", artifactPath)
	}

	artifact := &BuildArtifact{
		Platform:     PlatformAndroid,
		Architecture: "universal",
		FileName:     fileName,
		FilePath:     artifactPath,
	}

	return []*BuildArtifact{artifact}, nil
}

// buildIOS builds for iOS platform
func (bm *BuildManager) buildIOS(config *IOSBuildConfig) ([]*BuildArtifact, error) {
	var allArtifacts []*BuildArtifact

	executor := NewCommandExecutor(bm.ProjectPath, bm.Logger)

	for _, buildType := range config.BuildTypes {
		bm.Logger.Info("Building iOS %s (%s)", buildType.Type, buildType.Name)

		// Construct Flutter build command
		args := []string{"build", buildType.Type}

		if buildType.BuildMode != "" {
			args = append(args, "--"+buildType.BuildMode)
		}

		// Add export method for IPA builds
		if buildType.Type == "ipa" && buildType.ExportMethod != "" {
			args = append(args, "--export-method", buildType.ExportMethod)
		}

		// Add custom arguments
		args = append(args, buildType.CustomArgs...)

		// Execute build
		if err := executor.ExecuteFlutterBuild(args, PlatformIOS); err != nil {
			return allArtifacts, fmt.Errorf("iOS %s build failed: %w", buildType.Type, err)
		}

		// Collect artifacts
		artifacts, err := bm.collectIOSArtifacts(buildType.Type)
		if err != nil {
			bm.Logger.Warning("Failed to collect iOS artifacts: %v", err)
			continue
		}

		// Set build type for all artifacts
		for _, artifact := range artifacts {
			artifact.BuildType = buildType.Type
			artifact.BuildTime = time.Now()
		}

		allArtifacts = append(allArtifacts, artifacts...)
	}

	return allArtifacts, nil
}

// collectIOSArtifacts collects iOS artifacts
func (bm *BuildManager) collectIOSArtifacts(buildType string) ([]*BuildArtifact, error) {
	switch buildType {
	case "archive":
		return bm.collectXCArchive()
	case "ipa":
		return bm.collectIPA()
	default:
		return nil, fmt.Errorf("unknown iOS build type: %s", buildType)
	}
}

// collectXCArchive collects .xcarchive artifacts
func (bm *BuildManager) collectXCArchive() ([]*BuildArtifact, error) {
	archivePath := bm.findXCArchive()
	if archivePath == "" {
		return nil, fmt.Errorf("no .xcarchive found")
	}

	artifact := &BuildArtifact{
		Platform:     PlatformIOS,
		Architecture: "universal",
		FileName:     filepath.Base(archivePath),
		FilePath:     archivePath,
	}

	return []*BuildArtifact{artifact}, nil
}

// collectIPA collects .ipa artifacts
func (bm *BuildManager) collectIPA() ([]*BuildArtifact, error) {
	ipaPath := bm.findIPA()
	if ipaPath == "" {
		return nil, fmt.Errorf("no .ipa found")
	}

	artifact := &BuildArtifact{
		Platform:     PlatformIOS,
		Architecture: "universal",
		FileName:     filepath.Base(ipaPath),
		FilePath:     ipaPath,
	}

	return []*BuildArtifact{artifact}, nil
}

// findXCArchive finds the .xcarchive file
func (bm *BuildManager) findXCArchive() string {
	buildDir := filepath.Join(bm.ProjectPath, "build", "ios", "archive")
	
	entries, err := os.ReadDir(buildDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), ".xcarchive") {
			return filepath.Join(buildDir, entry.Name())
		}
	}

	return ""
}

// findIPA finds the .ipa file
func (bm *BuildManager) findIPA() string {
	buildDir := filepath.Join(bm.ProjectPath, "build", "ios", "ipa")
	
	entries, err := os.ReadDir(buildDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".ipa") {
			return filepath.Join(buildDir, entry.Name())
		}
	}

	return ""
}

// buildWeb builds for Web platform
func (bm *BuildManager) buildWeb(config *WebBuildConfig) ([]*BuildArtifact, error) {
	var allArtifacts []*BuildArtifact

	executor := NewCommandExecutor(bm.ProjectPath, bm.Logger)

	for _, buildType := range config.BuildTypes {
		bm.Logger.Info("Building Web %s (%s)", buildType.Type, buildType.Name)

		// Construct Flutter build command
		args := []string{"build", "web"}

		if buildType.BuildMode != "" {
			args = append(args, "--"+buildType.BuildMode)
		}

		// Add custom arguments
		args = append(args, buildType.CustomArgs...)

		// Execute build
		if err := executor.ExecuteFlutterBuild(args, PlatformWeb); err != nil {
			return allArtifacts, fmt.Errorf("Web build failed: %w", err)
		}

		// Collect artifacts
		artifacts, err := bm.collectWebArtifacts()
		if err != nil {
			bm.Logger.Warning("Failed to collect Web artifacts: %v", err)
			continue
		}

		// Set build type for all artifacts
		for _, artifact := range artifacts {
			artifact.BuildType = buildType.Type
			artifact.BuildTime = time.Now()
		}

		allArtifacts = append(allArtifacts, artifacts...)
	}

	return allArtifacts, nil
}

// collectWebArtifacts collects Web build artifacts
func (bm *BuildManager) collectWebArtifacts() ([]*BuildArtifact, error) {
	webBuildDir := filepath.Join(bm.ProjectPath, "build", "web")
	
	if _, err := os.Stat(webBuildDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("web build directory not found: %s", webBuildDir)
	}

	// Create a zip archive of the web build
	zipPath := filepath.Join(bm.ProjectPath, "build", "web.zip")
	if err := bm.createWebArchive(webBuildDir, zipPath); err != nil {
		return nil, fmt.Errorf("failed to create web archive: %w", err)
	}

	artifact := &BuildArtifact{
		Platform:     PlatformWeb,
		Architecture: "web",
		FileName:     "web.zip",
		FilePath:     zipPath,
	}

	return []*BuildArtifact{artifact}, nil
}

// createWebArchive creates a zip archive of the web build
func (bm *BuildManager) createWebArchive(sourceDir, zipPath string) error {
	// Simple implementation using system zip command
	// TODO: Implement proper zip creation in Go
	cmd := exec.Command("zip", "-r", zipPath, ".")
	cmd.Dir = sourceDir
	
	return cmd.Run()
}

// buildMacOS builds for macOS platform
func (bm *BuildManager) buildMacOS(config *MacOSBuildConfig) ([]*BuildArtifact, error) {
	return bm.buildDesktop(PlatformMacOS, config.BuildTypes[0].BuildMode, config.BuildTypes[0].CustomArgs)
}

// buildLinux builds for Linux platform
func (bm *BuildManager) buildLinux(config *LinuxBuildConfig) ([]*BuildArtifact, error) {
	return bm.buildDesktop(PlatformLinux, config.BuildTypes[0].BuildMode, config.BuildTypes[0].CustomArgs)
}

// buildWindows builds for Windows platform
func (bm *BuildManager) buildWindows(config *WindowsBuildConfig) ([]*BuildArtifact, error) {
	return bm.buildDesktop(PlatformWindows, config.BuildTypes[0].BuildMode, config.BuildTypes[0].CustomArgs)
}

// buildDesktop builds for desktop platforms (macOS, Linux, Windows)
func (bm *BuildManager) buildDesktop(platform Platform, buildMode string, customArgs []string) ([]*BuildArtifact, error) {
	executor := NewCommandExecutor(bm.ProjectPath, bm.Logger)

	bm.Logger.Info("Building %s", platform)

	// Construct Flutter build command
	args := []string{"build", string(platform)}

	if buildMode != "" {
		args = append(args, "--"+buildMode)
	}

	// Add custom arguments
	args = append(args, customArgs...)

	// Execute build
	if err := executor.ExecuteFlutterBuild(args, platform); err != nil {
		return nil, fmt.Errorf("%s build failed: %w", platform, err)
	}

	// Collect artifacts
	artifacts, err := bm.collectDesktopArtifacts(platform)
	if err != nil {
		return nil, fmt.Errorf("failed to collect %s artifacts: %w", platform, err)
	}

	// Set build time
	for _, artifact := range artifacts {
		artifact.BuildTime = time.Now()
	}

	return artifacts, nil
}

// collectDesktopArtifacts collects desktop platform artifacts
func (bm *BuildManager) collectDesktopArtifacts(platform Platform) ([]*BuildArtifact, error) {
	var buildDir string
	var expectedExtension string

	switch platform {
	case PlatformMacOS:
		buildDir = filepath.Join(bm.ProjectPath, "build", "macos", "Build", "Products", "Release")
		expectedExtension = ".app"
	case PlatformLinux:
		buildDir = filepath.Join(bm.ProjectPath, "build", "linux", "x64", "release", "bundle")
		expectedExtension = "" // Linux executables don't have extensions
	case PlatformWindows:
		buildDir = filepath.Join(bm.ProjectPath, "build", "windows", "runner", "Release")
		expectedExtension = ".exe"
	default:
		return nil, fmt.Errorf("unsupported desktop platform: %s", platform)
	}

	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("build directory not found: %s", buildDir)
	}

	entries, err := os.ReadDir(buildDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read build directory: %w", err)
	}

	var artifacts []*BuildArtifact

	for _, entry := range entries {
		if platform == PlatformMacOS && entry.IsDir() && strings.HasSuffix(entry.Name(), expectedExtension) {
			// macOS .app bundle
			appPath := filepath.Join(buildDir, entry.Name())
			arch := bm.detectMacOSArchitecture(appPath)
			
			artifact := &BuildArtifact{
				Platform:     platform,
				Architecture: arch,
				FileName:     entry.Name(),
				FilePath:     appPath,
			}
			artifacts = append(artifacts, artifact)
		} else if !entry.IsDir() && (expectedExtension == "" || strings.HasSuffix(entry.Name(), expectedExtension)) {
			// Linux/Windows executable
			artifact := &BuildArtifact{
				Platform:     platform,
				Architecture: bm.detectArchitecture(),
				FileName:     entry.Name(),
				FilePath:     filepath.Join(buildDir, entry.Name()),
			}
			artifacts = append(artifacts, artifact)
		}
	}

	if len(artifacts) == 0 {
		return nil, fmt.Errorf("no %s artifacts found in %s", platform, buildDir)
	}

	return artifacts, nil
}

// detectMacOSArchitecture detects the architecture of a macOS app bundle
func (bm *BuildManager) detectMacOSArchitecture(appPath string) string {
	// Use 'file' command to detect architecture
	executableName := strings.TrimSuffix(filepath.Base(appPath), ".app")
	executablePath := filepath.Join(appPath, "Contents", "MacOS", executableName)

	cmd := exec.Command("file", executablePath)
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "arm64") {
		return "arm64"
	} else if strings.Contains(outputStr, "x86_64") {
		return "x86_64"
	}

	return "universal"
}

// detectArchitecture detects the current system architecture
func (bm *BuildManager) detectArchitecture() string {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "x64" // Default fallback
	}

	arch := strings.TrimSpace(string(output))
	switch arch {
	case "x86_64", "amd64":
		return "x64"
	case "arm64", "aarch64":
		return "arm64"
	case "i386", "i686":
		return "x86"
	default:
		return arch
	}
}
