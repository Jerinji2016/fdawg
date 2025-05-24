package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/environment"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/localization"
)

// setupAPIRoutes sets up the API routes for the server
func setupAPIRoutes(project *flutter.ValidationResult) {
	// Localization API routes
	setupLocalizationAPIRoutes(project)
	// Environment API routes
	http.HandleFunc("/api/environment/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Get environment name
		envName := r.FormValue("env_name")
		if envName == "" {
			http.Error(w, "Environment name is required", http.StatusBadRequest)
			return
		}

		// Get copy from parameter
		copyFrom := r.FormValue("copy_from")

		// Create environment file
		var createErr error
		if copyFrom != "" {
			// Copy from existing environment file
			createErr = environment.CopyEnvFile(project.ProjectPath, copyFrom, envName)
		} else {
			// Create empty environment file
			createErr = environment.CreateEnvFile(project.ProjectPath, envName, make(map[string]interface{}))
		}

		if createErr != nil {
			http.Error(w, fmt.Sprintf("Failed to create environment file: %v", createErr), http.StatusInternalServerError)
			return
		}

		// Redirect back to the environment page
		http.Redirect(w, r, fmt.Sprintf("/environment?env=%s", envName), http.StatusSeeOther)
	})

	http.HandleFunc("/api/environment/add-variable", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Get environment name
		envName := r.FormValue("env_name")
		if envName == "" {
			http.Error(w, "Environment name is required", http.StatusBadRequest)
			return
		}

		// Get key and value
		key := r.FormValue("key")
		if key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}

		// Validate key format (must start with letter or underscore, and contain only letters, numbers, and underscores)
		keyRegex := regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
		if !keyRegex.MatchString(key) {
			if regexp.MustCompile("^\\d").MatchString(key) {
				http.Error(w, "Key must not start with a number (Dart variable naming convention)", http.StatusBadRequest)
			} else {
				http.Error(w, "Key must contain only letters, numbers, and underscores (no spaces or special characters)", http.StatusBadRequest)
			}
			return
		}

		valueStr := r.FormValue("value")

		// Parse value (try to convert to appropriate type)
		var value interface{} = valueStr

		// Try to parse as number or boolean
		if strings.EqualFold(valueStr, "true") {
			value = true
		} else if strings.EqualFold(valueStr, "false") {
			value = false
		} else if strings.Contains(valueStr, ".") {
			// Try to parse as float
			var f float64
			if _, err := fmt.Sscanf(valueStr, "%f", &f); err == nil {
				value = f
			}
		} else {
			// Try to parse as integer
			var i int64
			if _, err := fmt.Sscanf(valueStr, "%d", &i); err == nil {
				value = i
			}
		}

		// Add variable to environment file
		err = environment.AddVariable(project.ProjectPath, envName, key, value)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to add variable: %v", err), http.StatusInternalServerError)
			return
		}

		// Redirect back to the environment page
		http.Redirect(w, r, fmt.Sprintf("/environment?env=%s", envName), http.StatusSeeOther)
	})

	http.HandleFunc("/api/environment/delete-variable", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Get environment name
		envName := r.FormValue("env_name")
		if envName == "" {
			http.Error(w, "Environment name is required", http.StatusBadRequest)
			return
		}

		// Get key
		key := r.FormValue("key")
		if key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}

		// Delete variable from environment file
		err = environment.DeleteVariable(project.ProjectPath, envName, key)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete variable: %v", err), http.StatusInternalServerError)
			return
		}

		// Redirect back to the environment page
		http.Redirect(w, r, fmt.Sprintf("/environment?env=%s", envName), http.StatusSeeOther)
	})

	http.HandleFunc("/api/environment/delete-env", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Get environment name
		envName := r.FormValue("env_name")
		if envName == "" {
			http.Error(w, "Environment name is required", http.StatusBadRequest)
			return
		}

		// Delete environment file
		err = environment.DeleteEnvFile(project.ProjectPath, envName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete environment file: %v", err), http.StatusInternalServerError)
			return
		}

		// Redirect back to the environment page
		http.Redirect(w, r, "/environment", http.StatusSeeOther)
	})

	http.HandleFunc("/api/environment/download", func(w http.ResponseWriter, r *http.Request) {
		// Get environment name
		envName := r.URL.Query().Get("env_name")
		if envName == "" {
			http.Error(w, "Environment name is required", http.StatusBadRequest)
			return
		}

		// Get environment file
		envFile, err := environment.GetEnvFile(project.ProjectPath, envName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get environment file: %v", err), http.StatusInternalServerError)
			return
		}

		// Set headers for file download
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", envName))
		w.Header().Set("Content-Type", "application/json")

		// Marshal the variables to JSON
		data, err := json.MarshalIndent(envFile.Variables, "", "  ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Write(data)
	})
}

