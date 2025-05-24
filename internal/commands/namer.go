package commands

import (
	"fmt"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/namer"
	"github.com/Jerinji2016/fdawg/pkg/utils"
	"github.com/urfave/cli/v2"
)

// NamerCommand returns the CLI command for managing app names
func NamerCommand() *cli.Command {
	return &cli.Command{
		Name:        "namer",
		Usage:       "Manage app display names for Flutter projects",
		Description: "Commands for managing app display names across all platforms",
		Subcommands: []*cli.Command{
			{
				Name:        "get",
				Usage:       "Get current app names",
				Description: "Retrieves current app names for specified platforms",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "platforms",
						Aliases: []string{"p"},
						Usage:   "Platforms to get names from (android, ios, macos, linux, windows, web)",
					},
				},
				Action: getAppNames,
			},
			{
				Name:        "set",
				Usage:       "Set app names",
				Description: "Sets app names for specified platforms",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "value",
						Aliases: []string{"v"},
						Usage:   "App name to set universally across all platforms",
					},
					&cli.StringSliceFlag{
						Name:    "platforms",
						Aliases: []string{"p"},
						Usage:   "Platforms to set names for (android, ios, macos, linux, windows, web)",
					},
					&cli.StringFlag{
						Name:  "android",
						Usage: "App name for Android platform",
					},
					&cli.StringFlag{
						Name:  "ios",
						Usage: "App name for iOS platform",
					},
					&cli.StringFlag{
						Name:  "macos",
						Usage: "App name for macOS platform",
					},
					&cli.StringFlag{
						Name:  "linux",
						Usage: "App name for Linux platform",
					},
					&cli.StringFlag{
						Name:  "windows",
						Usage: "App name for Windows platform",
					},
					&cli.StringFlag{
						Name:  "web",
						Usage: "App name for Web platform",
					},
				},
				Action: setAppNames,
			},
			{
				Name:        "list",
				Usage:       "List current app names",
				Description: "Lists current app names for all available platforms (alias for get)",
				Action:      listAppNames,
			},
		},
	}
}

// validateFlutterProjectForNamer validates that the current directory is a Flutter project
func validateFlutterProjectForNamer() (*flutter.ValidationResult, error) {
	result, err := flutter.ValidateProject(".")
	if err != nil {
		utils.Error("Not a valid Flutter project: %v", err)
		return nil, err
	}

	if !result.IsValid {
		utils.Error("Not a valid Flutter project: %s", result.ErrorMessage)
		return nil, fmt.Errorf("not a valid Flutter project: %s", result.ErrorMessage)
	}

	return result, nil
}

// getAppNames retrieves app names for specified platforms
func getAppNames(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForNamer()
	if err != nil {
		return err
	}

	// Parse platforms
	platforms, err := parsePlatforms(c.StringSlice("platforms"))
	if err != nil {
		return err
	}

	// Get app names
	utils.Info("Retrieving app names...")
	result, err := namer.GetAppNames(project.ProjectPath, platforms)
	if err != nil {
		utils.Error("Failed to get app names: %v", err)
		return err
	}

	// Display results
	displayAppNames(result)
	return nil
}

// setAppNames sets app names for specified platforms
func setAppNames(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForNamer()
	if err != nil {
		return err
	}

	// Build request
	request := &namer.SetAppNameRequest{
		ProjectPath: project.ProjectPath,
		Universal:   c.String("value"),
		Platforms:   make(map[namer.Platform]string),
	}

	// Check for platform-specific names
	platformFlags := map[string]namer.Platform{
		"android": namer.PlatformAndroid,
		"ios":     namer.PlatformIOS,
		"macos":   namer.PlatformMacOS,
		"linux":   namer.PlatformLinux,
		"windows": namer.PlatformWindows,
		"web":     namer.PlatformWeb,
	}

	for flag, platform := range platformFlags {
		if value := c.String(flag); value != "" {
			request.Platforms[platform] = value
		}
	}

	// If platforms flag is specified with universal value, limit to those platforms
	if request.Universal != "" && c.IsSet("platforms") {
		platforms, err := parsePlatforms(c.StringSlice("platforms"))
		if err != nil {
			return err
		}
		
		// Override with limited platforms
		request.Platforms = make(map[namer.Platform]string)
		for _, platform := range platforms {
			request.Platforms[platform] = request.Universal
		}
		request.Universal = "" // Clear universal since we're using platform-specific
	}

	// Validate request
	if request.Universal == "" && len(request.Platforms) == 0 {
		utils.Error("Either --value or platform-specific flags must be provided")
		utils.Info("Usage examples:")
		utils.Info("  fdawg namer set --value \"My App\"")
		utils.Info("  fdawg namer set --android \"Android App\" --ios \"iOS App\"")
		utils.Info("  fdawg namer set --value \"My App\" --platforms android,ios")
		return fmt.Errorf("no app name provided")
	}

	// Set app names
	utils.Info("Setting app names...")
	if err := namer.SetAppNames(request); err != nil {
		utils.Error("Failed to set app names: %v", err)
		return err
	}

	utils.Success("App names updated successfully")
	
	// Show updated names
	utils.Info("Updated app names:")
	result, _ := namer.GetAppNames(project.ProjectPath, nil)
	if result != nil {
		displayAppNames(result)
	}

	return nil
}

// listAppNames lists all current app names (alias for get with no platforms)
func listAppNames(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForNamer()
	if err != nil {
		return err
	}

	// Get all app names
	utils.Info("Listing all app names...")
	result, err := namer.GetAppNames(project.ProjectPath, nil)
	if err != nil {
		utils.Error("Failed to list app names: %v", err)
		return err
	}

	// Display results
	displayAppNames(result)
	return nil
}

// parsePlatforms parses platform strings into Platform types
func parsePlatforms(platformStrs []string) ([]namer.Platform, error) {
	if len(platformStrs) == 0 {
		return nil, nil
	}

	var platforms []namer.Platform
	validPlatforms := map[string]namer.Platform{
		"android": namer.PlatformAndroid,
		"ios":     namer.PlatformIOS,
		"macos":   namer.PlatformMacOS,
		"linux":   namer.PlatformLinux,
		"windows": namer.PlatformWindows,
		"web":     namer.PlatformWeb,
	}

	for _, platformStr := range platformStrs {
		platform, exists := validPlatforms[strings.ToLower(platformStr)]
		if !exists {
			return nil, fmt.Errorf("invalid platform: %s. Valid platforms: android, ios, macos, linux, windows, web", platformStr)
		}
		platforms = append(platforms, platform)
	}

	return platforms, nil
}

// displayAppNames displays app names in a formatted way
func displayAppNames(result *namer.AppNameResult) {
	fmt.Println(utils.Separator("=", 60))
	utils.Success("App Names for Project: %s", result.ProjectPath)
	fmt.Println(utils.Separator("=", 60))

	availableCount := 0
	for _, appName := range result.AppNames {
		fmt.Println(utils.Separator("-", 60))
		
		platformName := strings.Title(string(appName.Platform))
		if appName.Available {
			availableCount++
			utils.Info("%s Platform", platformName)
			
			if appName.Error != "" {
				utils.Warning("Error: %s", appName.Error)
			} else {
				fmt.Printf("  Display Name: %s\n", appName.DisplayName)
				if appName.InternalName != "" && appName.InternalName != appName.DisplayName {
					fmt.Printf("  Internal Name: %s\n", appName.InternalName)
				}
			}
		} else {
			utils.Warning("%s Platform (Not Available)", platformName)
		}
	}

	fmt.Println(utils.Separator("=", 60))
	utils.Success("Available Platforms: %d", availableCount)
}
