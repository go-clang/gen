package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// A remapping of original source files and their translated files.
type Remapping struct {
	c C.CXRemapping
}

// Retrieve a remapping. \param path the path that contains metadata about remappings. \returns the requested remapping. This remapping must be freed via a call to \c clang_remap_dispose(). Can return NULL if an error occurred.
func NewRemappings(path string) Remapping {
	c_path := C.CString(path)
	defer C.free(unsafe.Pointer(c_path))

	return Remapping{C.clang_getRemappings(c_path)}
}

// Determine the number of remappings.
func (r Remapping) Remap_getNumFiles() uint16 {
	return uint16(C.clang_remap_getNumFiles(r.c))
}

// Get the original and the associated filename from the remapping. \param original If non-NULL, will be set to the original filename. \param transformed If non-NULL, will be set to the filename that the original is associated with.
func (r Remapping) Remap_getFilenames(index uint16) (string, string) {
	var original cxstring
	defer original.Dispose()
	var transformed cxstring
	defer transformed.Dispose()

	C.clang_remap_getFilenames(r.c, C.uint(index), &original.c, &transformed.c)

	return original.String(), transformed.String()
}

// Dispose the remapping.
func (r Remapping) Remap_Dispose() {
	C.clang_remap_dispose(r.c)
}
