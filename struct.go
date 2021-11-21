package gen

import (
	"fmt"
	"strings"

	"github.com/go-clang/bootstrap/clang"
)

// Struct represents a generation struct.
type Struct struct {
	IncludeFiles

	api *API

	Name           string
	CName          string
	CNameIsTypeDef bool
	Receiver       Receiver
	Comment        string

	IsPointerComposition bool

	Fields  []*StructField
	Methods []interface{}
}

// StructField field of Struct.
type StructField struct {
	CName   string
	Comment string
	Type    Type
}

// HandleStructCursor handles the struct cursor.
func HandleStructCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Struct {
	s := &Struct{
		IncludeFiles:   NewIncludeFiles(),
		Name:           TrimLanguagePrefix(cname),
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
	}
	s.Comment = CleanDoxygenComment(s.Name, cursor.RawCommentText())
	s.Receiver.Name = CommonReceiverName(s.Name)

	cursor.Visit(func(cursor, _ clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_FieldDecl:
			typ, err := TypeFromClangType(cursor.Type())
			if err != nil {
				panic(fmt.Errorf("unexpected error: %w, cursor.Type(): %#v", err, cursor.Type()))
			}

			if typ.IsFunctionPointer {
				return clang.ChildVisit_Continue
			}

			field := &StructField{
				CName: cursor.DisplayName(),
				Type:  typ,
			}
			field.Comment = CleanDoxygenComment(TrimCommonFunctionName(field.CName, typ), cursor.RawCommentText())
			s.Fields = append(s.Fields, field)
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

// Generate generates the struct.
func (s *Struct) Generate() error {
	f := NewFile(strings.ToLower(s.Name))
	f.Structs = append(f.Structs, s)

	return f.Generate()
}

// AddFieldGetters adds field getters to s.
func (s *Struct) AddFieldGetters() error {
	if s.api.PrepareStructFields != nil {
		s.api.PrepareStructFields(s)
	}

	// generate the getters we can handle
	for _, m := range s.Fields {
		if s.api.FilterStructFieldGetter != nil && !s.api.FilterStructFieldGetter(m) {
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
