package phoenix

// #include "go-clang.h"
import "C"

// Source location passed to index callbacks.
type IdxLoc struct {
	c C.CXIdxLoc
}
