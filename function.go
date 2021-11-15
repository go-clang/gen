package gen

import (
	"fmt"
	"strings"

	"github.com/go-clang/bootstrap/clang"
)

// Function represents a generation function.
type Function struct {
	IncludeFiles IncludeFiles

	Name    string
	CName   string
	Comment string

	Parameters []FunctionParameter
	ReturnType Type
	Receiver   Receiver
	Member     *FunctionParameter
}

// FunctionParameter represents a generation function parameter.
type FunctionParameter struct {
	Name  string
	CName string
	Type  Type
}

// NewFunction returns the initialized *Function.
func NewFunction(name, cname, comment, member string, typ Type) *Function {
	functionName := UpperFirstCharacter(name)
	receiverType := TrimLanguagePrefix(cname)
	receiverName := CommonReceiverName(receiverType)

	f := &Function{
		IncludeFiles: NewIncludeFiles(),
		Name:         functionName,
		CName:        cname,
		Comment:      comment,
		Parameters: []FunctionParameter{ // TODO(go-clang): this might not be needed if the receiver code is refactored: https://github.com/go-clang/gen/issues/52
			{
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

// HandleFunctionCursor handles function cursor.
func HandleFunctionCursor(cursor clang.Cursor) *Function {
	fname := cursor.Spelling()
	f := Function{
		IncludeFiles: NewIncludeFiles(),
		Name:         fname,
		CName:        fname,
		Comment:      CleanDoxygenComment(cursor.RawCommentText()),
	}

	typ, err := TypeFromClangType(cursor.ResultType())
	if err != nil {
		panic(fmt.Errorf("unexpected cursor.ResultType: %#v: %w", cursor.ResultType(), err))
	}
	f.ReturnType = typ

	numParam := int(cursor.NumArguments())
	f.Parameters = make([]FunctionParameter, 0, numParam)
	for i := 0; i < numParam; i++ {
		param := cursor.Argument(uint32(i))

		p := FunctionParameter{
			CName: param.DisplayName(),
		}

		typ, err := TypeFromClangType(param.Type())
		if err != nil {
			panic(fmt.Errorf("unexpected error: %w, param.Type(): %#v", err, param.Type()))
		}
		p.Type = typ

		p.Name = p.CName
		if p.Name == "" {
			p.Name = CommonReceiverName(p.Type.GoName)
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

// Generate generates the function.
func (f *Function) Generate() string {
	fa := NewASTFunc(f)
	fa.Generate()

	fStr := GenerateFunctionString(fa)

	// TODO(go-clang): find out how to position the comment correctly and do this using the AST: https://github.com/go-clang/gen/issues/54
	if f.Comment != "" {
		fStr = f.Comment + "\n" + fStr
	}

	return fStr
}
