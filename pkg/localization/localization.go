package localization

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/utils"
)

const (
	// TranslationsDir is the directory where translations are stored
	TranslationsDir = "assets/translations"

	// DefaultLanguage is the default language code
	DefaultLanguage = "en"

	// DefaultCountry is the default country code
	DefaultCountry = "US"
)

// TranslationFile represents a translation file
type TranslationFile struct {
	Path     string                 // Path to the translation file
	Language string                 // Language code
	Data     map[string]interface{} // Translation data
}

// InitializationStatus represents the localization initialization status
type InitializationStatus struct {
	IsInitialized       bool   `json:"isInitialized"`
	HasTranslationsDir  bool   `json:"hasTranslationsDir"`
	HasTranslationFiles bool   `json:"hasTranslationFiles"`
	HasEasyLocalization bool   `json:"hasEasyLocalization"`
	HasIOSConfig        bool   `json:"hasIOSConfig"`
	ErrorMessage        string `json:"errorMessage,omitempty"`
}

// InitLocalization initializes localization in a Flutter project
func InitLocalization(projectPath string) error {
	// Create translations directory
	translationsPath := filepath.Join(projectPath, TranslationsDir)
	if err := utils.EnsureDirExists(translationsPath); err != nil {
		return fmt.Errorf("failed to create translations directory: %v", err)
	}

	// Create default translation file (en_US.json)
	defaultLang := fmt.Sprintf("%s_%s", DefaultLanguage, DefaultCountry)
	defaultTranslationPath := filepath.Join(translationsPath, defaultLang+".json")

	// Check if the file already exists
	if _, err := os.Stat(defaultTranslationPath); os.IsNotExist(err) {
		// Create empty translation file with a welcome message
		defaultTranslations := map[string]interface{}{
			"app": map[string]interface{}{
				"title":   "Flutter App",
				"welcome": "Welcome to Flutter!",
			},
		}

		// Write the file
		if err := writeTranslationFile(defaultTranslationPath, defaultTranslations); err != nil {
			return fmt.Errorf("failed to create default translation file: %v", err)
		}
	}

	// Update pubspec.yaml to add easy_localization dependency and assets
	if err := updatePubspecForLocalization(projectPath); err != nil {
		return fmt.Errorf("failed to update pubspec.yaml: %v", err)
	}

	// Update iOS Info.plist for localization
	if err := updateIOSInfoPlistForLocalization(projectPath, defaultLang); err != nil {
		utils.Info("Warning: Failed to update iOS Info.plist: %v", err)
		// Don't fail the initialization if iOS config fails
	}

	return nil
}

// IsInitialized checks if localization is properly initialized in the project
func IsInitialized(projectPath string) InitializationStatus {
	status := InitializationStatus{}

	// Check if translations directory exists
	translationsPath := filepath.Join(projectPath, TranslationsDir)
	if _, err := os.Stat(translationsPath); err == nil {
		status.HasTranslationsDir = true
	}

	// Check if translation files exist
	if status.HasTranslationsDir {
		translationFiles, err := ListTranslationFiles(projectPath)
		if err == nil && len(translationFiles) > 0 {
			status.HasTranslationFiles = true
		}
	}

	// Check if pubspec.yaml has easy_localization dependency
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	if pubspecData, err := os.ReadFile(pubspecPath); err == nil {
		pubspecContent := string(pubspecData)
		if strings.Contains(pubspecContent, "easy_localization:") {
			status.HasEasyLocalization = true
		}
	}

	// Check if iOS Info.plist has localization configuration
	iosInfoPlistPath := filepath.Join(projectPath, "ios", "Runner", "Info.plist")
	if infoPlistData, err := os.ReadFile(iosInfoPlistPath); err == nil {
		infoPlistContent := string(infoPlistData)
		if strings.Contains(infoPlistContent, "CFBundleLocalizations") {
			status.HasIOSConfig = true
		}
	}

	// Determine if fully initialized (iOS config is optional)
	status.IsInitialized = status.HasTranslationsDir &&
		status.HasTranslationFiles &&
		status.HasEasyLocalization

	return status
}

