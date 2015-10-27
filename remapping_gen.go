package phoenix

// #include "go-clang.h"
import "C"

// A remapping of original source files and their translated files.
type Remapping struct {
	c C.CXRemapping
}

// Determine the number of remappings.
func (r Remapping) Remap_getNumFiles() uint16 {
	return uint16(C.clang_remap_getNumFiles(r.c))
}

// Dispose the remapping.
func (r Remapping) remap_Dispose() {
	C.clang_remap_dispose(r.c)
}
