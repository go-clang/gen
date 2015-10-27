package phoenix

// #include "go-clang.h"
import "C"

// A fast container representing a set of CXCursors.
type CursorSet struct {
	c C.CXCursorSet
}

// Creates an empty CXCursorSet.
func NewCursorSet() CursorSet {
	return CursorSet{C.clang_createCXCursorSet()}
}

// Disposes a CXCursorSet and releases its associated memory.
func (cs CursorSet) Dispose() {
	C.clang_disposeCXCursorSet(cs.c)
}
