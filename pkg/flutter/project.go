package flutter

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/Jerinji2016/fdawg/pkg/utils"
    "gopkg.in/yaml.v3"
)

// PubspecInfo represents the structure of a pubspec.yaml file
type PubspecInfo struct {
    Name            string                 `yaml:"name"`
    Description     string                 `yaml:"description"`
    Version         string                 `yaml:"version"`
    Environment     map[string]string      `yaml:"environment"`
    Dependencies    map[string]interface{} `yaml:"dependencies"`
    DevDependencies map[string]interface{} `yaml:"dev_dependencies"`
    Flutter         struct {
        Assets []string               `yaml:"assets"`
        Fonts  []map[string]interface{} `yaml:"fonts"`
    } `yaml:"flutter"`
}

// ValidationResult contains the result of a Flutter project validation
type ValidationResult struct {
    IsValid      bool
    ProjectPath  string
    PubspecInfo  *PubspecInfo
    MissingDirs  []string
    ErrorMessage string
}

// ValidateProject checks if the specified directory is a valid Flutter project
func ValidateProject(dir string) (*ValidationResult, error) {
    result := &ValidationResult{
        IsValid: false,
    }

    // Get absolute path for better error messages
    absPath, err := filepath.Abs(dir)
    if err != nil {
        return result, fmt.Errorf("failed to resolve path: %v", err)
    }
    result.ProjectPath = absPath

    // Check if directory exists
    dirInfo, err := os.Stat(dir)
    if os.IsNotExist(err) {
        result.ErrorMessage = fmt.Sprintf("directory does not exist: %s", absPath)
        return result, fmt.Errorf(result.ErrorMessage)
    }

    // Check if it's actually a directory
    if !dirInfo.IsDir() {
        result.ErrorMessage = fmt.Sprintf("not a directory: %s", absPath)
        return result, fmt.Errorf(result.ErrorMessage)
    }

    // Check for pubspec.yaml (required)
    pubspecPath := filepath.Join(dir, "pubspec.yaml")
    if _, err := os.Stat(pubspecPath); os.IsNotExist(err) {
        result.ErrorMessage = "not a Flutter project: pubspec.yaml not found"
        return result, fmt.Errorf(result.ErrorMessage)
    }

    // Parse pubspec.yaml
    pubspecData, err := os.ReadFile(pubspecPath)
    if err != nil {
        result.ErrorMessage = fmt.Sprintf("failed to read pubspec.yaml: %v", err)
        return result, fmt.Errorf(result.ErrorMessage)
    }

    var pubspec PubspecInfo
    if err := yaml.Unmarshal(pubspecData, &pubspec); err != nil {
        result.ErrorMessage = fmt.Sprintf("failed to parse pubspec.yaml: %v", err)
        return result, fmt.Errorf(result.ErrorMessage)
    }
    result.PubspecInfo = &pubspec

    // Additional checks (optional)
    requiredDirs := []string{"lib", "android", "ios"}
    missingDirs := []string{}

    for _, requiredDir := range requiredDirs {
        checkPath := filepath.Join(dir, requiredDir)
        if _, err := os.Stat(checkPath); os.IsNotExist(err) {
            missingDirs = append(missingDirs, requiredDir)
        }
    }
    result.MissingDirs = missingDirs

    // Project is valid if we got this far
    result.IsValid = true
    return result, nil
}

// DisplayProjectInfo formats and displays the pubspec information
func DisplayProjectInfo(result *ValidationResult) {
    if !result.IsValid || result.PubspecInfo == nil {
        utils.Error("Cannot display information for an invalid project")
        return
    }

    pubspec := result.PubspecInfo
    fmt.Println("\n" + utils.Separator("=", 50))
    utils.Success("Flutter Project Details")
    fmt.Println(utils.Separator("=", 50))

    // Basic info
    fmt.Printf("%-15s: %s\n", "Name", pubspec.Name)
    fmt.Printf("%-15s: %s\n", "Description", pubspec.Description)
    fmt.Printf("%-15s: %s\n", "Version", pubspec.Version)

    // SDK environment
    fmt.Println(utils.Separator("-", 50))
    utils.Info("SDK Environment")
    for env, version := range pubspec.Environment {
        fmt.Printf("%-15s: %s\n", env, version)
    }

    // Dependencies
    fmt.Println(utils.Separator("-", 50))
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
    fmt.Println(utils.Separator("-", 50))
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
    fmt.Println(utils.Separator("-", 50))
    utils.Info("Assets")
    if len(pubspec.Flutter.Assets) == 0 {
        fmt.Println("No assets defined")
    } else {
        for _, asset := range pubspec.Flutter.Assets {
            fmt.Printf("- %s\n", asset)
        }
    }

    // Fonts
    fmt.Println(utils.Separator("-", 50))
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

    fmt.Println(utils.Separator("=", 50))
}