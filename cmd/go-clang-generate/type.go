package main

import (
	"fmt"

	"github.com/sbinet/go-clang"
)

const (
	GoByte      = "byte"
	GoInt8      = "int8"
	GoUInt8     = "uint8"
	GoInt16     = "int16"
	GoUInt16    = "uint16"
	GoInt32     = "int32"
	GoUInt32    = "uint32"
	GoInt64     = "int64"
	GoUInt64    = "uint64"
	GoFloat32   = "float32"
	GoFloat64   = "float64"
	GoBool      = "bool"
	GoInterface = "interface"
	GoPointer   = "unsafe.Pointer"

	CChar      = "char"      // byte
	CSChar     = "schar"     // int8
	CUChar     = "uchar"     // uint8
	CShort     = "short"     // int16
	CUShort    = "ushort"    // uint16
	CInt       = "int"       // int
	CUInt      = "uint"      // uint32
	CLongInt   = "long"      // int32 or int64
	CULongInt  = "ulong"     // uint32 or uint64
	CLongLong  = "longlong"  // int64
	CULongLong = "ulonglong" // uint64
	CFloat     = "float"     // float32
	CDouble    = "double"    // float64

)

type Type struct {
	Original  string
	Name      string
	CName     string
	Primitive string

	PointerLevel      int
	IsPrimitive       bool
	IsArray           bool
	ArraySize         int64
	IsEnumLiteral     bool
	IsFunctionPointer bool
	IsReturnArgument  bool
	IsSlice           bool
	LengthOfSlice     string
}

func getType(cType clang.Type) (Type, error) {
	typ := Type{
		CName:    cType.TypeSpelling(),
		Original: cType.TypeSpelling(),

		PointerLevel:      0,
		IsPrimitive:       true,
		IsArray:           false,
		ArraySize:         -1,
		IsEnumLiteral:     false,
		IsFunctionPointer: false,
	}

	switch cType.Kind() {
	case clang.TK_Char_S:
		typ.CName = CSChar
		typ.Name = GoInt8
		typ.Primitive = "schar"
	case clang.TK_Char_U:
		typ.CName = CUChar
		typ.Name = GoUInt8
		typ.Primitive = "uchar"
	case clang.TK_Int:
		typ.CName = CInt
		typ.Name = GoInt16
		typ.Primitive = "int"
	case clang.TK_Short:
		typ.CName = CShort
		typ.Name = GoInt16
		typ.Primitive = "int"
	case clang.TK_UShort:
		typ.CName = CUShort
		typ.Name = GoUInt16
		typ.Primitive = "uint"
	case clang.TK_UInt:
		typ.CName = CUInt
		typ.Name = GoUInt16
		typ.Primitive = "uint"
	case clang.TK_Long:
		typ.CName = CLongInt
		typ.Name = GoInt32
		typ.Primitive = "long"
	case clang.TK_ULong:
		typ.CName = CULongInt
		typ.Name = GoUInt32
		typ.Primitive = "ulong"
	case clang.TK_LongLong:
		typ.CName = CLongLong
		typ.Name = GoInt64
		typ.Primitive = "longlong"
	case clang.TK_ULongLong:
		typ.CName = CULongLong
		typ.Name = GoUInt64
		typ.Primitive = "ulonglong"
	case clang.TK_Float:
		typ.CName = CFloat
		typ.Name = GoFloat32
		typ.Primitive = "float"
	case clang.TK_Double:
		typ.CName = CDouble
		typ.Name = GoFloat64
		typ.Primitive = "double"
	case clang.TK_Bool:
		typ.Name = GoBool
	case clang.TK_Void:
		typ.CName = "void"
		typ.Name = "void"
		typ.Primitive = "void"
	case clang.TK_ConstantArray:
		subTyp, err := getType(cType.ArrayElementType())
		if err != nil {
			return Type{}, err
		}

		typ.CName = subTyp.CName
		typ.Name = subTyp.Name
		typ.Primitive = subTyp.Primitive
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsArray = true
		typ.ArraySize = cType.ArraySize()
	case clang.TK_Typedef:
		typeStr := cType.TypeSpelling()
		if typeStr == "CXString" {
			typeStr = "cxstring"
		} else if typeStr == "time_t" {
			typ.Primitive = typeStr
			typeStr = "time.Time"
		} else {
			typeStr = trimClangPrefix(cType.Declaration().Type().TypeSpelling())
		}

		typ.CName = cType.Declaration().Type().TypeSpelling()
		typ.Name = typeStr
		typ.IsPrimitive = false

		if cType.CanonicalType().Kind() == clang.TK_Enum {
			typ.IsEnumLiteral = true
			typ.IsPrimitive = true
		}
	case clang.TK_Pointer:
		typ.PointerLevel++

		if cType.PointeeType().CanonicalType().Kind() == clang.TK_FunctionProto {
			typ.IsFunctionPointer = true
		}

		subTyp, err := getType(cType.PointeeType())
		if err != nil {
			return Type{}, err
		}

		typ.CName = subTyp.CName
		typ.Name = subTyp.Name
		typ.Primitive = subTyp.Primitive
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsPrimitive = subTyp.IsPrimitive
	case clang.TK_Record:
		typ.CName = cType.Declaration().Type().TypeSpelling()
		typ.Name = trimClangPrefix(typ.CName)
		typ.IsPrimitive = false
	case clang.TK_FunctionProto:
		typ.IsFunctionPointer = true
		typ.CName = cType.Declaration().Type().TypeSpelling()
		typ.Name = trimClangPrefix(typ.CName)
	case clang.TK_Enum:
		typ.Name = trimClangPrefix(cType.Declaration().DisplayName())
		typ.IsEnumLiteral = true
		typ.IsPrimitive = true
	case clang.TK_Unexposed: // there is a bug in clang for enums the kind is set to unexposed dunno why, bug persists since 2013
		subTyp, err := getType(cType.CanonicalType())
		if err != nil {
			return Type{}, err
		}

		typ.CName = subTyp.CName
		typ.Name = subTyp.Name
		typ.Primitive = subTyp.Primitive
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsPrimitive = subTyp.IsPrimitive
	default:
		return Type{}, fmt.Errorf("unhandled type %q of kind %q", cType.TypeSpelling(), cType.Kind().Spelling())
	}

	return typ, nil
}
