package generate

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"

	"github.com/zimmski/go-clang-phoenix-bootstrap/clang"
)

type HeaderFile struct {
	api *API

	name string
	dir  string

	enums     []*Enum
	functions []*Function
	structs   []*Struct

	lookupEnum        map[string]*Enum
	lookupNonTypedefs map[string]string
	lookupStruct      map[string]*Struct
}

func (h *HeaderFile) addMethod(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	fname = UpperFirstCharacter(fnamePrefix + UpperFirstCharacter(fname))

	if e, ok := h.lookupEnum[rt.Type.GoName]; ok {
		f.Name = fname

		e.Methods = append(e.Methods, f)

		return true
	} else if s, ok := h.lookupStruct[rt.Type.GoName]; ok && s.CName != "CXString" {
		f.Name = fname

		if !rt.Type.IsSlice && rt.Type.PointerLevel > 0 {
			s.IsPointerComposition = true
		}

		s.Methods = append(s.Methods, f)

		return true
	}

	return false
}

func (h *HeaderFile) addBasicMethods(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	if len(f.Parameters) == 0 && h.IsEnumOrStruct(f.ReturnType.GoName) {
		rtc := rt
		rtc.Type = f.ReturnType

		fname = TrimCommonFunctionNamePrefix(fname)
		if fname == "" {
			fname = f.ReturnType.GoName
		}

		if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") || strings.HasSuffix(f.CName, "_create") {
			fname = "New" + fname
		}

		return h.addMethod(f, fname, fnamePrefix, rtc)
	} else if (fname[0] == 'i' && fname[1] == 's' && unicode.IsUpper(rune(fname[2]))) || (fname[0] == 'h' && fname[1] == 'a' && fname[2] == 's' && unicode.IsUpper(rune(fname[3]))) {
		f.ReturnType.GoName = "bool"

		return h.addMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 1 && h.IsEnumOrStruct(f.Parameters[0].Type.GoName) && strings.HasPrefix(fname, "dispose") && f.ReturnType.GoName == "void" {
		fname = "Dispose"

		return h.addMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 2 && strings.HasPrefix(fname, "equal") && h.IsEnumOrStruct(f.Parameters[0].Type.GoName) && f.Parameters[0].Type == f.Parameters[1].Type {
		fname = "Equal"
		f.Parameters[0].Name = commonReceiverName(f.Parameters[0].Type.GoName)
		f.Parameters[1].Name = f.Parameters[0].Name + "2"

		f.ReturnType.GoName = "bool"

		return h.addMethod(f, fname, fnamePrefix, rt)
	}

	return false
}

func (h *HeaderFile) IsEnumOrStruct(name string) bool {
	if _, ok := h.lookupEnum[name]; ok {
		return true
	} else if _, ok := h.lookupStruct[name]; ok {
		return true
	}

	return false
}

func (h *HeaderFile) setIsPointerComposition(typ *Type) {
	if s, ok := h.lookupStruct[typ.GoName]; ok && s.IsPointerComposition {
		typ.IsPointerComposition = true
	}
}

func (h *HeaderFile) prepareFile() error {
	/*
		Hide all "void *" fields of structs by replacing the type with "uintptr_t".

		To paraphrase the original go-clang source code:
			Not hiding these fields confuses the Go GC during garbage collection and
			pointer scanning, making it think the heap/stack has been somehow corrupted.

		I do not know how the original author debugged this, but one thing: Thank you!
	*/
	findStructsRe := regexp.MustCompile(`(?s)struct[\s\w]+{.+?}`)
	f, err := ioutil.ReadFile(h.name)
	if err != nil {
		return fmt.Errorf("Cannot read Index.h: %v", err)
	}
	voidPointerReplacements := map[string]string{}
	findVoidPointerRe := regexp.MustCompile(`(?:const\s+)?void\s*\*\s*(\w+(\[\d+\])?;)`)
	for _, s := range findStructsRe.FindAll(f, -1) {
		s2 := findVoidPointerRe.ReplaceAll(s, []byte("uintptr_t $1"))
		if len(s) != len(s2) {
			voidPointerReplacements[string(s)] = string(s2)
		}
	}
	fs := string(f)
	for s, r := range voidPointerReplacements {
		fs = strings.Replace(fs, s, r, -1)
	}
	if incl := "#include <stdint.h>"; !strings.HasPrefix(fs, incl) { // Include for uintptr_t
		fs = "#include <stdint.h>\n\n" + fs
	}
	err = ioutil.WriteFile(h.name, []byte(fs), 0600)
	if err != nil {
		return fmt.Errorf("Cannot write %s: %v", h.name, err)
	}

	return nil
}

func (h *HeaderFile) handleFile(cursor clang.Cursor) {
	/*
		TODO mark the enum https://github.com/zimmski/go-clang-phoenix-gen/issues/40
			typedef enum CXChildVisitResult (*CXCursorVisitor)(CXCursor cursor, CXCursor parent, CXClientData client_data);
		as manually implemented
	*/
	// TODO report other enums like callbacks that they are not implemented https://github.com/zimmski/go-clang-phoenix-gen/issues/51

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		// Only handle code of the current file
		sourceFile, _, _, _ := cursor.Location().FileLocation()
		isCurrentFile := sourceFile.Name() == h.name

		if !strings.HasPrefix(sourceFile.Name(), h.dir) {
			return clang.ChildVisit_Continue
		}

		cname := cursor.Spelling()
		cnameIsTypeDef := false

		if parentCName := parent.Spelling(); parent.Kind() == clang.Cursor_TypedefDecl && parentCName != "" {
			cname = parentCName
			cnameIsTypeDef = true
		}

		switch cursor.Kind() {
		case clang.Cursor_EnumDecl:
			if cname == "" {
				break
			}

			e := handleEnumCursor(cursor, cname, cnameIsTypeDef)
			e.IncludeFiles.addIncludeFile(sourceFile.Name())

			if _, ok := h.lookupEnum[e.Name]; !ok {
				h.lookupEnum[e.Name] = e
				h.lookupNonTypedefs["enum "+e.CName] = e.Name
				h.lookupEnum[e.CName] = e

				h.enums = append(h.enums, e)
			}
		case clang.Cursor_FunctionDecl:
			if !isCurrentFile {
				return clang.ChildVisit_Continue
			}

			f := handleFunctionCursor(cursor)
			if f != nil {
				f.IncludeFiles.addIncludeFile(sourceFile.Name())

				h.functions = append(h.functions, f)
			}
		case clang.Cursor_StructDecl:
			if cname == "" {
				break
			}

			s := handleStructCursor(cursor, cname, cnameIsTypeDef)
			s.api = h.api
			s.IncludeFiles.addIncludeFile(sourceFile.Name())

			if _, ok := h.lookupStruct[s.Name]; !ok {
				h.lookupStruct[s.Name] = s
				h.lookupNonTypedefs["struct "+s.CName] = s.Name
				h.lookupStruct[s.CName] = s

				h.structs = append(h.structs, s)
			}
		case clang.Cursor_TypedefDecl:
			underlyingType := cursor.TypedefDeclUnderlyingType().Spelling()
			underlyingStructType := strings.TrimSuffix(strings.TrimPrefix(underlyingType, "struct "), " *")

			if s, ok := h.lookupStruct[underlyingStructType]; ok && !s.CNameIsTypeDef && strings.HasPrefix(underlyingType, "struct "+s.CName) {
				// Sometimes the typedef is not a parent of the struct but a sibling
				sn := handleStructCursor(cursor, cname, true)
				sn.api = h.api
				sn.IncludeFiles.addIncludeFile(sourceFile.Name())

				if sn.Comment == "" {
					sn.Comment = s.Comment
				}
				sn.Members = s.Members
				sn.Methods = s.Methods

				h.lookupStruct[sn.Name] = sn
				h.lookupNonTypedefs["struct "+sn.CName] = sn.Name
				h.lookupStruct[sn.CName] = sn

				// Update the lookups for the old struct
				h.lookupStruct[s.Name] = sn
				h.lookupStruct[s.CName] = sn

				for i, si := range h.structs {
					if si == s {
						h.structs[i] = sn

						break
					}
				}
			} else if underlyingType == "void *" {
				s := handleStructCursor(cursor, cname, true)
				s.api = h.api
				s.IncludeFiles.addIncludeFile(sourceFile.Name())

				if _, ok := h.lookupStruct[s.Name]; !ok {
					h.lookupStruct[s.Name] = s
					h.lookupNonTypedefs["struct "+s.CName] = s.Name
					h.lookupStruct[s.CName] = s

					h.structs = append(h.structs, s)
				}
			}
		}

		return clang.ChildVisit_Recurse
	})
}

