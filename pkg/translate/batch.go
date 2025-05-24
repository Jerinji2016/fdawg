package translate

import (
	"fmt"
	"strings"
)

// RowTranslationRequest represents a request to translate an entire row
type RowTranslationRequest struct {
	Key                string            `json:"key"`
	SourceLanguage     string            `json:"source_language"`
	TargetLanguages    []string          `json:"target_languages"`
	ExistingTranslations map[string]string `json:"existing_translations"`
}

// RowTranslationResponse represents the response for row translation
type RowTranslationResponse struct {
	Key          string                       `json:"key"`
	Translations map[string]TranslationResult `json:"translations"`
	Errors       map[string]string            `json:"errors,omitempty"`
}

// TranslationResult represents a single translation result
type TranslationResult struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

// CellTranslationRequest represents a request to translate a single cell
type CellTranslationRequest struct {
	Key                string            `json:"key"`
	TargetLanguage     string            `json:"target_language"`
	SourceLanguage     string            `json:"source_language,omitempty"`
	ExistingTranslations map[string]string `json:"existing_translations"`
}

// CellTranslationResponse represents the response for cell translation
type CellTranslationResponse struct {
	Key            string `json:"key"`
	TargetLanguage string `json:"target_language"`
	SourceLanguage string `json:"source_language"`
	TranslatedText string `json:"translated_text"`
	Success        bool   `json:"success"`
	Error          string `json:"error,omitempty"`
}

// TranslateRow translates an entire row from a source language to multiple target languages
func (s *Service) TranslateRow(req RowTranslationRequest) (*RowTranslationResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("translation service is not enabled")
	}

	// Get the source text
	sourceText, exists := req.ExistingTranslations[req.SourceLanguage]
	if !exists || strings.TrimSpace(sourceText) == "" {
		return nil, fmt.Errorf("no source text found for language: %s", req.SourceLanguage)
	}

	response := &RowTranslationResponse{
		Key:          req.Key,
		Translations: make(map[string]TranslationResult),
		Errors:       make(map[string]string),
	}

	// Filter target languages to only include those that don't already have content
	var languagesToTranslate []string
	for _, targetLang := range req.TargetLanguages {
		if targetLang == req.SourceLanguage {
			continue // Skip source language
		}
		
		existingText := req.ExistingTranslations[targetLang]
		if strings.TrimSpace(existingText) == "" {
			languagesToTranslate = append(languagesToTranslate, targetLang)
		}
	}

	if len(languagesToTranslate) == 0 {
		return response, nil // No languages need translation
	}

	// Translate to each target language
	for _, targetLang := range languagesToTranslate {
		translationReq := TranslationRequest{
			Text:       sourceText,
			SourceLang: req.SourceLanguage,
			TargetLang: targetLang,
		}

		translationResp, err := s.Translate(translationReq)
		if err != nil {
			response.Errors[targetLang] = err.Error()
			response.Translations[targetLang] = TranslationResult{
				Text:       "",
				SourceLang: req.SourceLanguage,
				TargetLang: targetLang,
				Success:    false,
				Error:      err.Error(),
			}
			continue
		}

		response.Translations[targetLang] = TranslationResult{
			Text:       translationResp.TranslatedText,
			SourceLang: req.SourceLanguage,
			TargetLang: targetLang,
			Success:    true,
		}
	}

	return response, nil
}

// TranslateCell translates a single cell
func (s *Service) TranslateCell(req CellTranslationRequest) (*CellTranslationResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("translation service is not enabled")
	}

	response := &CellTranslationResponse{
		Key:            req.Key,
		TargetLanguage: req.TargetLanguage,
		Success:        false,
	}

	// If source language is not specified, try to find one automatically
	sourceLanguage := req.SourceLanguage
	if sourceLanguage == "" {
		// Find the first language that has content
		for lang, text := range req.ExistingTranslations {
			if lang != req.TargetLanguage && strings.TrimSpace(text) != "" {
				sourceLanguage = lang
				break
			}
		}
		
		if sourceLanguage == "" {
			response.Error = "no source text available for translation"
			return response, nil
		}
	}

	// Get the source text
	sourceText, exists := req.ExistingTranslations[sourceLanguage]
	if !exists || strings.TrimSpace(sourceText) == "" {
		response.Error = fmt.Sprintf("no source text found for language: %s", sourceLanguage)
		return response, nil
	}

	response.SourceLanguage = sourceLanguage

	// Perform the translation
	translationReq := TranslationRequest{
		Text:       sourceText,
		SourceLang: sourceLanguage,
		TargetLang: req.TargetLanguage,
	}

	translationResp, err := s.Translate(translationReq)
	if err != nil {
		response.Error = err.Error()
		return response, nil
	}

	response.TranslatedText = translationResp.TranslatedText
	response.Success = true

	return response, nil
}

// GetAvailableSourceLanguages returns languages that have content for translation
func GetAvailableSourceLanguages(translations map[string]string, excludeLanguage string) []string {
	var availableLanguages []string
	
	for lang, text := range translations {
		if lang != excludeLanguage && strings.TrimSpace(text) != "" {
			availableLanguages = append(availableLanguages, lang)
		}
	}
	
	return availableLanguages
}

// HasTranslatableContent checks if there's any content available for translation
func HasTranslatableContent(translations map[string]string) bool {
	for _, text := range translations {
		if strings.TrimSpace(text) != "" {
			return true
		}
	}
	return false
}
