package gen

import (
	"regexp"
	"strings"
)

var (
	reReplaceCComments  = regexp.MustCompile(`[ \t]*\n[ \t]*\*[ \t]*`)
	reReplaceMultilines = regexp.MustCompile(`[ \t]*\n[ \t]*`)
)

// CleanDoxygenComment converts Clang Doxygen comment to Go comment.
func CleanDoxygenComment(comment string) string {
	// remove C style comment
	comment = strings.TrimPrefix(comment, "/**")
	comment = strings.TrimSuffix(comment, "*/")
	comment = strings.TrimPrefix(comment, "//")
	comment = reReplaceCComments.ReplaceAllString(comment, "\n")

	// replace some tags
	comment = strings.ReplaceAll(comment, "\\brief ", "")
	comment = strings.ReplaceAll(comment, "\\c ", "")
	comment = strings.ReplaceAll(comment, "\\param ", "Parameter ")
	comment = strings.ReplaceAll(comment, "\\returns ", "Returns ")

	// if there is no empty line we definitely do the comment in one line
	if !strings.Contains(comment, "\n\n") {
		comment = strings.ReplaceAll(comment, "\n", " ")
	}

	// replace spaces
	comment = strings.ReplaceAll(comment, "  ", " ")
	comment = strings.TrimSpace(comment)

	if comment == "" {
		return comment
	}

	// this might be a bug in the parsing by clang.
	// the nearest comment is added to an item as its comment even tough it is not directly the item's comment
	if strings.HasPrefix(comment, "\\defgroup") {
		return ""
	}

	// 	indent multiline comments
	// 	if strings.ContainsRune(comment, '\n') {
	// 		comment = reReplaceMultilines.ReplaceAllString(comment, "\n\t")
	// 		return "/*\n\t" + comment + "\n*/"
	// 	}
	//
	// 	return "// " + comment
	if strings.ContainsRune(comment, '\n') {
		comment = reReplaceMultilines.ReplaceAllString(comment, "\n\t")

		return "/*\n\t" + comment + "\n*/"
	} else {
		return "// " + comment
	}
}
