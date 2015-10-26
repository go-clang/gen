package phoenix

// #include "go-clang.h"
import "C"

// Identifies a half-open character range in the source code. Use clang_getRangeStart() and clang_getRangeEnd() to retrieve the starting and end locations from a source range, respectively.
type SourceRange struct {
	c C.CXSourceRange
}

// Retrieve a NULL (invalid) source range.
func NewNullRange() SourceRange {
	return SourceRange{C.clang_getNullRange()}
}

// Determine whether two ranges are equivalent. \returns non-zero if the ranges are the same, zero if they differ.
func EqualRanges(sr1, sr2 SourceRange) bool {
	o := C.clang_equalRanges(sr1.c, sr2.c)

	return o != C.uint(0)
}

// Returns non-zero if \p range is null.
func (sr SourceRange) Range_IsNull() bool {
	o := C.clang_Range_isNull(sr.c)

	return o != C.int(0)
}

// Retrieve a source location representing the first character within a source range.
func (sr SourceRange) RangeStart() SourceLocation {
	return SourceLocation{C.clang_getRangeStart(sr.c)}
}

// Retrieve a source location representing the last character within a source range.
func (sr SourceRange) RangeEnd() SourceLocation {
	return SourceLocation{C.clang_getRangeEnd(sr.c)}
}
