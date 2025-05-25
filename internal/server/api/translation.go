package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jerinji2016/fdawg/pkg/config"
	"github.com/Jerinji2016/fdawg/pkg/flutter"
	"github.com/Jerinji2016/fdawg/pkg/translate"
)

// TranslationAPI handles translation-related API endpoints
type TranslationAPI struct {
	project *flutter.ValidationResult
}

// NewTranslationAPI creates a new TranslationAPI instance
func NewTranslationAPI(project *flutter.ValidationResult) *TranslationAPI {
	return &TranslationAPI{
		project: project,
	}
}

// RegisterRoutes registers translation API routes
func (api *TranslationAPI) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/localizations/translate-config", api.handleTranslateConfig)
	mux.HandleFunc("/api/localizations/translate-cell", api.handleTranslateCell)
	mux.HandleFunc("/api/localizations/translate-row", api.handleTranslateRow)
	mux.HandleFunc("/api/localizations/update-api-key", api.handleUpdateAPIKey)
}

// handleTranslateConfig handles GET requests to get translation configuration
func (api *TranslationAPI) handleTranslateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load translation configuration
	config, err := translate.LoadConfig(api.project.ProjectPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load translation config: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"enabled":   config.IsEnabled(),
		"hasApiKey": config.APIKey != "",
	})
}

// handleTranslateCell handles POST requests to translate a single cell
func (api *TranslationAPI) handleTranslateCell(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load translation service
	config, err := translate.LoadConfig(api.project.ProjectPath)
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
}

// handleTranslateRow handles POST requests to translate an entire row
func (api *TranslationAPI) handleTranslateRow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load translation service
	config, err := translate.LoadConfig(api.project.ProjectPath)
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
}

// handleUpdateAPIKey handles POST requests to update the translation API key
func (api *TranslationAPI) handleUpdateAPIKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	apiKey := r.FormValue("api_key")

	// Check if this is an update (config already exists)
	existingConfig, err := config.GetTranslationConfig(api.project.ProjectPath)
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
	err = config.UpdateTranslationAPIKey(api.project.ProjectPath, apiKey)
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
}

// SetupTranslationAPIRoutes sets up translation API routes
func SetupTranslationAPIRoutes(project *flutter.ValidationResult) {
	translationAPI := NewTranslationAPI(project)
	translationAPI.RegisterRoutes(http.DefaultServeMux)
}
