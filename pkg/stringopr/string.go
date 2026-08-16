package stringopr

import (
	"strings"
	"unicode"
)

func ToCapitalCase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		runes := []rune(strings.ToLower(word))
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

func IsValidCapitalCase(input string) bool {
	words := strings.Fields(input)
	for _, word := range words {
		if len(word) == 0 {
			continue
		}
		if !unicode.IsUpper(rune(word[0])) {
			return false
		}
		for _, char := range word[1:] {
			if !unicode.IsLower(char) {
				return false
			}
		}
	}
	return true
}
