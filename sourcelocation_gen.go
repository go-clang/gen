package phoenix

// #include "go-clang.h"
import "C"

// Identifies a specific source location within a translation unit. Use clang_getExpansionLocation() or clang_getSpellingLocation() to map a source location to a particular file, line, and column.
type SourceLocation struct {
	c C.CXSourceLocation
}

// Determine whether two source locations, which must refer into the same translation unit, refer to exactly the same point in the source code. \returns non-zero if the source locations refer to the same location, zero if they refer to different locations.
func EqualLocations(sl1, sl2 SourceLocation) bool {
	o := C.clang_equalLocations(sl1.c, sl2.c)

	return o != C.uint(0)
}
