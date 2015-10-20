package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/sbinet/go-clang"
)

type Struct struct {
	Name    string
	CName   string
	Comment string
}

func handleStructCursor(cname string, cursor clang.Cursor) Struct {
	s := handleVoidStructCursor(cname, cursor)

	return s
}

func handleVoidStructCursor(cname string, cursor clang.Cursor) Struct {
	s := Struct{
		CName:   cname,
		Comment: cleanDoxygenComment(cursor.RawCommentText()),
	}

	s.Name = trimClangPrefix(s.CName)

	return s
}

var templateGenerateStruct = template.Must(template.New("go-clang-generate-struct").Parse(`package phoenix

// #include "go-clang.h"
import "C"

{{$.Comment}}
type {{$.Name}} struct {
	c C.{{$.CName}}
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
