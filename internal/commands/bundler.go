package commands

import (
	"fmt"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/bundler"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/utils"
	"github.com/urfave/cli/v2"
)

// BundlerCommand returns the CLI command for managing bundle IDs
func BundlerCommand() *cli.Command {
	return &cli.Command{
		Name:        "bundler",
		Usage:       "Manage bundle IDs for Flutter projects",
		Description: "Commands for managing bundle IDs across all platforms",
		Subcommands: []*cli.Command{
			{
				Name:        "get",
				Usage:       "Get current bundle IDs",
				Description: "Retrieves current bundle IDs for specified platforms",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "platforms",
						Aliases: []string{"p"},
						Usage:   "Platforms to get bundle IDs from (android, ios, macos, linux, windows, web)",
					},
				},
				Action: getBundleIDs,
			},
			{
				Name:        "set",
				Usage:       "Set bundle IDs",
				Description: "Sets bundle IDs universally across all platforms or for specific platforms",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "value",
						Aliases: []string{"v"},
						Usage:   "Universal bundle ID for all platforms",
					},
					&cli.StringSliceFlag{
						Name:    "platforms",
						Aliases: []string{"p"},
						Usage:   "Platforms to apply universal bundle ID to (android, ios, macos, linux, windows, web)",
					},
					&cli.StringFlag{
						Name:  "android",
						Usage: "Bundle ID for Android platform",
					},
					&cli.StringFlag{
						Name:  "ios",
						Usage: "Bundle ID for iOS platform",
					},
					&cli.StringFlag{
						Name:  "macos",
						Usage: "Bundle ID for macOS platform",
					},
					&cli.StringFlag{
						Name:  "linux",
						Usage: "Bundle ID for Linux platform",
					},
					&cli.StringFlag{
						Name:  "windows",
						Usage: "Bundle ID for Windows platform",
					},
					&cli.StringFlag{
						Name:  "web",
						Usage: "Bundle ID for Web platform",
					},
				},
				Action: setBundleIDs,
			},
			{
				Name:        "list",
				Usage:       "List all current bundle IDs",
				Description: "Lists current bundle IDs for all available platforms (alias for get with no platforms)",
				Action:      listBundleIDs,
			},
			{
				Name:        "validate",
				Usage:       "Validate bundle ID format",
				Description: "Validates bundle ID format and naming conventions",
				ArgsUsage:   "<bundle-id>",
				Action:      validateBundleID,
			},
		},
	}
}

// validateFlutterProjectForBundler validates that the current directory is a Flutter project
func validateFlutterProjectForBundler() (*flutter.ValidationResult, error) {
	result, err := flutter.ValidateProject(".")
	if err != nil {
		utils.Error("Failed to validate Flutter project: %v", err)
		return nil, err
	}

	if len(result.MissingDirs) > 0 {
		utils.Warning("Some Flutter directories are missing. This might affect bundle ID management.")
	}

	return result, nil
}

// parseBundlerPlatforms parses platform strings into bundler Platform enums
func parseBundlerPlatforms(platformStrings []string) ([]bundler.Platform, error) {
	if len(platformStrings) == 0 {
		return nil, nil
	}

	var platforms []bundler.Platform
	validPlatforms := map[string]bundler.Platform{
		"android": bundler.PlatformAndroid,
		"ios":     bundler.PlatformIOS,
		"macos":   bundler.PlatformMacOS,
		"linux":   bundler.PlatformLinux,
		"windows": bundler.PlatformWindows,
		"web":     bundler.PlatformWeb,
	}

	for _, platformStr := range platformStrings {
		platform, exists := validPlatforms[strings.ToLower(platformStr)]
		if !exists {
			return nil, fmt.Errorf("invalid platform: %s. Valid platforms are: android, ios, macos, linux, windows, web", platformStr)
		}
		platforms = append(platforms, platform)
	}

	return platforms, nil
}

// getBundleIDs retrieves bundle IDs for specified platforms
func getBundleIDs(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForBundler()
	if err != nil {
		return err
	}

	// Parse platforms
	platforms, err := parseBundlerPlatforms(c.StringSlice("platforms"))
	if err != nil {
		return err
	}

	// Get bundle IDs
	utils.Info("Retrieving bundle IDs...")
	result, err := bundler.GetBundleIDs(project.ProjectPath, platforms)
	if err != nil {
		utils.Error("Failed to get bundle IDs: %v", err)
		return err
	}

	// Display results
	displayBundleIDs(result)
	return nil
}

