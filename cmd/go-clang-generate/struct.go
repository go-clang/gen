package main

import (
	"bytes"
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
	ImportStdLib   bool

	Methods []string
}

func handleStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := handleVoidStructCursor(cursor, cname, cnameIsTypeDef)

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

		switch cursor.Kind() {
		case clang.CK_FieldDecl:
			conv := getTypeConversion(cursor.Type())

			if conv.IsFunctionPointer {
				return clang.CVR_Continue
			}

			comment := cleanDoxygenComment(cursor.RawCommentText())

			var method string

			if conv.PointerLevel == 2 {

				f := &FunctionSliceReturn{
					Function: *generateFunction(cursor.DisplayName(), cname, comment, cursor.DisplayName(), conv),

					ElementType:          "*" + conv.GoType,
					IsElementTypePointer: true,
					PointeeType:          conv.GoType,
					IsPrimitive:          conv.IsPrimitive,
				}

				method = generateFunctionSliceReturn(f)
				s.ImportStdLib = true

			} else if conv.PointerLevel < 2 {
				if conv.PointerLevel == 1 && conv.GoType == "void" {
					conv.GoType = GoPointer
					conv.PointerLevel = 0
					s.ImportUnsafe = true
				}

				f := generateFunction(cursor.DisplayName(), cname, comment, cursor.DisplayName(), conv)

				method = generateFunctionStructMemberGetter(f)
			} else {
				panic("Three pointers")
			}

			s.Methods = append(s.Methods, method)
		}

		return clang.CVR_Continue
	})

	return s
}

func generateFunction(name, cname, comment, member string, conv Conversion) *Function {
	receiverType := trimClangPrefix(cname)
	receiverName := strings.ToLower(string(receiverType[0]))
	functionName := upperFirstCharacter(name)

	rType := ""
	rTypePrimitive := ""

	if conv.IsPrimitive {
		rTypePrimitive = conv.GoType
	} else {
		rType = conv.GoType
	}

	f := &Function{
		Name:    functionName,
		CName:   cname,
		Comment: comment,

		Parameters: []FunctionParameter{},

		ReturnType:          rType,
		ReturnPrimitiveType: rTypePrimitive,
		IsReturnTypePointer: conv.PointerLevel > 0,
		IsReturnTypeArray:   conv.IsArray,

		Receiver:     receiverName,
		ReceiverType: receiverType,

		Member: member,
	}

	return f
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
{{if $.ImportStdLib}}
// #include <stdlib.h>{{end}}
// #include "go-clang.h"
import "C"
{{if $.ImportUnsafe}}
import "unsafe"{{end}}

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

	// TODO remove "_" from names for files here?

	return ioutil.WriteFile(strings.ToLower(s.Name)+"_gen.go", b.Bytes(), 0600)
}
