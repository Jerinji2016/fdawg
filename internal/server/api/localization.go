package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/config"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/localization"
	"github.com/Jerinji2016/fdawg/pkg/translate"
)

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

// SetupLocalizationAPIRoutes sets up the localization API routes
func SetupLocalizationAPIRoutes(project *flutter.ValidationResult) {
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

	// Translation configuration endpoint
	http.HandleFunc("/api/localizations/translate-config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Load translation configuration
		config, err := translate.LoadConfig(project.ProjectPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load translation config: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled":   config.IsEnabled(),
			"hasApiKey": config.APIKey != "",
		})
	})

	// Translate cell endpoint
	http.HandleFunc("/api/localizations/translate-cell", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Load translation service
		config, err := translate.LoadConfig(project.ProjectPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load translation config: %v", err), http.StatusInternalServerError)
			return
		}

		if !config.IsEnabled() {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Translation service is not enabled. Please configure Google Translate API key in the Web UI.",
			})
			return
		}

		service, err := translate.NewService(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create translation service: %v", err), http.StatusInternalServerError)
			return
		}

		// Parse request parameters
		translationKey := r.FormValue("translation_key")
		if translationKey == "" {
			http.Error(w, "Translation key is required", http.StatusBadRequest)
			return
		}

		targetLanguage := r.FormValue("target_language")
		if targetLanguage == "" {
			http.Error(w, "Target language is required", http.StatusBadRequest)
			return
		}

		sourceLanguage := r.FormValue("source_language") // Optional

		existingTranslationsJSON := r.FormValue("existing_translations")
		if existingTranslationsJSON == "" {
			http.Error(w, "Existing translations data is required", http.StatusBadRequest)
			return
		}

		// Parse existing translations JSON
		var existingTranslations map[string]string
		if err := json.Unmarshal([]byte(existingTranslationsJSON), &existingTranslations); err != nil {
			http.Error(w, "Invalid existing translations JSON", http.StatusBadRequest)
			return
		}

		// Create translation request
		req := translate.CellTranslationRequest{
			Key:                  translationKey,
			TargetLanguage:       targetLanguage,
			SourceLanguage:       sourceLanguage,
			ExistingTranslations: existingTranslations,
		}

		// Perform translation
		resp, err := service.TranslateCell(req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Translation failed: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Translate row endpoint
	http.HandleFunc("/api/localizations/translate-row", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Load translation service
		config, err := translate.LoadConfig(project.ProjectPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load translation config: %v", err), http.StatusInternalServerError)
			return
		}

		if !config.IsEnabled() {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Translation service is not enabled. Please configure Google Translate API key in the Web UI.",
			})
			return
		}

		service, err := translate.NewService(config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create translation service: %v", err), http.StatusInternalServerError)
			return
		}

		// Parse request parameters
		translationKey := r.FormValue("translation_key")
		if translationKey == "" {
			http.Error(w, "Translation key is required", http.StatusBadRequest)
			return
		}

		sourceLanguage := r.FormValue("source_language")
		if sourceLanguage == "" {
			http.Error(w, "Source language is required", http.StatusBadRequest)
			return
		}

		targetLanguagesJSON := r.FormValue("target_languages")
		if targetLanguagesJSON == "" {
			http.Error(w, "Target languages data is required", http.StatusBadRequest)
			return
		}

		existingTranslationsJSON := r.FormValue("existing_translations")
		if existingTranslationsJSON == "" {
			http.Error(w, "Existing translations data is required", http.StatusBadRequest)
			return
		}

		// Parse target languages JSON
		var targetLanguages []string
		if err := json.Unmarshal([]byte(targetLanguagesJSON), &targetLanguages); err != nil {
			http.Error(w, "Invalid target languages JSON", http.StatusBadRequest)
			return
		}

		// Parse existing translations JSON
		var existingTranslations map[string]string
		if err := json.Unmarshal([]byte(existingTranslationsJSON), &existingTranslations); err != nil {
			http.Error(w, "Invalid existing translations JSON", http.StatusBadRequest)
			return
		}

		// Create translation request
		req := translate.RowTranslationRequest{
			Key:                  translationKey,
			SourceLanguage:       sourceLanguage,
			TargetLanguages:      targetLanguages,
			ExistingTranslations: existingTranslations,
		}

		// Perform translation
		resp, err := service.TranslateRow(req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Translation failed: %v", err),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	// Update translation API key endpoint
	http.HandleFunc("/api/localizations/update-api-key", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		apiKey := r.FormValue("api_key")

		// Check if this is an update (config already exists)
		existingConfig, err := config.GetTranslationConfig(project.ProjectPath)
		isUpdate := err == nil && existingConfig.GoogleTranslateAPIKey != ""

		// For new configurations, API key is required
		// For updates, empty API key means keep the existing one
		if !isUpdate && apiKey == "" {
			http.Error(w, "API key is required", http.StatusBadRequest)
			return
		}

		// If updating and API key is empty, keep the existing one
		if isUpdate && apiKey == "" {
			apiKey = existingConfig.GoogleTranslateAPIKey
		}

		// Update the API key in the project config
		err = config.UpdateTranslationAPIKey(project.ProjectPath, apiKey)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to update API key: %v", err),
			})
			return
		}

		message := "Google Translate API key saved successfully"
		if isUpdate {
			message = "Google Translate API key updated successfully"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": message,
		})
	})
}