// AddLanguage adds a new language to the project
func AddLanguage(projectPath, languageCode string) error {
	// Validate language code
	langInfo, err := GetLanguageInfo(languageCode)
	if err != nil {
		return err
	}

	// Format the language code
	formattedCode := langInfo.String()

	// Check if translations directory exists
	translationsPath := filepath.Join(projectPath, TranslationsDir)
	if _, err := os.Stat(translationsPath); os.IsNotExist(err) {
		return fmt.Errorf("translations directory not found. Please run 'fdawg lang init' first to initialize localization")
	}

	// Check if the language already exists
	translationFilePath := filepath.Join(translationsPath, formattedCode+".json")

	if _, err := os.Stat(translationFilePath); err == nil {
		return fmt.Errorf("language %s (%s) already exists", langInfo.DisplayName(), formattedCode)
	}

	// Get the default translation file to copy structure
	defaultLang := fmt.Sprintf("%s_%s", DefaultLanguage, DefaultCountry)
	defaultTranslationPath := filepath.Join(translationsPath, defaultLang+".json")

	// Check if default translation file exists
	if _, err := os.Stat(defaultTranslationPath); os.IsNotExist(err) {
		return fmt.Errorf("default translation file (%s.json) not found. Please run 'fdawg lang init' first to initialize localization", defaultLang)
	}

	// Read default translations
	defaultTranslations, err := readTranslationFile(defaultTranslationPath)
	if err != nil {
		return fmt.Errorf("failed to read default translations: %v. Please run 'fdawg lang init' first to initialize localization", err)
	}

	// Create a new translation file with the same structure but empty values
	newTranslations := make(map[string]interface{})
	copyTranslationStructure(defaultTranslations, newTranslations, "")

	// Write the new translation file
	if err := writeTranslationFile(translationFilePath, newTranslations); err != nil {
		return fmt.Errorf("failed to create translation file: %v", err)
	}

	// Update iOS Info.plist to add the new language
	if err := addLanguageToIOSInfoPlist(projectPath, formattedCode); err != nil {
		utils.Info("Warning: Failed to update iOS Info.plist: %v", err)
		// Don't fail if iOS config fails
	}

	return nil
}

// RemoveLanguage removes a language from the project
func RemoveLanguage(projectPath, languageCode string) error {
	// Validate language code
	langInfo, err := GetLanguageInfo(languageCode)
	if err != nil {
		return err
	}

	// Format the language code
	formattedCode := langInfo.String()

	// Check if translations directory exists
	translationsPath := filepath.Join(projectPath, TranslationsDir)
	if _, err := os.Stat(translationsPath); os.IsNotExist(err) {
		return fmt.Errorf("translations directory not found. Please run 'fdawg lang init' first to initialize localization")
	}

	// Check if it's the default language
	defaultLang := fmt.Sprintf("%s_%s", DefaultLanguage, DefaultCountry)
	if formattedCode == defaultLang {
		return fmt.Errorf("cannot remove the default language (%s)", defaultLang)
	}

	// Check if the language exists
	translationFilePath := filepath.Join(translationsPath, formattedCode+".json")

	if _, err := os.Stat(translationFilePath); os.IsNotExist(err) {
		return fmt.Errorf("language %s (%s) does not exist. Use 'fdawg lang list' to see available languages", langInfo.DisplayName(), formattedCode)
	}

	// Remove the translation file
	if err := os.Remove(translationFilePath); err != nil {
		return fmt.Errorf("failed to remove translation file: %v", err)
	}

	// Update iOS Info.plist to remove the language
	if err := removeLanguageFromIOSInfoPlist(projectPath, formattedCode); err != nil {
		utils.Info("Warning: Failed to update iOS Info.plist: %v", err)
		// Don't fail if iOS config fails
	}

	utils.Info("Language %s removed from translation files", formattedCode)
	utils.Info("Note: You may want to manually remove the language from main.dart if no longer needed")

	return nil
}

