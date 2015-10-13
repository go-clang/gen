package main

import (
	"regexp"
	"strings"
)

var (
	reReplaceDoxygenTags = regexp.MustCompile(`\\brief`)
	reReplaceCComments   = regexp.MustCompile(`[ \t]*\n[ \t]*\*[ \t]*`)
	reReplaceMultilines  = regexp.MustCompile(`[ \t]*\n[ \t]*`)
)

func cleanDoxygenComment(comment string) string {
	// Do not touch the comment if there is a code example in it
	if strings.Contains(comment, "\\code") {
		// But at least make the comment more Go friendly
		if strings.HasPrefix(comment, "/**") {
			comment = "/*" + comment[3:]
		}

		return comment
	}

	comment = strings.TrimPrefix(comment, "/**")
	comment = strings.TrimSuffix(comment, "*/")

	comment = reReplaceDoxygenTags.ReplaceAllString(comment, "")
	comment = reReplaceCComments.ReplaceAllString(comment, " ")
	comment = reReplaceMultilines.ReplaceAllString(comment, " ")

	comment = strings.Replace(comment, "  ", " ", -1)

	return "// " + strings.TrimSpace(comment)
}

func trimClangPrefix(name string) string {
	name = strings.TrimPrefix(name, "CX_CXX")
	name = strings.TrimPrefix(name, "CXX")
	name = strings.TrimPrefix(name, "CX")

	return name
}
