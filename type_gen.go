package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// The type of an element in the abstract syntax tree.
type Type struct {
	c C.CXType
}

func (t Type) Kind() TypeKind {
	value := TypeKind(t.c.kind)
	return value
}

func (t Type) Data() []unsafe.Pointer {
	sc := []unsafe.Pointer{}

	length := 2
	goslice := (*[1 << 30]*C.void)(unsafe.Pointer(&t.c.data))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, unsafe.Pointer(goslice[is]))
	}

	return sc
}

// Determine whether two CXTypes represent the same type. \returns non-zero if the CXTypes represent the same type and zero otherwise.
func (t Type) EqualTypes(t2 Type) bool {
	o := C.clang_equalTypes(t.c, t2.c)

	return o != C.uint(0)
}

// Retrieve the type of an argument of a function type. If a non-function type is passed in or the function does not have enough parameters, an invalid type is returned.
func (t Type) ArgType(i uint16) Type {
	return Type{C.clang_getArgType(t.c, C.uint(i))}
}

// Return the offset of a field named S in a record of type T in bits as it would be returned by __offsetof__ as per C++11[18.2p4] If the cursor is not a record field declaration, CXTypeLayoutError_Invalid is returned. If the field's type declaration is an incomplete type, CXTypeLayoutError_Incomplete is returned. If the field's type declaration is a dependent type, CXTypeLayoutError_Dependent is returned. If the field's name S is not found, CXTypeLayoutError_InvalidFieldName is returned.
func (t Type) OffsetOf(S string) int64 {
	c_S := C.CString(S)
	defer C.free(unsafe.Pointer(c_S))

	return int64(C.clang_Type_getOffsetOf(t.c, c_S))
}
