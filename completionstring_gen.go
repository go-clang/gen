package phoenix

// #include "go-clang.h"
import "C"

// A semantic string that describes a code-completion result. A semantic string that describes the formatting of a code-completion result as a single "template" of text that should be inserted into the source buffer when a particular code-completion result is selected. Each semantic string is made up of some number of "chunks", each of which contains some text along with a description of what that text means, e.g., the name of the entity being referenced, whether the text chunk is part of the template, or whether it is a "placeholder" that the user should replace with actual code,of a specific kind. See \c CXCompletionChunkKind for a description of the different kinds of chunks.
type CompletionString struct {
	c C.CXCompletionString
}

// Retrieve the brief documentation comment attached to the declaration that corresponds to the given completion string.
func (cs CompletionString) CompletionBriefComment() string {
	o := cxstring{C.clang_getCompletionBriefComment(cs.c)}
	defer o.Dispose()

	return o.String()
}
