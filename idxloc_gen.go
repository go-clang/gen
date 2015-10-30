package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// Source location passed to index callbacks.
type IdxLoc struct {
	c C.CXIdxLoc
}

func (il IdxLoc) Ptr_data() []unsafe.Pointer {
	sc := []unsafe.Pointer{}

	length := 2
	goslice := (*[1 << 30]*C.void)(unsafe.Pointer(&il.c.ptr_data))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, unsafe.Pointer(goslice[is]))
	}

	return sc
}

func (il IdxLoc) Int_data() uint16 {
	value := uint16(il.c.int_data)
	return value
}

// Retrieve the CXSourceLocation represented by the given CXIdxLoc.
func (il IdxLoc) IndexLoc_getCXSourceLocation() SourceLocation {
	return SourceLocation{C.clang_indexLoc_getCXSourceLocation(il.c)}
}
