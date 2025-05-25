package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jerinji2016/fdawg/internal/server/helpers"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/localization"
)

// LocalizationAPI handles localization-related API endpoints
type LocalizationAPI struct {
	project *flutter.ValidationResult
}

// NewLocalizationAPI creates a new LocalizationAPI instance
func NewLocalizationAPI(project *flutter.ValidationResult) *LocalizationAPI {
	return &LocalizationAPI{
		project: project,
	}
}

// RegisterRoutes registers localization API routes
func (api *LocalizationAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/localizations/status", api.handleStatus)
	mux.HandleFunc("/api/localizations/init", api.handleInit)
	mux.HandleFunc("/api/localizations/data", api.handleData)
	mux.HandleFunc("/api/localizations/add-language", api.handleAddLanguage)
	mux.HandleFunc("/api/localizations/delete-language", api.handleDeleteLanguage)
	mux.HandleFunc("/api/localizations/download/", api.handleDownloadLanguage)
	mux.HandleFunc("/api/localizations/add-key", api.handleAddKey)
	mux.HandleFunc("/api/localizations/delete-key", api.handleDeleteKey)
	mux.HandleFunc("/api/localizations/update-translations", api.handleUpdateTranslations)
}

// handleStatus handles GET requests to check initialization status
func (api *LocalizationAPI) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := localization.IsInitialized(api.project.ProjectPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handleInit handles POST requests to initialize localization
func (api *LocalizationAPI) handleInit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize localization
	err := localization.InitLocalization(api.project.ProjectPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to initialize localization: %v", err),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Localization initialized successfully",
	})
}

// handleData handles GET requests to get localization data
func (api *LocalizationAPI) handleData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get translation files
	translationFiles, err := localization.ListTranslationFiles(api.project.ProjectPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list translation files: %v", err), http.StatusInternalServerError)
		return
	}

	// Build response data
	data := helpers.BuildLocalizationData(translationFiles)

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

// handleAddLanguage handles POST requests to add languages
func (api *LocalizationAPI) handleAddLanguage(w http.ResponseWriter, r *http.Request) {
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
	err := localization.AddLanguage(api.project.ProjectPath, languageCode)
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
}

// handleDeleteLanguage handles POST requests to delete languages
func (api *LocalizationAPI) handleDeleteLanguage(w http.ResponseWriter, r *http.Request) {
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
	err := localization.RemoveLanguage(api.project.ProjectPath, languageCode)
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
}

// handleDownloadLanguage handles GET requests to download language files
func (api *LocalizationAPI) handleDownloadLanguage(w http.ResponseWriter, r *http.Request) {
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
	translationsDir := filepath.Join(api.project.ProjectPath, "assets/translations")
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
}

// handleAddKey handles POST requests to add translation keys
func (api *LocalizationAPI) handleAddKey(w http.ResponseWriter, r *http.Request) {
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
	err := localization.InsertTranslationKey(api.project.ProjectPath, translationKey, make(map[string]string))
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
}

// handleDeleteKey handles POST requests to delete translation keys
func (api *LocalizationAPI) handleDeleteKey(w http.ResponseWriter, r *http.Request) {
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
	err := localization.DeleteTranslationKey(api.project.ProjectPath, translationKey)
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
}

// handleUpdateTranslations handles POST requests to update translations
func (api *LocalizationAPI) handleUpdateTranslations(w http.ResponseWriter, r *http.Request) {
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
	err := localization.InsertTranslationKey(api.project.ProjectPath, translationKey, translations)
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
}

// SetupLocalizationAPIRoutes sets up localization API routes
func SetupLocalizationAPIRoutes(project *flutter.ValidationResult) {
	localizationAPI := NewLocalizationAPI(project)
	localizationAPI.RegisterRoutes(http.DefaultServeMux)
}
