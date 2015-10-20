package phoenix

// #include "go-clang.h"
import "C"

// The memory usage of a CXTranslationUnit, broken into categories.
type TUResourceUsage struct {
	c C.CXTUResourceUsage
}
