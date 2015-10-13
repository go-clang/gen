package main

import (
	"bytes"
	"github.com/sbinet/go-clang"
	"io/ioutil"
	"strings"
	"text/template"
)

type enum struct {
	Name    string
	CName   string
	Comment string

	Items []enumerator
}

type enumerator struct {
	Name    string
	CName   string
	Comment string
}

func trimClangPrefix(name string) string {
	name = strings.TrimPrefix(name, "CX_CXX")
	name = strings.TrimPrefix(name, "CXX")
	name = strings.TrimPrefix(name, "CX")

	return name
}

func handleEnumCursor(cursor clang.Cursor) enum {
	e := enum{
		CName:   cursor.Spelling(),
		Comment: cursor.RawCommentText(),

		Items: []enumerator{},
	}

	e.Name = trimClangPrefix(e.CName)

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.CK_EnumConstantDecl:
			ei := enumerator{
				CName:   cursor.Spelling(),
				Comment: cursor.RawCommentText(),
			}

			ei.Name = trimClangPrefix(ei.CName)

			e.Items = append(e.Items, ei)
		default:
			panic(cursor.Kind())
		}

		return clang.CVR_Continue
	})

	return e
}

var templateGenerateEnum = template.Must(template.New("go-clang-generate-enum").Parse(`package phoenix

// #include "go-clang.h"
import "C"

{{$.Comment}}
type {{$.Name}} int

const (
{{range $i, $e := .Items}}	{{if $e.Comment}}{{$e.Comment}}
	{{end}}{{$e.Name}} {{$.Name}} = C.{{$e.CName}}
{{end}}
)
`))

func generateEnum(e enum) error {

	var b bytes.Buffer
	if err := templateGenerateEnum.Execute(&b, e); err != nil {
		return err
	}

	return ioutil.WriteFile(strings.ToLower(e.Name)+"_gen.go", b.Bytes(), 0600)
}
