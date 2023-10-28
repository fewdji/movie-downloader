package helper

import (
	"hash/fnv"
	"strconv"
	"strings"
)

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

func Hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.Itoa(int(h.Sum32()))
}
