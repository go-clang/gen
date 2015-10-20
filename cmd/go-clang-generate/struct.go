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
	Comment        string
}

func handleStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) Struct {
	s := handleVoidStructCursor(cursor, cname, cnameIsTypeDef)

	return s
}

func handleVoidStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) Struct {
	s := Struct{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),
	}

	s.Name = trimClangPrefix(s.CName)

	return s
}

var templateGenerateStruct = template.Must(template.New("go-clang-generate-struct").Parse(`package phoenix

// #include "go-clang.h"
import "C"

{{$.Comment}}
type {{$.Name}} struct {
	c C.{{if not $.CNameIsTypeDef}}struct_{{end}}{{$.CName}}
}

`))

func generateStruct(s Struct) error {
	var b bytes.Buffer
	if err := templateGenerateStruct.Execute(&b, s); err != nil {
		return err
	}

	// TODO remove "_" from names for files here?

	return ioutil.WriteFile(strings.ToLower(s.Name)+"_gen.go", b.Bytes(), 0600)
}
