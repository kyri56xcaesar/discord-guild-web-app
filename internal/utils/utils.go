package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("Invalid field %q: %s", e.Field, e.Message)
}

// Security related utils
func IsNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}

func IsAlphanumeric(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(s)
}

func IsAlphanumericPlus(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9,*?=/\\]+$`)
	return re.MatchString(s)
}

func VerifyStrSlice(s []string) bool {
	valid := true
	for _, str := range s {
		valid = IsAlphanumeric(str)
		if !valid {
			return valid
		}
	}

	return valid
}

func ToUpperFirstLetter(s string) string {
	if s == "" {
		return s
	} else if len(s) == 1 {
		return strings.ToUpper(s)
	}

	newStr := strings.ToLower(s)
	newStr = strings.Join([]string{
		strings.ToUpper(string(s[0])),
		s[1:],
	}, "")
	return newStr
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

// Datatype conversions and checks utils
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

func AppendKeys[K comparable, V any](mapslice []map[K]V) []K {
	var allKeys []K
	for _, v := range mapslice {
		allKeys = append(allKeys, KeysSliceFromMap(v)...)
	}

	return allKeys
}

func IsMapSliceEmpty[K comparable, V any](mapslice []map[K]V, isEmpty func(v V) bool) bool {
	for _, m := range mapslice {
		if !IsMapValuesEmpty(m, isEmpty) {
			return false
		}
	}
	return true
}

func IsMapValuesEmpty[K comparable, V any](m map[K]V, isEmpty func(v V) bool) bool {
	for _, v := range m {
		if !isEmpty(v) {
			return false
		}
	}
	return true
}

func IsSliceEmpty(v []string) bool {
	for _, str := range v {
		if len(str) != 0 {
			return false
		}
	}
	return true
}
