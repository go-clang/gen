package phoenix

// #include "go-clang.h"
import "C"

// A single result of code completion.
type CompletionResult struct {
	c C.CXCompletionResult
}

// The kind of entity that this completion refers to. The cursor kind will be a macro, keyword, or a declaration (one of the *Decl cursor kinds), describing the entity that the completion is referring to. \todo In the future, we would like to provide a full cursor, to allow the client to extract additional information from declaration.
func (cr CompletionResult) CursorKind() CursorKind {
	value := CursorKind(cr.c.CursorKind)
	return value
}

// The code-completion string that describes how to insert this code-completion result into the editing buffer.
func (cr CompletionResult) CompletionString() CompletionString {
	value := CompletionString{cr.c.CompletionString}
	return value
}
