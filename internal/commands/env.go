package commands

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/Jerinji2016/fdawg/pkg/environment"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/utils"
	"github.com/urfave/cli/v2"
)

// EnvCommand returns the CLI command for managing environment files
func EnvCommand() *cli.Command {
	return &cli.Command{
		Name:        "env",
		Usage:       "Manage environment files for Flutter projects",
		Description: "Commands for managing environment files in the .environment directory",
		Subcommands: []*cli.Command{
			{
				Name:        "list",
				Usage:       "List all environment files",
				Description: "Lists all environment files in the .environment directory",
				Action:      listEnvFiles,
			},
			{
				Name:        "show",
				Usage:       "Show variables in an environment file",
				Description: "Shows all variables in a specific environment file",
				ArgsUsage:   "<env-name>",
				Action:      showEnvVariables,
			},
			{
				Name:        "create",
				Usage:       "Create a new environment file",
				Description: "Creates a new environment file in the .environment directory",
				ArgsUsage:   "<env-name>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "copy",
						Aliases: []string{"c"},
						Usage:   "Copy variables from an existing environment file",
					},
				},
				Action: createEnvFile,
			},
			{
				Name:        "add",
				Usage:       "Add or update a variable in an environment file",
				Description: "Adds or updates a variable in a specific environment file",
				ArgsUsage:   "<key> <value>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "env",
						Aliases: []string{"e"},
						Usage:   "Environment file to add the variable to",
						Value:   "development",
					},
				},
				Action: addEnvVariable,
			},
			{
				Name:        "delete-env",
				Usage:       "Delete an environment file",
				Description: "Deletes an environment file from the .environment directory",
				ArgsUsage:   "<env-name>",
				Action:      deleteEnvFile,
			},
			{
				Name:        "delete-var",
				Usage:       "Delete a variable from an environment file",
				Description: "Deletes a variable from a specific environment file",
				ArgsUsage:   "<key>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "env",
						Aliases: []string{"e"},
						Usage:   "Environment file to delete the variable from",
						Value:   "development",
					},
				},
				Action: deleteEnvVariable,
			},
		},
	}
}

// validateFlutterProject checks if the current directory is a valid Flutter project
func validateFlutterProject() (*flutter.ValidationResult, error) {
	dir := "."
	utils.Info("Checking Flutter project in: %s", dir)

	result, err := flutter.ValidateProject(dir)
	if err != nil {
		utils.Error(err.Error())
		return nil, err
	}

	if !result.IsValid {
		utils.Error("Not a valid Flutter project")
		return nil, fmt.Errorf("not a valid Flutter project")
	}

	return result, nil
}

// listEnvFiles lists all environment files in the project
func listEnvFiles(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProject()
	if err != nil {
		return err
	}

	// Get environment files
	envFiles, err := environment.ListEnvFiles(project.ProjectPath)
	if err != nil {
		utils.Error("Failed to list environment files: %v", err)
		return err
	}

	if len(envFiles) == 0 {
		utils.Info("No environment files found in %s", environment.GetEnvDir(project.ProjectPath))
		utils.Info("Use 'fdawg env create <env-name>' to create a new environment file")
		return nil
	}

	// Display environment files
	utils.Success("Environment files:")
	for _, envFile := range envFiles {
		varCount := len(envFile.Variables)
		utils.Log("  - %s (%d variable%s)", envFile.Name, varCount, pluralize(varCount))
	}

	return nil
}

// showEnvVariables shows all variables in a specific environment file
func showEnvVariables(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProject()
	if err != nil {
		return err
	}

	// Check if environment name is provided
	if c.Args().Len() == 0 {
		utils.Error("Environment name is required")
		utils.Info("Usage: fdawg env show <env-name>")
		return fmt.Errorf("environment name is required")
	}

	envName := c.Args().First()

	// Get environment file
	envFile, err := environment.GetEnvFile(project.ProjectPath, envName)
	if err != nil {
		utils.Error("Failed to get environment file: %v", err)
		return err
	}

	// Display variables
	utils.Success("Variables in %s environment:", envName)

	if len(envFile.Variables) == 0 {
		utils.Info("No variables found")
		return nil
	}

	// Use tabwriter for aligned output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  KEY\tVALUE")
	fmt.Fprintln(w, "  ---\t-----")

	for key, value := range envFile.Variables {
		fmt.Fprintf(w, "  %s\t%v\n", key, value)
	}
	w.Flush()

	return nil
}

