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
	Receiver       Receiver
	Comment        string

	Imports map[string]struct{}

	Methods []string
}

func handleStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := handleVoidStructCursor(cursor, cname, cnameIsTypeDef)

	if false == true {
		cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

			switch cursor.Kind() {
			case clang.CK_FieldDecl:
				conv, err := getType(cursor.Type()) // TODO error handling
				if err != nil {
					return clang.CVR_Continue
				}

				if conv.IsFunctionPointer {
					return clang.CVR_Continue
				}

				fmt.Println(cursor.Type().TypeSpelling())

				comment := cleanDoxygenComment(cursor.RawCommentText())

				if conv.PointerLevel >= 1 && conv.GoType == "void" {
					conv.GoType = GoPointer
					conv.PointerLevel--

					s.Imports["unsafe"] = struct{}{}
				}

				var method string

				if conv.PointerLevel == 2 || conv.IsArray {
					f := &FunctionSliceReturn{
						Function: *generateFunction(cursor.DisplayName(), cname, comment, cursor.DisplayName(), conv),

						ElementType:     conv.GoType,
						IsPrimitive:     conv.IsPrimitive,
						ArrayDimensions: conv.PointerLevel,
					}

					method = generateFunctionSliceReturn(f)

				} else if conv.PointerLevel < 2 {

					f := generateFunction(cursor.DisplayName(), cname, comment, cursor.DisplayName(), conv)

					method = generateFunctionStructMemberGetter(f)

				} else {
					panic("Three pointers")
				}

				fmt.Println(method)

				s.Methods = append(s.Methods, method)
			}

			return clang.CVR_Continue
		})
	}
	return s
}

func handleVoidStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := Struct{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),

		Imports: map[string]struct{}{},
	}

	s.Name = trimClangPrefix(s.CName)
	s.Receiver.Name = receiverName(s.Name)

	return &s
}

var templateGenerateStruct = template.Must(template.New("go-clang-generate-struct").Parse(`package phoenix

// #include "go-clang.h"
import "C"
{{if $.Imports}}
import (
{{range $import, $empty := $.Imports}}	"{{$import}}"
{{end}}){{end}}

{{$.Comment}}
type {{$.Name}} struct {
	c C.{{if not $.CNameIsTypeDef}}struct_{{end}}{{$.CName}}
}
{{range $i, $m := $.Methods}}
{{$m}}
{{end}}
`))

func generateStruct(s *Struct) error {
	// TODO remove this hack
	for _, m := range s.Methods {
		if strings.Contains(m, "time.Time") {
			s.Imports["time"] = struct{}{}
		}
		if strings.Contains(m, "unsafe.") {
			s.Imports["unsafe"] = struct{}{}
		}
	}

	var b bytes.Buffer
	if err := templateGenerateStruct.Execute(&b, s); err != nil {
		return err
	}

	// TODO remove "_" from names for files here?

	return ioutil.WriteFile(strings.ToLower(s.Name)+"_gen.go", b.Bytes(), 0600)
}
