package gen

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/go-clang/bootstrap/clang"
)

// Enum represents a generation enum.
type Enum struct {
	IncludeFiles IncludeFiles

	Name           string
	CName          string
	CNameIsTypeDef bool
	Receiver       Receiver
	Comment        string
	UnderlyingType string

	Items []EnumItem

	Methods []interface{}
}

// EnumItem represents a generation enum item.
type EnumItem struct {
	Name    string
	CName   string
	Comment string
	Value   uint64
}

// HandleEnumCursor handles enum clang.Cursor and roterns the new *Enum.
func HandleEnumCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Enum {
	e := Enum{
		IncludeFiles:   NewIncludeFiles(),
		Name:           TrimLanguagePrefix(cname),
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        CleanDoxygenComment(cursor.RawCommentText()),
		Items:          []EnumItem{},
	}
	e.Receiver.Name = CommonReceiverName(e.Name)
	e.Receiver.Type.GoName = e.Name
	e.Receiver.Type.CGoName = e.CName

	if cnameIsTypeDef {
		e.Receiver.Type.CGoName = e.CName
	} else {
		e.Receiver.Type.CGoName = "enum_" + e.CName
	}

	enumNamePrefix := e.Name
	enumNamePrefix = strings.TrimSuffix(enumNamePrefix, "Kind")
	enumNamePrefix = strings.SplitN(enumNamePrefix, "_", 2)[0]

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		switch cursor.Kind() {
		case clang.Cursor_EnumConstantDecl:
			ei := EnumItem{
				CName: cursor.Spelling(),
				// TODO(go-clang): we are always using the same comment if there is none, see "TypeKind" https://github.com/go-clang/gen/issues/58
				Comment: CleanDoxygenComment(cursor.RawCommentText()),
				Value:   cursor.EnumConstantDeclUnsignedValue(),
			}
			ei.Name = TrimLanguagePrefix(ei.CName)

			// check if the first item has an enum prefix
			if len(e.Items) == 0 {
				eis := strings.SplitN(ei.Name, "_", 2)
				if len(eis) == 2 {
					enumNamePrefix = ""
				}
			}

			// add the enum prefix to the item
			if enumNamePrefix != "" {
				ei.Name = strings.TrimSuffix(ei.Name, enumNamePrefix)

				if !strings.HasPrefix(ei.Name, enumNamePrefix) {
					ei.Name = enumNamePrefix + "_" + ei.Name
				}
			}
			e.Items = append(e.Items, ei)

		default:
			panic(fmt.Errorf("unexpected cursor.Kind: %#v", cursor.Kind()))
		}

		return clang.ChildVisit_Continue
	})

	if strings.HasSuffix(e.Name, "Error") {
		e.UnderlyingType = "int32"
	} else {
		e.UnderlyingType = "uint32"
	}

	return &e
}

// ContainsMethod reports whether the contains name to Enum.Methods.
func (e *Enum) ContainsMethod(name string) bool {
	for _, m := range e.Methods {
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

// Generate generates enum.
func (e *Enum) Generate() error {
	f := NewFile(strings.ToLower(e.Name))
	f.Enums = append(f.Enums, e)

	return f.Generate()
}

// AddEnumStringMethods adds Enum String methods to e.
func (e *Enum) AddEnumStringMethods() error {
	if !e.ContainsMethod("Spelling") {
		if err := e.AddEnumSpellingMethod(); err != nil {
			return err
		}
	}

	if !e.ContainsMethod("String") {
		if err := e.AddSpellingMethodAlias("String"); err != nil {
			return err
		}
	}

	if strings.HasSuffix(e.Name, "Error") && !e.ContainsMethod("Error") {
		if err := e.AddSpellingMethodAlias("Error"); err != nil {
			return err
		}
	}

	return nil
}

// AddEnumSpellingMethod adds Enum spelling method to e.
func (e *Enum) AddEnumSpellingMethod() error {
	f := NewFunction("Spelling", e.Name, "", "", Type{GoName: "string"})
	fa := NewASTFunc(f)

	fa.GenerateReceiver()

	fa.AddReturnType("", Type{GoName: "string"})

	switchStmt := doSwitchStmt(&ast.Ident{Name: f.Receiver.Name})
	fa.Body.List = append(fa.Body.List, switchStmt)

	m := make(map[uint64]*ast.CaseClause)

	for _, enumerator := range e.Items {
		// EnumItems might have the same value, e.g.:
		//  enum Example {
		//  	aValue = 1
		//  	bValue = aValue
		//  }
		// Translating each EnumItem in its own case would result in a compliation error (https://golang.org/issues/4524).
		// Thus, all EnumItems with the same value need to be pooled.
		if caseClause, ok := m[enumerator.Value]; !ok {
			c := []ast.Expr{&ast.Ident{Name: enumerator.Name}}
			ret := &ast.ReturnStmt{
				Results: []ast.Expr{
					doStringLit(strings.Replace(enumerator.Name, "_", "=", 1)),
				},
			}
			b := []ast.Stmt{ret}

			caseClause = doCaseClause(c, b)
			switchStmt.Body.List = append(switchStmt.Body.List, caseClause)
			m[enumerator.Value] = caseClause
		} else {
			retStr := caseClause.Body[0].(*ast.ReturnStmt).Results[0].(*ast.BasicLit).Value
			retStr = retStr[0:len(retStr)-1] + ", " + enumerator.Name[strings.Index(enumerator.Name, "_")+1:] + "\""
			caseClause.Body[0].(*ast.ReturnStmt).Results[0].(*ast.BasicLit).Value = retStr
		}
	}

	fa.AddReturnItem(doCall(
		"fmt",
		"Sprintf",
		doStringLit(f.Receiver.Type.GoName+" unkown %d"),
		doCast("int", &ast.Ident{Name: f.Receiver.Name}),
	))

	fa.AddEmptyLine()
	fa.AddStatement(fa.ret)

	e.Methods = append(e.Methods, GenerateFunctionString(fa))

	return nil
}

// AddSpellingMethodAlias adds spelling method alias to e.
func (e *Enum) AddSpellingMethodAlias(name string) error {
	returnType := Type{
		GoName: "string",
	}

	f := NewFunction(name, e.Name, "", "", returnType)
	fa := NewASTFunc(f)

	fa.GenerateReceiver()

	fa.AddReturnType("", Type{GoName: "string"})

	fa.AddReturnItem(doCall(e.Receiver.Name, "Spelling"))

	fa.AddStatement(fa.ret)

	e.Methods = append(e.Methods, GenerateFunctionString(fa))

	return nil
}
