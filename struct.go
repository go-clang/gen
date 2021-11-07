package gen

import (
	"fmt"
	"strings"

	"github.com/go-clang/bootstrap/clang"
)

// Struct represents a generation struct.
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

// StructMember member of Struct.
type StructMember struct {
	CName   string
	Comment string

	Type Type
}

// HandleStructCursor handles the struct cursor.
func HandleStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := &Struct{
		IncludeFiles:   NewIncludeFiles(),
		Name:           TrimLanguagePrefix(cname),
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        CleanDoxygenComment(cursor.RawCommentText()),
	}
	s.Receiver.Name = CommonReceiverName(s.Name)

	cursor.Visit(func(cursor, _ clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_FieldDecl:
			typ, err := typeFromClangType(cursor.Type())
			if err != nil {
				panic(fmt.Errorf("unexpected error: %w, cursor.Type(): %#v", err, cursor.Type()))
			}

			if typ.IsFunctionPointer {
				return clang.ChildVisit_Continue
			}

			s.Members = append(s.Members, &StructMember{
				CName:   cursor.DisplayName(),
				Comment: CleanDoxygenComment(cursor.RawCommentText()),

				Type: typ,
			})
		}

		return clang.ChildVisit_Continue
	})

	return s
}

// ContainsMethod reports whether the contains name in Struct.
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

// Generate struct gereration.
func (s *Struct) Generate() error {
	f := NewFile(strings.ToLower(s.Name))
	f.Structs = append(f.Structs, s)

	return f.Generate()
}

// AddMemberGetters adds member getters to s.
func (s *Struct) AddMemberGetters() error {
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

		if m.Type.IsArray { // TODO(go-clang): generate arrays with the correct size and type https://github.com/go-clang/gen/issues/48
			continue
		}

		f := NewFunction(m.CName, s.CName, m.Comment, m.CName, m.Type)

		if !s.ContainsMethod(f.Name) {
			s.Methods = append(s.Methods, f)
		}
	}

	return nil
}
