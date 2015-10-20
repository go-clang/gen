package phoenix

// #include "go-clang.h"
import "C"

type CursorAndRangeVisitor struct {
	c C.CXCursorAndRangeVisitor
}