func (h *HeaderFile) parse(clangArguments []string) error {
	if err := h.prepareFile(); err != nil {
		return err
	}

	// Parse the header file to analyse everything we need to know
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit(h.name, clangArguments, nil, 0)
	defer tu.Dispose()

	if !tu.IsValid() {
		return fmt.Errorf("Cannot parse Index.h")
	}

	for _, diag := range tu.Diagnostics() {
		switch diag.Severity() {
		case clang.Diagnostic_Error:
			return fmt.Errorf("Diagnostic error in Index.h: %s", diag.Spelling())
		case clang.Diagnostic_Fatal:
			return fmt.Errorf("Diagnostic fatal in Index.h: %s", diag.Spelling())
		}
	}

	h.handleFile(tu.TranslationUnitCursor())

	return nil
}

func (h *HeaderFile) handle() error {
	// Prepare all functions
	clangFile := newFile("clang")

	for _, f := range h.functions {
		var fname string
		if h.api.PrepareFunctionName != nil {
			fname = h.api.PrepareFunctionName(h, f)
			f.Name = fname
		} else {
			fname = f.Name
		}

		if h.api.FilterFunction != nil && !h.api.FilterFunction(f) {
			continue
		}

		// Prepare the parameters
		for i := range f.Parameters {
			p := &f.Parameters[i]

			if n, ok := h.lookupNonTypedefs[p.Type.CGoName]; ok {
				p.Type.GoName = n
			}
			if e, ok := h.lookupEnum[p.Type.GoName]; ok {
				p.CName = e.Receiver.CName
				// TODO remove the receiver... and copy only names here to preserve the original pointers and so https://github.com/zimmski/go-clang-phoenix-gen/issues/52
				p.Type.GoName = e.Receiver.Type.GoName
				p.Type.CGoName = e.Receiver.Type.CGoName
				p.Type.CGoName = e.Receiver.Type.CGoName
			} else if _, ok := h.lookupStruct[p.Type.GoName]; ok {
			}
		}

		if h.api.PrepareFunction != nil {
			h.api.PrepareFunction(f)
		}

		// Prepare the return argument
		if n, ok := h.lookupNonTypedefs[f.ReturnType.CGoName]; ok {
			f.ReturnType.GoName = n
		}
		if e, ok := h.lookupEnum[f.ReturnType.GoName]; ok {
			f.ReturnType.CGoName = e.Receiver.Type.CGoName
		} else if _, ok := h.lookupStruct[f.ReturnType.GoName]; ok {
		}

		// Prepare the receiver
		var rt Receiver
		if len(f.Parameters) > 0 {
			rt.Name = commonReceiverName(f.Parameters[0].Type.GoName)
			rt.CName = f.Parameters[0].CName
			rt.Type = f.Parameters[0].Type
		} else {
			if e, ok := h.lookupEnum[f.ReturnType.GoName]; ok {
				rt.Type = e.Receiver.Type
			} else if s, ok := h.lookupStruct[f.ReturnType.GoName]; ok {
				rt.Type.GoName = s.Name
			}
		}

		// Check upfront if we can handle a function
		found := false

		for _, p := range f.Parameters {
			if h.api.FilterFunctionParameter != nil && !h.api.FilterFunctionParameter(p) {
				continue
			}

			// Return arguments and slices are always ok since we mark them earlier
			if p.Type.IsReturnArgument || p.Type.IsSlice {
				continue
			}

			if (!h.IsEnumOrStruct(p.Type.GoName) && !p.Type.IsPrimitive) || p.Type.PointerLevel != 0 {
				found = true

				fmt.Printf("Cannot handle parameter %s -> %#v\n", f.CName, p)

				break
			}
		}

		// If we find a heuristic to add the function, add it!
		added := false
		if !found {
			added = h.addBasicMethods(f, fname, "", rt)

			if !added {
				if s := strings.Split(f.Name, "_"); len(s) == 2 {
					if s[0] == rt.Type.GoName {
						rtc := rt
						rtc.Name = s[0]

						added = h.addBasicMethods(f, s[1], "", rtc)
					} else {
						if s[0] != "" {
							s[0] += "_"
						}

						added = h.addBasicMethods(f, strings.Join(s[1:], ""), s[0], rt)
					}
				}
			}

			if !added {
				fname = TrimCommonFunctionName(fname, rt.Type)
				if fn := strings.TrimPrefix(fname, f.ReturnType.GoName+"_"); len(fn) != len(fname) {
					fname = fn
				}

				added = h.addMethod(f, fname, "", rt)

				if !added && h.IsEnumOrStruct(f.ReturnType.GoName) {
					rtc := rt
					rtc.Type = f.ReturnType

					fname = TrimCommonFunctionNamePrefix(fname)
					if fname == "" {
						fname = f.ReturnType.GoName
					}

					if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") || strings.HasSuffix(f.CName, "_create") {
						fname = "New" + fname
					}

					added = h.addMethod(f, fname, "", rtc)
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

	for _, e := range h.enums {
		if err := e.addEnumStringMethods(); err != nil {
			return fmt.Errorf("Cannot generate enum string methods: %v", err)
		}

		for i, m := range e.Methods {
			if r := h.handleMethod(e.Name, m); r != "" {
				e.Methods[i] = r
			}
			switch m := m.(type) {
			case *Function:
				if len(m.Parameters) > 0 && !m.Parameters[0].Type.IsSlice && m.Parameters[0].Type.GoName == e.Name {
					m.Receiver = Receiver{
						Name: commonReceiverName(e.Name),
						Type: m.Parameters[0].Type,
					}

					h.setIsPointerComposition(&m.Receiver.Type)
				}
				for i := range m.Parameters {
					h.setIsPointerComposition(&m.Parameters[i].Type)
				}

				e.IncludeFiles.unifyIncludeFiles(m.IncludeFiles)

				e.Methods[i] = m.generate()
			}
		}

		if err := e.generate(); err != nil {
			return fmt.Errorf("Cannot generate enum: %v", err)
		}
	}

	for _, s := range h.structs {
		if err := s.addMemberGetters(); err != nil {
			return fmt.Errorf("Cannot generate struct member getters: %v", err)
		}

		for i, m := range s.Methods {
			if r := h.handleMethod(s.Name, m); r != "" {
				s.Methods[i] = r
			}
			switch m := m.(type) {
			case *Function:
				s.IncludeFiles.unifyIncludeFiles(m.IncludeFiles)
			}
		}

		if err := s.generate(); err != nil {
			return fmt.Errorf("Cannot generate struct: %v", err)
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
			if r := h.handleMethod("", m); r != "" {
				clangFile.Functions[i] = r
			}
		}

		if err := clangFile.generate(); err != nil {
			return fmt.Errorf("Cannot generate clang file: %v", err)
		}
	}

	return nil
}

func (h *HeaderFile) handleMethod(receiverName string, m interface{}) string {
	switch m := m.(type) {
	case *Function:
		if h.api.FixedFunctionName != nil {
			if fname := h.api.FixedFunctionName(m); fname != "" {
				m.Name = fname
			}
		}

		if len(m.Parameters) > 0 && !m.Parameters[0].Type.IsSlice && m.Parameters[0].Type.GoName == receiverName {
			m.Receiver = Receiver{
				Name: commonReceiverName(receiverName),
				Type: m.Parameters[0].Type,
			}

			h.setIsPointerComposition(&m.Receiver.Type)
		}
		for i := range m.Parameters {
			h.setIsPointerComposition(&m.Parameters[i].Type)
		}
		h.setIsPointerComposition(&m.ReturnType)

		return m.generate()
	}

	return ""
}
