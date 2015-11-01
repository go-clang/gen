package clang

// #include "go-clang.h"
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"
)

/**
 * \brief Retrieve the parent context of the given completion string.
 *
 * The parent context of a completion string is the semantic parent of
 * the declaration (if any) that the code completion represents. For example,
 * a code completion for an Objective-C method would have the method's class
 * or protocol as its context.
 *
 * \param completion_string The code completion string whose parent is
 * being queried.
 *
 * \param kind DEPRECATED: always set to CXCursor_NotImplemented if non-NULL.
 *
 * \returns The name of the completion parent, e.g., "NSObject" if
 * the completion string represents a method in the NSObject class.
 */
func (cs CompletionString) CompletionParent() string {
	o := cxstring{C.clang_getCompletionParent(cs.c, nil)}
	defer o.Dispose()
	return o.String()
}

/**
- * \brief Contains the results of code-completion.
- *
- * This data structure contains the results of code completion, as
- * produced by \c clang_codeCompleteAt(). Its contents must be freed by
- * \c clang_disposeCodeCompleteResults.
- */
type CodeCompleteResults struct {
	c *C.CXCodeCompleteResults
}

// TODO(): is there a better way to handle this?
func (ccr CodeCompleteResults) Results() (ret []CompletionResult) {
	header := (*reflect.SliceHeader)((unsafe.Pointer(&ret)))
	header.Cap = int(ccr.c.NumResults)
	header.Len = int(ccr.c.NumResults)
	header.Data = uintptr(unsafe.Pointer(ccr.c.Results))
	return
}

/**
 * \brief Free the given set of code-completion results.
 */
func (ccr CodeCompleteResults) Dispose() {
	C.clang_disposeCodeCompleteResults(ccr.c)
}
