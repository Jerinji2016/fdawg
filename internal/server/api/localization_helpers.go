package api

import (
	"fmt"
	"strings"

	"github.com/Jerinji2016/fdawg/pkg/localization"
)

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
					found = true
					value := getValueFromData(file.Data, key)
					if strings.TrimSpace(value) == "" {
						missingCount++
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
