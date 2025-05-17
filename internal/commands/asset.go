package commands

import (
	"fmt"
	"path/filepath"

	"github.com/Jerinji2016/fdawg/pkg/asset"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/utils"
	"github.com/urfave/cli/v2"
)

// AssetCommand returns the CLI command for managing assets
func AssetCommand() *cli.Command {
	return &cli.Command{
		Name:        "asset",
		Usage:       "Manage assets for Flutter projects",
		Description: "Commands for managing assets in the assets directory",
		Subcommands: []*cli.Command{
			{
				Name:        "add",
				Usage:       "Add an asset to the project",
				Description: "Adds an asset to the project and updates the pubspec.yaml file",
				ArgsUsage:   "<asset-path>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "Type of asset (images, animations, audio, videos, json, svgs, misc)",
					},
				},
				Action: addAsset,
			},
			{
				Name:        "remove",
				Usage:       "Remove an asset from the project",
				Description: "Removes an asset from the project and updates the pubspec.yaml file",
				ArgsUsage:   "<asset-name>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "type",
						Aliases: []string{"t"},
						Usage:   "Type of asset (images, animations, audio, videos, json, svgs, misc)",
					},
				},
				Action: removeAsset,
			},
			{
				Name:        "list",
				Usage:       "List all assets in the project",
				Description: "Lists all assets in the project by type",
				Action:      listAssets,
			},
			{
				Name:        "generate-dart",
				Usage:       "Generate Dart asset file",
				Description: "Generates a Dart asset file with all assets",
				Action:      generateDartAssetFile,
			},
			{
				Name:        "migrate",
				Usage:       "Migrate assets to organized folders by type",
				Description: "Migrates assets to organized folders by type and cleans up empty directories",
				Action:      migrateAssets,
			},
		},
	}
}

// validateFlutterProjectForAsset validates that the current directory is a Flutter project
// This is a separate function to avoid name conflicts with other commands
func validateFlutterProjectForAsset() (*flutter.ValidationResult, error) {
	// Check if the current directory is a Flutter project
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

// addAsset adds an asset to the project
func addAsset(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForAsset()
	if err != nil {
		return err
	}

	// Check if asset path is provided
	if c.Args().Len() == 0 {
		utils.Error("Asset path is required")
		utils.Info("Usage: fdawg asset add <asset-path> [--type <asset-type>]")
		return fmt.Errorf("asset path is required")
	}

	assetPath := c.Args().First()
	assetTypeStr := c.String("type")

	var assetType asset.AssetType
	if assetTypeStr != "" {
		assetType = asset.AssetType(assetTypeStr)
	} else {
		// Determine asset type from file extension
		assetType = asset.DetermineAssetType(assetPath)
	}

	// Add the asset
	utils.Info("Adding asset %s as type %s...", assetPath, assetType)

	err = asset.AddAsset(project.ProjectPath, assetPath, assetType)
	if err != nil {
		utils.Error("Failed to add asset: %v", err)
		return err
	}

	utils.Success("Asset added successfully")

	// Generate the Dart asset file
	utils.Info("Generating Dart asset file...")

	err = asset.GenerateDartAssetFile(project.ProjectPath)
	if err != nil {
		utils.Warning("Failed to generate Dart asset file: %v", err)
	} else {
		dartFilePath := filepath.Join(project.ProjectPath, "lib", "config", "asset.dart")
		utils.Success("Dart asset file generated successfully at %s", dartFilePath)
	}

	return nil
}

// removeAsset removes an asset from the project
func removeAsset(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForAsset()
	if err != nil {
		return err
	}

	// Check if asset name is provided
	if c.Args().Len() == 0 {
		utils.Error("Asset name is required")
		utils.Info("Usage: fdawg asset remove <asset-name> [--type <asset-type>]")
		return fmt.Errorf("asset name is required")
	}

	assetName := c.Args().First()
	assetTypeStr := c.String("type")

	var assetType asset.AssetType
	if assetTypeStr != "" {
		assetType = asset.AssetType(assetTypeStr)
	}

	// Remove the asset
	utils.Info("Removing asset %s...", assetName)

	err = asset.RemoveAsset(project.ProjectPath, assetName, assetType)
	if err != nil {
		utils.Error("Failed to remove asset: %v", err)
		return err
	}

	utils.Success("Asset removed successfully")

	// Generate the Dart asset file
	utils.Info("Generating Dart asset file...")

	err = asset.GenerateDartAssetFile(project.ProjectPath)
	if err != nil {
		utils.Warning("Failed to generate Dart asset file: %v", err)
	} else {
		dartFilePath := filepath.Join(project.ProjectPath, "lib", "config", "asset.dart")
		utils.Success("Dart asset file generated successfully at %s", dartFilePath)
	}

	return nil
}

// listAssets lists all assets in the project
func listAssets(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForAsset()
	if err != nil {
		return err
	}

	// List assets
	utils.Info("Listing assets...")

	assets, err := asset.ListAssets(project.ProjectPath)
	if err != nil {
		utils.Error("Failed to list assets: %v", err)
		return err
	}

	// Display assets by type
	fmt.Println(utils.Separator("=", 50))
	utils.Success("Project Assets")
	fmt.Println(utils.Separator("=", 50))

	totalAssets := 0
	for assetType, assetFiles := range assets {
		fmt.Println(utils.Separator("-", 50))
		utils.Info("%s Assets", assetType)

		if len(assetFiles) == 0 {
			fmt.Println("No assets found")
		} else {
			for _, assetFile := range assetFiles {
				fmt.Printf("- %s\n", assetFile)
			}
			totalAssets += len(assetFiles)
		}
	}

	fmt.Println(utils.Separator("=", 50))
	utils.Success("Total Assets: %d", totalAssets)

	return nil
}

// generateDartAssetFile generates a Dart asset file with all assets
func generateDartAssetFile(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForAsset()
	if err != nil {
		return err
	}

	// Generate the Dart asset file
	utils.Info("Generating Dart asset file...")

	err = asset.GenerateDartAssetFile(project.ProjectPath)
	if err != nil {
		utils.Error("Failed to generate Dart asset file: %v", err)
		return err
	}

	dartFilePath := filepath.Join(project.ProjectPath, "lib", "config", "asset.dart")
	utils.Success("Dart asset file generated successfully at %s", dartFilePath)
	return nil
}

// migrateAssets migrates assets from a flat structure to organized folders
func migrateAssets(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForAsset()
	if err != nil {
		return err
	}

	// Migrate assets
	utils.Info("Migrating assets to organized folders...")

	err = asset.MigrateAssets(project.ProjectPath)
	if err != nil {
		utils.Error("Failed to migrate assets: %v", err)
		return err
	}

	utils.Success("Assets migrated successfully")

	// Generate the Dart asset file
	utils.Info("Generating Dart asset file...")

	err = asset.GenerateDartAssetFile(project.ProjectPath)
	if err != nil {
		utils.Warning("Failed to generate Dart asset file: %v", err)
	} else {
		dartFilePath := filepath.Join(project.ProjectPath, "lib", "config", "asset.dart")
		utils.Success("Dart asset file generated successfully at %s", dartFilePath)
	}

	return nil
}
