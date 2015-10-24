package phoenix

// #include "go-clang.h"
import "C"

// Source location passed to index callbacks.
type IdxLoc struct {
	c C.CXIdxLoc
}

// Retrieve the CXSourceLocation represented by the given CXIdxLoc.
func (il IdxLoc) IndexLoc_getCXSourceLocation() SourceLocation {
	return SourceLocation{C.clang_indexLoc_getCXSourceLocation(il.c)}
}
