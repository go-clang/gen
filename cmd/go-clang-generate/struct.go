package main

import (
	"fmt"
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

	Members []*StructMember
	Methods []string
}

type StructMember struct {
	CName   string
	Comment string

	Type Type
}

func HandleStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := &Struct{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),
	}

	s.Name = trimLanguagePrefix(s.CName)
	s.Receiver.Name = receiverName(s.Name)

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

			s.Members = append(s.Members, &StructMember{
				CName:   cursor.DisplayName(),
				Comment: cleanDoxygenComment(cursor.RawCommentText()),

				Type: typ,
			})
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

func (s *Struct) AddMemberGetters() error {
	for _, m := range s.Members {
		if (m.Type.PointerLevel >= 1 && m.Type.GoName == "void") || m.Type.CGoName == "uintptr_t" {
			/*typ.CName = "void"
			typ.Name = GoPointer
			if typ.PointerLevel >= 1 {
				typ.PointerLevel--
			}
			typ.IsPrimitive = true*/

			continue
		}

		var fName string
		var method string

		if m.Type.PointerLevel == 2 || m.Type.IsArray {
			sizeMember := ""

			if m.Type.ArraySize == -1 {
				sizeMember = "num" + upperFirstCharacter(m.CName)
			}

			f := &FunctionSliceReturn{
				Function: NewFunction(m.CName, s.CName, m.Comment, m.CName, m.Type),

				SizeMember: sizeMember,

				CElementType:    m.Type.CGoName,
				ElementType:     m.Type.GoName,
				IsPrimitive:     m.Type.IsPrimitive,
				ArrayDimensions: m.Type.PointerLevel,
				ArraySize:       m.Type.ArraySize,
			}

			fName = f.Name
			method = generateFunctionSliceReturn(f)
		} else if m.Type.PointerLevel < 2 {
			f := NewFunction(m.CName, s.CName, m.Comment, m.CName, m.Type)

			fName = f.Name
			method = generateFunctionStructMemberGetter(f)
		} else {
			return fmt.Errorf("Three pointers")
		}

		if !s.ContainsMethod(fName) {
			s.Methods = append(s.Methods, method)
		}
	}

	return nil
}
