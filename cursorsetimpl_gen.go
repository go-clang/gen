package phoenix

// #include "go-clang.h"
import "C"

type CursorSetImpl struct {
	c C.CXCursorSetImpl
}
