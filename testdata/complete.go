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
 * \brief Retrieve the annotation associated with the given completion string.
 *
 * \param completion_string the completion string to query.
 *
 * \param annotation_number the 0-based index of the annotation of the
 * completion string.
 *
 * \returns annotation string associated with the completion at index
 * \c annotation_number, or a NULL string if that annotation is not available.
 */

func (cs CompletionString) Chunks() (ret []CompletionChunk) {
	ret = make([]CompletionChunk, C.clang_getNumCompletionChunks(cs.c))
	for i := range ret {
		ret[i].cs = cs.c
		ret[i].number = C.uint(i)
	}
	return
}

type CompletionChunk struct {
	cs     C.CXCompletionString
	number C.uint
}

func (cc CompletionChunk) String() string {
	return fmt.Sprintf("%s %s", cc.Kind(), cc.Text())
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

func (ccr CodeCompleteResults) IsValid() bool {
	return ccr.c != nil
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
 * \brief Sort the code-completion results in case-insensitive alphabetical
 * order.
 *
 * \param Results The set of results to sort.
 * \param NumResults The number of results in \p Results.
 */
func (ccr CodeCompleteResults) Sort() {
	C.clang_sortCodeCompletionResults(ccr.c.Results, ccr.c.NumResults)
}

/**
 * \brief Free the given set of code-completion results.
 */
func (ccr CodeCompleteResults) Dispose() {
	C.clang_disposeCodeCompleteResults(ccr.c)
}

/**
 * \brief Retrieve a diagnostic associated with the given code completion.
 *
 * \param Results the code completion results to query.
 * \param Index the zero-based diagnostic number to retrieve.
 *
 * \returns the requested diagnostic. This diagnostic must be freed
 * via a call to \c clang_disposeDiagnostic().
 */
func (ccr CodeCompleteResults) Diagnostics() (ret Diagnostics) {
	ret = make(Diagnostics, C.clang_codeCompleteGetNumDiagnostics(ccr.c))
	for i := range ret {
		ret[i].c = C.clang_codeCompleteGetDiagnostic(ccr.c, C.uint(i))
	}
	return
}
