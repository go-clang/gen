package phoenix

// #include "go-clang.h"
import "C"

// A single diagnostic, containing the diagnostic's severity, location, text, source ranges, and fix-it hints.
type Diagnostic struct {
	c C.CXDiagnostic
}

// Format the given diagnostic in a manner that is suitable for display. This routine will format the given diagnostic to a string, rendering the diagnostic according to the various options given. The \c clang_defaultDiagnosticDisplayOptions() function returns the set of options that most closely mimics the behavior of the clang compiler. \param Diagnostic The diagnostic to print. \param Options A set of options that control the diagnostic display, created by combining \c CXDiagnosticDisplayOptions values. \returns A new string containing for formatted diagnostic.
func (d Diagnostic) Format(Options uint16) string {
	o := cxstring{C.clang_formatDiagnostic(d.c, C.uint(Options))}
	defer o.Dispose()

	return o.String()
}

// Retrieve a source range associated with the diagnostic. A diagnostic's source ranges highlight important elements in the source code. On the command line, Clang displays source ranges by underlining them with '~' characters. \param Diagnostic the diagnostic whose range is being extracted. \param Range the zero-based index specifying which range to \returns the requested source range.
func (d Diagnostic) Range(Range uint16) SourceRange {
	return SourceRange{C.clang_getDiagnosticRange(d.c, C.uint(Range))}
}
