package main

import (
	"strings"

	"github.com/sbinet/go-clang"
)

// Struct represents a generation struct
type Struct struct {
	HeaderFile string

	Name           string
	CName          string
	CNameIsTypeDef bool
	Receiver       Receiver
	Comment        string

	Methods []string
}

func HandleVoidStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := Struct{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),
	}

	s.Name = trimLanguagePrefix(s.CName)
	s.Receiver.Name = receiverName(s.Name)

	return &s
}

func HandleStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := HandleVoidStructCursor(cursor, cname, cnameIsTypeDef)

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {

		switch cursor.Kind() {
		case clang.CK_FieldDecl:
			typ, err := TypeFromClangType(cursor.Type()) // TODO error handling
			if err != nil {
				return clang.CVR_Continue
			}

			if typ.IsFunctionPointer {
				return clang.CVR_Continue
			}

			comment := cleanDoxygenComment(cursor.RawCommentText())

			if (typ.PointerLevel >= 1 && typ.GoName == "void") || typ.CGoName == "uintptr_t" {
				/*typ.CName = "void"
				typ.Name = GoPointer
				if typ.PointerLevel >= 1 {
					typ.PointerLevel--
				}
				typ.IsPrimitive = true*/

				break
			}

			var method string

			var fName string

			if typ.PointerLevel == 2 || typ.IsArray {
				sizeMember := ""

				if typ.ArraySize == -1 {
					sizeMember = "num" + upperFirstCharacter(cursor.DisplayName())
				}

				f := &FunctionSliceReturn{
					Function: NewFunction(cursor.DisplayName(), cname, comment, cursor.DisplayName(), typ),

					SizeMember: sizeMember,

					CElementType:    typ.CGoName,
					ElementType:     typ.GoName,
					IsPrimitive:     typ.IsPrimitive,
					ArrayDimensions: typ.PointerLevel,
					ArraySize:       typ.ArraySize,
				}

				method = generateFunctionSliceReturn(f)
				fName = f.Name

			} else if typ.PointerLevel < 2 {

				f := NewFunction(cursor.DisplayName(), cname, comment, cursor.DisplayName(), typ)

				method = generateFunctionStructMemberGetter(f)
				fName = f.Name

			} else {
				panic("Three pointers")
			}

			if !s.ContainsMethod(fName) {
				s.Methods = append(s.Methods, method)
			}
		}

		return clang.CVR_Continue
	})
	return s
}

func (s *Struct) ContainsMethod(name string) bool {
	for _, m := range s.Methods {
		if strings.Contains(m, ") "+name+"()") {
			return true
		}
	}

	return false
}

func (s *Struct) Generate() error {
	f := NewFile(strings.ToLower(s.Name))
	f.Structs = append(f.Structs, s)

	return f.Generate()
}
