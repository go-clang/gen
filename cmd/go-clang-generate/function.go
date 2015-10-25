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

	Parameters          []FunctionParameter
	ReturnType          string
	ReturnPrimitiveType string
	IsReturnTypePointer bool
	IsReturnTypeArray   bool

	Receiver              string
	ReceiverType          string
	ReceiverPrimitiveType string

	Member string
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
		ReturnType: trimClangPrefix(cursor.ResultType().TypeSpelling()),
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

var templateGenerateFunctionGetter = template.Must(template.New("go-clang-generate-function-getter").Parse(`{{$.Comment}}
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() {{$.ReturnType}} {
	return {{$.ReturnType}}{{if $.ReturnPrimitiveType}}({{else}}{{"{"}}{{end}}C.{{$.CName}}({{if ne $.ReceiverPrimitiveType ""}}{{$.ReceiverPrimitiveType}}({{$.Receiver}}{{else}}{{$.Receiver}}.c{{end}}){{if $.ReturnPrimitiveType}}){{else}}{{"}"}}{{end}}
}
`))

func generateFunctionGetter(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateFunctionGetter.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
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

func generateFunctionIs(f *Function) string {
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

var templateGenerateFunctionEqual = template.Must(template.New("go-clang-generate-function-equal").Parse(`{{$.Comment}}
func {{$.Name}}({{$.Receiver}}1, {{$.Receiver}}2 {{$.ReceiverType}}) bool {
	o := C.{{$.CName}}({{$.Receiver}}1.c, {{$.Receiver}}2.c)

	return o != C.uint(0)
}
`))

func generateFunctionEqual(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateFunctionEqual.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}

var templateGenerateStructMemberGetter = template.Must(template.New("go-clang-generate-function-getter").Parse(`{{$.Comment}}
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() {{if $.IsReturnTypePointer}}*{{end}}{{if $.IsReturnTypeArray}}[]{{end}}{{if $.ReturnPrimitiveType}}{{$.ReturnPrimitiveType}}{{else}}{{$.ReturnType}}{{end}} {
	return {{if $.IsReturnTypePointer}}&{{end}}{{if $.IsReturnTypeArray}}[]{{end}}{{if $.ReturnPrimitiveType}}{{$.ReturnPrimitiveType}}{{else}}{{$.ReturnType}}{{end}}{{if $.ReturnPrimitiveType}}({{$.Receiver}}.c.{{$.Member}}){{else}}{{"{"}}{{$.Receiver}}.c.{{$.Member}}{{"}"}}{{end}}
}
`))

func generateFunctionStructMemberGetter(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateStructMemberGetter.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}

type FunctionSliceReturn struct {
	Function

	ElementType          string
	IsElementTypePointer bool
	PointeeType          string
	IsPrimitive          bool
}

var templateGenerateReturnSlice = template.Must(template.New("go-clang-generate-slice").Parse(`{{$.Comment}}
func ({{$.Name}} {{$.ReceiverType}}) {{$.Name}}() []{{$.ElementType}} {
	s := []{{$.ElementType}}{}
	length := C.sizeof({{$.Receiver}}.c.{{$.Member}}[0]) / C.sizeof({{$.Receiver}}.c.{{$.Member}}[0][0])

	for is := 0; is < length; is++ {
		s = append(s, {{if $.IsElementTypePointer}}&{{$.PointeeType}}{{else}}{{$.ElementType}}{{end}}{{if $.IsPrimitive}}({{$.Receiver}}.c.{{$.Member}}[is])){{else}}{"{"}{{$.Receiver}}.c.{{$.Member}}[is]){{"}"}}{{end}}
	}

	return s
}
`))

func generateFunctionSliceReturn(f *FunctionSliceReturn) string {
	var b bytes.Buffer
	if err := templateGenerateReturnSlice.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()

}