// setupLocalizationAPIRoutes sets up the localization API routes
func setupLocalizationAPIRoutes(project *flutter.ValidationResult) {
	// Get localization data
	http.HandleFunc("/api/localizations/data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get translation files
		translationFiles, err := localization.ListTranslationFiles(project.ProjectPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list translation files: %v", err), http.StatusInternalServerError)
			return
		}

		// Build response data
		data := buildLocalizationData(translationFiles)

		// Set content type
		w.Header().Set("Content-Type", "application/json")

		// Encode and send response
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	})

	// Add language endpoint
	http.HandleFunc("/api/localizations/add-language", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		languageCode := r.FormValue("language_code")
		if languageCode == "" {
			http.Error(w, "Language code is required", http.StatusBadRequest)
			return
		}

		// Add language using the localization package
		err := localization.AddLanguage(project.ProjectPath, languageCode)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to add language: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Language %s added successfully", languageCode),
		})
	})

	// Delete language endpoint
	http.HandleFunc("/api/localizations/delete-language", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		languageCode := r.FormValue("language_code")
		if languageCode == "" {
			http.Error(w, "Language code is required", http.StatusBadRequest)
			return
		}

		// Delete language using the localization package
		err := localization.RemoveLanguage(project.ProjectPath, languageCode)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to delete language: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Language %s deleted successfully", languageCode),
		})
	})

	// Download language file endpoint
	http.HandleFunc("/api/localizations/download/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract language code from URL path
		path := strings.TrimPrefix(r.URL.Path, "/api/localizations/download/")
		languageCode := strings.TrimSuffix(path, ".json")

		if languageCode == "" {
			http.Error(w, "Language code is required", http.StatusBadRequest)
			return
		}

		// Get translation file path
		translationsDir := filepath.Join(project.ProjectPath, "assets/translations")
		filePath := filepath.Join(translationsDir, languageCode+".json")

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "Translation file not found", http.StatusNotFound)
			return
		}

		// Set headers for file download
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", languageCode))
		w.Header().Set("Content-Type", "application/json")

		// Serve the file
		http.ServeFile(w, r, filePath)
	})

	// Add translation key endpoint
	http.HandleFunc("/api/localizations/add-key", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		translationKey := r.FormValue("translation_key")
		if translationKey == "" {
			http.Error(w, "Translation key is required", http.StatusBadRequest)
			return
		}

		// Add translation key using the localization package
		err := localization.InsertTranslationKey(project.ProjectPath, translationKey, make(map[string]string))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to add translation key: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Translation key %s added successfully", translationKey),
		})
	})

	// Delete translation key endpoint
	http.HandleFunc("/api/localizations/delete-key", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		translationKey := r.FormValue("translation_key")
		if translationKey == "" {
			http.Error(w, "Translation key is required", http.StatusBadRequest)
			return
		}

		// Delete translation key using the localization package
		err := localization.DeleteTranslationKey(project.ProjectPath, translationKey)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to delete translation key: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Translation key %s deleted successfully", translationKey),
		})
	})

	// Update translations endpoint
	http.HandleFunc("/api/localizations/update-translations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		translationKey := r.FormValue("translation_key")
		if translationKey == "" {
			http.Error(w, "Translation key is required", http.StatusBadRequest)
			return
		}

		translationsJSON := r.FormValue("translations")
		if translationsJSON == "" {
			http.Error(w, "Translations data is required", http.StatusBadRequest)
			return
		}

		// Parse translations JSON
		var translations map[string]string
		if err := json.Unmarshal([]byte(translationsJSON), &translations); err != nil {
			http.Error(w, "Invalid translations JSON", http.StatusBadRequest)
			return
		}

		// Update translations using the localization package
		err := localization.InsertTranslationKey(project.ProjectPath, translationKey, translations)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to update translations: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Translations for %s updated successfully", translationKey),
		})
	})
}

