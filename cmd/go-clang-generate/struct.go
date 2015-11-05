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

	IsPointerComposition bool

	Members []*StructMember
	Methods []interface{}
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
			typ, err := TypeFromClangType(cursor.Type())
			if err != nil {
				panic(err)
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
		switch m := m.(type) {
		case *Function:
			if m.Name == name {
				return true
			}
		case string:
			if strings.Contains(m, ") "+name+"()") {
				return true
			}
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
	// Prepare members
	for _, m := range s.Members {
		// TODO happy hack, if this is an array length parameter we need to find its partner https://github.com/zimmski/go-clang-phoenix/issues/40
		maCName := ArrayNameFromLength(m.CName)

		if maCName != "" {
			for _, ma := range s.Members {
				if strings.ToLower(ma.CName) == strings.ToLower(maCName) {
					m.Type.LengthOfSlice = ma.CName
					ma.Type.IsSlice = true
					ma.Type.LengthOfSlice = m.CName // TODO wrong usage but needed for the getter generation... maybe refactor this LengthOfSlice alltogether? https://github.com/zimmski/go-clang-phoenix/issues/49

					break
				}
			}
		}
	}

	// Generate the getters we can handle
	for _, m := range s.Members {
		// TODO happy hack, we do not want getters to *int_data members https://github.com/zimmski/go-clang-phoenix/issues/40
		if strings.HasSuffix(m.CName, "int_data") {
			continue
		}

		if m.Type.CGoName == "void" || m.Type.CGoName == "uintptr_t" {
			continue
		}

		if m.Type.IsArray { // TODO generate arrays with the correct size and type https://github.com/zimmski/go-clang-phoenix/issues/48
			continue
		}

		f := NewFunction(m.CName, s.CName, m.Comment, m.CName, m.Type)

		if !s.ContainsMethod(f.Name) {
			s.Methods = append(s.Methods, f)
		}
	}

	return nil
}
