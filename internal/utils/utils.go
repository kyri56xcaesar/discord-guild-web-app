package utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
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

func MergeMaps[K comparable, V any](maps []map[K]V) map[K]V {
	merged := make(map[K]V)
	for _, m := range maps {
		for key, value := range m {
			merged[key] = value
		}
	}

	return merged
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

func FilterStruct(data interface{}) (interface{}, error) {
	// Prepare the result
	result := []string{}

	// Get the reflection value
	value := reflect.ValueOf(data)

	// Ensure the input is either a struct or a slice of structs
	switch value.Kind() {
	case reflect.Ptr:
		value = value.Elem() // Dereference pointer
		if value.Kind() == reflect.Struct {
			// Single struct
			filtered, err := filterSingleStruct(value)
			if err != nil {
				return nil, err
			}
			return filtered, nil
		}
	case reflect.Struct:
		// Single struct
		filtered, err := filterSingleStruct(value)
		if err != nil {
			return nil, err
		}
		return filtered, nil
	case reflect.Slice:
		// Slice of structs
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem() // Dereference pointer
			}
			if elem.Kind() != reflect.Struct {
				return nil, fmt.Errorf("slice contains non-struct element")
			}
			filtered, err := filterSingleStruct(elem)
			if err != nil {
				return nil, err
			}
			result = append(result, filtered...)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("input must be a struct, pointer to struct, or slice of structs, got %s", value.Kind())
	}

	return nil, fmt.Errorf("unsupported input type")
}

func filterSingleStruct(value reflect.Value) ([]string, error) {
	filtered := []string{}

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)

		// Check if the field is a string and not empty
		if fieldValue.Kind() == reflect.String {
			strValue := fieldValue.String()
			if strValue != "" {
				filtered = append(filtered, strValue)
			}
		} else if fieldValue.Kind() == reflect.Int {
			intValue := fieldValue.Int()
			if intValue != 0 {
				filtered = append(filtered, strconv.Itoa(int(intValue)))
			}
		}
	}

	return filtered, nil
}