// InsertTranslationKey adds a new translation key to all language files
func InsertTranslationKey(projectPath, key string, values map[string]string) error {
	// Validate key format
	if !isValidTranslationKey(key) {
		return fmt.Errorf("invalid translation key format: %s (use dot notation, e.g., 'category.subcategory.key')", key)
	}

	// Check if translations directory exists
	translationsPath := filepath.Join(projectPath, TranslationsDir)
	if _, err := os.Stat(translationsPath); os.IsNotExist(err) {
		return fmt.Errorf("translations directory not found. Please run 'fdawg lang init' first to initialize localization")
	}

	// Get all translation files
	translationFiles, err := ListTranslationFiles(projectPath)
	if err != nil {
		return fmt.Errorf("failed to list translation files: %v", err)
	}

	if len(translationFiles) == 0 {
		return fmt.Errorf("no translation files found. Please run 'fdawg lang init' first to initialize localization")
	}

	// Add the key to each translation file
	for _, file := range translationFiles {
		// Get the value for this language
		value, ok := values[file.Language]
		if !ok {
			// If no value is provided for this language, use an empty string
			value = ""
		}

		// Add the key to the translation file
		if err := addKeyToTranslationFile(file.Path, key, value); err != nil {
			return fmt.Errorf("failed to add key to %s: %v", file.Language, err)
		}
	}

	return nil
}

// DeleteTranslationKey deletes a translation key from all language files
func DeleteTranslationKey(projectPath, key string) error {
	// Check if translations directory exists
	translationsPath := filepath.Join(projectPath, TranslationsDir)
	if _, err := os.Stat(translationsPath); os.IsNotExist(err) {
		return fmt.Errorf("translations directory not found. Please run 'fdawg lang init' first to initialize localization")
	}

	// Get all translation files
	translationFiles, err := ListTranslationFiles(projectPath)
	if err != nil {
		return fmt.Errorf("failed to list translation files: %v", err)
	}

	if len(translationFiles) == 0 {
		return fmt.Errorf("no translation files found. Please run 'fdawg lang init' first to initialize localization")
	}

	// Delete the key from each translation file
	for _, file := range translationFiles {
		if err := deleteKeyFromTranslationFile(file.Path, key); err != nil {
			return fmt.Errorf("failed to delete key from %s: %v", file.Language, err)
		}
	}

	return nil
}

// ListTranslationFiles lists all translation files in the project
func ListTranslationFiles(projectPath string) ([]TranslationFile, error) {
	translationsPath := filepath.Join(projectPath, TranslationsDir)

	// Check if the translations directory exists
	if _, err := os.Stat(translationsPath); os.IsNotExist(err) {
		return nil, nil
	}

	// Read the directory
	entries, err := os.ReadDir(translationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read translations directory: %v", err)
	}

	var files []TranslationFile

	// Process each JSON file
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		// Extract language code from filename
		langCode := strings.TrimSuffix(entry.Name(), ".json")

		// Create TranslationFile
		filePath := filepath.Join(translationsPath, entry.Name())
		data, err := readTranslationFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read translation file %s: %v", entry.Name(), err)
		}

		files = append(files, TranslationFile{
			Path:     filePath,
			Language: langCode,
			Data:     data,
		})
	}

	return files, nil
}

// Helper functions

// writeTranslationFile writes translation data to a JSON file
func writeTranslationFile(path string, data map[string]interface{}) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := utils.EnsureDirExists(dir); err != nil {
		return err
	}

	// Marshal the data with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write the file
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// readTranslationFile reads translation data from a JSON file
func readTranslationFile(path string) (map[string]interface{}, error) {
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Unmarshal the JSON
	var translations map[string]interface{}
	if err := json.Unmarshal(data, &translations); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return translations, nil
}

// copyTranslationStructure copies the structure of a translation map but with empty values
func copyTranslationStructure(src, dst map[string]interface{}, prefix string) {
	for key, value := range src {
		// Build the full key path
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		// Handle nested maps
		if nestedMap, ok := value.(map[string]interface{}); ok {
			// Create a new map for this key
			dst[key] = make(map[string]interface{})
			// Recursively copy the structure
			copyTranslationStructure(nestedMap, dst[key].(map[string]interface{}), fullKey)
		} else {
			// For leaf nodes, use an empty string
			dst[key] = ""
		}
	}
}

