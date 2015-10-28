package main

import (
	"errors"
	"fmt"

	"github.com/sbinet/go-clang"
)

const (
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
)

type Type struct {
	Name  string
	CName string

	PointerLevel      int
	IsPrimitive       bool
	IsArray           bool
	IsEnumLiteral     bool
	IsFunctionPointer bool
}

func getType(cType clang.Type) (Type, error) {
	typ := Type{
		CName:             cType.TypeSpelling(),
		PointerLevel:      0,
		IsPrimitive:       true,
		IsArray:           false,
		IsFunctionPointer: false,
	}

	switch cType.Kind() {
	case clang.TK_Char_S:
		typ.Name = string(GoInt8)
	case clang.TK_Char_U:
		typ.Name = GoUInt8
	case clang.TK_Int, clang.TK_Short:
		typ.Name = GoInt16
	case clang.TK_UInt, clang.TK_UShort:
		typ.Name = GoUInt16
	case clang.TK_Long:
		typ.Name = GoInt32
	case clang.TK_ULong:
		typ.Name = GoUInt32
	case clang.TK_LongLong:
		typ.Name = GoInt64
	case clang.TK_ULongLong:
		typ.Name = GoUInt64
	case clang.TK_Float:
		typ.Name = GoFloat32
	case clang.TK_Double:
		typ.Name = GoFloat64
	case clang.TK_Bool:
		typ.Name = GoBool
	case clang.TK_Void:
		typ.Name = "void"
	case clang.TK_ConstantArray:
		subConv, err := getType(cType.ArrayElementType())
		if err != nil {
			return Type{}, err
		}

		typ.Name = subConv.Name
		typ.PointerLevel += subConv.PointerLevel
		typ.IsArray = true

	case clang.TK_Typedef:
		typeStr := cType.TypeSpelling()
		if typeStr == "CXString" {
			typeStr = "cxstring"
		} else {
			typeStr = trimClangPrefix(cType.TypeSpelling())
		}

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

		subConv, err := getType(cType.PointeeType().Declaration().Type()) // ComplexTypes
		if err != nil {
			return Type{}, err
		}

		if subConv.Name == "" { // datatypes
			subConv, err = getType(cType.PointeeType())
			if err != nil {
				return Type{}, err
			}
		} else {
			typ.IsPrimitive = false
		}

		typ.Name = subConv.Name
		typ.PointerLevel += subConv.PointerLevel

	case clang.TK_Unexposed: // there is a bug in clang for enums the kind is set to unexposed dunno why, bug persists since 2013

		if cType.CanonicalType().Kind() == clang.TK_Enum {
			typ.Name = trimClangPrefix(cType.CanonicalType().Declaration().DisplayName())
			fmt.Println("blub" + typ.Name)
			typ.IsEnumLiteral = true
			typ.IsPrimitive = true
		} else {
			return Type{}, errors.New("unknown type")
		}

	}

	return typ, nil
}
