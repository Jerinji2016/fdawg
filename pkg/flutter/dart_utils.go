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
// - "my file name.jpg" becomes "myFileName"
// - "some weird@#$chars" becomes "someWeirdChars"
func FormatDartVariableName(name string) string {
	// First, replace all non-alphanumeric characters with spaces
	// This handles spaces, dots, underscores, hyphens, and any other special characters
	spaceReplacer := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return ' '
	}, name)

	// Split by spaces
	parts := strings.Fields(spaceReplacer)

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

	// If the result is empty, return a default name
	if len(result) == 0 {
		return "_asset"
	}

	// Ensure the name starts with a letter or underscore
	if !unicode.IsLetter(rune(result[0])) {
		result = "_" + result
	}

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
