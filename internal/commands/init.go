package commands

import (
    "fmt"

    "github.com/Jerinji2016/fdawg/pkg/flutter"
    "github.com/Jerinji2016/fdawg/pkg/utils"
    "github.com/urfave/cli/v2"
)

// InitCommand returns the CLI command for initializing or checking a Flutter project
func InitCommand() *cli.Command {
    return &cli.Command{
        Name:        "init",
        Usage:       "Check if the specified directory is a Flutter project",
        Description: "If no directory is provided, the current directory will be checked",
        ArgsUsage:   "[directory]",
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
    utils.Info("Checking Flutter project in: %s", dir)

    result, err := flutter.ValidateProject(dir)
    if err != nil {
        utils.Error(err.Error())
        return err
    }

    if len(result.MissingDirs) > 0 {
        utils.Warning("Some Flutter directories are missing:")
        for _, missingDir := range result.MissingDirs {
            fmt.Printf("  - %s\n", missingDir)
        }
        utils.Warning("This might be a partial or incomplete Flutter project.")
    } else {
        utils.Success("✓ Valid Flutter project detected!")
    }

    flutter.DisplayProjectInfo(result)
    return nil
}
