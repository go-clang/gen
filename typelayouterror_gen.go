package phoenix

// #include "go-clang.h"
import "C"

// List the possible error codes for \c clang_Type_getSizeOf, \c clang_Type_getAlignOf, \c clang_Type_getOffsetOf and \c clang_Cursor_getOffsetOf. A value of this enumeration type can be returned if the target type is not a valid argument to sizeof, alignof or offsetof.
type TypeLayoutError int32

const (
	// Type is of kind CXType_Invalid.
	TypeLayoutError_Invalid TypeLayoutError = C.CXTypeLayoutError_Invalid
	// The type is an incomplete Type.
	TypeLayoutError_Incomplete = C.CXTypeLayoutError_Incomplete
	// The type is a dependent Type.
	TypeLayoutError_Dependent = C.CXTypeLayoutError_Dependent
	// The type is not a constant size type.
	TypeLayoutError_NotConstantSize = C.CXTypeLayoutError_NotConstantSize
	// The Field name is not valid for this record.
	TypeLayoutError_InvalidFieldName = C.CXTypeLayoutError_InvalidFieldName
)
