package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/sbinet/go-clang"
)

type Enum struct {
	Name           string
	CName          string
	CNameIsTypeDef bool
	Receiver       Receiver
	Comment        string
	UnderlyingType string

	Imports map[string]struct{}

	Items []Enumerator

	Methods []string
}

type Enumerator struct {
	Name    string
	CName   string
	Comment string
}

func handleEnumCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Enum {
	e := Enum{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),

		Items: []Enumerator{},
	}

	e.Name = trimClangPrefix(e.CName)
	e.Receiver.Name = receiverName(e.Name)
	e.Receiver.Type = e.Name
	e.Receiver.CType = e.CName
	if cnameIsTypeDef {
		e.Receiver.PrimitiveType = e.CName
	} else {
		e.Receiver.PrimitiveType = "enum_" + e.CName
	}

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.CK_EnumConstantDecl:
			ei := Enumerator{
				CName:   cursor.Spelling(),
				Comment: cleanDoxygenComment(cursor.RawCommentText()), // TODO We are always using the same comment if there is none, see "TypeKind"
			}

			ei.Name = trimClangPrefix(ei.CName)
			// TODO remove underlines to make the names more Go idiomatic e.g. "C.CXComment_InlineCommand" should be "CommentInlineCommand"

			e.Items = append(e.Items, ei)
		default:
			panic(cursor.Kind())
		}

		return clang.CVR_Continue
	})

	if strings.HasSuffix(e.Name, "Error") {
		e.UnderlyingType = "int32"
	} else {
		e.UnderlyingType = "uint32"
	}

	return &e
}

var templateGenerateEnum = template.Must(template.New("go-clang-generate-enum").Parse(`package phoenix

// #include "go-clang.h"
import "C"
{{if $.Imports}}
import (
{{range $import, $empty := $.Imports}}	"{{$import}}"
{{end}}){{end}}

{{$.Comment}}
type {{$.Name}} {{$.UnderlyingType}}

const (
{{range $i, $e := .Items}}	{{if $e.Comment}}{{$e.Comment}}
	{{end}}{{$e.Name}}{{if eq $i 0}} {{$.Name}}{{end}} = C.{{$e.CName}}
{{end}}
)

{{range $i, $m := $.Methods}}
{{$m}}
{{end}}
`))

func generateEnum(e *Enum) error {
	// TODO remove this hack
	for _, m := range e.Methods {
		if strings.Contains(m, "unsafe.") {
			e.Imports["unsafe"] = struct{}{}

			break
		}
	}

	var b bytes.Buffer
	if err := templateGenerateEnum.Execute(&b, e); err != nil {
		return err
	}

	// TODO remove "_" from names for files here?

	return ioutil.WriteFile(strings.ToLower(e.Name)+"_gen.go", b.Bytes(), 0600)
}
