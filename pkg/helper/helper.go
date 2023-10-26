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
