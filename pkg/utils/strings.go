package utils

// Separator creates a string of repeated characters
func Separator(char string, length int) string {
    result := ""
    for i := 0; i < length; i++ {
        result += char
    }
    return result
}