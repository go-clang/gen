package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
	"unicode"

	"github.com/sbinet/go-clang"
)

type Struct struct {
	Name    string
	CName   string
	Comment string
}

type funcDef struct {
	Name       string
	CType      string
	FuncName   string
	ReturnType string
}

func handleStructCursor(cname string, cursor clang.Cursor) Struct {
	s := handleVoidStructCursor(cname, cursor)
	t := trimClangPrefix(cname)

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.CK_FieldDecl:
			if cursor.Type().ArraySize() == 0 {

			} else {
				fmt.Println("Original : ")
				fmt.Println(cursor.Type().TypeSpelling() + " " + cursor.DisplayName())
				fmt.Println()
				fmt.Println()

				f := funcDef{
					Name:       strings.ToLower(string(t[0])),
					CType:      t,
					FuncName:   upperFirstCharacter(cursor.DisplayName()),
					ReturnType: getTypeString(cursor),
				}

				//TODO:CXString

				/*if cursor.Type().Kind() == clang.TK_Pointer {

				} else {
					cursor.Type().Kind()
					f.Type = trimClangPrefix(cursor.Type().Declaration().DisplayName())
				}*/

				var b bytes.Buffer
				if err := templateGenerateGetter.Execute(&b, f); err != nil {
					fmt.Println(err.Error())
				}

				fmt.Println("Generated: " + b.String())

			}
			/*fmt.Println("-- " + cursor.DisplayName())
			fmt.Println("  +const " + fmt.Sprintf("%t", cursor.Type().IsConstQualified()))
			fmt.Println("  +type " + cursor.Type().PointeeType().TypeSpelling())
			fmt.Println("  +type " + cursor.Type().ClassType().TypeSpelling())
			fmt.Println("  +type " + cursor.Type().Declaration().Type().TypeSpelling())
			if cursor.Type().ArraySize() > 0 {
				fmt.Println("  +arrtype " + cursor.Type().ArrayElementType().TypeSpelling())
			}
			//	fmt.Println("  +type1 " + cursor.Type().ResultType().TypeSpelling())
			//	fmt.Println("  +decl " + cursor.Type().Declaration().Spelling())
			fmt.Println("  +pointer " + cursor.Type().Kind().Spelling())
			fmt.Println("")
			*/
			//fmt.Println("--+pointer" + cursor
		}

		return clang.CVR_Continue
	})

	return s
}

func getTypeString(cursor clang.Cursor) string {

	switch cursor.Type().Kind() {
	case clang.TK_Char_S:
		return "int8"
	case clang.TK_Char_U:
		return "uint8"
	case clang.TK_Int, clang.TK_Short:
		return "int16"
	case clang.TK_UInt, clang.TK_UShort:
		return "uint16"
	case clang.TK_Long:
		return "int32"
	case clang.TK_ULong:
		return "uint32"
	case clang.TK_Float:
		return "float32"
	case clang.TK_Double:
		return "float64"
	case clang.TK_Bool:
		return "bool"
	case clang.TK_Typedef:
		return trimClangPrefix(cursor.Type().TypeSpelling())
	case clang.TK_Pointer:
		pointedType := cursor.Referenced().Type().PointeeType().Declaration().Type().TypeSpelling()
		if pointedType == "" {
			pointedType = getTypeString(cursor.Referenced().Type().Declaration().Type().Declaration())
		}
		return pointedType
	//	case clang.TK_Enum:
	//	f.ReturnType = trimClangPrefix(cursor.Type().TypeSpelling())

	default:
		fmt.Println("+++++++++++++++" + cursor.Type().Kind().Spelling())
		return ""
	}
}

func upperFirstCharacter(s string) string {
	a := []rune(s)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}

func handleVoidStructCursor(cname string, cursor clang.Cursor) Struct {
	s := Struct{
		CName:   cname,
		Comment: cleanDoxygenComment(cursor.RawCommentText()),
	}

	s.Name = trimClangPrefix(s.CName)

	return s
}

var templateGenerateStruct = template.Must(template.New("go-clang-generate-struct").Parse(`package phoenix

// #include "go-clang.h"
import "C"

{{$.Comment}}
type {{$.Name}} struct {
	c C.{{$.CName}}
}

`))

var templateGenerateGetter = template.Must(template.New("go-clang-generate-struct").Parse(`
	func ({{$.Name}} {{$.CType}}) {{$.FuncName}}() {{$.ReturnType}} {

	}
`))

func generateStruct(s Struct) error {
	var b bytes.Buffer
	if err := templateGenerateStruct.Execute(&b, s); err != nil {
		return err
	}

	// TODO remove "_" from names for files here?

	return ioutil.WriteFile(strings.ToLower(s.Name)+"_gen.ago", b.Bytes(), 0600)
}
