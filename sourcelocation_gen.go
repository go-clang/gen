package phoenix

// #include "go-clang.h"
import "C"

// Identifies a specific source location within a translation unit. Use clang_getExpansionLocation() or clang_getSpellingLocation() to map a source location to a particular file, line, and column.
type SourceLocation struct {
	c C.CXSourceLocation
}

// Retrieve a NULL (invalid) source location.
func NewNullLocation() SourceLocation {
	return SourceLocation{C.clang_getNullLocation()}
}

// Determine whether two source locations, which must refer into the same translation unit, refer to exactly the same point in the source code. \returns non-zero if the source locations refer to the same location, zero if they refer to different locations.
func (sl SourceLocation) EqualLocations(sl2 SourceLocation) bool {
	o := C.clang_equalLocations(sl.c, sl2.c)

	return o != C.uint(0)
}

// Returns non-zero if the given source location is in a system header.
func (sl SourceLocation) Location_IsInSystemHeader() bool {
	o := C.clang_Location_isInSystemHeader(sl.c)

	return o != C.int(0)
}

// Returns non-zero if the given source location is in the main file of the corresponding translation unit.
func (sl SourceLocation) Location_IsFromMainFile() bool {
	o := C.clang_Location_isFromMainFile(sl.c)

	return o != C.int(0)
}

// Retrieve a source range given the beginning and ending source locations.
func (sl SourceLocation) Range(end SourceLocation) SourceRange {
	return SourceRange{C.clang_getRange(sl.c, end.c)}
}
