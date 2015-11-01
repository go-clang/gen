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

// Free the given set of code-completion results.
func (ccr *CodeCompleteResults) Dispose() {
	C.clang_disposeCodeCompleteResults(&ccr.c)
}

// Determine the number of diagnostics produced prior to the location where code completion was performed.
func (ccr *CodeCompleteResults) CodeCompleteGetNumDiagnostics() uint16 {
	return uint16(C.clang_codeCompleteGetNumDiagnostics(&ccr.c))
}

// Retrieve a diagnostic associated with the given code completion. \param Results the code completion results to query. \param Index the zero-based diagnostic number to retrieve. \returns the requested diagnostic. This diagnostic must be freed via a call to \c clang_disposeDiagnostic().
func (ccr *CodeCompleteResults) CodeCompleteGetDiagnostic(Index uint16) Diagnostic {
	return Diagnostic{C.clang_codeCompleteGetDiagnostic(&ccr.c, C.uint(Index))}
}

// Determines what compeltions are appropriate for the context the given code completion. \param Results the code completion results to query \returns the kinds of completions that are appropriate for use along with the given code completion results.
func (ccr *CodeCompleteResults) CodeCompleteGetContexts() uint64 {
	return uint64(C.clang_codeCompleteGetContexts(&ccr.c))
}

// Returns the cursor kind for the container for the current code completion context. The container is only guaranteed to be set for contexts where a container exists (i.e. member accesses or Objective-C message sends); if there is not a container, this function will return CXCursor_InvalidCode. \param Results the code completion results to query \param IsIncomplete on return, this value will be false if Clang has complete information about the container. If Clang does not have complete information, this value will be true. \returns the container kind, or CXCursor_InvalidCode if there is not a container
func (ccr *CodeCompleteResults) CodeCompleteGetContainerKind() (uint16, CursorKind) {
	var IsIncomplete C.uint

	o := CursorKind(C.clang_codeCompleteGetContainerKind(&ccr.c, &IsIncomplete))

	return uint16(IsIncomplete), o
}

// Returns the USR for the container for the current code completion context. If there is not a container for the current context, this function will return the empty string. \param Results the code completion results to query \returns the USR for the container
func (ccr *CodeCompleteResults) CodeCompleteGetContainerUSR() string {
	o := cxstring{C.clang_codeCompleteGetContainerUSR(&ccr.c)}
	defer o.Dispose()

	return o.String()
}

// Returns the currently-entered selector for an Objective-C message send, formatted like "initWithFoo:bar:". Only guaranteed to return a non-empty string for CXCompletionContext_ObjCInstanceMessage and CXCompletionContext_ObjCClassMessage. \param Results the code completion results to query \returns the selector (or partial selector) that has been entered thus far for an Objective-C message send.
func (ccr *CodeCompleteResults) CodeCompleteGetObjCSelector() string {
	o := cxstring{C.clang_codeCompleteGetObjCSelector(&ccr.c)}
	defer o.Dispose()

	return o.String()
}
