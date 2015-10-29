package phoenix

// #include "go-clang.h"
import "C"

// A semantic string that describes a code-completion result. A semantic string that describes the formatting of a code-completion result as a single "template" of text that should be inserted into the source buffer when a particular code-completion result is selected. Each semantic string is made up of some number of "chunks", each of which contains some text along with a description of what that text means, e.g., the name of the entity being referenced, whether the text chunk is part of the template, or whether it is a "placeholder" that the user should replace with actual code,of a specific kind. See \c CXCompletionChunkKind for a description of the different kinds of chunks.
type CompletionString struct {
	c C.CXCompletionString
}

// Determine the kind of a particular chunk within a completion string. \param completion_string the completion string to query. \param chunk_number the 0-based index of the chunk in the completion string. \returns the kind of the chunk at the index \c chunk_number.
func (cs CompletionString) CompletionChunkKind(chunk_number uint16) CompletionChunkKind {
	return CompletionChunkKind(C.clang_getCompletionChunkKind(cs.c, C.uint(chunk_number)))
}

// Retrieve the text associated with a particular chunk within a completion string. \param completion_string the completion string to query. \param chunk_number the 0-based index of the chunk in the completion string. \returns the text associated with the chunk at index \c chunk_number.
func (cs CompletionString) CompletionChunkText(chunk_number uint16) string {
	o := cxstring{C.clang_getCompletionChunkText(cs.c, C.uint(chunk_number))}
	defer o.Dispose()

	return o.String()
}

// Retrieve the completion string associated with a particular chunk within a completion string. \param completion_string the completion string to query. \param chunk_number the 0-based index of the chunk in the completion string. \returns the completion string associated with the chunk at index \c chunk_number.
func (cs CompletionString) CompletionChunkCompletionString(chunk_number uint16) CompletionString {
	return CompletionString{C.clang_getCompletionChunkCompletionString(cs.c, C.uint(chunk_number))}
}

// Retrieve the annotation associated with the given completion string. \param completion_string the completion string to query. \param annotation_number the 0-based index of the annotation of the completion string. \returns annotation string associated with the completion at index \c annotation_number, or a NULL string if that annotation is not available.
func (cs CompletionString) CompletionAnnotation(annotation_number uint16) string {
	o := cxstring{C.clang_getCompletionAnnotation(cs.c, C.uint(annotation_number))}
	defer o.Dispose()

	return o.String()
}