// LocalizationData represents the data structure for localization API responses
type LocalizationData struct {
	Languages       []LanguageInfo    `json:"languages"`
	TranslationKeys []TranslationKey  `json:"translationKeys"`
	Stats           LocalizationStats `json:"stats"`
}

// LanguageInfo represents information about a supported language
type LanguageInfo struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	Flag           string `json:"flag"`
	CompletionRate int    `json:"completionRate"`
	MissingKeys    int    `json:"missingKeys"`
}

// TranslationKey represents a translation key with its values in different languages
type TranslationKey struct {
	Key          string            `json:"key"`
	Translations map[string]string `json:"translations"`
}

// LocalizationStats represents statistics about the localization
type LocalizationStats struct {
	SupportedLanguages  int `json:"supportedLanguages"`
	TranslationKeys     int `json:"translationKeys"`
	MissingTranslations int `json:"missingTranslations"`
	CompletionRate      int `json:"completionRate"`
}

// buildLocalizationData builds the localization data from translation files
func buildLocalizationData(translationFiles []localization.TranslationFile) LocalizationData {
	data := LocalizationData{
		Languages:       []LanguageInfo{},
		TranslationKeys: []TranslationKey{},
		Stats:           LocalizationStats{},
	}

	if len(translationFiles) == 0 {
		return data
	}

	// Build language info
	languageMap := make(map[string]*LanguageInfo)
	allKeys := make(map[string]bool)

	// First pass: collect all keys and initialize languages
	for _, file := range translationFiles {
		langInfo := &LanguageInfo{
			Code:           file.Language,
			Name:           getLanguageName(file.Language),
			Flag:           getLanguageFlag(file.Language),
			CompletionRate: 0,
			MissingKeys:    0,
		}
		languageMap[file.Language] = langInfo

		// Collect all keys from this file
		collectKeys(file.Data, "", allKeys)
	}

	// Convert allKeys map to slice
	var keysList []string
	for key := range allKeys {
		keysList = append(keysList, key)
	}

	// Build translation keys data
	for _, key := range keysList {
		translationKey := TranslationKey{
			Key:          key,
			Translations: make(map[string]string),
		}

		for _, file := range translationFiles {
			value := getValueFromData(file.Data, key)
			translationKey.Translations[file.Language] = value
		}

		data.TranslationKeys = append(data.TranslationKeys, translationKey)
	}

	// Calculate completion rates and missing keys
	totalKeys := len(keysList)
	totalMissingTranslations := 0

	for _, langInfo := range languageMap {
		missingCount := 0
		for _, key := range keysList {
			found := false
			for _, file := range translationFiles {
				if file.Language == langInfo.Code {
					value := getValueFromData(file.Data, key)
					if strings.TrimSpace(value) == "" {
						missingCount++
					} else {
						found = true
					}
					break
				}
			}
			if !found {
				missingCount++
			}
		}

		langInfo.MissingKeys = missingCount
		if totalKeys > 0 {
			completedKeys := totalKeys - missingCount
			if completedKeys < 0 {
				completedKeys = 0
			}
			langInfo.CompletionRate = (completedKeys * 100) / totalKeys
		}
		totalMissingTranslations += missingCount

		data.Languages = append(data.Languages, *langInfo)
	}

	// Calculate overall stats
	data.Stats.SupportedLanguages = len(translationFiles)
	data.Stats.TranslationKeys = totalKeys
	data.Stats.MissingTranslations = totalMissingTranslations

	if len(translationFiles) > 0 && totalKeys > 0 {
		totalPossibleTranslations := len(translationFiles) * totalKeys
		completedTranslations := totalPossibleTranslations - totalMissingTranslations
		if completedTranslations < 0 {
			completedTranslations = 0
		}
		data.Stats.CompletionRate = (completedTranslations * 100) / totalPossibleTranslations
	}

	return data
}

