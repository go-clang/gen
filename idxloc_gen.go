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

// Retrieve the CXIdxFile, file, line, column, and offset represented by the given CXIdxLoc. If the location refers into a macro expansion, retrieves the location of the macro expansion and if it refers into a macro argument retrieves the location of the argument.
func (il IdxLoc) IndexLoc_getFileLocation() (IdxClientFile, File, uint16, uint16, uint16) {
	var indexFile IdxClientFile
	var file File
	var line C.uint
	var column C.uint
	var offset C.uint

	C.clang_indexLoc_getFileLocation(il.c, &indexFile.c, &file.c, &line, &column, &offset)

	return indexFile, file, uint16(line), uint16(column), uint16(offset)
}

// Retrieve the CXSourceLocation represented by the given CXIdxLoc.
func (il IdxLoc) IndexLoc_getCXSourceLocation() SourceLocation {
	return SourceLocation{C.clang_indexLoc_getCXSourceLocation(il.c)}
}
