package build

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// BuildConfig represents the complete build configuration
type BuildConfig struct {
	Metadata  MetadataConfig  `yaml:"metadata"`
	PreBuild  PreBuildConfig  `yaml:"pre_build"`
	Platforms PlatformsConfig `yaml:"platforms"`
	Artifacts ArtifactsConfig `yaml:"artifacts"`
	Execution ExecutionConfig `yaml:"execution"`
}

// MetadataConfig contains build metadata configuration
type MetadataConfig struct {
	AppNameSource string `yaml:"app_name_source"` // namer, pubspec, custom
	CustomAppName string `yaml:"custom_app_name"` // used if source is custom
	VersionSource string `yaml:"version_source"`  // pubspec, custom
	CustomVersion string `yaml:"custom_version"`  // used if source is custom
}

// PreBuildConfig contains pre-build step configurations
type PreBuildConfig struct {
	Global  []BuildStep `yaml:"global"`
	Android []BuildStep `yaml:"android"`
	IOS     []BuildStep `yaml:"ios"`
	Web     []BuildStep `yaml:"web"`
}

// BuildStep represents a single build step
type BuildStep struct {
	Name        string            `yaml:"name"`
	Command     string            `yaml:"command"`
	Required    bool              `yaml:"required"`
	Timeout     int               `yaml:"timeout"` // seconds
	WorkingDir  string            `yaml:"working_dir"`
	Environment map[string]string `yaml:"env"`
	Condition   string            `yaml:"condition"`
}

// PlatformsConfig contains platform-specific configurations
type PlatformsConfig struct {
	Android AndroidBuildConfig `yaml:"android"`
	IOS     IOSBuildConfig     `yaml:"ios"`
	Web     WebBuildConfig     `yaml:"web"`
	MacOS   MacOSBuildConfig   `yaml:"macos"`
	Linux   LinuxBuildConfig   `yaml:"linux"`
	Windows WindowsBuildConfig `yaml:"windows"`
}

// AndroidBuildConfig contains Android build configuration
type AndroidBuildConfig struct {
	Enabled     bool               `yaml:"enabled"`
	BuildTypes  []AndroidBuildType `yaml:"build_types"`
	Environment map[string]string  `yaml:"environment"`
}

// AndroidBuildType represents an Android build type
type AndroidBuildType struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`       // apk, appbundle
	BuildMode   string   `yaml:"build_mode"` // release, debug, profile
	SplitPerABI bool     `yaml:"split_per_abi"`
	CustomArgs  []string `yaml:"custom_args"`
}

// IOSBuildConfig contains iOS build configuration
type IOSBuildConfig struct {
	Enabled    bool           `yaml:"enabled"`
	BuildTypes []IOSBuildType `yaml:"build_types"`
}

// IOSBuildType represents an iOS build type
type IOSBuildType struct {
	Name         string   `yaml:"name"`
	Type         string   `yaml:"type"`          // ipa, archive
	BuildMode    string   `yaml:"build_mode"`    // release, debug, profile
	ExportMethod string   `yaml:"export_method"` // app-store, development, ad-hoc
	CustomArgs   []string `yaml:"custom_args"`
}

// WebBuildConfig contains Web build configuration
type WebBuildConfig struct {
	Enabled    bool           `yaml:"enabled"`
	BuildTypes []WebBuildType `yaml:"build_types"`
}

// WebBuildType represents a Web build type
type WebBuildType struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`       // web
	BuildMode  string   `yaml:"build_mode"` // release, debug, profile
	PWA        bool     `yaml:"pwa"`
	CustomArgs []string `yaml:"custom_args"`
}

// MacOSBuildConfig contains macOS build configuration
type MacOSBuildConfig struct {
	Enabled    bool             `yaml:"enabled"`
	BuildTypes []MacOSBuildType `yaml:"build_types"`
}

