package translate

import (
	"fmt"
	"strings"
)

// LanguageMapping maps easy_localization language codes to Google Translate language codes
var LanguageMapping = map[string]string{
	// Common languages
	"en":    "en",                // English
	"en_US": "en",                // English (US)
	"en_GB": "en",                // English (UK)
	"es":    "es",                // Spanish
	"es_ES": "es",                // Spanish (Spain)
	"es_MX": "es",                // Spanish (Mexico)
	"fr":    "fr",                // French
	"fr_FR": "fr",                // French (France)
	"fr_CA": "fr",                // French (Canada)
	"de":    "de",                // German
	"de_DE": "de",                // German (Germany)
	"it":    "it",                // Italian
	"it_IT": "it",                // Italian (Italy)
	"pt":    "pt",                // Portuguese
	"pt_PT": "pt",                // Portuguese (Portugal)
	"pt_BR": "pt",                // Portuguese (Brazil)
	"ru":    "ru",                // Russian
	"ru_RU": "ru",                // Russian (Russia)
	"zh":    "zh",                // Chinese
	"zh_CN": "zh-cn",             // Chinese (Simplified)
	"zh_TW": "zh-tw",             // Chinese (Traditional)
	"ja":    "ja",                // Japanese
	"ja_JP": "ja",                // Japanese (Japan)
	"ko":    "ko",                // Korean
	"ko_KR": "ko",                // Korean (Korea)
	"ar":    "ar",                // Arabic
	"ar_SA": "ar",                // Arabic (Saudi Arabia)
	"hi":    "hi",                // Hindi
	"hi_IN": "hi",                // Hindi (India)
	"th":    "th",                // Thai
	"th_TH": "th",                // Thai (Thailand)
	"vi":    "vi",                // Vietnamese
	"vi_VN": "vi",                // Vietnamese (Vietnam)
	"tr":    "tr",                // Turkish
	"tr_TR": "tr",                // Turkish (Turkey)
	"pl":    "pl",                // Polish
	"pl_PL": "pl",                // Polish (Poland)
	"nl":    "nl",                // Dutch
	"nl_NL": "nl",                // Dutch (Netherlands)
	"sv":    "sv",                // Swedish
	"sv_SE": "sv",                // Swedish (Sweden)
	"da":    "da",                // Danish
	"da_DK": "da",                // Danish (Denmark)
	"no":    "no",                // Norwegian
	"no_NO": "no",                // Norwegian (Norway)
	"fi":    "fi",                // Finnish
	"fi_FI": "fi",                // Finnish (Finland)
	"he":    "he",                // Hebrew
	"he_IL": "he",                // Hebrew (Israel)
	"cs":    "cs",                // Czech
	"cs_CZ": "cs",                // Czech (Czech Republic)
	"sk":    "sk",                // Slovak
	"sk_SK": "sk",                // Slovak (Slovakia)
	"hu":    "hu",                // Hungarian
	"hu_HU": "hu",                // Hungarian (Hungary)
	"ro":    "ro",                // Romanian
	"ro_RO": "ro",                // Romanian (Romania)
	"bg":    "bg",                // Bulgarian
	"bg_BG": "bg",                // Bulgarian (Bulgaria)
	"hr":    "hr",                // Croatian
	"hr_HR": "hr",                // Croatian (Croatia)
	"sr":    "sr",                // Serbian
	"sr_RS": "sr",                // Serbian (Serbia)
	"sl":    "sl",                // Slovenian
	"sl_SI": "sl",                // Slovenian (Slovenia)
	"et":    "et",                // Estonian
	"et_EE": "et",                // Estonian (Estonia)
	"lv":    "lv",                // Latvian
	"lv_LV": "lv",                // Latvian (Latvia)
	"lt":    "lt",                // Lithuanian
	"lt_LT": "lt",                // Lithuanian (Lithuania)
	"uk":    "uk",                // Ukrainian
	"uk_UA": "uk",                // Ukrainian (Ukraine)
	"be":    "be",                // Belarusian
	"be_BY": "be",                // Belarusian (Belarus)
	"mk":    "mk",                // Macedonian
	"mk_MK": "mk",                // Macedonian (Macedonia)
	"mt":    "mt",                // Maltese
	"mt_MT": "mt",                // Maltese (Malta)
	"is":    "is",                // Icelandic
	"is_IS": "is",                // Icelandic (Iceland)
	"ga":    "ga",                // Irish
	"ga_IE": "ga",                // Irish (Ireland)
	"cy":    "cy",                // Welsh
	"cy_GB": "cy",                // Welsh (Wales)
	"eu":    "eu",                // Basque
	"eu_ES": "eu",                // Basque (Spain)
	"ca":    "ca",                // Catalan
	"ca_ES": "ca",                // Catalan (Spain)
	"gl":    "gl",                // Galician
	"gl_ES": "gl",                // Galician (Spain)
	"af":    "af",                // Afrikaans
	"af_ZA": "af",                // Afrikaans (South Africa)
	"sw":    "sw",                // Swahili
	"sw_KE": "sw",                // Swahili (Kenya)
	"am":    "am",                // Amharic
	"am_ET": "am",                // Amharic (Ethiopia)
	"bn":    "bn",                // Bengali
	"bn_BD": "bn",                // Bengali (Bangladesh)
	"gu":    "gu",                // Gujarati
	"gu_IN": "gu",                // Gujarati (India)
	"kn":    "kn",                // Kannada
	"kn_IN": "kn",                // Kannada (India)
	"ml":    "ml",                // Malayalam
	"ml_IN": "ml",                // Malayalam (India)
	"mr":    "mr",                // Marathi
	"mr_IN": "mr",                // Marathi (India)
	"ne":    "ne",                // Nepali
	"ne_NP": "ne",                // Nepali (Nepal)
	"pa":    "pa",                // Punjabi
	"pa_IN": "pa",                // Punjabi (India)
	"si":    "si",                // Sinhala
	"si_LK": "si",                // Sinhala (Sri Lanka)
	"ta":    "ta",                // Tamil
	"ta_IN": "ta",                // Tamil (India)
	"te":    "te",                // Telugu
	"te_IN": "te",                // Telugu (India)
	"ur":    "ur",                // Urdu
	"ur_PK": "ur",                // Urdu (Pakistan)
	"fa":    "fa",                // Persian
	"fa_IR": "fa",                // Persian (Iran)
	"ps":    "ps",                // Pashto
	"ps_AF": "ps",                // Pashto (Afghanistan)
	"sd":    "sd",                // Sindhi
	"sd_PK": "sd",                // Sindhi (Pakistan)
	"ky":    "ky",                // Kyrgyz
	"ky_KG": "ky",                // Kyrgyz (Kyrgyzstan)
	"kk":    "kk",                // Kazakh
	"kk_KZ": "kk",                // Kazakh (Kazakhstan)
	"uz":    "uz", "uz_UZ": "uz", // Uzbek
	"tg": "tg", "tg_TJ": "tg", // Tajik
	"mn": "mn", "mn_MN": "mn", // Mongolian
	"my": "my", "my_MM": "my", // Myanmar (Burmese)
	"km": "km", "km_KH": "km", // Khmer
	"lo": "lo", "lo_LA": "lo", // Lao
	"ka": "ka", "ka_GE": "ka", // Georgian
	"hy": "hy", "hy_AM": "hy", // Armenian
	"az": "az", "az_AZ": "az", // Azerbaijani
	"sq": "sq", "sq_AL": "sq", // Albanian
	"bs": "bs", "bs_BA": "bs", // Bosnian
	"me": "me", "me_ME": "me", // Montenegrin
	"id": "id", "id_ID": "id", // Indonesian
	"ms": "ms", "ms_MY": "ms", // Malay
	"tl": "tl", "tl_PH": "tl", // Filipino
	"haw": "haw", // Hawaiian
	"mi":  "mi",  // Maori
	"sm":  "sm",  // Samoan
	"to":  "to",  // Tongan
	"fj":  "fj",  // Fijian
	"mg":  "mg",  // Malagasy
	"ny":  "ny",  // Chichewa
	"sn":  "sn",  // Shona
	"st":  "st",  // Sesotho
	"su":  "su",  // Sundanese
	"jw":  "jw",  // Javanese
	"ceb": "ceb", // Cebuano
	"hmn": "hmn", // Hmong
	"co":  "co",  // Corsican
	"eo":  "eo",  // Esperanto
	"fy":  "fy",  // Frisian
	"gd":  "gd",  // Scottish Gaelic
	"ig":  "ig",  // Igbo
	"ku":  "ku",  // Kurdish
	"la":  "la",  // Latin
	"lb":  "lb",  // Luxembourgish
	"xh":  "xh",  // Xhosa
	"yi":  "yi",  // Yiddish
	"yo":  "yo",  // Yoruba
	"zu":  "zu",  // Zulu
}

// GetGoogleLanguageCode converts an easy_localization language code to Google Translate language code
func GetGoogleLanguageCode(easyLocalizationCode string) (string, error) {
	// First try exact match
	if googleCode, exists := LanguageMapping[easyLocalizationCode]; exists {
		return googleCode, nil
	}

	// If no exact match, try the base language (before underscore)
	if strings.Contains(easyLocalizationCode, "_") {
		baseCode := strings.Split(easyLocalizationCode, "_")[0]
		if googleCode, exists := LanguageMapping[baseCode]; exists {
			return googleCode, nil
		}
	}

	return "", fmt.Errorf("unsupported language code: %s", easyLocalizationCode)
}

// IsLanguageSupported checks if a language code is supported for translation
func IsLanguageSupported(easyLocalizationCode string) bool {
	_, err := GetGoogleLanguageCode(easyLocalizationCode)
	return err == nil
}

// GetSupportedLanguages returns all supported easy_localization language codes
func GetSupportedLanguages() []string {
	languages := make([]string, 0, len(LanguageMapping))
	for code := range LanguageMapping {
		languages = append(languages, code)
	}
	return languages
}
