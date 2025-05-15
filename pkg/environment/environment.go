package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/Jerinji2016/fdawg/pkg/utils"
)

const (
	// EnvDirName is the name of the directory where environment files are stored
	EnvDirName = ".environment"
)

// EnvFile represents an environment file
type EnvFile struct {
	Name      string                 `json:"-"`
	Path      string                 `json:"-"`
	Variables map[string]interface{} `json:"-"`
}

// EnvVariable represents a key-value pair in an environment file
type EnvVariable struct {
	Key   string
	Value interface{}
}

// GetEnvDir returns the path to the environment directory for a Flutter project
func GetEnvDir(projectPath string) string {
	return filepath.Join(projectPath, EnvDirName)
}

// EnsureEnvDirExists creates the environment directory if it doesn't exist
func EnsureEnvDirExists(projectPath string) error {
	envDir := GetEnvDir(projectPath)
	return utils.EnsureDirExists(envDir)
}

// ListEnvFiles returns a list of all environment files in the project
func ListEnvFiles(projectPath string) ([]EnvFile, error) {
	envDir := GetEnvDir(projectPath)

	// Check if environment directory exists
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if err := EnsureEnvDirExists(projectPath); err != nil {
			return nil, fmt.Errorf("failed to create environment directory: %v", err)
		}
		return []EnvFile{}, nil
	}

	// Read all files in the environment directory
	files, err := os.ReadDir(envDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment directory: %v", err)
	}

	var envFiles []EnvFile
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Check if the file is a JSON file
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(envDir, file.Name())
		envName := strings.TrimSuffix(file.Name(), ".json")

		// Read the file content
		variables, err := readEnvFile(filePath)
		if err != nil {
			utils.Warning("Failed to read environment file %s: %v", file.Name(), err)
			continue
		}

		envFiles = append(envFiles, EnvFile{
			Name:      envName,
			Path:      filePath,
			Variables: variables,
		})
	}

	// Sort environment files by name
	sort.Slice(envFiles, func(i, j int) bool {
		return envFiles[i].Name < envFiles[j].Name
	})

	return envFiles, nil
}

// GetEnvFile returns a specific environment file by name
func GetEnvFile(projectPath, envName string) (*EnvFile, error) {
	envDir := GetEnvDir(projectPath)
	filePath := filepath.Join(envDir, envName+".json")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("environment file %s does not exist", envName)
	}

	// Read the file content
	variables, err := readEnvFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment file: %v", err)
	}

	return &EnvFile{
		Name:      envName,
		Path:      filePath,
		Variables: variables,
	}, nil
}

// CreateEnvFile creates a new environment file
func CreateEnvFile(projectPath, envName string, variables map[string]interface{}) error {
	// Ensure the environment directory exists
	if err := EnsureEnvDirExists(projectPath); err != nil {
		return fmt.Errorf("failed to create environment directory: %v", err)
	}

	envDir := GetEnvDir(projectPath)
	filePath := filepath.Join(envDir, envName+".json")

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("environment file %s already exists", envName)
	}

	// Create the file
	if err := writeEnvFile(filePath, variables); err != nil {
		return err
	}

	// Generate the Dart environment file
	if err := GenerateDartEnvironmentFile(projectPath); err != nil {
		utils.Warning("Failed to generate Dart environment file: %v", err)
	}

	return nil
}

// CopyEnvFile copies an existing environment file to a new one
func CopyEnvFile(projectPath, sourceEnvName, targetEnvName string) error {
	// Get the source environment file
	sourceEnv, err := GetEnvFile(projectPath, sourceEnvName)
	if err != nil {
		return fmt.Errorf("failed to get source environment file: %v", err)
	}

	// Create a new environment file with the same variables
	return CreateEnvFile(projectPath, targetEnvName, sourceEnv.Variables)
}

// AddVariable adds or updates a variable in an environment file
func AddVariable(projectPath, envName, key string, value interface{}) error {
	// Get the environment file
	envFile, err := GetEnvFile(projectPath, envName)
	if err != nil {
		return fmt.Errorf("failed to get environment file: %v", err)
	}

	// Add or update the variable
	envFile.Variables[key] = value

	// Write the updated variables back to the file
	if err := writeEnvFile(envFile.Path, envFile.Variables); err != nil {
		return err
	}

	// Generate the Dart environment file
	if err := GenerateDartEnvironmentFile(projectPath); err != nil {
		utils.Warning("Failed to generate Dart environment file: %v", err)
	}

	return nil
}

