package commands

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/Jerinji2016/fdawg/pkg/utils"
    "github.com/urfave/cli/v2"
)

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
        
        // Print some project info
        if pubspecData, err := os.ReadFile(pubspecPath); err == nil {
            utils.Info("Project information:")
            // Safely print a portion of the pubspec
            dataLen := len(pubspecData)
            previewLen := 100
            if dataLen < previewLen {
                previewLen = dataLen
            }
            fmt.Printf("%s\n", string(pubspecData[:previewLen]))
            if dataLen > previewLen {
                fmt.Println("...")
            }
        }
    }

    return nil
}
