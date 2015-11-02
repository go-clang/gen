package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/sbinet/go-clang"
)

// Struct represents a generation struct
type Struct struct {
	HeaderFile string

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

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

		switch cursor.Kind() {
		case clang.CK_FieldDecl:
			typ, err := getType(cursor.Type()) // TODO error handling
			if err != nil {
				return clang.CVR_Continue
			}

			if typ.IsFunctionPointer {
				return clang.CVR_Continue
			}

			comment := cleanDoxygenComment(cursor.RawCommentText())

			if (typ.PointerLevel >= 1 && typ.GoName == "void") || typ.CGoName == "uintptr_t" {
				/*typ.CName = "void"
				typ.Name = GoPointer
				if typ.PointerLevel >= 1 {
					typ.PointerLevel--
				}
				typ.IsPrimitive = true

				s.Imports["unsafe"] = struct{}{}*/

				break
			}

			var method string

			var fName string

			if typ.PointerLevel == 2 || typ.IsArray {
				sizeMember := ""

				if typ.ArraySize == -1 {
					sizeMember = "num" + upperFirstCharacter(cursor.DisplayName())
				}

				f := &FunctionSliceReturn{
					Function: *generateFunction(cursor.DisplayName(), cname, comment, cursor.DisplayName(), typ),

					SizeMember: sizeMember,

					CElementType:    typ.CGoName,
					ElementType:     typ.GoName,
					IsPrimitive:     typ.IsPrimitive,
					ArrayDimensions: typ.PointerLevel,
					ArraySize:       typ.ArraySize,
				}

				method = generateFunctionSliceReturn(f)
				fName = f.Name

			} else if typ.PointerLevel < 2 {

				f := generateFunction(cursor.DisplayName(), cname, comment, cursor.DisplayName(), typ)

				method = generateFunctionStructMemberGetter(f)
				fName = f.Name

			} else {
				panic("Three pointers")
			}

			if !containsMethod(s.Methods, fName) {
				s.Methods = append(s.Methods, method)
			}
		}

		return clang.CVR_Continue
	})
	return s
}

func containsMethod(methods []string, fName string) bool {
	idx := -1
	for i, mem := range methods {
		if strings.Contains(mem, ") "+fName+"()") {
			idx = i
		}
	}

	if idx != -1 {
		return true
	}

	return false
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

{{if $.HeaderFile}}// #include "{{$.HeaderFile}}"
{{end}}// #include "go-clang.h"
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
		if strings.Contains(m, "reflect.") {
			s.Imports["reflect"] = struct{}{}
		}
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

	return ioutil.WriteFile(strings.ToLower(s.Name)+"_gen.go", b.Bytes(), 0600)
}
