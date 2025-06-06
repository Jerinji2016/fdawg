package build

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/utils"
)

// CommandExecutor handles execution of build commands
type CommandExecutor struct {
	WorkingDir  string
	Environment map[string]string
	Logger      *utils.Logger
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(workingDir string, logger *utils.Logger) *CommandExecutor {
	return &CommandExecutor{
		WorkingDir:  workingDir,
		Environment: make(map[string]string),
		Logger:      logger,
	}
}

// ExecuteStep executes a single build step
func (ce *CommandExecutor) ExecuteStep(step BuildStep) error {
	// Check conditions before execution
	if !ce.shouldExecuteStep(step) {
		ce.Logger.Info("Skipping step '%s' - condition not met: %s", step.Name, step.Condition)
		return nil
	}

	ce.Logger.Info("Executing: %s", step.Name)
	ce.Logger.Debug("Command: %s", step.Command)

	// Setup command context with timeout
	timeout := step.GetTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create command
	cmd := exec.CommandContext(ctx, "sh", "-c", step.Command)
	cmd.Dir = ce.resolveWorkingDir(step.WorkingDir)
	cmd.Env = ce.buildEnvironment(step.Environment)

	// Execute with real-time output
	return ce.executeWithOutput(cmd, step.Name)
}

// ExecuteFlutterBuild executes a Flutter build command
func (ce *CommandExecutor) ExecuteFlutterBuild(args []string, platform Platform) error {
	ce.Logger.Info("Executing Flutter build for %s", platform)
	ce.Logger.Debug("Command: flutter %s", strings.Join(args, " "))

	// Create command
	cmd := exec.Command("flutter", args...)
	cmd.Dir = ce.WorkingDir
	cmd.Env = ce.buildEnvironment(nil)

	// Execute with real-time output
	return ce.executeWithOutput(cmd, fmt.Sprintf("Flutter build %s", platform))
}

// shouldExecuteStep evaluates whether a step should be executed based on conditions
func (ce *CommandExecutor) shouldExecuteStep(step BuildStep) bool {
	if step.Condition == "" {
		return true
	}

	// Parse and evaluate conditions
	switch {
	case strings.HasPrefix(step.Condition, "file_exists:"):
		file := strings.TrimPrefix(step.Condition, "file_exists:")
		return ce.fileExists(file)

	case strings.HasPrefix(step.Condition, "platform_available:"):
		platform := strings.TrimPrefix(step.Condition, "platform_available:")
		return ce.isPlatformAvailable(platform)

	case strings.HasPrefix(step.Condition, "env_set:"):
		envVar := strings.TrimPrefix(step.Condition, "env_set:")
		return os.Getenv(envVar) != ""

	case strings.HasPrefix(step.Condition, "command_exists:"):
		command := strings.TrimPrefix(step.Condition, "command_exists:")
		return ce.commandExists(command)

	case strings.HasPrefix(step.Condition, "dir_exists:"):
		dir := strings.TrimPrefix(step.Condition, "dir_exists:")
		return ce.dirExists(dir)

	default:
		ce.Logger.Warning("Unknown condition: %s", step.Condition)
		return true
	}
}

// executeWithOutput executes a command and streams output in real-time
func (ce *CommandExecutor) executeWithOutput(cmd *exec.Cmd, stepName string) error {
	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Stream output in real-time
	go ce.streamOutput(stdout, stepName, "stdout")
	go ce.streamOutput(stderr, stepName, "stderr")

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// streamOutput streams command output in real-time
func (ce *CommandExecutor) streamOutput(reader io.Reader, stepName, streamType string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if streamType == "stderr" {
			ce.Logger.Warning("[%s] %s", stepName, line)
		} else {
			ce.Logger.Debug("[%s] %s", stepName, line)
		}
	}
}

// resolveWorkingDir resolves the working directory for a command
func (ce *CommandExecutor) resolveWorkingDir(workingDir string) string {
	if workingDir == "" {
		return ce.WorkingDir
	}

	if filepath.IsAbs(workingDir) {
		return workingDir
	}

	return filepath.Join(ce.WorkingDir, workingDir)
}

// buildEnvironment builds the environment variables for a command
func (ce *CommandExecutor) buildEnvironment(stepEnv map[string]string) []string {
	// Start with system environment
	env := os.Environ()

	// Add executor environment
	for key, value := range ce.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	// Add step-specific environment
	for key, value := range stepEnv {
		// Expand environment variables in values
		expandedValue := os.ExpandEnv(value)
		env = append(env, fmt.Sprintf("%s=%s", key, expandedValue))
	}

	return env
}

// fileExists checks if a file exists
func (ce *CommandExecutor) fileExists(filename string) bool {
	path := filename
	if !filepath.IsAbs(path) {
		path = filepath.Join(ce.WorkingDir, filename)
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists
func (ce *CommandExecutor) dirExists(dirname string) bool {
	path := dirname
	if !filepath.IsAbs(path) {
		path = filepath.Join(ce.WorkingDir, dirname)
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// isPlatformAvailable checks if a platform is available in the project
func (ce *CommandExecutor) isPlatformAvailable(platform string) bool {
	platformPath := filepath.Join(ce.WorkingDir, platform)
	return ce.dirExists(platformPath)
}

// commandExists checks if a command exists in PATH
func (ce *CommandExecutor) commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// SetEnvironment sets environment variables for the executor
func (ce *CommandExecutor) SetEnvironment(env map[string]string) {
	for key, value := range env {
		ce.Environment[key] = value
	}
}

// ExecuteWithProgress executes a command with progress indication
func (ce *CommandExecutor) ExecuteWithProgress(cmd *exec.Cmd, stepName string, progressCallback func(string)) error {
	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	// Stream output with progress callback
	go ce.streamOutputWithProgress(stdout, stepName, "stdout", progressCallback)
	go ce.streamOutputWithProgress(stderr, stepName, "stderr", progressCallback)

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// streamOutputWithProgress streams output with progress callback
func (ce *CommandExecutor) streamOutputWithProgress(reader io.Reader, stepName, streamType string, progressCallback func(string)) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		// Call progress callback
		if progressCallback != nil {
			progressCallback(line)
		}

		// Log the output
		if streamType == "stderr" {
			ce.Logger.Warning("[%s] %s", stepName, line)
		} else {
			ce.Logger.Debug("[%s] %s", stepName, line)
		}
	}
}

// ValidateFlutterInstallation validates that Flutter is installed and available
func (ce *CommandExecutor) ValidateFlutterInstallation() error {
	if !ce.commandExists("flutter") {
		return fmt.Errorf("Flutter is not installed or not in PATH")
	}

	// Check Flutter version
	cmd := exec.Command("flutter", "--version")
	cmd.Dir = ce.WorkingDir

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Flutter version: %w", err)
	}

	ce.Logger.Info("Flutter installation validated")
	ce.Logger.Debug("Flutter version: %s", strings.TrimSpace(string(output)))

	return nil
}

// ValidatePlatformRequirements validates platform-specific requirements
func (ce *CommandExecutor) ValidatePlatformRequirements(platform Platform) error {
	switch platform {
	case PlatformAndroid:
		return ce.validateAndroidRequirements()
	case PlatformIOS:
		return ce.validateIOSRequirements()
	case PlatformWeb:
		return ce.validateWebRequirements()
	case PlatformMacOS:
		return ce.validateMacOSRequirements()
	case PlatformLinux:
		return ce.validateLinuxRequirements()
	case PlatformWindows:
		return ce.validateWindowsRequirements()
	default:
		return fmt.Errorf("unknown platform: %s", platform)
	}
}

// validateAndroidRequirements validates Android build requirements
func (ce *CommandExecutor) validateAndroidRequirements() error {
	// Check if Android SDK is available
	if os.Getenv("ANDROID_HOME") == "" && os.Getenv("ANDROID_SDK_ROOT") == "" {
		return fmt.Errorf("Android SDK not found (ANDROID_HOME or ANDROID_SDK_ROOT not set)")
	}

	// Check if Java is available
	if !ce.commandExists("java") {
		return fmt.Errorf("Java is not installed or not in PATH")
	}

	return nil
}

// validateIOSRequirements validates iOS build requirements
func (ce *CommandExecutor) validateIOSRequirements() error {
	// Check if running on macOS
	if !ce.commandExists("xcodebuild") {
		return fmt.Errorf("Xcode is not installed (iOS builds require macOS with Xcode)")
	}

	return nil
}

// validateWebRequirements validates Web build requirements
func (ce *CommandExecutor) validateWebRequirements() error {
	// Web builds don't have special requirements beyond Flutter
	return nil
}

// validateMacOSRequirements validates macOS build requirements
func (ce *CommandExecutor) validateMacOSRequirements() error {
	// Check if running on macOS
	if !ce.commandExists("xcodebuild") {
		return fmt.Errorf("Xcode is not installed (macOS builds require macOS with Xcode)")
	}

	return nil
}

// validateLinuxRequirements validates Linux build requirements
func (ce *CommandExecutor) validateLinuxRequirements() error {
	// Check for common Linux build tools
	requiredTools := []string{"cmake", "ninja-build"}

	for _, tool := range requiredTools {
		if !ce.commandExists(tool) {
			ce.Logger.Warning("Recommended tool '%s' not found", tool)
		}
	}

	return nil
}

// validateWindowsRequirements validates Windows build requirements
func (ce *CommandExecutor) validateWindowsRequirements() error {
	// Check for Visual Studio or Build Tools
	if !ce.commandExists("msbuild") {
		ce.Logger.Warning("MSBuild not found - Visual Studio or Build Tools may be required")
	}

	return nil
}
