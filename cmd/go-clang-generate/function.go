package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/sbinet/go-clang"
)

type Function struct {
	Name    string
	CName   string
	Comment string

	Parameters []FunctionParameter
	ReturnType string

	Receiver              string
	ReceiverType          string
	ReceiverPrimitiveType string
}

type FunctionParameter struct {
	Name  string
	CName string
	Type  string
}

func handleFunctionCursor(cursor clang.Cursor) *Function {
	f := Function{
		CName:   cursor.Spelling(),
		Comment: cleanDoxygenComment(cursor.RawCommentText()),

		Parameters: []FunctionParameter{},
		ReturnType: cursor.ResultType().TypeSpelling(),
	}

	f.Name = strings.TrimPrefix(f.CName, "clang_")

	numParam := uint(cursor.NumArguments())
	for i := uint(0); i < numParam; i++ {
		param := cursor.Argument(i)

		p := FunctionParameter{
			CName: param.DisplayName(),
			Type:  param.Type().TypeSpelling(),
		}

		f.Parameters = append(f.Parameters, p)
	}

	return &f
}

var templateGenerateFunctionStringGetter = template.Must(template.New("go-clang-generate-function-string-getter").Parse(`{{$.Comment}}
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() string {
	o := cxstring{C.{{$.CName}}({{if ne $.ReceiverPrimitiveType ""}}{{$.ReceiverPrimitiveType}}({{$.Receiver}}){{else}}{{$.Receiver}}.c{{end}})}
	defer o.Dispose()

	return o.String()
}
`))

func generateFunctionStringGetter(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateFunctionStringGetter.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}

var templateGenerateFunctionIs = template.Must(template.New("go-clang-generate-function-is").Parse(`{{$.Comment}}
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() bool {
	o := C.{{$.CName}}({{if ne $.ReceiverPrimitiveType ""}}{{$.ReceiverPrimitiveType}}({{$.Receiver}}){{else}}{{$.Receiver}}.c{{end}})

	return o != C.uint(0)
}
`))

func generateGenerateFunctionIs(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateFunctionIs.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}

var templateGenerateFunctionVoidMethod = template.Must(template.New("go-clang-generate-function-void-method").Parse(`{{$.Comment}}
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() {
	C.{{$.CName}}({{if ne $.ReceiverPrimitiveType ""}}{{$.ReceiverPrimitiveType}}({{$.Receiver}}){{else}}{{$.Receiver}}.c{{end}})
}
`))

func generateFunctionVoidMethod(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateFunctionVoidMethod.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}
