package gen

import (
	"fmt"
	"strings"
	"unicode"
)

// Generation represents a generation entrypoint.
type Generation struct {
	Lookup

	api *API

	enums     []*Enum
	functions []*Function
	structs   []*Struct
}

// NewGeneration returns the new *Generation from a.
func NewGeneration(a *API) *Generation {
	gen := &Generation{
		Lookup: NewLookup(),
		api:    a,
	}

	return gen
}

// AddHeaderFiles adds headerFiles to g.
func (g *Generation) AddHeaderFiles(headerFiles []*HeaderFile) {
	for _, h := range headerFiles {
		for _, e := range h.Enums {
			g.enums = append(g.enums, e)
			g.RegisterEnum(e)
		}

		for _, s := range h.Structs {
			g.structs = append(g.structs, s)
			g.RegisterStruct(s)
		}

		g.functions = append(g.functions, h.Functions...)
	}
}

// Generate Clang bindings generation.
func (g *Generation) Generate() error {
	// prepare all functions
	clangFile := NewFile("clang")

	for _, f := range g.functions {
		fname := f.Name
		if g.api.PrepareFunctionName != nil {
			fname = g.api.PrepareFunctionName(g, f)
			f.Name = fname
		}

		if g.api.FilterFunction != nil && !g.api.FilterFunction(f) {
			continue
		}

		// prepare the parameters
		for i := range f.Parameters {
			p := &f.Parameters[i]

			if n, ok := g.LookupNonTypedef(p.Type.CGoName); ok {
				p.Type.GoName = n
			}

			if e, ok := g.HasEnum(p.Type.GoName); ok {
				p.CName = e.Receiver.CName
				// TODO(go-clang): remove the receiver... and copy only names here to preserve the original pointers and so https://github.com/go-clang/gen/issues/52
				p.Type.GoName = e.Receiver.Type.GoName
				p.Type.CGoName = e.Receiver.Type.CGoName
				p.Type.CGoName = e.Receiver.Type.CGoName
			}
		}

		if g.api.PrepareFunction != nil {
			g.api.PrepareFunction(f)
		}

		// prepare the return argument
		if n, ok := g.LookupNonTypedef(f.ReturnType.CGoName); ok {
			f.ReturnType.GoName = n
		}
		if e, ok := g.HasEnum(f.ReturnType.GoName); ok {
			f.ReturnType.CGoName = e.Receiver.Type.CGoName
		}

		// prepare the receiver
		var rt Receiver
		if len(f.Parameters) > 0 {
			rt.Name = CommonReceiverName(f.Parameters[0].Type.GoName)
			rt.CName = f.Parameters[0].CName
			rt.Type = f.Parameters[0].Type
		} else {
			if e, ok := g.HasEnum(f.ReturnType.GoName); ok {
				rt.Type = e.Receiver.Type
			} else if s, ok := g.HasStruct(f.ReturnType.GoName); ok {
				rt.Type.GoName = s.Name
			}
		}

		// check upfront if we can handle a function
		found := false

		for _, p := range f.Parameters {
			if g.api.FilterFunctionParameter != nil && !g.api.FilterFunctionParameter(p) {
				continue
			}

			// return arguments and slices are always ok since we mark them earlier
			if p.Type.IsReturnArgument || p.Type.IsSlice {
				continue
			}

			if (!g.IsEnumOrStruct(p.Type.GoName) && !p.Type.IsPrimitive) || p.Type.PointerLevel != 0 {
				found = true

				fmt.Printf("Cannot handle parameter %s -> %#v\n", f.CName, p)

				break
			}
		}

		// if we find a heuristic to add the function, add it!
		added := false
		if !found {
			added = g.AddBasicMethods(f, fname, "", rt)

			if !added {
				if s := strings.Split(f.Name, "_"); len(s) == 2 {
					if s[0] == rt.Type.GoName {
						rtc := rt
						rtc.Name = s[0]

						added = g.AddBasicMethods(f, s[1], "", rtc)
					} else {
						if s[0] != "" {
							s[0] += "_"
						}

						added = g.AddBasicMethods(f, strings.Join(s[1:], ""), s[0], rt)
					}
				}
			}

			if !added {
				fname = TrimCommonFunctionName(fname, rt.Type)
				if fn := strings.TrimPrefix(fname, f.ReturnType.GoName+"_"); len(fn) != len(fname) {
					fname = fn
				}

				added = g.AddMethod(f, fname, "", rt)

				if !added && g.IsEnumOrStruct(f.ReturnType.GoName) {
					rtc := rt
					rtc.Type = f.ReturnType

					fname = TrimCommonFunctionNamePrefix(fname)
					if fname == "" {
						fname = f.ReturnType.GoName
					}

					if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") || strings.HasSuffix(f.CName, "_create") {
						fname = "New" + fname
					}

					added = g.AddMethod(f, fname, "", rtc)
				}

				if !added {
					f.Name = UpperFirstCharacter(f.Name)

					clangFile.Functions = append(clangFile.Functions, f)

					added = true
				}
			}
		}

		if !added {
			fmt.Println("Unused function:", f.Name)
		}
	}

	for _, e := range g.enums {
		if err := e.AddEnumStringMethods(); err != nil {
			return fmt.Errorf("cannot generate enum string methods: %w", err)
		}

		for i, m := range e.Methods {
			e.Methods[i] = g.GenerateMethod(e.Name, m)

			switch m := m.(type) {
			case *Function:
				if len(m.Parameters) > 0 && !m.Parameters[0].Type.IsSlice && m.Parameters[0].Type.GoName == e.Name {
					m.Receiver = Receiver{
						Name: CommonReceiverName(e.Name),
						Type: m.Parameters[0].Type,
					}

					g.SetIsPointerComposition(&m.Receiver.Type)
				}

				for i := range m.Parameters {
					g.SetIsPointerComposition(&m.Parameters[i].Type)
				}

				e.IncludeFiles.unifyIncludeFiles(m.IncludeFiles)

				e.Methods[i] = m.Generate()
			}
		}

		if err := e.Generate(); err != nil {
			return fmt.Errorf("cannot generate enum: %w", err)
		}
	}

	for _, s := range g.structs {
		if err := s.AddFieldGetters(); err != nil {
			return fmt.Errorf("cannot generate struct member getters: %w", err)
		}

		for i, m := range s.Methods {
			s.Methods[i] = g.GenerateMethod(s.Name, m)

			switch m := m.(type) {
			case *Function:
				s.IncludeFiles.unifyIncludeFiles(m.IncludeFiles)
			}
		}

		if err := s.Generate(); err != nil {
			return fmt.Errorf("cannot generate struct: %w", err)
		}
	}

	if len(clangFile.Functions) > 0 {
		for _, m := range clangFile.Functions {
			switch m := m.(type) {
			case *Function:
				clangFile.IncludeFiles.unifyIncludeFiles(m.IncludeFiles)
			}
		}

		for i, m := range clangFile.Functions {
			clangFile.Functions[i] = g.GenerateMethod("", m)
		}

		if err := clangFile.Generate(); err != nil {
			return fmt.Errorf("cannot generate clang file: %w", err)
		}
	}

	return nil
}

