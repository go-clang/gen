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
	IsReturnTypeEnumLit bool

	Receiver Receiver

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
func ({{$.Receiver.Name}} {{$.Receiver.Type}}) {{$.Name}}() {{$.ReturnType}} {
	return {{$.ReturnType}}{{if $.ReturnPrimitiveType}}({{else}}{{"{"}}{{end}}C.{{$.CName}}({{if ne $.Receiver.PrimitiveType ""}}{{$.Receiver.PrimitiveType}}({{$.Receiver.Name}}){{else}}{{$.Receiver.Name}}.c{{end}}){{if $.ReturnPrimitiveType}}){{else}}{{"}"}}{{end}}
}
`))

func generateFunctionGetter(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateFunctionGetter.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}

var templateGenerateFunctionGetterPrimitive = template.Must(template.New("go-clang-generate-function-getter-primitive").Parse(`{{$.Comment}}
func ({{$.Receiver.Name}} {{$.Receiver.Type}}) {{$.Name}}() {{$.ReturnType}} {
	return {{$.ReturnType}}(C.{{$.CName}}({{if ne $.Receiver.PrimitiveType ""}}{{$.Receiver.PrimitiveType}}({{$.Receiver.Name}}){{else}}{{$.Receiver.Name}}.c{{end}}))
}
`))

func generateFunctionGetterPrimitive(f *Function) string {
	var b bytes.Buffer
	if err := templateGenerateFunctionGetterPrimitive.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()
}

var templateGenerateFunctionStringGetter = template.Must(template.New("go-clang-generate-function-string-getter").Parse(`{{$.Comment}}
func ({{$.Receiver.Name}} {{$.Receiver.Type}}) {{$.Name}}() string {
	o := cxstring{C.{{$.CName}}({{if ne $.Receiver.PrimitiveType ""}}{{$.Receiver.PrimitiveType}}({{$.Receiver.Name}}){{else}}{{$.Receiver.Name}}.c{{end}})}
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
func ({{$.Receiver.Name}} {{$.Receiver.Type}}) {{$.Name}}() bool {
	o := C.{{$.CName}}({{if ne $.Receiver.PrimitiveType ""}}C.{{$.Receiver.CName}}({{$.Receiver.Name}}){{else}}{{$.Receiver.Name}}.c{{end}})

	return o != C.{{if eq $.ReturnType "int"}}int{{else}}uint{{end}}(0)
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
func ({{$.Receiver.Name}} {{$.Receiver.Type}}) {{$.Name}}() {
	C.{{$.CName}}({{if ne $.Receiver.PrimitiveType ""}}{{$.Receiver.PrimitiveType}}({{$.Receiver.Name}}){{else}}{{$.Receiver.Name}}.c{{end}})
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
func {{$.Name}}({{$.Receiver.Name}}1, {{$.Receiver.Name}}2 {{$.Receiver.Type}}) bool {
	o := C.{{$.CName}}({{$.Receiver.Name}}1.c, {{$.Receiver.Name}}2.c)

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
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() {{if $.IsReturnTypePointer}}*{{end}}{{if $.ReturnPrimitiveType}}{{$.ReturnPrimitiveType}}{{else}}{{$.ReturnType}}{{end}} {
	return {{if $.IsReturnTypePointer}}&{{end}}{{if $.ReturnPrimitiveType}}{{$.ReturnPrimitiveType}}{{else}}{{$.ReturnType}}{{end}}{{if $.ReturnPrimitiveType}}({{if $.IsReturnTypePointer}}*{{end}}{{$.Receiver}}.c.{{$.Member}}){{else}}{{"{"}}{{if $.IsReturnTypePointer}}*{{end}}{{$.Receiver}}.c.{{$.Member}}{{"}"}}{{end}}
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

	ElementType     string
	IsPrimitive     bool
	ArrayDimensions int
}

var templateGenerateReturnSlice = template.Must(template.New("go-clang-generate-slice").Parse(`{{$.Comment}}
func ({{$.Receiver}} {{$.ReceiverType}}) {{$.Name}}() []{{if eq $.ArrayDimensions 2 }}*{{end}}{{$.ElementType}} {
	sc := []{{$.ElementType}}{}
	{{if eq $.ArrayDimensions 2 }}
	length := int(C.sizeof({{$.Receiver}}.c.{{$.Member}}[0])) / int(sizeof({{$.Receiver}}.c.{{$.Member}}[0][0]))
	{{else}}
	length := int(sizeof({{$.Receiver}}.c.{{$.Member}}))
	{{end}}
	for is := 0; is < length; is++ {
		sc = append(sc, {{if eq $.ArrayDimensions 2}}&{{$.ElementType}}{{else}}{{$.ElementType}}{{end}}{{if $.IsPrimitive}}({{$.Receiver}}.c.{{$.Member}}[is])){{else}}{"{"}{{$.Receiver}}.c.{{$.Member}}[is]){{"}"}}{{end}}
	}

	return sc
}
`))

func generateFunctionSliceReturn(f *FunctionSliceReturn) string {
	var b bytes.Buffer
	if err := templateGenerateReturnSlice.Execute(&b, f); err != nil {
		panic(err)
	}

	return b.String()

}

func generateFunction(name, cname, comment, member string, conv Conversion) *Function {
	receiverType := trimClangPrefix(cname)
	receiverName := receiverName(string(receiverType[0]))
	functionName := upperFirstCharacter(name)

	rType := ""
	rTypePrimitive := ""

	if conv.IsPrimitive {
		rTypePrimitive = conv.GoType
	} else {
		rType = conv.GoType
	}

	f := &Function{
		Name:    functionName,
		CName:   cname,
		Comment: comment,

		Parameters: []FunctionParameter{},

		ReturnType:          rType,
		ReturnPrimitiveType: rTypePrimitive,
		IsReturnTypePointer: conv.PointerLevel > 0,
		IsReturnTypeEnumLit: conv.IsEnumLiteral,

		Receiver: Receiver{
			Name: receiverName,
			Type: receiverType,
		},

		Member: member,
	}

	return f
}
