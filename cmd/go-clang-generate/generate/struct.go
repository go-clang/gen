package generate

import (
	"strings"

	"github.com/zimmski/go-clang-phoenix"
)

// Struct represents a generation struct
type Struct struct {
	api *API

	IncludeFiles includeFiles

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

func handleStructCursor(cursor phoenix.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := &Struct{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        CleanDoxygenComment(cursor.RawCommentText()),

		IncludeFiles: newIncludeFiles(),
	}

	s.Name = TrimLanguagePrefix(s.CName)
	s.Receiver.Name = commonReceiverName(s.Name)

	cursor.Visit(func(cursor, parent phoenix.Cursor) phoenix.ChildVisitResult {
		switch cursor.Kind() {
		case phoenix.Cursor_FieldDecl:
			typ, err := typeFromClangType(cursor.Type())
			if err != nil {
				panic(err)
			}

			if typ.IsFunctionPointer {
				return phoenix.ChildVisit_Continue
			}

			s.Members = append(s.Members, &StructMember{
				CName:   cursor.DisplayName(),
				Comment: CleanDoxygenComment(cursor.RawCommentText()),

				Type: typ,
			})
		}

		return phoenix.ChildVisit_Continue
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

func (s *Struct) generate() error {
	f := newFile(strings.ToLower(s.Name))
	f.Structs = append(f.Structs, s)

	return f.generate()
}

func (s *Struct) addMemberGetters() error {
	if s.api.PrepareStructMembers != nil {
		s.api.PrepareStructMembers(s)
	}

	// Generate the getters we can handle
	for _, m := range s.Members {
		if s.api.FilterStructMemberGetter != nil && !s.api.FilterStructMemberGetter(m) {
			continue
		}

		if m.Type.CGoName == "void" || m.Type.CGoName == "uintptr_t" {
			continue
		}

		if m.Type.IsArray { // TODO generate arrays with the correct size and type https://github.com/zimmski/go-clang-phoenix/issues/48
			continue
		}

		f := newFunction(m.CName, s.CName, m.Comment, m.CName, m.Type)

		if !s.ContainsMethod(f.Name) {
			s.Methods = append(s.Methods, f)
		}
	}

	return nil
}
