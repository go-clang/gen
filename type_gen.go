package phoenix

// #include "go-clang.h"
import "C"

// The type of an element in the abstract syntax tree.
type Type struct {
	c C.CXType
}

// Pretty-print the underlying type using the rules of the language of the translation unit from which it came. If the type is invalid, an empty string is returned.
func (t Type) Spelling() string {
	o := cxstring{C.clang_getTypeSpelling(t.c)}
	defer o.Dispose()

	return o.String()
}

// Determine whether a CXType has the "const" qualifier set, without looking through typedefs that may have added "const" at a different level.
func (t Type) IsConstQualifiedType() bool {
	o := C.clang_isConstQualifiedType(t.c)

	return o != C.uint(0)
}

// Determine whether a CXType has the "volatile" qualifier set, without looking through typedefs that may have added "volatile" at a different level.
func (t Type) IsVolatileQualifiedType() bool {
	o := C.clang_isVolatileQualifiedType(t.c)

	return o != C.uint(0)
}

// Determine whether a CXType has the "restrict" qualifier set, without looking through typedefs that may have added "restrict" at a different level.
func (t Type) IsRestrictQualifiedType() bool {
	o := C.clang_isRestrictQualifiedType(t.c)

	return o != C.uint(0)
}

// Return 1 if the CXType is a variadic function type, and 0 otherwise.
func (t Type) IsFunctionTypeVariadic() bool {
	o := C.clang_isFunctionTypeVariadic(t.c)

	return o != C.uint(0)
}

// Return 1 if the CXType is a POD (plain old data) type, and 0 otherwise.
func (t Type) IsPODType() bool {
	o := C.clang_isPODType(t.c)

	return o != C.uint(0)
}
