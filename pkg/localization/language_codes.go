package localization

import (
	"fmt"
	"strings"
)

// LanguageCode represents a language code with country variant
type LanguageCode struct {
	Code        string // The language code (e.g., "en", "es", "fr")
	CountryCode string // The country code (e.g., "US", "GB", "ES")
	Name        string // The language name (e.g., "English", "Spanish", "French")
}

// String returns the string representation of the language code
// If CountryCode is empty, returns just the Code
// Otherwise, returns Code_CountryCode (e.g., "en_US")
func (lc LanguageCode) String() string {
	if lc.CountryCode == "" {
		return lc.Code
	}
	return fmt.Sprintf("%s_%s", lc.Code, lc.CountryCode)
}

// DisplayName returns a user-friendly display name for the language
// If CountryCode is empty, returns just the Name
// Otherwise, returns Name (CountryCode) (e.g., "English (US)")
func (lc LanguageCode) DisplayName() string {
	if lc.CountryCode == "" {
		return lc.Name
	}
	return fmt.Sprintf("%s (%s)", lc.Name, lc.CountryCode)
}

// CommonLanguageCodes is a list of common language codes
var CommonLanguageCodes = []LanguageCode{
	{"en", "US", "English"},
	{"en", "GB", "English"},
	{"es", "", "Spanish"},
	{"fr", "", "French"},
	{"de", "", "German"},
	{"it", "", "Italian"},
	{"pt", "BR", "Portuguese"},
	{"pt", "PT", "Portuguese"},
	{"ru", "", "Russian"},
	{"zh", "CN", "Chinese"},
	{"zh", "TW", "Chinese"},
	{"ja", "", "Japanese"},
	{"ko", "", "Korean"},
	{"ar", "", "Arabic"},
	{"hi", "", "Hindi"},
	{"bn", "", "Bengali"},
	{"pa", "", "Punjabi"},
	{"ta", "", "Tamil"},
	{"te", "", "Telugu"},
	{"ml", "", "Malayalam"},
	{"th", "", "Thai"},
	{"vi", "", "Vietnamese"},
	{"id", "", "Indonesian"},
	{"ms", "", "Malay"},
	{"tr", "", "Turkish"},
	{"nl", "", "Dutch"},
	{"sv", "", "Swedish"},
	{"no", "", "Norwegian"},
	{"da", "", "Danish"},
	{"fi", "", "Finnish"},
	{"pl", "", "Polish"},
	{"cs", "", "Czech"},
	{"sk", "", "Slovak"},
	{"hu", "", "Hungarian"},
	{"ro", "", "Romanian"},
	{"bg", "", "Bulgarian"},
	{"hr", "", "Croatian"},
	{"sr", "", "Serbian"},
	{"uk", "", "Ukrainian"},
	{"el", "", "Greek"},
	{"he", "", "Hebrew"},
	{"fa", "", "Persian"},
}

// IsValidLanguageCode checks if a language code is valid
func IsValidLanguageCode(code string) bool {
	normalizedCode := strings.ToLower(code)
	
	// Check if it's a simple language code (e.g., "en")
	if len(normalizedCode) == 2 {
		for _, lc := range CommonLanguageCodes {
			if strings.ToLower(lc.Code) == normalizedCode {
				return true
			}
		}
		return false
	}
	
	// Check if it's a language code with country variant (e.g., "en_US")
	parts := strings.Split(normalizedCode, "_")
	if len(parts) == 2 {
		langCode := parts[0]
		countryCode := strings.ToUpper(parts[1])
		
		for _, lc := range CommonLanguageCodes {
			if strings.ToLower(lc.Code) == langCode && 
			   (lc.CountryCode == "" || lc.CountryCode == countryCode) {
				return true
			}
		}
	}
	
	return false
}

// GetLanguageInfo returns information about a language code
func GetLanguageInfo(code string) (LanguageCode, error) {
	normalizedCode := strings.TrimSpace(code)
	
	// Handle language code with country variant (e.g., "en_US")
	parts := strings.Split(normalizedCode, "_")
	langCode := strings.ToLower(parts[0])
	countryCode := ""
	
	if len(parts) > 1 {
		countryCode = strings.ToUpper(parts[1])
	}
	
	// Look for an exact match first
	for _, lc := range CommonLanguageCodes {
		if strings.ToLower(lc.Code) == langCode && 
		   (countryCode == "" || lc.CountryCode == countryCode) {
			return lc, nil
		}
	}
	
	// If no exact match, look for a language match without country code
	if countryCode != "" {
		for _, lc := range CommonLanguageCodes {
			if strings.ToLower(lc.Code) == langCode {
				// Return with the requested country code
				return LanguageCode{lc.Code, countryCode, lc.Name}, nil
			}
		}
	}
	
	return LanguageCode{}, fmt.Errorf("invalid language code: %s", code)
}

// FormatLanguageCode formats a language code to the standard format
func FormatLanguageCode(code string) (string, error) {
	info, err := GetLanguageInfo(code)
	if err != nil {
		return "", err
	}
	return info.String(), nil
}

// ListLanguageCodes returns a list of all supported language codes
func ListLanguageCodes() []LanguageCode {
	return CommonLanguageCodes
}