// DeleteVariable deletes a variable from an environment file
func DeleteVariable(projectPath, envName, key string) error {
	// Get the environment file
	envFile, err := GetEnvFile(projectPath, envName)
	if err != nil {
		return fmt.Errorf("failed to get environment file: %v", err)
	}

	// Check if the variable exists
	if _, exists := envFile.Variables[key]; !exists {
		return fmt.Errorf("variable %s does not exist in environment file %s", key, envName)
	}

	// Delete the variable
	delete(envFile.Variables, key)

	// Write the updated variables back to the file
	if err := writeEnvFile(envFile.Path, envFile.Variables); err != nil {
		return err
	}

	// Generate the Dart environment file
	if err := GenerateDartEnvironmentFile(projectPath); err != nil {
		utils.Warning("Failed to generate Dart environment file: %v", err)
	}

	return nil
}

// readEnvFile reads and parses an environment file
func readEnvFile(filePath string) (map[string]interface{}, error) {
	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Parse the JSON content
	var variables map[string]interface{}
	if err := json.Unmarshal(data, &variables); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return variables, nil
}

// DeleteEnvFile deletes an environment file
func DeleteEnvFile(projectPath, envName string) error {
	envDir := GetEnvDir(projectPath)
	filePath := filepath.Join(envDir, envName+".json")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("environment file %s does not exist", envName)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete environment file: %v", err)
	}

	// Generate the Dart environment file
	if err := GenerateDartEnvironmentFile(projectPath); err != nil {
		utils.Warning("Failed to generate Dart environment file: %v", err)
	}

	return nil
}

// writeEnvFile writes variables to an environment file
func writeEnvFile(filePath string, variables map[string]interface{}) error {
	// Marshal the variables to JSON
	data, err := json.MarshalIndent(variables, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write the JSON to the file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// GenerateDartEnvironmentFile generates a Dart environment file with all environment variables
func GenerateDartEnvironmentFile(projectPath string) error {
	// Get all environment files
	envFiles, err := ListEnvFiles(projectPath)
	if err != nil {
		return fmt.Errorf("failed to list environment files: %v", err)
	}

	if len(envFiles) == 0 {
		return fmt.Errorf("no environment files found")
	}

	// Create a map to store all unique variables
	allVariables := make(map[string]struct{})

	// Collect all unique variable names from all environment files
	for _, envFile := range envFiles {
		for key := range envFile.Variables {
			allVariables[key] = struct{}{}
		}
	}

	// Sort the variable names for consistent output
	var sortedVars []string
	for key := range allVariables {
		sortedVars = append(sortedVars, key)
	}
	sort.Strings(sortedVars)

	// Create the Dart file content
	var content strings.Builder

	// Add file header
	content.WriteString(`// GENERATED CODE - DO NOT MODIFY BY HAND
// Generated by fdawg

import 'dart:core';

/// Environment configuration
///
/// This class provides access to environment variables defined in the .environment directory.
/// It is automatically generated by fdawg and should not be modified manually.
class Environment {
  // Private constructor to prevent instantiation
  Environment._();

`)

	// Add static constants for each variable
	for _, key := range sortedVars {
		// Determine the type of the variable based on the first environment file that has it
		var varType string
		var defaultValue string

		for _, envFile := range envFiles {
			if value, exists := envFile.Variables[key]; exists {
				switch v := value.(type) {
				case bool:
					varType = "bool"
					if v {
						defaultValue = "true"
					} else {
						defaultValue = "false"
					}
				case float64:
					// Check if it's an integer
					if v == float64(int(v)) {
						varType = "int"
						defaultValue = fmt.Sprintf("%d", int(v))
					} else {
						varType = "double"
						defaultValue = fmt.Sprintf("%g", v)
					}
				case int:
					varType = "int"
					defaultValue = fmt.Sprintf("%d", v)
				case int64:
					varType = "int"
					defaultValue = fmt.Sprintf("%d", v)
				default:
					varType = "String"
					defaultValue = fmt.Sprintf("\"%v\"", v)
				}
				break
			}
		}

		// Add the static constant
		content.WriteString(fmt.Sprintf("  /// %s environment variable\n", key))
		content.WriteString(fmt.Sprintf("  static const %s %s = %s.fromEnvironment('%s', defaultValue: %s);\n\n",
			varType, formatDartVariableName(key), varType, key, defaultValue))
	}

	// Close the class
	content.WriteString("}\n")

	// Write the file
	dartFilePath := filepath.Join(projectPath, "lib", "config")

	// Ensure the directory exists
	if err := os.MkdirAll(dartFilePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	dartFilePath = filepath.Join(dartFilePath, "environment.dart")

	if err := os.WriteFile(dartFilePath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("failed to write Dart file: %v", err)
	}

	return nil
}

// formatDartVariableName converts an environment variable name to a Dart variable name
// For example, API_URL becomes apiUrl
func formatDartVariableName(name string) string {
	parts := strings.Split(name, "_")
	for i := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(parts[i])
		} else {
			// Capitalize the first letter of each word
			if len(parts[i]) > 0 {
				r := []rune(strings.ToLower(parts[i]))
				r[0] = unicode.ToUpper(r[0])
				parts[i] = string(r)
			}
		}
	}
	return strings.Join(parts, "")
}
