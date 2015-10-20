package main

import (
	"strings"

	"github.com/sbinet/go-clang"
)

type Function struct {
	Name    string
	CName   string
	Comment string

	Parameters []FunctionParameter
	ReturnType string
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
