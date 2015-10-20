package main

import (
	"bytes"
	"github.com/sbinet/go-clang"
	"io/ioutil"
	"strings"
	"text/template"
)

type enum struct { // TODO make public
	Name           string
	CName          string
	CNameIsTypeDef bool
	Comment        string

	Items []enumerator
}

type enumerator struct { // TODO make public
	Name    string
	CName   string
	Comment string
}

func handleEnumCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) enum {
	e := enum{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),

		Items: []enumerator{},
	}

	e.Name = trimClangPrefix(e.CName)

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.CK_EnumConstantDecl:
			ei := enumerator{
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

	return e
}

var templateGenerateEnum = template.Must(template.New("go-clang-generate-enum").Parse(`package phoenix

// #include "go-clang.h"
import "C"

{{$.Comment}}
type {{$.Name}} int

const (
{{range $i, $e := .Items}}	{{if $e.Comment}}{{$e.Comment}}
	{{end}}{{$e.Name}}{{if eq $i 0}} {{$.Name}}{{end}} = C.{{$e.CName}}
{{end}}
)
`))

func generateEnum(e enum) error {
	var b bytes.Buffer
	if err := templateGenerateEnum.Execute(&b, e); err != nil {
		return err
	}

	// TODO remove "_" from names for files here?

	return ioutil.WriteFile(strings.ToLower(e.Name)+"_gen.go", b.Bytes(), 0600)
}
