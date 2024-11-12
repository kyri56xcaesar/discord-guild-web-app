package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
)

type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("Invalid field %q: %s", e.Field, e.Message)
}

// utils
func IsNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}

func IsAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}

func IsValidUTF8String(s string) bool {
	// Updated regex to include space (\s) and new line (\n) characters
	re := regexp.MustCompile(`^[\p{L}\p{N}\s\n!@#\$%\^&\*\(\):\?><\.\-]+$`)

	return re.MatchString(s)
}

func IsValidColor(s string, allowedColors map[string]bool) bool {
	hexColorRe := regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)
	rgbaColorRe := regexp.MustCompile(`^rgba?\(\d{1,3},\d{1,3},\d{1,3}(?:,\d?(\.\d+)?)?\)$`)
	// Check if the color is a valid hex or rgba color
	return hexColorRe.MatchString(s) || rgbaColorRe.MatchString(s) || allowedColors[s]
}

func IsValidURLOrBase64(s string) bool {
	if isValidURL(s) || isValidBase64(s) || s == "None" {
		return true
	}
	return false
}

func isValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func isValidBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func InterfaceSlice(slice []string) []interface{} {
	interfaces := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaces[i] = v
	}
	return interfaces
}

func KeysSliceFromMap[K comparable, V any](mymap map[K]V) []K {
	keys := make([]K, 0, len(mymap))
	for k := range mymap {
		keys = append(keys, k)
	}
	return keys
}
