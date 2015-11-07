package generate

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"unicode"

	"github.com/sbinet/go-clang"
)

func TrimCommonFunctionName(name string, typ Type) string {
	name = TrimCommonFunctionNamePrefix(name)

	if fn := strings.TrimPrefix(name, typ.GoName+"_"); len(fn) != len(name) {
		name = fn
	} else if fn := strings.TrimPrefix(name, typ.GoName); len(fn) != len(name) {
		name = fn
	}
	if tkn := strings.TrimSuffix(typ.GoName, "Kind"); len(tkn) != len(typ.GoName) {
		if fn := strings.TrimPrefix(name, tkn+"_"); len(fn) != len(name) {
			name = fn
		} else if fn := strings.TrimPrefix(name, tkn); len(fn) != len(name) {
			name = fn
		}
	}

	name = TrimCommonFunctionNamePrefix(name)

	// If the function name is empty at this point, it is a constructor
	if name == "" {
		name = typ.GoName
	}

	return name
}

func TrimCommonFunctionNamePrefix(name string) string {
	name = strings.TrimPrefix(name, "create")
	name = strings.TrimPrefix(name, "get")
	if len(name) > 4 && unicode.IsUpper(rune(name[3])) {
		name = strings.TrimPrefix(name, "Get")
	}

	name = TrimLanguagePrefix(name)

	return name
}

func TrimLanguagePrefix(name string) string {
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

	Member *FunctionParameter
}

// FunctionParameter represents a generation function parameter
type FunctionParameter struct {
	Name  string
	CName string
	Type  Type
}

func newFunction(name, cname, comment, member string, typ Type) *Function {
	receiverType := TrimLanguagePrefix(cname)
	receiverName := receiverName(receiverType)
	functionName := UpperFirstCharacter(name)

	if (strings.HasPrefix(name, "has") || strings.HasPrefix(name, "is")) && typ.GoName == GoInt16 {
		typ.GoName = GoBool
	}

	f := &Function{
		Name:    functionName,
		CName:   cname,
		Comment: comment,

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

func handleFunctionCursor(cursor clang.Cursor) *Function {
	f := Function{
		CName:   cursor.Spelling(),
		Comment: CleanDoxygenComment(cursor.RawCommentText()),

		Parameters: []FunctionParameter{},
	}

	f.Name = strings.TrimPrefix(f.CName, "clang_")

	typ, err := typeFromClangType(cursor.ResultType())
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

		typ, err := typeFromClangType(param.Type())
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
