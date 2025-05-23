package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/localization"
	"github.com/Jerinji2016/fdawg/pkg/utils"
	"github.com/urfave/cli/v2"
)

// LangCommand returns the CLI command for managing localizations
func LangCommand() *cli.Command {
	return &cli.Command{
		Name:        "lang",
		Usage:       "Manage localizations for Flutter projects",
		Description: "Commands for managing localizations using easy_localization",
		Subcommands: []*cli.Command{
			{
				Name:        "init",
				Usage:       "Initialize localization in the project",
				Description: "Initializes localization using easy_localization package",
				Action:      initLocalization,
			},
			{
				Name:        "add",
				Usage:       "Add a new language",
				Description: "Adds support for a new language",
				ArgsUsage:   "<language-code>",
				Action:      addLanguage,
			},
			{
				Name:        "remove",
				Usage:       "Remove a language",
				Description: "Removes support for a language",
				ArgsUsage:   "<language-code>",
				Action:      removeLanguage,
			},
			{
				Name:        "insert",
				Usage:       "Add a new translation key",
				Description: "Adds a new translation key to all languages",
				ArgsUsage:   "<key>",
				Action:      insertTranslationKey,
			},
			{
				Name:        "delete",
				Usage:       "Delete a translation key",
				Description: "Deletes a translation key from all languages",
				ArgsUsage:   "<key>",
				Action:      deleteTranslationKey,
			},
			{
				Name:        "list",
				Usage:       "List supported languages",
				Description: "Lists all supported languages in the project",
				Action:      listLanguages,
			},
		},
	}
}

// validateFlutterProjectForLocalization validates that the current directory is a Flutter project
func validateFlutterProjectForLocalization() (*flutter.ValidationResult, error) {
	// Check if the current directory is a Flutter project
	result, err := flutter.ValidateProject(".")
	if err != nil {
		utils.Error("Not a valid Flutter project: %v", err)
		return nil, err
	}

	if !result.IsValid {
		utils.Error("Not a valid Flutter project: %s", result.ErrorMessage)
		return nil, fmt.Errorf("not a valid Flutter project: %s", result.ErrorMessage)
	}

	return result, nil
}

// initLocalization initializes localization in the project
func initLocalization(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForLocalization()
	if err != nil {
		return err
	}

	// Initialize localization
	utils.Info("Initializing localization using easy_localization...")

	err = localization.InitLocalization(project.ProjectPath)
	if err != nil {
		utils.Error("Failed to initialize localization: %v", err)
		return err
	}

	utils.Success("Localization initialized successfully")
	utils.Info("Default language: en_US")
	utils.Info("Translations directory: %s", localization.TranslationsDir)
	utils.Info("Run 'flutter pub get' to install the easy_localization package")

	return nil
}

// addLanguage adds a new language to the project
func addLanguage(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForLocalization()
	if err != nil {
		return err
	}

	// Check if language code is provided
	if c.Args().Len() == 0 {
		utils.Error("Language code is required")
		utils.Info("Usage: fdawg lang add <language-code>")
		utils.Info("Example: fdawg lang add es (Spanish)")
		utils.Info("Example: fdawg lang add fr_FR (French, France)")
		return fmt.Errorf("language code is required")
	}

	langCode := c.Args().First()

	// Validate language code
	if !localization.IsValidLanguageCode(langCode) {
		utils.Error("Invalid language code: %s", langCode)
		utils.Info("Use a valid ISO language code (e.g., 'en', 'es', 'fr')")
		utils.Info("For country-specific variants, use format like 'en_US', 'pt_BR'")
		return fmt.Errorf("invalid language code: %s", langCode)
	}

	// Get language info
	langInfo, err := localization.GetLanguageInfo(langCode)
	if err != nil {
		utils.Error("Invalid language code: %v", err)
		return err
	}

	// Add the language
	utils.Info("Adding support for %s (%s)...", langInfo.DisplayName(), langInfo.String())

	err = localization.AddLanguage(project.ProjectPath, langCode)
	if err != nil {
		utils.Error("Failed to add language: %v", err)
		return err
	}

	utils.Success("Language %s (%s) added successfully", langInfo.DisplayName(), langInfo.String())
	return nil
}

