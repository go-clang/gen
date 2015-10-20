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

	Receiver     string
	ReceiverType string
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
	cstr := cxstring{C.{{$.CName}}({{$.Receiver}}.c)}
	defer cstr.Dispose()

	return cstr.String()
}
`))

func generateFunctionStringGetter(f *Function) (string, error) {
	var b bytes.Buffer
	if err := templateGenerateFunctionStringGetter.Execute(&b, f); err != nil {
		return "", err
	}

	return b.String(), nil
}
