package phoenix

// #include "go-clang.h"
import "C"

// Identifies a half-open character range in the source code. Use clang_getRangeStart() and clang_getRangeEnd() to retrieve the starting and end locations from a source range, respectively.
type SourceRange struct {
	c C.CXSourceRange
}

// Determine whether two ranges are equivalent. \returns non-zero if the ranges are the same, zero if they differ.
func EqualRanges(sr1, sr2 SourceRange) bool {
	o := C.clang_equalRanges(sr1.c, sr2.c)

	return o != C.uint(0)
}