// isValidTranslationKey checks if a translation key is valid
func isValidTranslationKey(key string) bool {
	// Key should not be empty
	if key == "" {
		return false
	}

	// Key should not start or end with a dot
	if strings.HasPrefix(key, ".") || strings.HasSuffix(key, ".") {
		return false
	}

	// Key should not contain consecutive dots
	if strings.Contains(key, "..") {
		return false
	}

	// Key should only contain alphanumeric characters, underscores, and dots
	for _, char := range key {
		if !isAlphaNumeric(char) && char != '.' && char != '_' {
			return false
		}
	}

	return true
}

// isAlphaNumeric checks if a character is alphanumeric
func isAlphaNumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

// addKeyToTranslationFile adds a key to a translation file
func addKeyToTranslationFile(path, key, value string) error {
	// Read the translation file
	translations, err := readTranslationFile(path)
	if err != nil {
		return err
	}

	// Split the key into parts
	parts := strings.Split(key, ".")

	// Navigate to the correct location in the nested map
	current := translations
	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part, set the value
			current[part] = value
		} else {
			// Not the last part, ensure the nested map exists
			if _, ok := current[part]; !ok {
				current[part] = make(map[string]interface{})
			}

			// Move to the next level
			if nestedMap, ok := current[part].(map[string]interface{}); ok {
				current = nestedMap
			} else {
				// Convert a leaf node to a map if needed
				current[part] = make(map[string]interface{})
				current = current[part].(map[string]interface{})
			}
		}
	}

	// Write the updated translations back to the file
	return writeTranslationFile(path, translations)
}

// deleteKeyFromTranslationFile deletes a key from a translation file
func deleteKeyFromTranslationFile(path, key string) error {
	// Read the translation file
	translations, err := readTranslationFile(path)
	if err != nil {
		return err
	}

	// Split the key into parts
	parts := strings.Split(key, ".")

	// Navigate to the correct location in the nested map
	current := translations
	parents := make([]map[string]interface{}, 0, len(parts)-1)

	// Navigate to the parent of the key to delete
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]

		// Check if this part exists
		if nestedMap, ok := current[part].(map[string]interface{}); ok {
			parents = append(parents, current)
			current = nestedMap
		} else {
			// Key doesn't exist, nothing to delete
			return nil
		}
	}

	// Delete the key
	lastPart := parts[len(parts)-1]
	delete(current, lastPart)

	// Clean up empty parent maps
	for i := len(parents) - 1; i >= 0; i-- {
		parentMap := parents[i]
		parentKey := parts[i]

		// Check if the child map is empty
		if childMap, ok := parentMap[parentKey].(map[string]interface{}); ok && len(childMap) == 0 {
			// Child map is empty, remove it
			delete(parentMap, parentKey)
		}
	}

	// Write the updated translations back to the file
	return writeTranslationFile(path, translations)
}