// createEnvFile creates a new environment file
func createEnvFile(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProject()
	if err != nil {
		return err
	}

	// Check if environment name is provided
	if c.Args().Len() == 0 {
		utils.Error("Environment name is required")
		utils.Info("Usage: fdawg env create <env-name>")
		return fmt.Errorf("environment name is required")
	}

	envName := c.Args().First()
	copyFrom := c.String("copy")

	// Create environment file
	if copyFrom != "" {
		// Copy from existing environment file
		utils.Info("Creating %s environment by copying from %s...", envName, copyFrom)

		err := environment.CopyEnvFile(project.ProjectPath, copyFrom, envName)
		if err != nil {
			utils.Error("Failed to create environment file: %v", err)
			return err
		}
	} else {
		// Create empty environment file
		utils.Info("Creating empty %s environment...", envName)

		err := environment.CreateEnvFile(project.ProjectPath, envName, make(map[string]interface{}))
		if err != nil {
			utils.Error("Failed to create environment file: %v", err)
			return err
		}
	}

	utils.Success("Environment file %s created successfully", envName)
	return nil
}

// addEnvVariable adds a variable to an environment file
func addEnvVariable(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProject()
	if err != nil {
		return err
	}

	// Check if key and value are provided
	if c.Args().Len() < 2 {
		utils.Error("Key and value are required")
		utils.Info("Usage: fdawg env add <key> <value> [--env <env-name>]")
		return fmt.Errorf("key and value are required")
	}

	key := c.Args().Get(0)
	valueStr := c.Args().Get(1)
	envName := c.String("env")

	// Parse value (try to convert to appropriate type)
	var value interface{} = valueStr

	// Try to parse as number or boolean
	if strings.EqualFold(valueStr, "true") {
		value = true
	} else if strings.EqualFold(valueStr, "false") {
		value = false
	} else if strings.Contains(valueStr, ".") {
		// Try to parse as float
		if floatVal, err := parseFloat(valueStr); err == nil {
			value = floatVal
		}
	} else {
		// Try to parse as integer
		if intVal, err := parseInt(valueStr); err == nil {
			value = intVal
		}
	}

	// Add variable to environment file
	utils.Info("Adding %s=%v to %s environment...", key, value, envName)

	err = environment.AddVariable(project.ProjectPath, envName, key, value)
	if err != nil {
		utils.Error("Failed to add variable: %v", err)
		return err
	}

	utils.Success("Variable added successfully")
	return nil
}

// deleteEnvFile deletes an environment file
func deleteEnvFile(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProject()
	if err != nil {
		return err
	}

	// Check if environment name is provided
	if c.Args().Len() == 0 {
		utils.Error("Environment name is required")
		utils.Info("Usage: fdawg env delete-env <env-name>")
		return fmt.Errorf("environment name is required")
	}

	envName := c.Args().First()

	// Confirm deletion
	utils.Warning("Are you sure you want to delete the %s environment file? (y/N): ", envName)
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		utils.Info("Deletion cancelled")
		return nil
	}

	// Delete environment file
	utils.Info("Deleting %s environment...", envName)

	err = environment.DeleteEnvFile(project.ProjectPath, envName)
	if err != nil {
		utils.Error("Failed to delete environment file: %v", err)
		return err
	}

	utils.Success("Environment file %s deleted successfully", envName)
	return nil
}

// deleteEnvVariable deletes a variable from an environment file
func deleteEnvVariable(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProject()
	if err != nil {
		return err
	}

	// Check if key is provided
	if c.Args().Len() == 0 {
		utils.Error("Variable key is required")
		utils.Info("Usage: fdawg env delete-var <key> [--env <env-name>]")
		return fmt.Errorf("variable key is required")
	}

	key := c.Args().First()
	envName := c.String("env")

	// Confirm deletion
	utils.Warning("Are you sure you want to delete the variable %s from the %s environment? (y/N): ", key, envName)
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		utils.Info("Deletion cancelled")
		return nil
	}

	// Delete variable
	utils.Info("Deleting variable %s from %s environment...", key, envName)

	err = environment.DeleteVariable(project.ProjectPath, envName, key)
	if err != nil {
		utils.Error("Failed to delete variable: %v", err)
		return err
	}

	utils.Success("Variable %s deleted successfully from %s environment", key, envName)
	return nil
}

// Helper functions

// pluralize returns "s" if count is not 1, otherwise returns empty string
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

// parseFloat tries to parse a string as a float64
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// parseInt tries to parse a string as an int64
func parseInt(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}