// MacOSBuildType represents a macOS build type
type MacOSBuildType struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`       // macos
	BuildMode  string   `yaml:"build_mode"` // release, debug, profile
	CustomArgs []string `yaml:"custom_args"`
}

// LinuxBuildConfig contains Linux build configuration
type LinuxBuildConfig struct {
	Enabled    bool             `yaml:"enabled"`
	BuildTypes []LinuxBuildType `yaml:"build_types"`
}

// LinuxBuildType represents a Linux build type
type LinuxBuildType struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`       // linux
	BuildMode  string   `yaml:"build_mode"` // release, debug, profile
	CustomArgs []string `yaml:"custom_args"`
}

// WindowsBuildConfig contains Windows build configuration
type WindowsBuildConfig struct {
	Enabled    bool               `yaml:"enabled"`
	BuildTypes []WindowsBuildType `yaml:"build_types"`
}

// WindowsBuildType represents a Windows build type
type WindowsBuildType struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`       // windows
	BuildMode  string   `yaml:"build_mode"` // release, debug, profile
	CustomArgs []string `yaml:"custom_args"`
}

// ArtifactsConfig contains artifact management configuration
type ArtifactsConfig struct {
	BaseOutputDir string             `yaml:"base_output_dir"`
	Organization  OrganizationConfig `yaml:"organization"`
	Naming        NamingConfig       `yaml:"naming"`
	Cleanup       CleanupConfig      `yaml:"cleanup"`
}

// OrganizationConfig contains artifact organization settings
type OrganizationConfig struct {
	ByDate      bool   `yaml:"by_date"`
	DateFormat  string `yaml:"date_format"`
	ByPlatform  bool   `yaml:"by_platform"`
	ByBuildType bool   `yaml:"by_build_type"`
}

// NamingConfig contains artifact naming settings
type NamingConfig struct {
	Pattern         string `yaml:"pattern"`
	FallbackAppName string `yaml:"fallback_app_name"`
}

// CleanupConfig contains artifact cleanup settings
type CleanupConfig struct {
	Enabled        bool `yaml:"enabled"`
	KeepLastBuilds int  `yaml:"keep_last_builds"`
	MaxAgeDays     int  `yaml:"max_age_days"`
}

// ExecutionConfig contains build execution settings
type ExecutionConfig struct {
	ParallelBuilds  bool   `yaml:"parallel_builds"`
	MaxParallel     int    `yaml:"max_parallel"`
	ContinueOnError bool   `yaml:"continue_on_error"`
	SaveLogs        bool   `yaml:"save_logs"`
	LogLevel        string `yaml:"log_level"`
}

// LoadBuildConfig loads build configuration from file
func LoadBuildConfig(projectPath, configPath string) (*BuildConfig, error) {
	// Make path relative to project if not absolute
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(projectPath, configPath)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	var config BuildConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate and set defaults
	if err := validateAndSetDefaults(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// SaveBuildConfig saves build configuration to file
func SaveBuildConfig(projectPath, configPath string, config *BuildConfig) error {
	// Make path relative to project if not absolute
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(projectPath, configPath)
	}

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// DefaultBuildConfig returns a default build configuration
func DefaultBuildConfig() *BuildConfig {
	return &BuildConfig{
		Metadata: MetadataConfig{
			AppNameSource: "namer",
			VersionSource: "pubspec",
		},
		PreBuild: PreBuildConfig{
			Global: []BuildStep{
				{
					Name:     "Install dependencies",
					Command:  "flutter pub get",
					Required: true,
					Timeout:  300,
				},
				{
					Name:      "Generate code",
					Command:   "dart run build_runner build --delete-conflicting-outputs",
					Required:  false,
					Timeout:   600,
					Condition: "file_exists:build.yaml",
				},
			},
		},
		Platforms: PlatformsConfig{
			Android: AndroidBuildConfig{
				Enabled: true,
				BuildTypes: []AndroidBuildType{
					{
						Name:        "release_apk",
						Type:        "apk",
						BuildMode:   "release",
						SplitPerABI: true,
						CustomArgs:  []string{"--obfuscate", "--split-debug-info=build/debug-info"},
					},
				},
			},
			IOS: IOSBuildConfig{
				Enabled: true,
				BuildTypes: []IOSBuildType{
					{
						Name:         "archive",
						Type:         "archive",
						BuildMode:    "release",
						ExportMethod: "development",
						CustomArgs:   []string{"--no-codesign"},
					},
				},
			},
			Web: WebBuildConfig{
				Enabled: true,
				BuildTypes: []WebBuildType{
					{
						Name:       "release",
						Type:       "web",
						BuildMode:  "release",
						PWA:        true,
						CustomArgs: []string{},
					},
				},
			},
			MacOS: MacOSBuildConfig{
				Enabled: true,
				BuildTypes: []MacOSBuildType{
					{
						Name:      "release",
						Type:      "macos",
						BuildMode: "release",
					},
				},
			},
			Linux: LinuxBuildConfig{
				Enabled: true,
				BuildTypes: []LinuxBuildType{
					{
						Name:      "release",
						Type:      "linux",
						BuildMode: "release",
					},
				},
			},
			Windows: WindowsBuildConfig{
				Enabled: true,
				BuildTypes: []WindowsBuildType{
					{
						Name:      "release",
						Type:      "windows",
						BuildMode: "release",
					},
				},
			},
		},
		Artifacts: ArtifactsConfig{
			BaseOutputDir: "build/fdawg-outputs",
			Organization: OrganizationConfig{
				ByDate:      true,
				DateFormat:  "January-2",
				ByPlatform:  true,
				ByBuildType: true,
			},
			Naming: NamingConfig{
				Pattern:         "{app_name}_{version}_{arch}",
				FallbackAppName: "flutter_app",
			},
			Cleanup: CleanupConfig{
				Enabled:        true,
				KeepLastBuilds: 10,
				MaxAgeDays:     30,
			},
		},
		Execution: ExecutionConfig{
			ParallelBuilds:  false,
			MaxParallel:     2,
			ContinueOnError: false,
			SaveLogs:        true,
			LogLevel:        "info",
		},
	}
}

// validateAndSetDefaults validates configuration and sets defaults
func validateAndSetDefaults(config *BuildConfig) error {
	// Set default metadata values
	if config.Metadata.AppNameSource == "" {
		config.Metadata.AppNameSource = "namer"
	}
	if config.Metadata.VersionSource == "" {
		config.Metadata.VersionSource = "pubspec"
	}

	// Set default artifact values
	if config.Artifacts.BaseOutputDir == "" {
		config.Artifacts.BaseOutputDir = "build/fdawg-outputs"
	}
	if config.Artifacts.Organization.DateFormat == "" {
		config.Artifacts.Organization.DateFormat = "January-2"
	}
	if config.Artifacts.Naming.Pattern == "" {
		config.Artifacts.Naming.Pattern = "{app_name}_{version}_{arch}"
	}
	if config.Artifacts.Naming.FallbackAppName == "" {
		config.Artifacts.Naming.FallbackAppName = "flutter_app"
	}

	// Set default execution values
	if config.Execution.LogLevel == "" {
		config.Execution.LogLevel = "info"
	}
	if config.Execution.MaxParallel == 0 {
		config.Execution.MaxParallel = 2
	}

	// Set default timeouts for build steps
	for i := range config.PreBuild.Global {
		if config.PreBuild.Global[i].Timeout == 0 {
			config.PreBuild.Global[i].Timeout = 300 // 5 minutes default
		}
	}

	// Validate app name source
	validAppNameSources := []string{"namer", "pubspec", "custom"}
	if !contains(validAppNameSources, config.Metadata.AppNameSource) {
		return fmt.Errorf("invalid app_name_source: %s (must be one of: %v)",
			config.Metadata.AppNameSource, validAppNameSources)
	}

	// Validate version source
	validVersionSources := []string{"pubspec", "custom"}
	if !contains(validVersionSources, config.Metadata.VersionSource) {
		return fmt.Errorf("invalid version_source: %s (must be one of: %v)",
			config.Metadata.VersionSource, validVersionSources)
	}

	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetTimeout returns the timeout duration for a build step
func (bs *BuildStep) GetTimeout() time.Duration {
	if bs.Timeout <= 0 {
		return 5 * time.Minute // default timeout
	}
	return time.Duration(bs.Timeout) * time.Second
}
