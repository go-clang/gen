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

	CChar      = "char"
	CSChar     = "schar"
	CUChar     = "uchar"
	CShort     = "short"
	CUShort    = "ushort"
	CInt       = "int"
	CUInt      = "uint"
	CLongInt   = "long"
	CULongInt  = "ulong"
	CLongLong  = "longlong"
	CULongLong = "ulonglong"
	CFloat     = "float"
	CDouble    = "double"
)

type Type struct {
	CName   string
	CGoName string
	GoName  string

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
		CName: cType.TypeSpelling(),

		PointerLevel:      0,
		IsPrimitive:       true,
		IsArray:           false,
		ArraySize:         -1,
		IsEnumLiteral:     false,
		IsFunctionPointer: false,
	}

	switch cType.Kind() {
	case clang.TK_Char_S:
		typ.CGoName = CSChar
		typ.GoName = GoInt8
	case clang.TK_Char_U:
		typ.CGoName = CUChar
		typ.GoName = GoUInt8
	case clang.TK_Int:
		typ.CGoName = CInt
		typ.GoName = GoInt16
	case clang.TK_Short:
		typ.CGoName = CShort
		typ.GoName = GoInt16
	case clang.TK_UShort:
		typ.CGoName = CUShort
		typ.GoName = GoUInt16
	case clang.TK_UInt:
		typ.CGoName = CUInt
		typ.GoName = GoUInt16
	case clang.TK_Long:
		typ.CGoName = CLongInt
		typ.GoName = GoInt32
	case clang.TK_ULong:
		typ.CGoName = CULongInt
		typ.GoName = GoUInt32
	case clang.TK_LongLong:
		typ.CGoName = CLongLong
		typ.GoName = GoInt64
	case clang.TK_ULongLong:
		typ.CGoName = CULongLong
		typ.GoName = GoUInt64
	case clang.TK_Float:
		typ.CGoName = CFloat
		typ.GoName = GoFloat32
	case clang.TK_Double:
		typ.CGoName = CDouble
		typ.GoName = GoFloat64
	case clang.TK_Bool:
		typ.GoName = GoBool
	case clang.TK_Void:
		// TODO Does not exist in Go
		typ.CGoName = "void"
		typ.GoName = "void"
	case clang.TK_ConstantArray:
		subTyp, err := getType(cType.ArrayElementType())
		if err != nil {
			return Type{}, err
		}

		typ.CGoName = subTyp.CGoName
		typ.GoName = subTyp.GoName
		typ.CGoName = subTyp.CGoName
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsArray = true
		typ.ArraySize = cType.ArraySize()
	case clang.TK_Typedef:
		typ.IsPrimitive = false

		typeStr := cType.TypeSpelling()
		if typeStr == "CXString" {
			typeStr = "cxstring"
		} else if typeStr == "time_t" {
			typ.CGoName = typeStr
			typeStr = "time.Time"

			typ.IsPrimitive = true
		} else {
			typeStr = trimClangPrefix(cType.Declaration().Type().TypeSpelling())
		}

		typ.CGoName = cType.Declaration().Type().TypeSpelling()
		typ.GoName = typeStr

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

		typ.CGoName = subTyp.CGoName
		typ.GoName = subTyp.GoName
		typ.CGoName = subTyp.CGoName
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsPrimitive = subTyp.IsPrimitive
	case clang.TK_Record:
		typ.CGoName = cType.Declaration().Type().TypeSpelling()
		typ.GoName = trimClangPrefix(typ.CGoName)
		typ.IsPrimitive = false
	case clang.TK_FunctionProto:
		typ.IsFunctionPointer = true
		typ.CGoName = cType.Declaration().Type().TypeSpelling()
		typ.GoName = trimClangPrefix(typ.CGoName)
	case clang.TK_Enum:
		typ.GoName = trimClangPrefix(cType.Declaration().DisplayName())
		typ.IsEnumLiteral = true
		typ.IsPrimitive = true
	case clang.TK_Unexposed: // there is a bug in clang for enums the kind is set to unexposed dunno why, bug persists since 2013
		subTyp, err := getType(cType.CanonicalType())
		if err != nil {
			return Type{}, err
		}

		typ.CGoName = subTyp.CGoName
		typ.GoName = subTyp.GoName
		typ.CGoName = subTyp.CGoName
		typ.PointerLevel += subTyp.PointerLevel
		typ.IsPrimitive = subTyp.IsPrimitive
	default:
		return Type{}, fmt.Errorf("unhandled type %q of kind %q", cType.TypeSpelling(), cType.Kind().Spelling())
	}

	return typ, nil
}