// updatePubspecForLocalization updates the pubspec.yaml file to add easy_localization dependency and assets
func updatePubspecForLocalization(projectPath string) error {
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")

	// Read the pubspec.yaml file
	pubspecData, err := os.ReadFile(pubspecPath)
	if err != nil {
		return fmt.Errorf("failed to read pubspec.yaml: %v", err)
	}

	pubspecContent := string(pubspecData)

	// Check if easy_localization is already added
	if strings.Contains(pubspecContent, "easy_localization:") {
		utils.Info("easy_localization dependency already exists in pubspec.yaml")
	} else {
		// Add easy_localization dependency
		dependenciesIndex := strings.Index(pubspecContent, "dependencies:")
		if dependenciesIndex == -1 {
			return fmt.Errorf("could not find dependencies section in pubspec.yaml")
		}

		// Find the end of the dependencies section
		nextSectionIndex := strings.Index(pubspecContent[dependenciesIndex:], "\n\n")
		if nextSectionIndex == -1 {
			// If there's no clear end, append to the end of the file
			pubspecContent += "\n  easy_localization: ^3.0.3\n"
		} else {
			// Insert before the next section
			insertPos := dependenciesIndex + nextSectionIndex
			pubspecContent = pubspecContent[:insertPos] + "\n  easy_localization: ^3.0.3" + pubspecContent[insertPos:]
		}
	}

	// Check if translations assets are already added
	if strings.Contains(pubspecContent, "assets/translations") {
		utils.Info("translations assets already exist in pubspec.yaml")
	} else {
		// Add translations assets
		flutterIndex := strings.Index(pubspecContent, "flutter:")
		if flutterIndex == -1 {
			// If flutter section doesn't exist, add it
			pubspecContent += "\nflutter:\n  assets:\n    - assets/translations/\n"
		} else {
			// Check if assets section exists (commented or uncommented)
			assetsIndex := strings.Index(pubspecContent[flutterIndex:], "assets:")
			commentedAssetsIndex := strings.Index(pubspecContent[flutterIndex:], "# assets:")

			if assetsIndex == -1 && commentedAssetsIndex == -1 {
				// No assets section at all, add one
				usesIndex := strings.Index(pubspecContent[flutterIndex:], "uses-material-design:")
				if usesIndex == -1 {
					// Add after flutter:
					pubspecContent = pubspecContent[:flutterIndex+8] + "\n  assets:\n    - assets/translations/\n" + pubspecContent[flutterIndex+8:]
				} else {
					// Add after uses-material-design section
					insertPos := flutterIndex + usesIndex
					endOfLine := strings.Index(pubspecContent[insertPos:], "\n")
					if endOfLine == -1 {
						pubspecContent += "\n  assets:\n    - assets/translations/\n"
					} else {
						insertPos += endOfLine + 1
						pubspecContent = pubspecContent[:insertPos] + "\n  assets:\n    - assets/translations/\n" + pubspecContent[insertPos:]
					}
				}
			} else if commentedAssetsIndex != -1 && (assetsIndex == -1 || commentedAssetsIndex < assetsIndex) {
				// There's a commented assets section, uncomment it and add our asset
				commentPos := flutterIndex + commentedAssetsIndex
				// Find the line with "# assets:"
				lineStart := strings.LastIndex(pubspecContent[:commentPos], "\n") + 1
				lineEnd := strings.Index(pubspecContent[commentPos:], "\n")
				if lineEnd == -1 {
					lineEnd = len(pubspecContent) - commentPos
				}
				lineEnd += commentPos

				// Replace "# assets:" with "assets:" and add our translation
				oldLine := pubspecContent[lineStart:lineEnd]
				newLine := strings.Replace(oldLine, "# assets:", "assets:", 1)
				pubspecContent = pubspecContent[:lineStart] + newLine + "\n    - assets/translations/" + pubspecContent[lineEnd:]
			} else {
				// Add to existing assets section
				insertPos := flutterIndex + assetsIndex + 8 // 8 is the length of "assets:"
				endOfLine := strings.Index(pubspecContent[insertPos:], "\n")
				if endOfLine == -1 {
					pubspecContent += "\n    - assets/translations/\n"
				} else {
					insertPos += endOfLine + 1
					pubspecContent = pubspecContent[:insertPos] + "    - assets/translations/\n" + pubspecContent[insertPos:]
				}
			}
		}
	}

	// Write the updated pubspec.yaml file
	if err := os.WriteFile(pubspecPath, []byte(pubspecContent), 0644); err != nil {
		return fmt.Errorf("failed to write pubspec.yaml: %v", err)
	}

	return nil
}

