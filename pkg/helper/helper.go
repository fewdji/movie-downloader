package helper

import "strings"

func ContainsAny(s string, contains []string) bool {
	for _, c := range contains {
		if strings.Contains(strings.ToLower(s), strings.ToLower(c)) {
			return true
		}
	}
	return false
}

func ContainsAll(s string, contains []string) bool {
	for _, c := range contains {
		if !strings.Contains(strings.ReplaceAll(strings.ToLower(s), "ё", "е"), strings.ReplaceAll(strings.ToLower(c), "ё", "е")) {
			return false
		}
	}
	return true
}
