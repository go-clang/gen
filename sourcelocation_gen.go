package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// Identifies a specific source location within a translation unit. Use clang_getExpansionLocation() or clang_getSpellingLocation() to map a source location to a particular file, line, and column.
type SourceLocation struct {
	c C.CXSourceLocation
}

func (sl SourceLocation) Ptr_data() []unsafe.Pointer {
	sc := []unsafe.Pointer{}

	length := 2
	goslice := (*[1 << 30]*C.void)(unsafe.Pointer(&sl.c.ptr_data))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, unsafe.Pointer(goslice[is]))
	}

	return sc
}

func (sl SourceLocation) Int_data() uint16 {
	value := uint16(sl.c.int_data)
	return value
}

// Determine whether two source locations, which must refer into the same translation unit, refer to exactly the same point in the source code. \returns non-zero if the source locations refer to the same location, zero if they refer to different locations.
func (sl SourceLocation) EqualLocations(sl2 SourceLocation) bool {
	o := C.clang_equalLocations(sl.c, sl2.c)

	return o != C.uint(0)
}

// Retrieve a source range given the beginning and ending source locations.
func (sl SourceLocation) Range(end SourceLocation) SourceRange {
	return SourceRange{C.clang_getRange(sl.c, end.c)}
}
