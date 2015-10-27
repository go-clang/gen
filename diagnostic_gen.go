package phoenix

// #include "go-clang.h"
import "C"

// A single diagnostic, containing the diagnostic's severity, location, text, source ranges, and fix-it hints.
type Diagnostic struct {
	c C.CXDiagnostic
}

// Retrieve the child diagnostics of a CXDiagnostic. This CXDiagnosticSet does not need to be released by clang_disposeDiagnosticSet.
func (d Diagnostic) ChildDiagnostics() DiagnosticSet {
	return DiagnosticSet{C.clang_getChildDiagnostics(d.c)}
}

// Destroy a diagnostic.
func (d Diagnostic) Dispose() {
	C.clang_disposeDiagnostic(d.c)
}

// Determine the severity of the given diagnostic.
func (d Diagnostic) Severity() DiagnosticSeverity {
	return DiagnosticSeverity(C.clang_getDiagnosticSeverity(d.c))
}

// Retrieve the source location of the given diagnostic. This location is where Clang would print the caret ('^') when displaying the diagnostic on the command line.
func (d Diagnostic) Location() SourceLocation {
	return SourceLocation{C.clang_getDiagnosticLocation(d.c)}
}

// Retrieve the text of the given diagnostic.
func (d Diagnostic) Spelling() string {
	o := cxstring{C.clang_getDiagnosticSpelling(d.c)}
	defer o.Dispose()

	return o.String()
}

// Retrieve the category number for this diagnostic. Diagnostics can be categorized into groups along with other, related diagnostics (e.g., diagnostics under the same warning flag). This routine retrieves the category number for the given diagnostic. \returns The number of the category that contains this diagnostic, or zero if this diagnostic is uncategorized.
func (d Diagnostic) Category() uint16 {
	return uint16(C.clang_getDiagnosticCategory(d.c))
}

// Retrieve the diagnostic category text for a given diagnostic. \returns The text of the given diagnostic category.
func (d Diagnostic) CategoryText() string {
	o := cxstring{C.clang_getDiagnosticCategoryText(d.c)}
	defer o.Dispose()

	return o.String()
}

// Determine the number of source ranges associated with the given diagnostic.
func (d Diagnostic) NumRanges() uint16 {
	return uint16(C.clang_getDiagnosticNumRanges(d.c))
}

// Determine the number of fix-it hints associated with the given diagnostic.
func (d Diagnostic) NumFixIts() uint16 {
	return uint16(C.clang_getDiagnosticNumFixIts(d.c))
}
