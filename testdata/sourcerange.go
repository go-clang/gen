package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

// SourceRange identifies a half-open character range in the source code.
//
// Use clang_getRangeStart() and clang_getRangeEnd() to retrieve the
// starting and end locations from a source range, respectively.
type SourceRange struct {
	c C.CXSourceRange
}

// NewNullRange creates a NULL (invalid) source range.
func NewNullRange() SourceRange {
	return SourceRange{C.clang_getNullRange()}
}

// NewRange creates a source range given the beginning and ending source
// locations.
func NewRange(beg, end SourceLocation) SourceRange {
	o := C.clang_getRange(beg.c, end.c)
	return SourceRange{o}
}
