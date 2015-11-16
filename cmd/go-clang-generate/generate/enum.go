package generate

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/zimmski/go-clang-phoenix"
)

// Enum represents a generation enum
type Enum struct {
	IncludeFiles includeFiles

	Name           string
	CName          string
	CNameIsTypeDef bool
	Receiver       Receiver
	Comment        string
	UnderlyingType string

	Items []EnumItem

	Methods []interface{}
}

// EnumItem represents a generation enum item
type EnumItem struct {
	Name    string
	CName   string
	Comment string
	Value   uint64
}

func handleEnumCursor(cursor phoenix.Cursor, cname string, cnameIsTypeDef bool) *Enum {
	e := Enum{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        CleanDoxygenComment(cursor.RawCommentText()),

		IncludeFiles: newIncludeFiles(),

		Items: []EnumItem{},
	}

	e.Name = TrimLanguagePrefix(e.CName)

	e.Receiver.Name = commonReceiverName(e.Name)
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

	cursor.Visit(func(cursor, parent phoenix.Cursor) phoenix.ChildVisitResult {
		switch cursor.Kind() {
		case phoenix.Cursor_EnumConstantDecl:
			ei := EnumItem{
				CName:   cursor.Spelling(),
				Comment: CleanDoxygenComment(cursor.RawCommentText()), // TODO We are always using the same comment if there is none, see "TypeKind" https://github.com/zimmski/go-clang-phoenix/issues/58
				Value:   cursor.EnumConstantDeclUnsignedValue(),
			}
			ei.Name = TrimLanguagePrefix(ei.CName)

			// Check if the first item has an enum prefix
			if len(e.Items) == 0 {
				eis := strings.SplitN(ei.Name, "_", 2)
				if len(eis) == 2 {
					enumNamePrefix = ""
				}
			}

			// Add the enum prefix to the item
			if enumNamePrefix != "" {
				ei.Name = strings.TrimSuffix(ei.Name, enumNamePrefix)

				if !strings.HasPrefix(ei.Name, enumNamePrefix) {
					ei.Name = enumNamePrefix + "_" + ei.Name
				}
			}

			e.Items = append(e.Items, ei)
		default:
			panic(cursor.Kind())
		}

		return phoenix.ChildVisit_Continue
	})

	if strings.HasSuffix(e.Name, "Error") {
		e.UnderlyingType = "int32"
	} else {
		e.UnderlyingType = "uint32"
	}

	return &e
}

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

func (e *Enum) generate() error {
	f := newFile(strings.ToLower(e.Name))
	f.Enums = append(f.Enums, e)

	return f.generate()
}

var templateGenerateEnumString = template.Must(template.New("go-clang-generate-enum-string").Parse(`
func ({{$.Receiver.Name}} {{$.Receiver.Type.GoName}}) String() string {
	return {{$.Receiver.Name}}.Spelling()
}
`))

var templateGenerateEnumError = template.Must(template.New("go-clang-generate-enum-error").Parse(`
func ({{$.Receiver.Name}} {{$.Receiver.Type.GoName}}) Error() string {
	return {{$.Receiver.Name}}.Spelling()
}
`))

func (e *Enum) addEnumStringMethods() error {
	if !e.ContainsMethod("Spelling") {
		if err := e.addEnumSpellingMethod(); err != nil {
			return err
		}
	}
	if !e.ContainsMethod("String") {
		if err := e.addEnumMethod(templateGenerateEnumString); err != nil {
			return err
		}
	}
	if strings.HasSuffix(e.Name, "Error") && !e.ContainsMethod("Error") {
		if err := e.addEnumMethod(templateGenerateEnumError); err != nil {
			return err
		}
	}

	return nil
}

var templateGenerateEnumSpelling = template.Must(template.New("go-clang-generate-enum-spelling").Parse(`
func ({{$.Receiver}} {{$.ReceiverType}}) Spelling() string {
	switch {{$.Receiver}} {
		{{range $en := $.Cases}}case {{range $sn := $en.CaseStr}}{{$sn}}{{end}}: return "{{$en.PrettyStr}}"
		{{end}}
	}

	return fmt.Sprintf("{{$.ReceiverType}} unkown %d", int({{$.Receiver}}))
}
`))

func (e *Enum) addEnumSpellingMethod() error {
	var b bytes.Buffer

	type Case struct {
		CaseStr   []string
		PrettyStr string
	}

	type Switch struct {
		Receiver     string
		ReceiverType string
		Cases        []Case
	}

	s := &Switch{
		Receiver:     e.Receiver.Name,
		ReceiverType: e.Receiver.Type.GoName,
		Cases:        []Case{},
	}

	m := make(map[string]struct{})

	for i, enumerator := range e.Items {
		if _, ok := m[enumerator.Name]; ok {
			continue
		}

		c := Case{
			CaseStr:   []string{enumerator.Name},
			PrettyStr: strings.Replace(enumerator.Name, "_", "=", 1),
		}
		m[enumerator.Name] = struct{}{}

		for s := i + 1; s < len(e.Items); s++ {
			if e.Items[s].Value == enumerator.Value {
				c.CaseStr = append(c.CaseStr, "/*"+e.Items[s].Name+"*/")
				c.PrettyStr = c.PrettyStr + ", " + e.Items[s].Name[strings.Index(e.Items[s].Name, "_")+1:]
				m[e.Items[s].Name] = struct{}{}
			}
		}

		s.Cases = append(s.Cases, c)
	}

	if err := templateGenerateEnumSpelling.Execute(&b, s); err != nil {
		return err
	}

	e.Methods = append(e.Methods, b.String())

	return nil
}

func (e *Enum) addEnumMethod(tmpl *template.Template) error {
	var b bytes.Buffer
	if err := tmpl.Execute(&b, e); err != nil {
		return err
	}

	e.Methods = append(e.Methods, b.String())

	return nil
}
