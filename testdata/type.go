package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// Type represents the type of an element in the abstract syntax tree.
type Type struct {
	c C.CXType
}

func (c Type) Kind() TypeKind {
	return TypeKind(c.c.kind)
}