// Helper functions for localization data processing

// getLanguageName returns the display name for a language code
func getLanguageName(code string) string {
	languageNames := map[string]string{
		"en":    "English",
		"en_US": "English (US)",
		"en_GB": "English (UK)",
		"es":    "Spanish",
		"es_ES": "Spanish (Spain)",
		"es_MX": "Spanish (Mexico)",
		"fr":    "French",
		"fr_FR": "French (France)",
		"fr_CA": "French (Canada)",
		"de":    "German",
		"de_DE": "German (Germany)",
		"it":    "Italian",
		"pt":    "Portuguese",
		"pt_BR": "Portuguese (Brazil)",
		"pt_PT": "Portuguese (Portugal)",
		"ru":    "Russian",
		"zh":    "Chinese",
		"zh_CN": "Chinese (Simplified)",
		"zh_TW": "Chinese (Traditional)",
		"ja":    "Japanese",
		"ko":    "Korean",
		"ar":    "Arabic",
		"hi":    "Hindi",
		"th":    "Thai",
		"vi":    "Vietnamese",
		"nl":    "Dutch",
		"sv":    "Swedish",
		"da":    "Danish",
		"no":    "Norwegian",
		"fi":    "Finnish",
		"pl":    "Polish",
		"tr":    "Turkish",
		"he":    "Hebrew",
		"cs":    "Czech",
		"sk":    "Slovak",
		"hu":    "Hungarian",
		"ro":    "Romanian",
		"bg":    "Bulgarian",
		"hr":    "Croatian",
		"sr":    "Serbian",
		"sl":    "Slovenian",
		"et":    "Estonian",
		"lv":    "Latvian",
		"lt":    "Lithuanian",
		"uk":    "Ukrainian",
		"be":    "Belarusian",
		"mk":    "Macedonian",
		"mt":    "Maltese",
		"is":    "Icelandic",
		"ga":    "Irish",
		"cy":    "Welsh",
		"eu":    "Basque",
		"ca":    "Catalan",
		"gl":    "Galician",
		"af":    "Afrikaans",
		"sq":    "Albanian",
		"az":    "Azerbaijani",
		"hy":    "Armenian",
		"ka":    "Georgian",
		"kk":    "Kazakh",
		"ky":    "Kyrgyz",
		"mn":    "Mongolian",
		"ne":    "Nepali",
		"si":    "Sinhala",
		"ta":    "Tamil",
		"te":    "Telugu",
		"ml":    "Malayalam",
		"kn":    "Kannada",
		"gu":    "Gujarati",
		"pa":    "Punjabi",
		"bn":    "Bengali",
		"or":    "Odia",
		"as":    "Assamese",
		"ur":    "Urdu",
		"fa":    "Persian",
		"ps":    "Pashto",
		"my":    "Myanmar",
		"km":    "Khmer",
		"lo":    "Lao",
		"bo":    "Tibetan",
		"id":    "Indonesian",
		"ms":    "Malay",
		"tl":    "Tagalog",
		"haw":   "Hawaiian",
		"mi":    "Maori",
		"to":    "Tongan",
		"fj":    "Fijian",
		"sm":    "Samoan",
		"kl":    "Kalaallisut",
		"fo":    "Faroese",
		"gd":    "Scottish Gaelic",
		"br":    "Breton",
		"co":    "Corsican",
		"sc":    "Sardinian",
		"rm":    "Romansh",
		"la":    "Latin",
		"eo":    "Esperanto",
	}

	if name, exists := languageNames[code]; exists {
		return name
	}

	// If not found, return the code itself with proper formatting
	parts := strings.Split(code, "_")
	if len(parts) == 2 {
		return strings.Title(parts[0]) + " (" + strings.ToUpper(parts[1]) + ")"
	}
	return strings.Title(code)
}