// updateIOSInfoPlistForLocalization updates the iOS Info.plist file to add localization support
func updateIOSInfoPlistForLocalization(projectPath, defaultLang string) error {
	iosInfoPlistPath := filepath.Join(projectPath, "ios", "Runner", "Info.plist")

	// Check if iOS directory exists
	if _, err := os.Stat(filepath.Join(projectPath, "ios")); os.IsNotExist(err) {
		return fmt.Errorf("iOS directory not found")
	}

	// Read the Info.plist file
	infoPlistData, err := os.ReadFile(iosInfoPlistPath)
	if err != nil {
		return fmt.Errorf("failed to read Info.plist: %v", err)
	}

	infoPlistContent := string(infoPlistData)

	// Check if CFBundleLocalizations is already present
	if strings.Contains(infoPlistContent, "CFBundleLocalizations") {
		utils.Info("CFBundleLocalizations already exists in Info.plist")
		return nil
	}

	// Find the closing </dict> tag before </plist>
	plistEndIndex := strings.LastIndex(infoPlistContent, "</plist>")
	if plistEndIndex == -1 {
		return fmt.Errorf("could not find </plist> tag in Info.plist")
	}

	dictEndIndex := strings.LastIndex(infoPlistContent[:plistEndIndex], "</dict>")
	if dictEndIndex == -1 {
		return fmt.Errorf("could not find </dict> tag in Info.plist")
	}

	// Extract language code from defaultLang (e.g., "en_US" -> "en")
	langCode := strings.Split(defaultLang, "_")[0]

	// Create the CFBundleLocalizations entry
	localizationsEntry := fmt.Sprintf("\t<key>CFBundleLocalizations</key>\n\t<array>\n\t\t<string>%s</string>\n\t</array>\n", langCode)

	// Insert before the closing </dict> tag
	infoPlistContent = infoPlistContent[:dictEndIndex] + localizationsEntry + infoPlistContent[dictEndIndex:]

	// Write the updated Info.plist file
	if err := os.WriteFile(iosInfoPlistPath, []byte(infoPlistContent), 0644); err != nil {
		return fmt.Errorf("failed to write Info.plist: %v", err)
	}

	return nil
}

// addLanguageToIOSInfoPlist adds a language to the iOS Info.plist file
func addLanguageToIOSInfoPlist(projectPath, langCode string) error {
	iosInfoPlistPath := filepath.Join(projectPath, "ios", "Runner", "Info.plist")

	// Check if iOS directory exists
	if _, err := os.Stat(filepath.Join(projectPath, "ios")); os.IsNotExist(err) {
		return fmt.Errorf("iOS directory not found")
	}

	// Read the Info.plist file
	infoPlistData, err := os.ReadFile(iosInfoPlistPath)
	if err != nil {
		return fmt.Errorf("failed to read Info.plist: %v", err)
	}

	infoPlistContent := string(infoPlistData)

	// Extract language code from langCode (e.g., "en_US" -> "en")
	langCodeParts := strings.Split(langCode, "_")
	localeCode := langCodeParts[0]

	// Check if CFBundleLocalizations exists
	if !strings.Contains(infoPlistContent, "CFBundleLocalizations") {
		// If it doesn't exist, create it with the new language
		return updateIOSInfoPlistForLocalization(projectPath, langCode)
	}

	// Check if the language is already in the array
	if strings.Contains(infoPlistContent, fmt.Sprintf("<string>%s</string>", localeCode)) {
		utils.Info("Language %s is already in CFBundleLocalizations", localeCode)
		return nil
	}

	// Find the CFBundleLocalizations array
	arrayStartIndex := strings.Index(infoPlistContent, "<key>CFBundleLocalizations</key>")
	if arrayStartIndex == -1 {
		return fmt.Errorf("could not find CFBundleLocalizations key in Info.plist")
	}

	arrayStartIndex = strings.Index(infoPlistContent[arrayStartIndex:], "<array>") + arrayStartIndex
	if arrayStartIndex == -1 {
		return fmt.Errorf("could not find CFBundleLocalizations array in Info.plist")
	}

	arrayEndIndex := strings.Index(infoPlistContent[arrayStartIndex:], "</array>") + arrayStartIndex
	if arrayEndIndex == -1 {
		return fmt.Errorf("could not find end of CFBundleLocalizations array in Info.plist")
	}

	// Add the new language before the closing </array> tag
	newLanguageEntry := fmt.Sprintf("\t\t<string>%s</string>\n\t", localeCode)
	infoPlistContent = infoPlistContent[:arrayEndIndex] + newLanguageEntry + infoPlistContent[arrayEndIndex:]

	// Write the updated Info.plist file
	if err := os.WriteFile(iosInfoPlistPath, []byte(infoPlistContent), 0644); err != nil {
		return fmt.Errorf("failed to write Info.plist: %v", err)
	}

	return nil
}

