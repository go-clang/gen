package phoenix

// #include "go-clang.h"
import "C"

// A single diagnostic, containing the diagnostic's severity, location, text, source ranges, and fix-it hints.
type Diagnostic struct {
	c C.CXDiagnostic
}

// Destroy a diagnostic.
func (d Diagnostic) Dispose() {
	C.clang_disposeDiagnostic(d.c)
}

// Retrieve the text of the given diagnostic.
func (d Diagnostic) Spelling() string {
	o := cxstring{C.clang_getDiagnosticSpelling(d.c)}
	defer o.Dispose()

	return o.String()
}

// Retrieve the diagnostic category text for a given diagnostic. \returns The text of the given diagnostic category.
func (d Diagnostic) CategoryText() string {
	o := cxstring{C.clang_getDiagnosticCategoryText(d.c)}
	defer o.Dispose()

	return o.String()
}
