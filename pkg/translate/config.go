package translate

import (
	"fmt"
	"os"

	"github.com/Jerinji2016/fdawg/pkg/config"
)

// Config holds the configuration for the translation service
type Config struct {
	APIKey  string
	Enabled bool
}

// LoadConfig loads the translation configuration from project config or environment variables
func LoadConfig(projectPath string) (*Config, error) {
	// First try to load from project config
	if projectPath != "" {
		projectConfig, err := config.GetTranslationConfig(projectPath)
		if err == nil && projectConfig.GoogleTranslateAPIKey != "" {
			return &Config{
				APIKey:  projectConfig.GoogleTranslateAPIKey,
				Enabled: projectConfig.Enabled,
			}, nil
		}
	}

	// Fallback to environment variable
	apiKey := os.Getenv("GOOGLE_TRANSLATE_API_KEY")
	enabled := apiKey != ""

	return &Config{
		APIKey:  apiKey,
		Enabled: enabled,
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
