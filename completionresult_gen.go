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

// Sort the code-completion results in case-insensitive alphabetical order. \param Results The set of results to sort. \param NumResults The number of results in \p Results.
func SortCodeCompletion(Results []CompletionResult) {
	ca_Results := make([]C.CXCompletionResult, len(Results))
	var cp_Results *C.CXCompletionResult
	if len(Results) > 0 {
		cp_Results = &ca_Results[0]
	}
	for i := range Results {
		ca_Results[i] = Results[i].c
	}

	C.clang_sortCodeCompletionResults(cp_Results, C.uint(len(Results)))
}
