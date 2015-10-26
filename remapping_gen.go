package phoenix

// #include "go-clang.h"
import "C"

// A remapping of original source files and their translated files.
type Remapping struct {
	c C.CXRemapping
}

// Dispose the remapping.
func (r Remapping) remap_Dispose() {
	C.clang_remap_dispose(r.c)
}
