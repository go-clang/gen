package phoenix

// #include "go-clang.h"
import "C"

// A fast container representing a set of CXCursors.
type CursorSet struct {
	c C.CXCursorSet
}

// Queries a CXCursorSet to see if it contains a specific CXCursor. \returns non-zero if the set contains the specified cursor.
func (cs CursorSet) Contains(cursor Cursor) uint16 {
	return uint16(C.clang_CXCursorSet_contains(cs.c, cursor.c))
}

// Inserts a CXCursor into a CXCursorSet. \returns zero if the CXCursor was already in the set, and non-zero otherwise.
func (cs CursorSet) Insert(cursor Cursor) uint16 {
	return uint16(C.clang_CXCursorSet_insert(cs.c, cursor.c))
}
