package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"github.com/sbinet/go-clang"
)

func trimCommonFunctionName(name string, typ Type) string {
	name = trimCommonFunctionNamePrefix(name)

	if fn := strings.TrimPrefix(name, typ.GoName+"_"); len(fn) != len(name) {
		name = fn
	} else if fn := strings.TrimPrefix(name, typ.GoName); len(fn) != len(name) {
		name = fn
	}

	name = trimCommonFunctionNamePrefix(name)

	// If the function name is empty at this point, it is a constructor
	if name == "" {
		name = typ.GoName
	}

	return name
}

func trimCommonFunctionNamePrefix(name string) string {
	name = strings.TrimPrefix(name, "create")
	name = strings.TrimPrefix(name, "get")

	name = trimLanguagePrefix(name)

	return name
}

func trimLanguagePrefix(name string) string {
	name = strings.TrimPrefix(name, "CX_CXX")
	name = strings.TrimPrefix(name, "CXX")
	name = strings.TrimPrefix(name, "CX")
	name = strings.TrimPrefix(name, "ObjC")

	return name
}

// Function represents a generation function
type Function struct {
	Name    string
	CName   string
	Comment string

	Parameters []FunctionParameter
	ReturnType Type

	Receiver Receiver

	Member string
}

// FunctionParameter represents a generation function parameter
type FunctionParameter struct {
	Name  string
	CName string
	Type  Type
}

func NewFunction(name, cname, comment, member string, typ Type) *Function {
	receiverType := trimLanguagePrefix(cname)
	receiverName := receiverName(receiverType)
	functionName := upperFirstCharacter(name)

	if (strings.HasPrefix(name, "has") || strings.HasPrefix(name, "is")) && typ.GoName == GoInt16 {
		typ.GoName = GoBool
	}

	f := &Function{
		Name:    functionName,
		CName:   cname,
		Comment: comment,

		Parameters: []FunctionParameter{ // TODO ... do a correct version here...
			FunctionParameter{
				Name:  receiverName,
				CName: cname,
				Type: Type{
					GoName: receiverType,
				},
			},
		},

		ReturnType: typ,
		Receiver: Receiver{
			Name: receiverName,
			Type: Type{
				GoName: receiverType,
			},
		},

		Member: member,
	}

	return f
}

func HandleFunctionCursor(cursor clang.Cursor) *Function {
	f := Function{
		CName:   cursor.Spelling(),
		Comment: cleanDoxygenComment(cursor.RawCommentText()),

		Parameters: []FunctionParameter{},
	}

	f.Name = strings.TrimPrefix(f.CName, "clang_")

	typ, err := TypeFromClangType(cursor.ResultType())
	if err != nil {
		panic(err)
	}
	f.ReturnType = typ

	numParam := uint(cursor.NumArguments())
	for i := uint(0); i < numParam; i++ {
		param := cursor.Argument(i)

		p := FunctionParameter{
			CName: param.DisplayName(),
		}

		typ, err := TypeFromClangType(param.Type())
		if err != nil {
			panic(err)
		}
		p.Type = typ

		p.Name = p.CName
		if p.Name == "" {
			p.Name = receiverName(p.Type.GoName)
		} else {
			pns := strings.Split(p.Name, "_")
			for i := range pns {
				pns[i] = upperFirstCharacter(pns[i])
			}
			p.Name = lowerFirstCharacter(strings.Join(pns, ""))
		}
		if r := ReplaceGoKeywords(p.Name); r != "" {
			p.Name = r
		}

		f.Parameters = append(f.Parameters, p)
	}

	return &f
}

func (f *Function) Generate() string {
	fa := NewASTFunc(f)
	fa.generate()

	var b bytes.Buffer
	err := format.Node(&b, token.NewFileSet(), []ast.Decl{fa.FuncDecl})
	if err != nil {
		panic(err)
	}

	sss := b.String()

	// TODO hack to make new lines...
	sss = strings.Replace(sss, "REMOVE()", "", -1)

	// TODO find out how to position the comment correctly and do this using the AST
	if f.Comment != "" {
		sss = f.Comment + "\n" + sss
	}

	return sss
}
