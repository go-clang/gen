package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/sbinet/go-clang"
)

type Struct struct {
	Name           string
	CName          string
	CNameIsTypeDef bool
	Receiver       string
	Comment        string
	ImportUnsafe   bool

	Methods []string
}

type FuncDef struct {
	Comment    string
	Name       string
	GType      string
	FuncName   string
	ReturnType string

	ReturnString string
}

var returnPrimitive = "return %s(%s)"
var returnComplex = "return %s{%s}"

var templateGenerateReturnSlice = template.Must(template.New("go-clang-generate-slice").Parse(`
	{{$.Comment}}
	func ({{$.Name}} {{$.GType}}) {{$.FuncName}}() {{$.ReturnType}} {
 		s := {{$.ReturnType}}{}
 		length := C.sizeof(%s[0]) / C.sizeof(%s[0][0])
		for is:=0;is < length; is++ {
     		s = append(s, %s{%s[is]})
		}
		return s
	}
`))

var templateGenerateGetter = template.Must(template.New("go-clang-generate-getter").Parse(`
	{{$.Comment}}
	func ({{$.Name}} {{$.GType}}) {{$.FuncName}}() {{$.ReturnType}} {
		{{$.ReturnString}}
	}
`))

func handleStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := handleVoidStructCursor(cursor, cname, cnameIsTypeDef)
	gType := trimClangPrefix(cname)

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

		switch cursor.Kind() {
		case clang.CK_FieldDecl:
			conv := getTypeConversion(cursor.Type())

			if conv.FunctionPointer {
				return clang.CVR_Continue
			}

			receiver := strings.ToLower(string(gType[0]))

			f := FuncDef{
				Comment:    cleanDoxygenComment(cursor.RawCommentText()),
				Name:       receiver,
				GType:      gType,
				FuncName:   upperFirstCharacter(cursor.DisplayName()),
				ReturnType: conv.GType,
			}

			callToC := receiver + ".c." + cursor.DisplayName()

			var method string

			if conv.Pointer == 2 {
				elemType := "&" + f.ReturnType
				f.ReturnType = "[]*" + f.ReturnType

				var b bytes.Buffer
				if err := templateGenerateReturnSlice.Execute(&b, f); err != nil {
					fmt.Println(err.Error())
				}

				method = fmt.Sprintf(b.String(), callToC, callToC, elemType, callToC)
			} else if conv.Pointer < 2 {
				if conv.Pointer == 1 {
					f.ReturnType = "*" + f.ReturnType
				}

				if f.ReturnType == "*void" {
					f.ReturnType = GoPointer
					s.ImportUnsafe = true
				}

				if conv.IsArray {
					f.ReturnType = "[]" + f.ReturnType
				}

				if conv.Primitive {
					f.ReturnString = fmt.Sprintf(returnPrimitive, strings.Replace(f.ReturnType, "*", "&", -1), callToC)
				} else {
					f.ReturnString = fmt.Sprintf(returnComplex, strings.Replace(f.ReturnType, "*", "&", -1), callToC)
				}

				var b bytes.Buffer
				if err := templateGenerateGetter.Execute(&b, f); err != nil {
					fmt.Println(err.Error())
				}

				method = b.String()
			} // fuck you, three levels of pointers or more

			s.Methods = append(s.Methods, method)
		}

		return clang.CVR_Continue
	})

	return s
}

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
	GType           string
	Pointer         int
	Primitive       bool
	IsArray         bool
	FunctionPointer bool
}

func getTypeConversion(cType clang.Type) Conversion {
	conv := Conversion{
		Pointer:   0,
		Primitive: true,
		IsArray:   false,
	}

	switch cType.Kind() {
	case clang.TK_Char_S:
		conv.GType = string(GoInt8)
	case clang.TK_Char_U:
		conv.GType = GoUInt8
	case clang.TK_Int, clang.TK_Short:
		conv.GType = GoInt16
	case clang.TK_UInt, clang.TK_UShort:
		conv.GType = GoUInt16
	case clang.TK_Long:
		conv.GType = GoInt32
	case clang.TK_ULong:
		conv.GType = GoUInt32
	case clang.TK_LongLong:
		conv.GType = GoInt64
	case clang.TK_ULongLong:
		conv.GType = GoUInt64
	case clang.TK_Float:
		conv.GType = GoFloat32
	case clang.TK_Double:
		conv.GType = GoFloat64
	case clang.TK_Bool:
		conv.GType = GoBool
	case clang.TK_Void:
		conv.GType = "void"
	case clang.TK_ConstantArray:
		subConv := getTypeConversion(cType.ArrayElementType())
		conv.GType = subConv.GType
		conv.Pointer += subConv.Pointer
		conv.IsArray = true
	case clang.TK_Typedef:
		typeStr := cType.TypeSpelling()
		if typeStr == "CXString" {
			typeStr = "cxstring"
		} else {
			typeStr = trimClangPrefix(cType.TypeSpelling())
		}
		conv.GType = typeStr
		conv.Primitive = false
	case clang.TK_Pointer:
		conv.Pointer++

		if cType.PointeeType().CanonicalType().Kind() == clang.TK_FunctionProto {
			conv.FunctionPointer = true
		}

		subConv := getTypeConversion(cType.PointeeType().Declaration().Type()) // ComplexTypes
		if subConv.GType == "" {                                               // datatypes
			subConv = getTypeConversion(cType.PointeeType())
		} else {
			conv.Primitive = false
		}
		conv.GType = subConv.GType
		conv.Pointer += subConv.Pointer
	case clang.TK_Unexposed: // there is a bug in clang for enums the kind is set to unexposed dunno why, bug persists since 2013
		enumStr := cType.CanonicalType().TypeSpelling()
		if strings.Contains(enumStr, "enum") {
			enumStr = trimClangPrefix(cType.CanonicalType().Declaration().DisplayName())
		} else {
			enumStr = ""
		}
		conv.GType = enumStr
	}

	return conv
}

func handleVoidStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := Struct{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),
	}

	s.Name = trimClangPrefix(s.CName)
	s.Receiver = receiverName(s.Name)

	return &s
}

var templateGenerateStruct = template.Must(template.New("go-clang-generate-struct").Parse(`package phoenix
// #include "go-clang.h"
import "C"
%s

{{$.Comment}}
type {{$.Name}} struct {
	c C.{{if not $.CNameIsTypeDef}}struct_{{end}}{{$.CName}}
}
{{range $i, $m := .Methods}}
{{$m}}
{{end}}
`))

func generateStruct(s *Struct) error {
	var b bytes.Buffer
	if err := templateGenerateStruct.Execute(&b, s); err != nil {
		return err
	}

	file := b.String()

	if s.ImportUnsafe {
		file = fmt.Sprintf(file, "import \"unsafe\"")
	} else {
		file = fmt.Sprintf(file, "")
	}
	// TODO remove "_" from names for files here?

	return ioutil.WriteFile(strings.ToLower(s.Name)+"_gen.go", []byte(file), 0600)
}