// getLanguageFlag returns the flag emoji for a language code
func getLanguageFlag(code string) string {
	languageFlags := map[string]string{
		"en":    "🇺🇸", // Default to US flag for English
		"en_US": "🇺🇸",
		"en_GB": "🇬🇧",
		"en_CA": "🇨🇦",
		"en_AU": "🇦🇺",
		"en_NZ": "🇳🇿",
		"en_IE": "🇮🇪",
		"en_ZA": "🇿🇦",
		"es":    "🇪🇸",
		"es_ES": "🇪🇸",
		"es_MX": "🇲🇽",
		"es_AR": "🇦🇷",
		"es_CO": "🇨🇴",
		"es_CL": "🇨🇱",
		"es_PE": "🇵🇪",
		"es_VE": "🇻🇪",
		"fr":    "🇫🇷",
		"fr_FR": "🇫🇷",
		"fr_CA": "🇨🇦",
		"fr_BE": "🇧🇪",
		"fr_CH": "🇨🇭",
		"de":    "🇩🇪",
		"de_DE": "🇩🇪",
		"de_AT": "🇦🇹",
		"de_CH": "🇨🇭",
		"it":    "🇮🇹",
		"it_IT": "🇮🇹",
		"it_CH": "🇨🇭",
		"pt":    "🇵🇹",
		"pt_PT": "🇵🇹",
		"pt_BR": "🇧🇷",
		"ru":    "🇷🇺",
		"zh":    "🇨🇳",
		"zh_CN": "🇨🇳",
		"zh_TW": "🇹🇼",
		"zh_HK": "🇭🇰",
		"ja":    "🇯🇵",
		"ko":    "🇰🇷",
		"ar":    "🇸🇦", // Default to Saudi Arabia for Arabic
		"ar_SA": "🇸🇦",
		"ar_EG": "🇪🇬",
		"ar_AE": "🇦🇪",
		"ar_JO": "🇯🇴",
		"ar_LB": "🇱🇧",
		"ar_SY": "🇸🇾",
		"ar_IQ": "🇮🇶",
		"ar_KW": "🇰🇼",
		"ar_QA": "🇶🇦",
		"ar_BH": "🇧🇭",
		"ar_OM": "🇴🇲",
		"ar_YE": "🇾🇪",
		"ar_MA": "🇲🇦",
		"ar_TN": "🇹🇳",
		"ar_DZ": "🇩🇿",
		"ar_LY": "🇱🇾",
		"ar_SD": "🇸🇩",
		"hi":    "🇮🇳",
		"th":    "🇹🇭",
		"vi":    "🇻🇳",
		"nl":    "🇳🇱",
		"nl_BE": "🇧🇪",
		"sv":    "🇸🇪",
		"da":    "🇩🇰",
		"no":    "🇳🇴",
		"fi":    "🇫🇮",
		"pl":    "🇵🇱",
		"tr":    "🇹🇷",
		"he":    "🇮🇱",
		"cs":    "🇨🇿",
		"sk":    "🇸🇰",
		"hu":    "🇭🇺",
		"ro":    "🇷🇴",
		"bg":    "🇧🇬",
		"hr":    "🇭🇷",
		"sr":    "🇷🇸",
		"sl":    "🇸🇮",
		"et":    "🇪🇪",
		"lv":    "🇱🇻",
		"lt":    "🇱🇹",
		"uk":    "🇺🇦",
		"be":    "🇧🇾",
		"mk":    "🇲🇰",
		"mt":    "🇲🇹",
		"is":    "🇮🇸",
		"ga":    "🇮🇪",
		"cy":    "🏴󠁧󠁢󠁷󠁬󠁳󠁿", // Wales flag
		"eu":    "🏴",       // Basque flag (generic)
		"ca":    "🏴",       // Catalonia flag (generic)
		"gl":    "🏴",       // Galicia flag (generic)
		"af":    "🇿🇦",
		"sq":    "🇦🇱",
		"az":    "🇦🇿",
		"hy":    "🇦🇲",
		"ka":    "🇬🇪",
		"kk":    "🇰🇿",
		"ky":    "🇰🇬",
		"mn":    "🇲🇳",
		"ne":    "🇳🇵",
		"si":    "🇱🇰",
		"ta":    "🇮🇳",      // Tamil (India)
		"te":    "🇮🇳",      // Telugu (India)
		"ml":    "🇮🇳",      // Malayalam (India)
		"kn":    "🇮🇳",      // Kannada (India)
		"gu":    "🇮🇳",      // Gujarati (India)
		"pa":    "🇮🇳",      // Punjabi (India)
		"bn":    "🇧🇩",      // Bengali (Bangladesh)
		"or":    "🇮🇳",      // Odia (India)
		"as":    "🇮🇳",      // Assamese (India)
		"ur":    "🇵🇰",      // Urdu (Pakistan)
		"fa":    "🇮🇷",      // Persian (Iran)
		"ps":    "🇦🇫",      // Pashto (Afghanistan)
		"my":    "🇲🇲",      // Myanmar
		"km":    "🇰🇭",      // Khmer (Cambodia)
		"lo":    "🇱🇦",      // Lao
		"bo":    "🇨🇳",      // Tibetan (China)
		"id":    "🇮🇩",      // Indonesian
		"ms":    "🇲🇾",      // Malay (Malaysia)
		"tl":    "🇵🇭",      // Tagalog (Philippines)
		"haw":   "🇺🇸",      // Hawaiian (US)
		"mi":    "🇳🇿",      // Maori (New Zealand)
		"to":    "🇹🇴",      // Tongan
		"fj":    "🇫🇯",      // Fijian
		"sm":    "🇼🇸",      // Samoan
		"kl":    "🇬🇱",      // Kalaallisut (Greenland)
		"fo":    "🇫🇴",      // Faroese
		"gd":    "🏴󠁧󠁢󠁳󠁣󠁴󠁿", // Scottish Gaelic
		"br":    "🇫🇷",      // Breton (France)
		"co":    "🇫🇷",      // Corsican (France)
		"sc":    "🇮🇹",      // Sardinian (Italy)
		"rm":    "🇨🇭",      // Romansh (Switzerland)
		"la":    "🇻🇦",      // Latin (Vatican)
		"eo":    "🌍",       // Esperanto (global)
	}

	if flag, exists := languageFlags[code]; exists {
		return flag
	}

	// Extract country code from language_COUNTRY format
	parts := strings.Split(code, "_")
	if len(parts) == 2 {
		countryCode := strings.ToLower(parts[1])
		if flag, exists := languageFlags[countryCode]; exists {
			return flag
		}
	}

	// Default flag
	return "🌐"
}

// collectKeys recursively collects all keys from a nested map
func collectKeys(data map[string]interface{}, prefix string, keys map[string]bool) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if nestedMap, ok := value.(map[string]interface{}); ok {
			collectKeys(nestedMap, fullKey, keys)
		} else {
			keys[fullKey] = true
		}
	}
}

// getValueFromData retrieves a value from nested map using dot notation
func getValueFromData(data map[string]interface{}, key string) string {
	parts := strings.Split(key, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part, get the value
			if value, exists := current[part]; exists {
				if str, ok := value.(string); ok {
					return str
				}
				return fmt.Sprintf("%v", value)
			}
			return ""
		} else {
			// Navigate deeper
			if nestedMap, ok := current[part].(map[string]interface{}); ok {
				current = nestedMap
			} else {
				return ""
			}
		}
	}

	return ""
}
