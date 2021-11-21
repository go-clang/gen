package gen

import (
	"regexp"
	"strings"
)

var (
	reReplaceCComments = regexp.MustCompile(`[ \t]*\n[ \t]*\*[ \t]*`)
)

// CleanDoxygenComment converts Clang Doxygen comment to Go comment.
func CleanDoxygenComment(name, comment string) string {
	// remove C style comment
	comment = strings.TrimPrefix(comment, "/**")
	comment = strings.TrimSuffix(comment, "*/")
	comment = strings.TrimPrefix(comment, "//")
	comment = reReplaceCComments.ReplaceAllString(comment, "\n")

	// replace some tags
	comment = strings.ReplaceAll(comment, "\\brief ", "")
	comment = strings.ReplaceAll(comment, "\\c ", "")
	comment = strings.ReplaceAll(comment, "\\p ", "")                       // \p
	comment = strings.ReplaceAll(comment, "\\param ", "Parameter ")         // parameter
	comment = strings.ReplaceAll(comment, "\\Parameter ", "parameter ")     // parameter
	comment = strings.ReplaceAll(comment, "\\paragraph ", "paragraph ")     // parameter
	comment = strings.ReplaceAll(comment, "\\return ", "Return ")           // return
	comment = strings.ReplaceAll(comment, "\\returns ", "Returns ")         // returns
	comment = strings.ReplaceAll(comment, "\\Returns ", "returns ")         // returns
	comment = strings.ReplaceAll(comment, "\\li ", "The ")                  // args
	comment = strings.ReplaceAll(comment, "\\\\arg ", "arg ")               // args
	comment = strings.ReplaceAll(comment, "\\\\tparam ", "template param ") // args

	comment = strings.ReplaceAll(comment, "\\@try", "try")          // \@try
	comment = strings.ReplaceAll(comment, "\\@catch", "catch")      // \@catch
	comment = strings.ReplaceAll(comment, "\\@finally ", "finally") // \@finally

	comment = strings.ReplaceAll(comment, "\\todo", "TODO:") // \todo
	comment = strings.ReplaceAll(comment, "Note :", "NOTE:") // Note

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

	// indent multiline comments
	if strings.ContainsRune(comment, '\n') {
		const cmPrefix = "// "
		var s strings.Builder
		comments := strings.Split(comment, "\n")
		s.Grow(len(comments) + (len(comments) * len(cmPrefix)))

		indent := false
		for i, cm := range comments {
			if i == 0 && name != "" {
				s.WriteString(cmPrefix)
				s.WriteString(name + " ")
				s.WriteString(LowerFirstCharacter(cm))
				s.WriteByte('\n')
				continue
			}

			// skip \code or \verbatim and add tab indent
			if strings.Contains(cm, "\\code") || strings.Contains(cm, "\\verbatim") {
				indent = true
				continue
			}
			// skip \endcode or \endverbatim and un-indent to next comments
			if strings.Contains(cm, "\\endcode") || strings.Contains(cm, "\\endverbatim") {
				indent = false
				continue
			}

			s.WriteString(cmPrefix)
			if indent {
				s.WriteString(" ")
			}
			s.WriteString(cm)
			if i != len(comments)-1 {
				s.WriteByte('\n')
			}
		}

		comment := s.String()
		comment = strings.ReplaceAll(comment, "// \n", "//\n")
		return strings.TrimSpace(comment)
	}

	if name != "" {
		comment = name + " " + LowerFirstCharacter(comment)
	}

	return "// " + comment
}
