package phoenix

// #include "go-clang.h"
import "C"

type IdxEntityInfo struct {
	c C.CXIdxEntityInfo
}
