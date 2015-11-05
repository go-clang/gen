package main

import (
	"unicode"
)

// Receiver TODO refactor https://github.com/zimmski/go-clang-phoenix/issues/52
type Receiver struct {
	Name  string
	CName string
	Type  Type
}

func receiverName(s string) string {
	var n []rune

	for _, c := range s {
		if unicode.IsUpper(c) {
			n = append(n, unicode.ToLower(c))
		}
	}

	return string(n)
}
