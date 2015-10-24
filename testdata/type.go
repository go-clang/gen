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

/**
 * \brief Return the the array size of a constant array.
 *
 * If a non-array type is passed in, -1 is returned.
 */
func (t Type) ArraySize() int64 {
	o := C.clang_getArraySize(t.c)
	return int64(o)
}

/**
 * \brief Return the alignment of a type in bytes as per C++[expr.alignof]
 *   standard.
 *
 * If the type declaration is invalid, CXTypeLayoutError_Invalid is returned.
 * If the type declaration is an incomplete type, CXTypeLayoutError_Incomplete
 *   is returned.
 * If the type declaration is a dependent type, CXTypeLayoutError_Dependent is
 *   returned.
 * If the type declaration is not a constant size type,
 *   CXTypeLayoutError_NotConstantSize is returned.
 */
func (t Type) AlignOf() (int, error) {
	o := C.clang_Type_getAlignOf(t.c)
	if o < 0 {
		return int(o), TypeLayoutError(o)
	}
	return int(o), nil
}

/**
 * \brief Return the size of a type in bytes as per C++[expr.sizeof] standard.
 *
 * If the type declaration is invalid, CXTypeLayoutError_Invalid is returned.
 * If the type declaration is an incomplete type, CXTypeLayoutError_Incomplete
 *   is returned.
 * If the type declaration is a dependent type, CXTypeLayoutError_Dependent is
 *   returned.
 */
func (t Type) SizeOf() (int, error) {
	o := C.clang_Type_getSizeOf(t.c)
	if o < 0 {
		return int(o), TypeLayoutError(o)
	}
	return int(o), nil
}

/**
 * \brief Return the offset of a field named S in a record of type T in bits
 *   as it would be returned by __offsetof__ as per C++11[18.2p4]
 *
 * If the cursor is not a record field declaration, CXTypeLayoutError_Invalid
 *   is returned.
 * If the field's type declaration is an incomplete type,
 *   CXTypeLayoutError_Incomplete is returned.
 * If the field's type declaration is a dependent type,
 *   CXTypeLayoutError_Dependent is returned.
 * If the field's name S is not found,
 *   CXTypeLayoutError_InvalidFieldName is returned.
 */
func (t Type) OffsetOf(s string) (int, error) {
	c_str := C.CString(s)
	defer C.free(unsafe.Pointer(c_str))
	o := C.clang_Type_getOffsetOf(t.c, c_str)
	if o < 0 {
		return int(o), TypeLayoutError(o)
	}
	return int(o), nil
}
