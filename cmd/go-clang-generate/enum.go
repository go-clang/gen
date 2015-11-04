package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/sbinet/go-clang"
)

// Enum represents a generation enum
type Enum struct {
	HeaderFile string

	Name           string
	CName          string
	CNameIsTypeDef bool
	Receiver       Receiver
	Comment        string
	UnderlyingType string

	Items []Enumerator

	Methods []string
}

// Enumerator represents a generation enum item
type Enumerator struct {
	Name    string
	CName   string
	Comment string
	Value   uint64
}

func handleEnumCursor(cursor clang.Cursor, cname string, cnameIsTypeDef bool) *Enum {
	e := Enum{
		CName:          cname,
		CNameIsTypeDef: cnameIsTypeDef,
		Comment:        cleanDoxygenComment(cursor.RawCommentText()),

		Items: []Enumerator{},
	}

	e.Name = trimLanguagePrefix(e.CName)

	e.Receiver.Name = receiverName(e.Name)
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
		case clang.CK_EnumConstantDecl:
			ei := Enumerator{
				CName:   cursor.Spelling(),
				Comment: cleanDoxygenComment(cursor.RawCommentText()), // TODO We are always using the same comment if there is none, see "TypeKind"
				Value:   cursor.EnumConstantDeclUnsignedValue(),
			}
			ei.Name = trimLanguagePrefix(ei.CName)

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

		return clang.CVR_Continue
	})

	if strings.HasSuffix(e.Name, "Error") {
		e.UnderlyingType = "int32"
	} else {
		e.UnderlyingType = "uint32"
	}

	return &e
}

func generateEnum(e *Enum) error {
	generateEnumStringMethods(e)

	f := NewFile(strings.ToLower(e.Name))
	f.Enums = append(f.Enums, e)

	return f.Generate()
}

func generateEnumStringMethods(e *Enum) {
	hasSpelling := false
	hasString := false
	hasError := false

	for _, fStr := range e.Methods {
		if strings.Contains(fStr, ") Spelling() ") {
			hasSpelling = true
		}

		if strings.Contains(fStr, ") String() ") {
			hasString = true
		}

		if strings.Contains(fStr, ") Error() ") {
			hasError = true
		}
	}

	if !hasSpelling {
		err := generateEnumSpellingMethod(e, templateGenerateEnumSpelling)
		if err != nil {
			fmt.Println(err)
		}
	}
	if !hasString {
		if err := generateEnumMethod(e, templateGenerateEnumString); err != nil {
			panic(err)
		}
	}
	if strings.HasSuffix(e.Name, "Error") && !hasError {
		if err := generateEnumMethod(e, templateGenerateEnumError); err != nil {
			panic(err)
		}
	}
}

func generateEnumSpellingMethod(e *Enum, tmpl *template.Template) error {
	var b bytes.Buffer
	var err error

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

	if err = tmpl.Execute(&b, s); err != nil {
		return err
	}

	e.Methods = append(e.Methods, b.String())

	return nil
}

func generateEnumMethod(e *Enum, tmpl *template.Template) error {
	var b bytes.Buffer
	var err error
	if err = tmpl.Execute(&b, e); err != nil {
		return err
	}

	e.Methods = append(e.Methods, b.String())

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
