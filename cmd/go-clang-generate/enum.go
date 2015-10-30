package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/sbinet/go-clang"
)

// Enum represents a generation enum
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

// Enumerator represents a generation enum item
type Enumerator struct {
	Name    string
	CName   string
	Comment string
	Value   int64
}

func handleEnumCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Enum {
	e := Enum{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),

		Imports: map[string]struct{}{},

		Items: []Enumerator{},
	}

	e.Name = trimClangPrefix(e.CName)
	e.Receiver.Name = receiverName(e.Name)
	e.Receiver.Type.GoName = e.Name
	e.Receiver.Type.CGoName = e.CName
	if cnameIsTypeDef {
		e.Receiver.Type.CGoName = e.CName
	} else {
		e.Receiver.Type.CGoName = "enum_" + e.CName
	}

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.CK_EnumConstantDecl:
			ei := Enumerator{
				CName:   cursor.Spelling(),
				Comment: cleanDoxygenComment(cursor.RawCommentText()), // TODO We are always using the same comment if there is none, see "TypeKind"
				Value:   cursor.EnumConstantDeclValue(),
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
		if strings.Contains(m, "time.Time") {
			e.Imports["time"] = struct{}{}
		}
		if strings.Contains(m, "unsafe.") {
			e.Imports["unsafe"] = struct{}{}
		}
	}

	generateEnumStringMethods(e)

	var b bytes.Buffer
	if err := templateGenerateEnum.Execute(&b, e); err != nil {
		return err
	}

	// TODO remove "_" from names for files here?

	return ioutil.WriteFile(strings.ToLower(e.Name)+"_gen.go", b.Bytes(), 0600)
}

func generateEnumStringMethods(e *Enum) {
	hasSpelling := false
	hasString := false
	hasError := false

	for _, fStr := range e.Methods {
		if strings.Contains(fStr, ") Spelling() ") {
			hasSpelling = true
		}

		if strings.Contains(fStr, ") String() ") {
			hasString = true
		}

		if strings.Contains(fStr, ") Error() ") {
			hasError = true
		}
	}

	if !hasSpelling {
		err := generateEnumSpellingMethod(e, templateGenerateEnumSpelling)
		if err != nil {
			fmt.Println(err)
		}
	}
	if !hasString {
		generateEnumMethod(e, templateGenerateEnumString)
	}
	if strings.HasSuffix(e.Name, "Error") && !hasError {
		generateEnumMethod(e, templateGenerateEnumError)
	}
}

func generateEnumSpellingMethod(e *Enum, tmpl *template.Template) error {
	var b bytes.Buffer
	var err error

	e.Imports["fmt"] = struct{}{}

	type Case struct {
		CaseStr   []string
		PrettyStr string
	}

	type Switch struct {
		Receiver     string
		ReceiverType string
		Cases        []Case
	}

	s := &Switch{
		Receiver:     e.Receiver.Name,
		ReceiverType: e.Receiver.Type.GoName,
		Cases:        []Case{},
	}

	m := make(map[string]struct{})

	for i, enumerator := range e.Items {
		if _, ok := m[enumerator.Name]; ok {
			continue
		}

		c := Case{
			CaseStr:   []string{enumerator.Name},
			PrettyStr: strings.Replace(enumerator.Name, "_", "=", 1),
		}

		m[enumerator.Name] = struct{}{}

		for s := i + 1; s < len(e.Items); s++ {
			if e.Items[s].Value == enumerator.Value {
				c.CaseStr = append(c.CaseStr, "/*"+e.Items[s].Name+"*/")
				c.PrettyStr = c.PrettyStr + ", " + e.Items[s].Name[strings.Index(e.Items[s].Name, "_")+1:]
				m[e.Items[s].Name] = struct{}{}
			}
		}

		s.Cases = append(s.Cases, c)
	}

	if err = tmpl.Execute(&b, s); err != nil {
		return err
	}

	e.Methods = append(e.Methods, b.String())

	return nil
}

func generateEnumMethod(e *Enum, tmpl *template.Template) error {
	var b bytes.Buffer
	var err error
	if err = tmpl.Execute(&b, e); err != nil {
		return err
	}

	e.Methods = append(e.Methods, b.String())

	return nil
}

var templateGenerateEnumSpelling = template.Must(template.New("go-clang-generate-enum-spelling").Parse(`
func ({{$.Receiver}} {{$.ReceiverType}}) Spelling() string {
	switch {{$.Receiver}} {
		{{range $en := $.Cases}}case {{range $sn := $en.CaseStr}}{{$sn}}{{end}}: return "{{$en.PrettyStr}}"
		{{end}}
	}

	return fmt.Sprintf("{{$.ReceiverType}} unkown %d", int({{$.Receiver}}))
}
`))

var templateGenerateEnumString = template.Must(template.New("go-clang-generate-enum-string").Parse(`
func ({{$.Receiver.Name}} {{$.Receiver.Type.GoName}}) String() string {
	return {{$.Receiver.Name}}.Spelling()
}
`))

var templateGenerateEnumError = template.Must(template.New("go-clang-generate-enum-error").Parse(`
func ({{$.Receiver.Name}} {{$.Receiver.Type.GoName}}) Error() string {
	return {{$.Receiver.Name}}.Spelling()
}
`))
