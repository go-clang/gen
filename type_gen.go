package phoenix

// #include "go-clang.h"
import "C"

// The type of an element in the abstract syntax tree.
type Type struct {
	c C.CXType
}

// Pretty-print the underlying type using the rules of the language of the translation unit from which it came. If the type is invalid, an empty string is returned.
func (t Type) Spelling() string {
	cstr := cxstring{C.clang_getTypeSpelling(t.c)}
	defer cstr.Dispose()

	return cstr.String()
}
