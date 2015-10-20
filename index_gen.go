package phoenix

// #include "go-clang.h"
import "C"

// An "index" that consists of a set of translation units that would typically be linked together into an executable or library.
type Index struct {
	c C.CXIndex
}