// removeLanguageFromIOSInfoPlist removes a language from the iOS Info.plist file
func removeLanguageFromIOSInfoPlist(projectPath, langCode string) error {
	iosInfoPlistPath := filepath.Join(projectPath, "ios", "Runner", "Info.plist")

	// Check if iOS directory exists
	if _, err := os.Stat(filepath.Join(projectPath, "ios")); os.IsNotExist(err) {
		return fmt.Errorf("iOS directory not found")
	}

	// Read the Info.plist file
	infoPlistData, err := os.ReadFile(iosInfoPlistPath)
	if err != nil {
		return fmt.Errorf("failed to read Info.plist: %v", err)
	}

	infoPlistContent := string(infoPlistData)

	// Extract language code from langCode (e.g., "en_US" -> "en")
	langCodeParts := strings.Split(langCode, "_")
	localeCode := langCodeParts[0]

	// Check if CFBundleLocalizations exists
	if !strings.Contains(infoPlistContent, "CFBundleLocalizations") {
		utils.Info("CFBundleLocalizations not found in Info.plist")
		return nil
	}

	// Check if the language is in the array
	languageEntry := fmt.Sprintf("<string>%s</string>", localeCode)
	if !strings.Contains(infoPlistContent, languageEntry) {
		utils.Info("Language %s is not in CFBundleLocalizations", localeCode)
		return nil
	}

	// Find and remove the language entry
	// Look for the entry with proper indentation and newlines
	patterns := []string{
		fmt.Sprintf("\t\t\t<string>%s</string>\n", localeCode),
		fmt.Sprintf("\t\t<string>%s</string>\n", localeCode),
		fmt.Sprintf("\t<string>%s</string>\n", localeCode),
		fmt.Sprintf("		<string>%s</string>\n", localeCode),  // tab characters
		fmt.Sprintf("			<string>%s</string>\n", localeCode), // more tabs
	}

	removed := false
	for _, pattern := range patterns {
		if strings.Contains(infoPlistContent, pattern) {
			infoPlistContent = strings.Replace(infoPlistContent, pattern, "", 1)
			removed = true
			break
		}
	}

	// If none of the patterns matched, try a more flexible approach
	if !removed {
		// Find the CFBundleLocalizations array
		arrayStartIndex := strings.Index(infoPlistContent, "<key>CFBundleLocalizations</key>")
		if arrayStartIndex == -1 {
			return fmt.Errorf("could not find CFBundleLocalizations key in Info.plist")
		}

		arrayStartIndex = strings.Index(infoPlistContent[arrayStartIndex:], "<array>") + arrayStartIndex
		if arrayStartIndex == -1 {
			return fmt.Errorf("could not find CFBundleLocalizations array in Info.plist")
		}

		arrayEndIndex := strings.Index(infoPlistContent[arrayStartIndex:], "</array>") + arrayStartIndex
		if arrayEndIndex == -1 {
			return fmt.Errorf("could not find end of CFBundleLocalizations array in Info.plist")
		}

		// Extract the array content
		arrayContent := infoPlistContent[arrayStartIndex+7 : arrayEndIndex] // +7 for "<array>"

		// Split by lines and filter out the language
		lines := strings.Split(arrayContent, "\n")
		var newLines []string

		for _, line := range lines {
			if !strings.Contains(line, fmt.Sprintf("<string>%s</string>", localeCode)) {
				newLines = append(newLines, line)
			}
		}

		// Reconstruct the array content
		newArrayContent := strings.Join(newLines, "\n")

		// Replace the old array content with the new one
		infoPlistContent = infoPlistContent[:arrayStartIndex+7] + newArrayContent + infoPlistContent[arrayEndIndex:]
	}

	// Write the updated Info.plist file
	if err := os.WriteFile(iosInfoPlistPath, []byte(infoPlistContent), 0644); err != nil {
		return fmt.Errorf("failed to write Info.plist: %v", err)
	}

	utils.Info("Language %s removed from iOS Info.plist", localeCode)
	return nil
}
