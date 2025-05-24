package api

import "strings"

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
