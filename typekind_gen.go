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
	Type_Unexposed TypeKind = C.CXType_Unexposed
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Void TypeKind = C.CXType_Void
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Bool TypeKind = C.CXType_Bool
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Char_U TypeKind = C.CXType_Char_U
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_UChar TypeKind = C.CXType_UChar
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Char16 TypeKind = C.CXType_Char16
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Char32 TypeKind = C.CXType_Char32
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_UShort TypeKind = C.CXType_UShort
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_UInt TypeKind = C.CXType_UInt
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ULong TypeKind = C.CXType_ULong
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ULongLong TypeKind = C.CXType_ULongLong
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_UInt128 TypeKind = C.CXType_UInt128
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Char_S TypeKind = C.CXType_Char_S
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_SChar TypeKind = C.CXType_SChar
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_WChar TypeKind = C.CXType_WChar
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Short TypeKind = C.CXType_Short
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Int TypeKind = C.CXType_Int
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Long TypeKind = C.CXType_Long
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_LongLong TypeKind = C.CXType_LongLong
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Int128 TypeKind = C.CXType_Int128
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Float TypeKind = C.CXType_Float
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Double TypeKind = C.CXType_Double
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_LongDouble TypeKind = C.CXType_LongDouble
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_NullPtr TypeKind = C.CXType_NullPtr
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Overload TypeKind = C.CXType_Overload
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Dependent TypeKind = C.CXType_Dependent
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCId TypeKind = C.CXType_ObjCId
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCClass TypeKind = C.CXType_ObjCClass
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCSel TypeKind = C.CXType_ObjCSel
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_FirstBuiltin TypeKind = C.CXType_FirstBuiltin
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_LastBuiltin TypeKind = C.CXType_LastBuiltin
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Complex TypeKind = C.CXType_Complex
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Pointer TypeKind = C.CXType_Pointer
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_BlockPointer TypeKind = C.CXType_BlockPointer
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_LValueReference TypeKind = C.CXType_LValueReference
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_RValueReference TypeKind = C.CXType_RValueReference
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Record TypeKind = C.CXType_Record
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Enum TypeKind = C.CXType_Enum
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Typedef TypeKind = C.CXType_Typedef
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCInterface TypeKind = C.CXType_ObjCInterface
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ObjCObjectPointer TypeKind = C.CXType_ObjCObjectPointer
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_FunctionNoProto TypeKind = C.CXType_FunctionNoProto
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_FunctionProto TypeKind = C.CXType_FunctionProto
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_ConstantArray TypeKind = C.CXType_ConstantArray
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_Vector TypeKind = C.CXType_Vector
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_IncompleteArray TypeKind = C.CXType_IncompleteArray
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_VariableArray TypeKind = C.CXType_VariableArray
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_DependentSizedArray TypeKind = C.CXType_DependentSizedArray
	/**
	 * \brief A type whose specific kind is not exposed via this
	 * interface.
	 */
	Type_MemberPointer TypeKind = C.CXType_MemberPointer
)
