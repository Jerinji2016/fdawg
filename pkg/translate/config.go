package translate

import (
	"fmt"

	"github.com/Jerinji2016/fdawg/pkg/config"
)

// Config holds the configuration for the translation service
type Config struct {
	APIKey  string
	Enabled bool
}

// LoadConfig loads the translation configuration from project config only
func LoadConfig(projectPath string) (*Config, error) {
	// Load from project config only
	if projectPath == "" {
		return &Config{
			APIKey:  "",
			Enabled: false,
		}, nil
	}

	projectConfig, err := config.GetTranslationConfig(projectPath)
	if err != nil {
		// If config doesn't exist or has error, return disabled config
		return &Config{
			APIKey:  "",
			Enabled: false,
		}, nil
	}

	return &Config{
		APIKey:  projectConfig.GoogleTranslateAPIKey,
		Enabled: projectConfig.Enabled,
	}, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Enabled && c.APIKey == "" {
		return fmt.Errorf("Google Translate API key is required when translation is enabled")
	}
	return nil
}

// IsEnabled returns whether translation is enabled
func (c *Config) IsEnabled() bool {
	return c.Enabled && c.APIKey != ""
}
