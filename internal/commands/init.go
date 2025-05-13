package commands

import (
    "fmt"
    "os"

    "github.com/urfave/cli/v2"
)

// InitCommand returns the CLI command for initializing or checking a Flutter project
func InitCommand() *cli.Command {
    return &cli.Command{
        Name:  "init",
        Usage: "Check if the current directory is a Flutter project",
        Action: func(c *cli.Context) error {
            return checkFlutterProject()
        },
    }
}

// checkFlutterProject verifies if the current directory contains Flutter project files
func checkFlutterProject() error {
    // Check for pubspec.yaml (required)
    if _, err := os.Stat("pubspec.yaml"); os.IsNotExist(err) {
        return fmt.Errorf("not a Flutter project: pubspec.yaml not found")
    }

    // Additional checks (optional)
    requiredDirs := []string{"lib", "android", "ios"}
    missingDirs := []string{}

    for _, dir := range requiredDirs {
        if _, err := os.Stat(dir); os.IsNotExist(err) {
            missingDirs = append(missingDirs, dir)
        }
    }

    if len(missingDirs) > 0 {
        fmt.Println("Warning: Some Flutter directories are missing:")
        for _, dir := range missingDirs {
            fmt.Printf("  - %s\n", dir)
        }
        fmt.Println("This might be a partial or incomplete Flutter project.")
    } else {
        fmt.Println("✓ Valid Flutter project detected!")
        
        // Print some project info
        if pubspecData, err := os.ReadFile("pubspec.yaml"); err == nil {
            fmt.Println("\nProject information:")
            fmt.Printf("%s\n", string(pubspecData[:100])) // Print first 100 chars of pubspec
            fmt.Println("...")
        }
    }

    return nil
}