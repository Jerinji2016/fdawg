package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FdawgConfig represents the fdawg configuration
type FdawgConfig struct {
	Translation TranslationConfig `json:"translation"`
}

// TranslationConfig holds translation-related configuration
type TranslationConfig struct {
	GoogleTranslateAPIKey string `json:"google_translate_api_key"`
	Enabled               bool   `json:"enabled"`
}

// GetConfigPath returns the path to the .fdawg-config file
func GetConfigPath(projectPath string) string {
	return filepath.Join(projectPath, ".fdawg-config")
}

// LoadConfig loads the fdawg configuration from the .fdawg-config file
func LoadConfig(projectPath string) (*FdawgConfig, error) {
	configPath := GetConfigPath(projectPath)
	
	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &FdawgConfig{
			Translation: TranslationConfig{
				GoogleTranslateAPIKey: "",
				Enabled:               false,
			},
		}, nil
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the JSON
	var config FdawgConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Update enabled status based on API key
	config.Translation.Enabled = config.Translation.GoogleTranslateAPIKey != ""

	return &config, nil
}

// SaveConfig saves the fdawg configuration to the .fdawg-config file
func SaveConfig(projectPath string, config *FdawgConfig) error {
	configPath := GetConfigPath(projectPath)

	// Update enabled status based on API key
	config.Translation.Enabled = config.Translation.GoogleTranslateAPIKey != ""

	// Convert to JSON with proper formatting
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// UpdateTranslationAPIKey updates the Google Translate API key
func UpdateTranslationAPIKey(projectPath, apiKey string) error {
	config, err := LoadConfig(projectPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	config.Translation.GoogleTranslateAPIKey = apiKey
	config.Translation.Enabled = apiKey != ""

	return SaveConfig(projectPath, config)
}

// GetTranslationConfig returns the translation configuration
func GetTranslationConfig(projectPath string) (*TranslationConfig, error) {
	config, err := LoadConfig(projectPath)
	if err != nil {
		return nil, err
	}

	return &config.Translation, nil
}

// IsConfigFileExists checks if the .fdawg-config file exists
func IsConfigFileExists(projectPath string) bool {
	configPath := GetConfigPath(projectPath)
	_, err := os.Stat(configPath)
	return err == nil
}
