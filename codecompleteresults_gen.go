package phoenix

// #include "go-clang.h"
import "C"

// Contains the results of code-completion. This data structure contains the results of code completion, as produced by \c clang_codeCompleteAt(). Its contents must be freed by \c clang_disposeCodeCompleteResults.
type CodeCompleteResults struct {
	c C.CXCodeCompleteResults
}
