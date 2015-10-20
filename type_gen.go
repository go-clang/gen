package phoenix

// #include "go-clang.h"
import "C"

// The type of an element in the abstract syntax tree.
type Type struct {
	c C.CXType
}
