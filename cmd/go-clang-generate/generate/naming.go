package generate

import (
	"unicode"
)

func LowerFirstCharacter(s string) string {
	r := []rune(s)

	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func UpperFirstCharacter(s string) string {
	r := []rune(s)

	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

var goKeywordReplacements = map[string]string{
	"range": "r",
}

func ReplaceGoKeywords(s string) string {
	if r, ok := goKeywordReplacements[s]; ok {
		return r
	}

	return ""
}
