package generate

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"github.com/zimmski/go-clang-phoenix"
)

// Function represents a generation function
type Function struct {
	IncludeFiles includeFiles

	Name    string
	CName   string
	Comment string

	Parameters []FunctionParameter
	ReturnType Type

	Receiver Receiver

	Member *FunctionParameter
}

// FunctionParameter represents a generation function parameter
type FunctionParameter struct {
	Name  string
	CName string
	Type  Type
}

func newFunction(name, cname, comment, member string, typ Type) *Function {
	functionName := UpperFirstCharacter(name)
	receiverType := TrimLanguagePrefix(cname)
	receiverName := commonReceiverName(receiverType)

	f := &Function{
		Name:    functionName,
		CName:   cname,
		Comment: comment,

		IncludeFiles: newIncludeFiles(),

		Parameters: []FunctionParameter{ // TODO this might not be needed if the receiver code is refactored https://github.com/zimmski/go-clang-phoenix/issues/52
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

		Member: &FunctionParameter{
			Name: member,
			Type: typ,
		},
	}

	return f
}

func handleFunctionCursor(cursor phoenix.Cursor) *Function {
	fname := cursor.Spelling()
	f := Function{
		Name:    fname,
		CName:   fname,
		Comment: CleanDoxygenComment(cursor.RawCommentText()),

		IncludeFiles: newIncludeFiles(),

		Parameters: []FunctionParameter{},
	}

	typ, err := typeFromClangType(cursor.ResultType())
	if err != nil {
		panic(err)
	}
	f.ReturnType = typ

	numParam := int(cursor.NumArguments())
	for i := 0; i < numParam; i++ {
		param := cursor.Argument(uint16(i))

		p := FunctionParameter{
			CName: param.DisplayName(),
		}

		typ, err := typeFromClangType(param.Type())
		if err != nil {
			panic(err)
		}
		p.Type = typ

		p.Name = p.CName
		if p.Name == "" {
			p.Name = commonReceiverName(p.Type.GoName)
		} else {
			pns := strings.Split(p.Name, "_")
			for i := range pns {
				pns[i] = UpperFirstCharacter(pns[i])
			}
			p.Name = LowerFirstCharacter(strings.Join(pns, ""))
		}
		if r := ReplaceGoKeywords(p.Name); r != "" {
			p.Name = r
		}

		f.Parameters = append(f.Parameters, p)
	}

	return &f
}

func (f *Function) generate() string {
	fa := newASTFunc(f)
	fa.generate()

	var b bytes.Buffer
	err := format.Node(&b, token.NewFileSet(), []ast.Decl{fa.FuncDecl})
	if err != nil {
		panic(err)
	}

	sss := b.String()

	// TODO hack to make new lines... https://github.com/zimmski/go-clang-phoenix/issues/53
	sss = strings.Replace(sss, "REMOVE()", "", -1)

	// TODO find out how to position the comment correctly and do this using the AST https://github.com/zimmski/go-clang-phoenix/issues/54
	if f.Comment != "" {
		sss = f.Comment + "\n" + sss
	}

	return sss
}
