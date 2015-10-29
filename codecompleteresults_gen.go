package phoenix

// #include "go-clang.h"
import "C"

// Contains the results of code-completion. This data structure contains the results of code completion, as produced by \c clang_codeCompleteAt(). Its contents must be freed by \c clang_disposeCodeCompleteResults.
type CodeCompleteResults struct {
	c C.CXCodeCompleteResults
}

// The code-completion results.
func (ccr CodeCompleteResults) Results() *CompletionResult {
	value := CompletionResult{*ccr.c.Results}
	return &value
}

// The number of code-completion results stored in the \c Results array.
func (ccr CodeCompleteResults) NumResults() uint16 {
	value := uint16(ccr.c.NumResults)
	return value
}
