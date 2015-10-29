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
