package flutter

import (
	"strings"
	"unicode"
)

// FormatDartVariableName converts a string to a valid Dart variable name in camelCase format
// For example:
// - "API_URL" becomes "apiUrl"
// - "user-name" becomes "userName"
// - "image.png" becomes "imagePng"
// - "123test" becomes "_123test"
func FormatDartVariableName(name string) string {
	// Split by common separators
	var parts []string
	
	// Check if the name contains underscores
	if strings.Contains(name, "_") {
		parts = strings.Split(name, "_")
	} else if strings.Contains(name, "-") {
		// If no underscores, try splitting by hyphens
		parts = strings.Split(name, "-")
	} else if strings.Contains(name, ".") {
		// If no underscores or hyphens, try splitting by dots (for filenames)
		parts = strings.Split(name, ".")
	} else {
		// If no common separators, treat the whole string as one part
		parts = []string{name}
	}
	
	// Process each part
	for i := range parts {
		if len(parts[i]) == 0 {
			continue
		}
		
		// Convert to lowercase
		parts[i] = strings.ToLower(parts[i])
		
		// Capitalize the first letter of each part except the first one
		if i > 0 && len(parts[i]) > 0 {
			r := []rune(parts[i])
			r[0] = unicode.ToUpper(r[0])
			parts[i] = string(r)
		}
	}
	
	// Join the parts
	result := strings.Join(parts, "")
	
	// Ensure the name starts with a letter or underscore
	if len(result) > 0 && unicode.IsDigit(rune(result[0])) {
		result = "_" + result
	}
	
	// Replace any remaining invalid characters with underscores
	result = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}
		return '_'
	}, result)
	
	return result
}

// EnsureValidDartIdentifier ensures that a string is a valid Dart identifier
// It replaces invalid characters with underscores and ensures the name starts with a letter or underscore
func EnsureValidDartIdentifier(name string) string {
	if len(name) == 0 {
		return "_empty"
	}
	
	// Replace invalid characters with underscores
	result := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			return r
		}
		return '_'
	}, name)
	
	// Ensure the name starts with a letter or underscore
	if unicode.IsDigit(rune(result[0])) {
		result = "_" + result
	}
	
	return result
}
