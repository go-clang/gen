package phoenix

// #include "go-clang.h"
import "C"

// A semantic string that describes a code-completion result. A semantic string that describes the formatting of a code-completion result as a single "template" of text that should be inserted into the source buffer when a particular code-completion result is selected. Each semantic string is made up of some number of "chunks", each of which contains some text along with a description of what that text means, e.g., the name of the entity being referenced, whether the text chunk is part of the template, or whether it is a "placeholder" that the user should replace with actual code,of a specific kind. See \c CXCompletionChunkKind for a description of the different kinds of chunks.
type CompletionString struct {
	c C.CXCompletionString
}

// Retrieve the number of chunks in the given code-completion string.
func (cs CompletionString) NumCompletionChunks() uint16 {
	return uint16(C.clang_getNumCompletionChunks(cs.c))
}

// Determine the priority of this code completion. The priority of a code completion indicates how likely it is that this particular completion is the completion that the user will select. The priority is selected by various internal heuristics. \param completion_string The completion string to query. \returns The priority of this completion string. Smaller values indicate higher-priority (more likely) completions.
func (cs CompletionString) CompletionPriority() uint16 {
	return uint16(C.clang_getCompletionPriority(cs.c))
}

// Determine the availability of the entity that this code-completion string refers to. \param completion_string The completion string to query. \returns The availability of the completion string.
func (cs CompletionString) CompletionAvailability() AvailabilityKind {
	return AvailabilityKind(C.clang_getCompletionAvailability(cs.c))
}

// Retrieve the number of annotations associated with the given completion string. \param completion_string the completion string to query. \returns the number of annotations associated with the given completion string.
func (cs CompletionString) CompletionNumAnnotations() uint16 {
	return uint16(C.clang_getCompletionNumAnnotations(cs.c))
}

// Retrieve the brief documentation comment attached to the declaration that corresponds to the given completion string.
func (cs CompletionString) CompletionBriefComment() string {
	o := cxstring{C.clang_getCompletionBriefComment(cs.c)}
	defer o.Dispose()

	return o.String()
}
