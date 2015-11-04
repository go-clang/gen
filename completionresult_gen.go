package phoenix

// #include "go-clang.h"
import "C"

// A single result of code completion.
type CompletionResult struct {
	c C.CXCompletionResult
}

/*
	The kind of entity that this completion refers to.

	The cursor kind will be a macro, keyword, or a declaration (one of the
	*Decl cursor kinds), describing the entity that the completion is
	referring to.

	\todo In the future, we would like to provide a full cursor, to allow
	the client to extract additional information from declaration.
*/
func (cr CompletionResult) CursorKind() CursorKind {
	value := CursorKind(cr.c.CursorKind)
	return value
}

// The code-completion string that describes how to insert this code-completion result into the editing buffer.
func (cr CompletionResult) CompletionString() CompletionString {
	value := CompletionString{cr.c.CompletionString}
	return value
}

/*
	Sort the code-completion results in case-insensitive alphabetical
	order.

	Parameter Results The set of results to sort.
	Parameter NumResults The number of results in \p Results.
*/
func SortCodeCompletionResults(results []CompletionResult) {
	ca_results := make([]C.CXCompletionResult, len(results))
	var cp_results *C.CXCompletionResult
	if len(results) > 0 {
		cp_results = &ca_results[0]
	}
	for i := range results {
		ca_results[i] = results[i].c
	}

	C.clang_sortCodeCompletionResults(cp_results, C.uint(len(results)))
}
