package gen

import (
	"strings"
	"unicode"
)

// LowerFirstCharacter converts first s character to lower.
func LowerFirstCharacter(s string) string {
	r := []rune(s)

	r[0] = unicode.ToLower(r[0])

	return string(r)
}

// UpperFirstCharacter converts first s character to upper.
func UpperFirstCharacter(s string) string {
	r := []rune(s)

	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

var goKeywordReplacements = map[string]string{
	"range": "r",
}

// ReplaceGoKeywords replaces s to Go keyword.
func ReplaceGoKeywords(s string) string {
	if r, ok := goKeywordReplacements[s]; ok {
		return r
	}

	return ""
}

// TrimCommonFunctionName trims common function-name from name.
func TrimCommonFunctionName(name string, typ Type) string {
	name = TrimCommonFunctionNamePrefix(name)

	switch {
	case strings.HasPrefix(name, typ.GoName+"_"):
		name = strings.TrimPrefix(name, typ.GoName+"_")

	case strings.HasPrefix(name, typ.GoName):
		name = strings.TrimPrefix(name, typ.GoName)
	}

	if strings.HasSuffix(typ.GoName, "Kind") {
		// trim Kind *suffix*
		tkn := strings.TrimSuffix(typ.GoName, "Kind")

		switch {
		case strings.HasPrefix(name, tkn+"_"):
			name = strings.TrimPrefix(name, tkn+"_")

		case strings.HasPrefix(name, tkn):
			name = strings.TrimPrefix(name, tkn)
		}
	}

	name = TrimCommonFunctionNamePrefix(name)

	// if the function name is empty at this point, it is a constructor
	if name == "" {
		name = typ.GoName
	}

	return name
}

// TrimCommonFunctionNamePrefix trims common function-name prefix from name.
func TrimCommonFunctionNamePrefix(name string) string {
	name = strings.TrimPrefix(name, "create")
	name = strings.TrimPrefix(name, "get")

	if len(name) > 4 && unicode.IsUpper(rune(name[3])) {
		name = strings.TrimPrefix(name, "Get")
	}

	switch name {
	case "CXXManglings", "ObjCManglings":
		// conflict if trims language prefix

	default:
		name = TrimLanguagePrefix(name)
	}

	return name
}

// TrimLanguagePrefix trims Language prefix from name.
func TrimLanguagePrefix(name string) string {
	name = strings.TrimPrefix(name, "CX_CXX")
	name = strings.TrimPrefix(name, "CXX")
	name = strings.TrimPrefix(name, "CX")
	name = strings.TrimPrefix(name, "ObjC")
	name = strings.TrimPrefix(name, "_")

	return name
}

// CommonReceiverName returns the common function receiver name.
func CommonReceiverName(s string) string {
	var n []rune

	for _, c := range s {
		if unicode.IsUpper(c) {
			n = append(n, unicode.ToLower(c))
		}
	}

	return string(n)
}