// GenerateMethod method generation.
func (g *Generation) GenerateMethod(receiverName string, m interface{}) string {
	switch m := m.(type) {
	case *Function:
		if g.api.FixFunctionName != nil {
			if fname := g.api.FixFunctionName(m); fname != "" {
				m.Name = fname
			}
		}

		if len(m.Parameters) > 0 && !m.Parameters[0].Type.IsSlice && m.Parameters[0].Type.GoName == receiverName {
			m.Receiver = Receiver{
				Name: CommonReceiverName(receiverName),
				Type: m.Parameters[0].Type,
			}

			g.SetIsPointerComposition(&m.Receiver.Type)
		}

		for i := range m.Parameters {
			g.SetIsPointerComposition(&m.Parameters[i].Type)
		}
		g.SetIsPointerComposition(&m.ReturnType)

		m.Comment = strings.ReplaceAll(m.Comment, strings.TrimPrefix(m.CName, "clang_"), m.Name)

		return m.Generate()

	case string:
		return m

	default:
		panic(fmt.Sprintf("Cannot handle %v in handleMethod", m))
	}
}

// SetIsPointerComposition sets IsPointerComposition if given.
func (g *Generation) SetIsPointerComposition(typ *Type) {
	if s, ok := g.HasStruct(typ.GoName); ok && s.IsPointerComposition {
		typ.IsPointerComposition = true
	}
}

// AddMethod adds method to g.
func (g *Generation) AddMethod(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	fname = UpperFirstCharacter(fnamePrefix + UpperFirstCharacter(fname))

	if e, ok := g.HasEnum(rt.Type.GoName); ok {
		f.Name = fname
		e.Methods = append(e.Methods, f)

		return true
	}

	if s, ok := g.HasStruct(rt.Type.GoName); ok && s.CName != "CXString" {
		f.Name = fname

		if !rt.Type.IsSlice && rt.Type.PointerLevel > 0 {
			s.IsPointerComposition = true
		}

		s.Methods = append(s.Methods, f)

		return true
	}

	return false
}

// AddBasicMethods adds basic methods.
func (g *Generation) AddBasicMethods(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	switch {
	case len(f.Parameters) == 0 && g.IsEnumOrStruct(f.ReturnType.GoName):
		rtc := rt
		rtc.Type = f.ReturnType

		fname = TrimCommonFunctionNamePrefix(fname)
		if fname == "" {
			fname = f.ReturnType.GoName
		}

		if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") || strings.HasSuffix(f.CName, "_create") {
			fname = "New" + fname
		}

		return g.AddMethod(f, fname, fnamePrefix, rtc)

	case (fname[0] == 'i' && fname[1] == 's' && unicode.IsUpper(rune(fname[2]))) || (fname[0] == 'h' && fname[1] == 'a' && fname[2] == 's' && unicode.IsUpper(rune(fname[3]))):
		f.ReturnType.GoName = "bool"

		return g.AddMethod(f, fname, fnamePrefix, rt)

	case len(f.Parameters) == 1 && g.IsEnumOrStruct(f.Parameters[0].Type.GoName) && strings.HasPrefix(fname, "dispose") && f.ReturnType.GoName == "void":
		fname = "Dispose"

		return g.AddMethod(f, fname, fnamePrefix, rt)

	case len(f.Parameters) == 2 && strings.HasPrefix(fname, "equal") && g.IsEnumOrStruct(f.Parameters[0].Type.GoName) && f.Parameters[0].Type == f.Parameters[1].Type:
		fname = "Equal"
		f.Parameters[0].Name = CommonReceiverName(f.Parameters[0].Type.GoName)
		f.Parameters[1].Name = f.Parameters[0].Name + "2"

		f.ReturnType.GoName = "bool"

		return g.AddMethod(f, fname, fnamePrefix, rt)
	}

	return false
}
