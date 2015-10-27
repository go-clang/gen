package main

import (
	"unicode"
)

type Receiver struct {
	Name          string
	CName         string
	Type          string
	CType         string
	PrimitiveType string
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
