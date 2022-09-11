package convert

import "strings"

func StringToBool(str string) bool {
	if strings.TrimSpace(strings.ToLower(str)) == "true" ||
		strings.TrimSpace(str) == "1" {
		return true
	}
	return false
}

func StringToPtrBool(str string) *bool {
	var result bool
	if str == "" {
		return nil
	}
	if strings.TrimSpace(strings.ToLower(str)) == "true" ||
		strings.TrimSpace(str) == "1" {
		result = true
		return &result
	}
	return &result
}