// removeLanguage removes a language from the project
func removeLanguage(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForLocalization()
	if err != nil {
		return err
	}

	// Check if language code is provided
	if c.Args().Len() == 0 {
		utils.Error("Language code is required")
		utils.Info("Usage: fdawg lang remove <language-code>")
		utils.Info("Example: fdawg lang remove es")
		return fmt.Errorf("language code is required")
	}

	langCode := c.Args().First()

	// Validate language code
	if !localization.IsValidLanguageCode(langCode) {
		utils.Error("Invalid language code: %s", langCode)
		utils.Info("Use a valid ISO language code (e.g., 'en', 'es', 'fr')")
		return fmt.Errorf("invalid language code: %s", langCode)
	}

	// Get language info
	langInfo, err := localization.GetLanguageInfo(langCode)
	if err != nil {
		utils.Error("Invalid language code: %v", err)
		return err
	}

	// Confirm deletion
	utils.Warning("Are you sure you want to remove support for %s (%s)? (y/N): ", langInfo.DisplayName(), langInfo.String())
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		utils.Info("Removal cancelled")
		return nil
	}

	// Remove the language
	utils.Info("Removing support for %s (%s)...", langInfo.DisplayName(), langInfo.String())

	err = localization.RemoveLanguage(project.ProjectPath, langCode)
	if err != nil {
		utils.Error("Failed to remove language: %v", err)
		return err
	}

	utils.Success("Language %s (%s) removed successfully", langInfo.DisplayName(), langInfo.String())
	return nil
}

// insertTranslationKey adds a new translation key to all languages
func insertTranslationKey(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForLocalization()
	if err != nil {
		return err
	}

	// Check if key is provided
	if c.Args().Len() == 0 {
		utils.Error("Translation key is required")
		utils.Info("Usage: fdawg lang insert <key>")
		utils.Info("Example: fdawg lang insert app.welcome")
		return fmt.Errorf("translation key is required")
	}

	key := c.Args().First()

	// Get all translation files
	translationFiles, err := localization.ListTranslationFiles(project.ProjectPath)
	if err != nil {
		utils.Error("Failed to list translation files: %v", err)
		return err
	}

	if len(translationFiles) == 0 {
		utils.Error("No translation files found")
		utils.Info("Run 'fdawg lang init' to initialize localization first")
		return fmt.Errorf("no translation files found")
	}

	// Collect values for each language
	values := make(map[string]string)
	reader := bufio.NewReader(os.Stdin)

	utils.Info("Enter translation values for each language:")
	for _, file := range translationFiles {
		fmt.Printf("%s (%s): ", file.Language, key)
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)
		values[file.Language] = value
	}

	// Add the key to all languages
	utils.Info("Adding translation key %s to all languages...", key)

	err = localization.InsertTranslationKey(project.ProjectPath, key, values)
	if err != nil {
		utils.Error("Failed to add translation key: %v", err)
		return err
	}

	utils.Success("Translation key %s added successfully", key)
	return nil
}

// deleteTranslationKey deletes a translation key from all languages
func deleteTranslationKey(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForLocalization()
	if err != nil {
		return err
	}

	// Check if key is provided
	if c.Args().Len() == 0 {
		utils.Error("Translation key is required")
		utils.Info("Usage: fdawg lang delete <key>")
		utils.Info("Example: fdawg lang delete app.welcome")
		return fmt.Errorf("translation key is required")
	}

	key := c.Args().First()

	// Confirm deletion
	utils.Warning("Are you sure you want to delete the translation key %s from all languages? (y/N): ", key)
	var confirm string
	fmt.Scanln(&confirm)

	if strings.ToLower(confirm) != "y" {
		utils.Info("Deletion cancelled")
		return nil
	}

	// Delete the key from all languages
	utils.Info("Deleting translation key %s from all languages...", key)

	err = localization.DeleteTranslationKey(project.ProjectPath, key)
	if err != nil {
		utils.Error("Failed to delete translation key: %v", err)
		return err
	}

	utils.Success("Translation key %s deleted successfully", key)
	return nil
}

// listLanguages lists all supported languages in the project
func listLanguages(c *cli.Context) error {
	// Validate Flutter project
	project, err := validateFlutterProjectForLocalization()
	if err != nil {
		return err
	}

	// Check if translations directory exists
	translationsPath := filepath.Join(project.ProjectPath, localization.TranslationsDir)
	if _, err := os.Stat(translationsPath); os.IsNotExist(err) {
		utils.Error("Translations directory not found")
		utils.Info("Run 'fdawg lang init' to initialize localization first")
		return nil
	}

	// Get all translation files
	translationFiles, err := localization.ListTranslationFiles(project.ProjectPath)
	if err != nil {
		utils.Error("Failed to list translation files: %v", err)
		return err
	}

	if len(translationFiles) == 0 {
		utils.Info("No languages configured")
		utils.Info("Run 'fdawg lang init' to initialize localization first")
		return nil
	}

	// Display languages
	fmt.Println(utils.Separator("=", 50))
	utils.Success("Supported Languages")
	fmt.Println(utils.Separator("=", 50))

	for _, file := range translationFiles {
		langInfo, err := localization.GetLanguageInfo(file.Language)
		if err == nil {
			fmt.Printf("- %s (%s)\n", langInfo.DisplayName(), file.Language)
		} else {
			fmt.Printf("- %s\n", file.Language)
		}
	}

	fmt.Println(utils.Separator("=", 50))
	utils.Success("Total Languages: %d", len(translationFiles))

	return nil
}
