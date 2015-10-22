package phoenix

// #include "go-clang.h"
import "C"

// An "index" that consists of a set of translation units that would typically be linked together into an executable or library.
type Index struct {
	c C.CXIndex
}

// Destroy the given index. The index must not be destroyed until all of the translation units created within that index have been destroyed.
func (i Index) Dispose() {
	C.clang_disposeIndex(i.c)
}
