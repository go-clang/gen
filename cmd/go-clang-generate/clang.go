package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reReplaceCComments  = regexp.MustCompile(`[ \t]*\n[ \t]*\*[ \t]*`)
	reReplaceMultilines = regexp.MustCompile(`[ \t]*\n[ \t]*`)
)

func cleanDoxygenComment(comment string) string {
	// Remove C style comment
	comment = strings.TrimPrefix(comment, "/**")
	comment = strings.TrimSuffix(comment, "*/")
	comment = reReplaceCComments.ReplaceAllString(comment, "\n")

	// Replace some tags
	comment = strings.Replace(comment, "\\brief ", "", -1)
	comment = strings.Replace(comment, "\\c ", "", -1)
	comment = strings.Replace(comment, "\\param ", "Parameter ", -1)
	comment = strings.Replace(comment, "\\returns ", "Returns ", -1)

	// If there is no empty line we definitely do the comment in one line
	if !strings.Contains(comment, "\n\n") {
		comment = strings.Replace(comment, "\n", " ", -1)
	}

	// Replace spaces
	comment = strings.Replace(comment, "  ", " ", -1)
	comment = strings.TrimSpace(comment)

	if comment == "" {
		return comment
	}

	// This might be a bug in the parsing by clang. The nearest comment is added to an item as its comment even tough it is not directly the item's comment
	if strings.HasPrefix(comment, "\\defgroup") {
		return ""
	}

	// Indent multiline comments
	if strings.ContainsRune(comment, '\n') {
		comment = reReplaceMultilines.ReplaceAllString(comment, "\n\t")

		return "/*\n\t" + comment + "\n*/"
	} else {
		return "// " + comment
	}
}

func trimClangPrefix(name string) string {
	name = strings.TrimPrefix(name, "CX_CXX")
	name = strings.TrimPrefix(name, "CXX")
	name = strings.TrimPrefix(name, "CX")
	name = strings.TrimPrefix(name, "ObjC")

	return name
}

type Version struct {
	Major    int
	Minor    int
	Subminor int
}

func ParseVersion(s []byte) *Version {
	m := regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?`).FindSubmatch(s)
	if m == nil {
		return nil
	}

	var err error
	var v Version

	if v.Major, err = strconv.Atoi(string(m[1])); err != nil {
		return nil
	}
	if v.Minor, err = strconv.Atoi(string(m[2])); err != nil {
		return nil
	}
	if len(m[3]) != 0 {
		if v.Subminor, err = strconv.Atoi(string(m[3])); err != nil {
			return nil
		}
	} else {
		v.Subminor = 0
	}

	return &v
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Subminor)
}

func (v Version) StringMinor() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}
