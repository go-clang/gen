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

// Format the given diagnostic in a manner that is suitable for display. This routine will format the given diagnostic to a string, rendering the diagnostic according to the various options given. The \c clang_defaultDiagnosticDisplayOptions() function returns the set of options that most closely mimics the behavior of the clang compiler. \param Diagnostic The diagnostic to print. \param Options A set of options that control the diagnostic display, created by combining \c CXDiagnosticDisplayOptions values. \returns A new string containing for formatted diagnostic.
func (d Diagnostic) Format(Options uint16) string {
	o := cxstring{C.clang_formatDiagnostic(d.c, C.uint(Options))}
	defer o.Dispose()

	return o.String()
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

// Retrieve the name of the command-line option that enabled this diagnostic. \param Diag The diagnostic to be queried. \param Disable If non-NULL, will be set to the option that disables this diagnostic (if any). \returns A string that contains the command-line option used to enable this warning, such as "-Wconversion" or "-pedantic".
func (d Diagnostic) Option() (string, string) {
	var Disable cxstring
	defer Disable.Dispose()

	o := cxstring{C.clang_getDiagnosticOption(d.c, &Disable.c)}
	defer o.Dispose()

	return Disable.String(), o.String()
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

// Retrieve a source range associated with the diagnostic. A diagnostic's source ranges highlight important elements in the source code. On the command line, Clang displays source ranges by underlining them with '~' characters. \param Diagnostic the diagnostic whose range is being extracted. \param Range the zero-based index specifying which range to \returns the requested source range.
func (d Diagnostic) Range(Range uint16) SourceRange {
	return SourceRange{C.clang_getDiagnosticRange(d.c, C.uint(Range))}
}

// Determine the number of fix-it hints associated with the given diagnostic.
func (d Diagnostic) NumFixIts() uint16 {
	return uint16(C.clang_getDiagnosticNumFixIts(d.c))
}

// Retrieve the replacement information for a given fix-it. Fix-its are described in terms of a source range whose contents should be replaced by a string. This approach generalizes over three kinds of operations: removal of source code (the range covers the code to be removed and the replacement string is empty), replacement of source code (the range covers the code to be replaced and the replacement string provides the new code), and insertion (both the start and end of the range point at the insertion location, and the replacement string provides the text to insert). \param Diagnostic The diagnostic whose fix-its are being queried. \param FixIt The zero-based index of the fix-it. \param ReplacementRange The source range whose contents will be replaced with the returned replacement string. Note that source ranges are half-open ranges [a, b), so the source code should be replaced from a and up to (but not including) b. \returns A string containing text that should be replace the source code indicated by the \c ReplacementRange.
func (d Diagnostic) FixIt(FixIt uint16) (SourceRange, string) {
	var ReplacementRange SourceRange

	o := cxstring{C.clang_getDiagnosticFixIt(d.c, C.uint(FixIt), &ReplacementRange.c)}
	defer o.Dispose()

	return ReplacementRange, o.String()
}
