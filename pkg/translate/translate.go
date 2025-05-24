package translate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Service provides translation functionality using Google Translate API
type Service struct {
	config     *Config
	httpClient *http.Client
}

// TranslationRequest represents a translation request
type TranslationRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
	Format     string `json:"format,omitempty"`
}

// TranslationResponse represents a translation response
type TranslationResponse struct {
	TranslatedText string `json:"translated_text"`
	SourceLang     string `json:"source_lang"`
	TargetLang     string `json:"target_lang"`
}

// BatchTranslationRequest represents a batch translation request
type BatchTranslationRequest struct {
	Texts      []string `json:"texts"`
	SourceLang string   `json:"source_lang"`
	TargetLang string   `json:"target_lang"`
	Format     string   `json:"format,omitempty"`
}

// BatchTranslationResponse represents a batch translation response
type BatchTranslationResponse struct {
	Translations []TranslationResponse `json:"translations"`
}

// GoogleTranslateResponse represents the response from Google Translate API
type GoogleTranslateResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText         string `json:"translatedText"`
			DetectedSourceLanguage string `json:"detectedSourceLanguage,omitempty"`
		} `json:"translations"`
	} `json:"data"`
}

// DetectLanguageResponse represents language detection response
type DetectLanguageResponse struct {
	Language   string  `json:"language"`
	Confidence float64 `json:"confidence"`
}

// NewService creates a new translation service
func NewService(config *Config) (*Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// IsEnabled returns whether the translation service is enabled
func (s *Service) IsEnabled() bool {
	return s.config.IsEnabled()
}

// Translate translates text from source language to target language
func (s *Service) Translate(req TranslationRequest) (*TranslationResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("translation service is not enabled")
	}

	// Convert language codes to Google Translate format
	sourceLang, err := GetGoogleLanguageCode(req.SourceLang)
	if err != nil {
		return nil, fmt.Errorf("unsupported source language: %v", err)
	}

	targetLang, err := GetGoogleLanguageCode(req.TargetLang)
	if err != nil {
		return nil, fmt.Errorf("unsupported target language: %v", err)
	}

	// Prepare the request to Google Translate API
	apiURL := "https://translation.googleapis.com/language/translate/v2"

	data := url.Values{}
	data.Set("key", s.config.APIKey)
	data.Set("q", req.Text)
	data.Set("source", sourceLang)
	data.Set("target", targetLang)
	data.Set("format", "text")

	resp, err := s.httpClient.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Translate API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google Translate API error (status %d): %s", resp.StatusCode, string(body))
	}

	var googleResp GoogleTranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, fmt.Errorf("failed to decode Google Translate response: %v", err)
	}

	if len(googleResp.Data.Translations) == 0 {
		return nil, fmt.Errorf("no translations returned from Google Translate API")
	}

	translation := googleResp.Data.Translations[0]

	return &TranslationResponse{
		TranslatedText: translation.TranslatedText,
		SourceLang:     req.SourceLang,
		TargetLang:     req.TargetLang,
	}, nil
}

// BatchTranslate translates multiple texts from source language to target language
func (s *Service) BatchTranslate(req BatchTranslationRequest) (*BatchTranslationResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("translation service is not enabled")
	}

	if len(req.Texts) == 0 {
		return &BatchTranslationResponse{Translations: []TranslationResponse{}}, nil
	}

	// Convert language codes to Google Translate format
	sourceLang, err := GetGoogleLanguageCode(req.SourceLang)
	if err != nil {
		return nil, fmt.Errorf("unsupported source language: %v", err)
	}

	targetLang, err := GetGoogleLanguageCode(req.TargetLang)
	if err != nil {
		return nil, fmt.Errorf("unsupported target language: %v", err)
	}

	// Prepare the request to Google Translate API
	apiURL := "https://translation.googleapis.com/language/translate/v2"

	data := url.Values{}
	data.Set("key", s.config.APIKey)
	for _, text := range req.Texts {
		data.Add("q", text)
	}
	data.Set("source", sourceLang)
	data.Set("target", targetLang)
	data.Set("format", "text")

	resp, err := s.httpClient.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Translate API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google Translate API error (status %d): %s", resp.StatusCode, string(body))
	}

	var googleResp GoogleTranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&googleResp); err != nil {
		return nil, fmt.Errorf("failed to decode Google Translate response: %v", err)
	}

	if len(googleResp.Data.Translations) != len(req.Texts) {
		return nil, fmt.Errorf("unexpected number of translations returned: expected %d, got %d",
			len(req.Texts), len(googleResp.Data.Translations))
	}

	translations := make([]TranslationResponse, len(googleResp.Data.Translations))
	for i, translation := range googleResp.Data.Translations {
		translations[i] = TranslationResponse{
			TranslatedText: translation.TranslatedText,
			SourceLang:     req.SourceLang,
			TargetLang:     req.TargetLang,
		}
	}

	return &BatchTranslationResponse{
		Translations: translations,
	}, nil
}

// DetectLanguage detects the language of the given text
func (s *Service) DetectLanguage(text string) (*DetectLanguageResponse, error) {
	if !s.IsEnabled() {
		return nil, fmt.Errorf("translation service is not enabled")
	}

	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Prepare the request to Google Translate API
	apiURL := "https://translation.googleapis.com/language/translate/v2/detect"

	data := url.Values{}
	data.Set("key", s.config.APIKey)
	data.Set("q", text)

	resp, err := s.httpClient.PostForm(apiURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to call Google Translate API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google Translate API error (status %d): %s", resp.StatusCode, string(body))
	}

	var detectResp struct {
		Data struct {
			Detections [][]struct {
				Language   string  `json:"language"`
				Confidence float64 `json:"confidence"`
			} `json:"detections"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&detectResp); err != nil {
		return nil, fmt.Errorf("failed to decode language detection response: %v", err)
	}

	if len(detectResp.Data.Detections) == 0 || len(detectResp.Data.Detections[0]) == 0 {
		return nil, fmt.Errorf("no language detection results returned")
	}

	detection := detectResp.Data.Detections[0][0]

	return &DetectLanguageResponse{
		Language:   detection.Language,
		Confidence: detection.Confidence,
	}, nil
}
