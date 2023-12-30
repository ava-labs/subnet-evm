package utils

import (
	"strings"
	"unicode"
)

func ContainsString(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func RemoveSpacesAndSpecialChars(str string) string {
	var builder strings.Builder
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
