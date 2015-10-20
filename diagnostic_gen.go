package phoenix

// #include "go-clang.h"
import "C"

// A single diagnostic, containing the diagnostic's severity, location, text, source ranges, and fix-it hints.
type Diagnostic struct {
	c C.CXDiagnostic
}

// Retrieve the text of the given diagnostic.
func (d Diagnostic) Spelling() string {
	cstr := cxstring{C.clang_getDiagnosticSpelling(d.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// Retrieve the diagnostic category text for a given diagnostic. \returns The text of the given diagnostic category.
func (d Diagnostic) CategoryText() string {
	cstr := cxstring{C.clang_getDiagnosticCategoryText(d.c)}
	defer cstr.Dispose()

	return cstr.String()
}
