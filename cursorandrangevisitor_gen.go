package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

type CursorAndRangeVisitor struct {
	c C.CXCursorAndRangeVisitor
}

func (carv CursorAndRangeVisitor) Context() unsafe.Pointer {
	value := unsafe.Pointer(carv.c.context)
	return value
}
