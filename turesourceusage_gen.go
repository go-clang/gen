package phoenix

// #include "go-clang.h"
import "C"

// The memory usage of a CXTranslationUnit, broken into categories.
type TUResourceUsage struct {
	c C.CXTUResourceUsage
}

func (turu TUResourceUsage) Dispose() {
	C.clang_disposeCXTUResourceUsage(turu.c)
}
