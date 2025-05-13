package commands

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/Jerinji2016/fdawg/pkg/utils"
    "github.com/urfave/cli/v2"
    "gopkg.in/yaml.v3"
)

// PubspecInfo represents the structure of a pubspec.yaml file
type PubspecInfo struct {
    Name        string                 `yaml:"name"`
    Description string                 `yaml:"description"`
    Version     string                 `yaml:"version"`
    Environment map[string]string      `yaml:"environment"`
    Dependencies map[string]interface{} `yaml:"dependencies"`
    DevDependencies map[string]interface{} `yaml:"dev_dependencies"`
    Flutter     struct {
        Assets []string               `yaml:"assets"`
        Fonts  []map[string]interface{} `yaml:"fonts"`
    } `yaml:"flutter"`
}

// InitCommand returns the CLI command for initializing or checking a Flutter project
func InitCommand() *cli.Command {
    return &cli.Command{
        Name:  "init",
        Usage: "Check if the specified directory is a Flutter project",
        Description: "If no directory is provided, the current directory will be checked",
        ArgsUsage: "[directory]",
        Action: func(c *cli.Context) error {
            // Get directory from arguments or use current directory
            dir := "."
            if c.Args().Len() > 0 {
                dir = c.Args().Get(0)
            }
            
            return checkFlutterProject(dir)
        },
    }
}

// checkFlutterProject verifies if the specified directory contains Flutter project files
func checkFlutterProject(dir string) error {
    // Get absolute path for better error messages
    absPath, err := filepath.Abs(dir)
    if err != nil {
        utils.Error("Failed to resolve path: %v", err)
        return err
    }
    
    // Check if directory exists
    dirInfo, err := os.Stat(dir)
    if os.IsNotExist(err) {
        utils.Error("Directory does not exist: %s", absPath)
        return fmt.Errorf("directory does not exist: %s", absPath)
    }
    
    // Check if it's actually a directory
    if !dirInfo.IsDir() {
        utils.Error("Not a directory: %s", absPath)
        return fmt.Errorf("not a directory: %s", absPath)
    }
    
    utils.Info("Checking Flutter project in: %s", absPath)
    
    // Check for pubspec.yaml (required)
    pubspecPath := filepath.Join(dir, "pubspec.yaml")
    if _, err := os.Stat(pubspecPath); os.IsNotExist(err) {
        utils.Error("Not a Flutter project: pubspec.yaml not found")
        return fmt.Errorf("not a Flutter project: pubspec.yaml not found")
    }

    // Additional checks (optional)
    requiredDirs := []string{"lib", "android", "ios"}
    missingDirs := []string{}

    for _, requiredDir := range requiredDirs {
        checkPath := filepath.Join(dir, requiredDir)
        if _, err := os.Stat(checkPath); os.IsNotExist(err) {
            missingDirs = append(missingDirs, requiredDir)
        }
    }

    if len(missingDirs) > 0 {
        utils.Warning("Some Flutter directories are missing:")
        for _, missingDir := range missingDirs {
            fmt.Printf("  - %s\n", missingDir)
        }
        utils.Warning("This might be a partial or incomplete Flutter project.")
    } else {
        utils.Success("✓ Valid Flutter project detected!")
    }
    
    // Parse and display pubspec information
    pubspecData, err := os.ReadFile(pubspecPath)
    if err != nil {
        utils.Error("Failed to read pubspec.yaml: %v", err)
        return err
    }
    
    var pubspec PubspecInfo
    if err := yaml.Unmarshal(pubspecData, &pubspec); err != nil {
        utils.Error("Failed to parse pubspec.yaml: %v", err)
        return err
    }
    
    displayPubspecInfo(pubspec)
    
    return nil
}

// displayPubspecInfo formats and displays the pubspec information
func displayPubspecInfo(pubspec PubspecInfo) {
    fmt.Println("\n" + strings.Repeat("=", 50))
    utils.Success("Flutter Project Details")
    fmt.Println(strings.Repeat("=", 50))
    
    // Basic info
    fmt.Printf("%-15s: %s\n", "Name", pubspec.Name)
    fmt.Printf("%-15s: %s\n", "Description", pubspec.Description)
    fmt.Printf("%-15s: %s\n", "Version", pubspec.Version)
    
    // SDK environment
    fmt.Println(strings.Repeat("-", 50))
    utils.Info("SDK Environment")
    for env, version := range pubspec.Environment {
        fmt.Printf("%-15s: %s\n", env, version)
    }
    
    // Dependencies
    fmt.Println(strings.Repeat("-", 50))
    utils.Info("Dependencies")
    if len(pubspec.Dependencies) == 0 {
        fmt.Println("No dependencies found")
    } else {
        for dep, version := range pubspec.Dependencies {
            switch v := version.(type) {
            case string:
                fmt.Printf("%-20s: %s\n", dep, v)
            case map[string]interface{}:
                fmt.Printf("%-20s: (complex dependency)\n", dep)
            default:
                fmt.Printf("%-20s: %v\n", dep, v)
            }
        }
    }
    
    // Dev Dependencies
    fmt.Println(strings.Repeat("-", 50))
    utils.Info("Dev Dependencies")
    if len(pubspec.DevDependencies) == 0 {
        fmt.Println("No dev dependencies found")
    } else {
        for dep, version := range pubspec.DevDependencies {
            switch v := version.(type) {
            case string:
                fmt.Printf("%-20s: %s\n", dep, v)
            case map[string]interface{}:
                fmt.Printf("%-20s: (complex dependency)\n", dep)
            default:
                fmt.Printf("%-20s: %v\n", dep, v)
            }
        }
    }
    
    // Assets
    fmt.Println(strings.Repeat("-", 50))
    utils.Info("Assets")
    if len(pubspec.Flutter.Assets) == 0 {
        fmt.Println("No assets defined")
    } else {
        for _, asset := range pubspec.Flutter.Assets {
            fmt.Printf("- %s\n", asset)
        }
    }
    
    // Fonts
    fmt.Println(strings.Repeat("-", 50))
    utils.Info("Fonts")
    if len(pubspec.Flutter.Fonts) == 0 {
        fmt.Println("No custom fonts defined")
    } else {
        for _, font := range pubspec.Flutter.Fonts {
            family, _ := font["family"].(string)
            fmt.Printf("Family: %s\n", family)
            
            if fonts, ok := font["fonts"].([]interface{}); ok {
                for _, f := range fonts {
                    if fontMap, ok := f.(map[string]interface{}); ok {
                        if asset, ok := fontMap["asset"].(string); ok {
                            fmt.Printf("  - Asset: %s\n", asset)
                        }
                        if style, ok := fontMap["style"].(string); ok {
                            fmt.Printf("    Style: %s\n", style)
                        }
                        if weight, ok := fontMap["weight"].(int); ok {
                            fmt.Printf("    Weight: %d\n", weight)
                        }
                    }
                }
            }
        }
    }
    
    fmt.Println(strings.Repeat("=", 50))
}
