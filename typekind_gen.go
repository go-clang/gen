package phoenix

// #include "go-clang.h"
import "C"

/**
 * \brief Describes the kind of type
 */
type TypeKind int

const (
	/**
	 * \brief Reprents an invalid type (e.g., where no type is available).
	 */
	Type_Invalid TypeKind = C.CXType_Invalid
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Unexposed = C.CXType_Unexposed
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Void = C.CXType_Void
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Bool = C.CXType_Bool
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Char_U = C.CXType_Char_U
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_UChar = C.CXType_UChar
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Char16 = C.CXType_Char16
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Char32 = C.CXType_Char32
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_UShort = C.CXType_UShort
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_UInt = C.CXType_UInt
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ULong = C.CXType_ULong
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ULongLong = C.CXType_ULongLong
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_UInt128 = C.CXType_UInt128
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Char_S = C.CXType_Char_S
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_SChar = C.CXType_SChar
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_WChar = C.CXType_WChar
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Short = C.CXType_Short
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Int = C.CXType_Int
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Long = C.CXType_Long
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_LongLong = C.CXType_LongLong
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Int128 = C.CXType_Int128
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Float = C.CXType_Float
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Double = C.CXType_Double
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_LongDouble = C.CXType_LongDouble
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_NullPtr = C.CXType_NullPtr
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Overload = C.CXType_Overload
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Dependent = C.CXType_Dependent
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCId = C.CXType_ObjCId
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCClass = C.CXType_ObjCClass
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCSel = C.CXType_ObjCSel
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_FirstBuiltin = C.CXType_FirstBuiltin
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_LastBuiltin = C.CXType_LastBuiltin
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Complex = C.CXType_Complex
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Pointer = C.CXType_Pointer
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_BlockPointer = C.CXType_BlockPointer
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_LValueReference = C.CXType_LValueReference
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_RValueReference = C.CXType_RValueReference
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Record = C.CXType_Record
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Enum = C.CXType_Enum
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Typedef = C.CXType_Typedef
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCInterface = C.CXType_ObjCInterface
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCObjectPointer = C.CXType_ObjCObjectPointer
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_FunctionNoProto = C.CXType_FunctionNoProto
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_FunctionProto = C.CXType_FunctionProto
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ConstantArray = C.CXType_ConstantArray
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Vector = C.CXType_Vector
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_IncompleteArray = C.CXType_IncompleteArray
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_VariableArray = C.CXType_VariableArray
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_DependentSizedArray = C.CXType_DependentSizedArray
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_MemberPointer = C.CXType_MemberPointer
)