// setBundleIDs sets bundle IDs for platforms
func setBundleIDs(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForBundler()
	if err != nil {
		return err
	}

	// Build request
	request := &bundler.BundleIDRequest{
		ProjectPath: project.ProjectPath,
		Universal:   c.String("value"),
		Platforms:   make(map[bundler.Platform]string),
	}

	// Add platform-specific bundle IDs
	platformFlags := map[string]bundler.Platform{
		"android": bundler.PlatformAndroid,
		"ios":     bundler.PlatformIOS,
		"macos":   bundler.PlatformMacOS,
		"linux":   bundler.PlatformLinux,
		"windows": bundler.PlatformWindows,
		"web":     bundler.PlatformWeb,
	}

	for flag, platform := range platformFlags {
		if bundleID := c.String(flag); bundleID != "" {
			request.Platforms[platform] = bundleID
		}
	}

	// If universal value is provided with platforms flag, apply to specific platforms only
	if request.Universal != "" && len(c.StringSlice("platforms")) > 0 {
		platforms, err := parseBundlerPlatforms(c.StringSlice("platforms"))
		if err != nil {
			return err
		}

		// Clear universal and set platform-specific
		universalValue := request.Universal
		request.Universal = ""
		for _, platform := range platforms {
			request.Platforms[platform] = universalValue
		}
	}

	// Validate request
	if request.Universal == "" && len(request.Platforms) == 0 {
		utils.Error("Either --value or platform-specific flags must be provided")
		utils.Info("Usage examples:")
		utils.Info("  fdawg bundler set --value \"com.mycompany.myapp\"")
		utils.Info("  fdawg bundler set --android \"com.mycompany.myapp.android\" --ios \"com.mycompany.myapp.ios\"")
		utils.Info("  fdawg bundler set --value \"com.mycompany.myapp\" --platforms android,ios")
		return fmt.Errorf("no bundle ID provided")
	}

	// Validate bundle ID format
	bundleIDsToValidate := []string{}
	if request.Universal != "" {
		bundleIDsToValidate = append(bundleIDsToValidate, request.Universal)
	}
	for _, bundleID := range request.Platforms {
		bundleIDsToValidate = append(bundleIDsToValidate, bundleID)
	}

	for _, bundleID := range bundleIDsToValidate {
		if err := validateBundleIDFormat(bundleID); err != nil {
			utils.Error("Invalid bundle ID format: %v", err)
			return err
		}
	}

	// Set bundle IDs
	utils.Info("Setting bundle IDs...")
	if err := bundler.SetBundleIDs(request); err != nil {
		utils.Error("Failed to set bundle IDs: %v", err)
		return err
	}

	utils.Success("Bundle IDs updated successfully")

	// Show updated bundle IDs
	utils.Info("Updated bundle IDs:")
	result, _ := bundler.GetBundleIDs(project.ProjectPath, nil)
	if result != nil {
		displayBundleIDs(result)
	}

	return nil
}

// listBundleIDs lists all current bundle IDs (alias for get with no platforms)
func listBundleIDs(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForBundler()
	if err != nil {
		return err
	}

	// Get all bundle IDs
	utils.Info("Listing all bundle IDs...")
	result, err := bundler.GetBundleIDs(project.ProjectPath, nil)
	if err != nil {
		utils.Error("Failed to list bundle IDs: %v", err)
		return err
	}

	// Display results
	displayBundleIDs(result)
	return nil
}

// validateBundleID validates a bundle ID format
func validateBundleID(c *cli.Context) error {
	if c.Args().Len() == 0 {
		utils.Error("Bundle ID argument is required")
		utils.Info("Usage: fdawg bundler validate <bundle-id>")
		return fmt.Errorf("bundle ID argument missing")
	}

	bundleID := c.Args().Get(0)

	if err := validateBundleIDFormat(bundleID); err != nil {
		utils.Error("Invalid bundle ID: %v", err)
		return err
	}

	utils.Success("Bundle ID '%s' is valid", bundleID)
	return nil
}

// displayBundleIDs displays bundle ID results in a formatted way
func displayBundleIDs(result *bundler.BundleIDResult) {
	if len(result.BundleIDs) == 0 {
		utils.Warning("No platforms found in project")
		return
	}

	utils.Info("Current Bundle IDs:")

	availableCount := 0
	for _, bundleInfo := range result.BundleIDs {
		if bundleInfo.Available {
			availableCount++
			platformName := strings.ToUpper(string(bundleInfo.Platform)[:1]) + string(bundleInfo.Platform)[1:]
			if bundleInfo.Error != "" {
				utils.Error("✗ %s: %s", platformName, bundleInfo.Error)
			} else {
				if bundleInfo.Namespace != "" && bundleInfo.Namespace != bundleInfo.BundleID {
					utils.Success("✓ %s: %s (namespace: %s)", platformName, bundleInfo.BundleID, bundleInfo.Namespace)
				} else {
					utils.Success("✓ %s: %s", platformName, bundleInfo.BundleID)
				}
			}
		}
	}

	utils.Info("")
	utils.Info("Platforms found: %d/%d", availableCount, len(result.BundleIDs))
}

// validateBundleIDFormat validates bundle ID format and naming conventions
func validateBundleIDFormat(bundleID string) error {
	if bundleID == "" {
		return fmt.Errorf("bundle ID cannot be empty")
	}

	// Check for reverse domain notation (at least one dot)
	if !strings.Contains(bundleID, ".") {
		return fmt.Errorf("bundle ID should follow reverse domain notation (e.g., com.company.app)")
	}

	// Check for valid characters (alphanumeric, dots, hyphens, underscores)
	for _, char := range bundleID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '.' || char == '-' || char == '_') {
			return fmt.Errorf("bundle ID contains invalid character: %c. Only alphanumeric characters, dots, hyphens, and underscores are allowed", char)
		}
	}

	// Check that it doesn't start or end with a dot
	if strings.HasPrefix(bundleID, ".") || strings.HasSuffix(bundleID, ".") {
		return fmt.Errorf("bundle ID cannot start or end with a dot")
	}

	// Check for consecutive dots
	if strings.Contains(bundleID, "..") {
		return fmt.Errorf("bundle ID cannot contain consecutive dots")
	}

	// Check minimum length
	if len(bundleID) < 3 {
		return fmt.Errorf("bundle ID is too short (minimum 3 characters)")
	}

	// Check maximum length (reasonable limit)
	if len(bundleID) > 255 {
		return fmt.Errorf("bundle ID is too long (maximum 255 characters)")
	}

	// Split by dots and validate each segment
	segments := strings.Split(bundleID, ".")
	for i, segment := range segments {
		if segment == "" {
			return fmt.Errorf("bundle ID segment %d is empty", i+1)
		}

		// Each segment should not start with a number (Java package naming convention)
		if len(segment) > 0 && segment[0] >= '0' && segment[0] <= '9' {
			return fmt.Errorf("bundle ID segment '%s' cannot start with a number", segment)
		}
	}

	return nil
}
