package main

import (
	"strings"

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

type Conversion struct {
	GoType            string
	PointerLevel      int
	IsPrimitive       bool
	IsArray           bool
	IsFunctionPointer bool
}

func getTypeConversion(cType clang.Type) Conversion {
	conv := Conversion{
		PointerLevel:      0,
		IsPrimitive:       true,
		IsArray:           false,
		IsFunctionPointer: false,
	}

	switch cType.Kind() {
	case clang.TK_Char_S:
		conv.GoType = string(GoInt8)
	case clang.TK_Char_U:
		conv.GoType = GoUInt8
	case clang.TK_Int, clang.TK_Short:
		conv.GoType = GoInt16
	case clang.TK_UInt, clang.TK_UShort:
		conv.GoType = GoUInt16
	case clang.TK_Long:
		conv.GoType = GoInt32
	case clang.TK_ULong:
		conv.GoType = GoUInt32
	case clang.TK_LongLong:
		conv.GoType = GoInt64
	case clang.TK_ULongLong:
		conv.GoType = GoUInt64
	case clang.TK_Float:
		conv.GoType = GoFloat32
	case clang.TK_Double:
		conv.GoType = GoFloat64
	case clang.TK_Bool:
		conv.GoType = GoBool
	case clang.TK_Void:
		conv.GoType = "void"
	case clang.TK_ConstantArray:
		subConv := getTypeConversion(cType.ArrayElementType())

		conv.GoType = subConv.GoType
		conv.PointerLevel += subConv.PointerLevel
		conv.IsArray = true

	case clang.TK_Typedef:
		typeStr := cType.TypeSpelling()
		if typeStr == "CXString" {
			typeStr = "cxstring"
		} else {
			typeStr = trimClangPrefix(cType.TypeSpelling())
		}

		conv.GoType = typeStr
		conv.IsPrimitive = false

	case clang.TK_Pointer:
		conv.PointerLevel++

		if cType.PointeeType().CanonicalType().Kind() == clang.TK_FunctionProto {
			conv.IsFunctionPointer = true
		}

		subConv := getTypeConversion(cType.PointeeType().Declaration().Type()) // ComplexTypes
		if subConv.GoType == "" {                                              // datatypes
			subConv = getTypeConversion(cType.PointeeType())
		} else {
			conv.IsPrimitive = false
		}

		conv.GoType = subConv.GoType
		conv.PointerLevel += subConv.PointerLevel

	case clang.TK_Unexposed: // there is a bug in clang for enums the kind is set to unexposed dunno why, bug persists since 2013

		enumStr := cType.CanonicalType().TypeSpelling()
		if strings.Contains(enumStr, "enum") {
			enumStr = trimClangPrefix(cType.CanonicalType().Declaration().DisplayName())
		} else {
			enumStr = ""
		}

		conv.GoType = enumStr
	}

	return conv
}
