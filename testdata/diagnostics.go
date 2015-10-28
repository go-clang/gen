package clang

// #include "go-clang.h"
import "C"

/**
 * \brief A single diagnostic, containing the diagnostic's severity,
 * location, text, source ranges, and fix-it hints.
 */
type Diagnostic struct {
	c C.CXDiagnostic
}

type Diagnostics []Diagnostic

func (d Diagnostics) Dispose() {
	for _, di := range d {
		di.Dispose()
	}
}

/**
 * \brief Retrieve a source range associated with the diagnostic.
 *
 * A diagnostic's source ranges highlight important elements in the source
 * code. On the command line, Clang displays source ranges by
 * underlining them with '~' characters.
 *
 * \param Diagnostic the diagnostic whose range is being extracted.
 *
 * \param Range the zero-based index specifying which range to
 *
 * \returns the requested source range.
 */
func (d Diagnostic) Ranges() (ret []SourceRange) {
	ret = make([]SourceRange, C.clang_getDiagnosticNumRanges(d.c))
	for i := range ret {
		ret[i].c = C.clang_getDiagnosticRange(d.c, C.uint(i))
	}
	return
}

type FixIt struct {
	Data             string
	ReplacementRange SourceRange
}

/**
 * \brief Retrieve the replacement information for a given fix-it.
 *
 * Fix-its are described in terms of a source range whose contents
 * should be replaced by a string. This approach generalizes over
 * three kinds of operations: removal of source code (the range covers
 * the code to be removed and the replacement string is empty),
 * replacement of source code (the range covers the code to be
 * replaced and the replacement string provides the new code), and
 * insertion (both the start and end of the range point at the
 * insertion location, and the replacement string provides the text to
 * insert).
 *
 * \param Diagnostic The diagnostic whose fix-its are being queried.
 *
 * \param FixIt The zero-based index of the fix-it.
 *
 * \param ReplacementRange The source range whose contents will be
 * replaced with the returned replacement string. Note that source
 * ranges are half-open ranges [a, b), so the source code should be
 * replaced from a and up to (but not including) b.
 *
 * \returns A string containing text that should be replace the source
 * code indicated by the \c ReplacementRange.
 */
func (d Diagnostic) FixIts() (ret []FixIt) {
	ret = make([]FixIt, C.clang_getDiagnosticNumFixIts(d.c))
	for i := range ret {
		cx := cxstring{C.clang_getDiagnosticFixIt(d.c, C.uint(i), &ret[i].ReplacementRange.c)}
		defer cx.Dispose()
		ret[i].Data = cx.String()
	}
	return
}

func (d Diagnostic) String() string {
	return d.Format(DefaultDiagnosticDisplayOptions())
}

/**
 * \brief Retrieve the set of display options most similar to the
 * default behavior of the clang compiler.
 *
 * \returns A set of display options suitable for use with \c
 * clang_displayDiagnostic().
 */
func DefaultDiagnosticDisplayOptions() DiagnosticDisplayOptions {
	return DiagnosticDisplayOptions(C.clang_defaultDiagnosticDisplayOptions())
}
