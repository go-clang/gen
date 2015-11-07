package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"

	"github.com/sbinet/go-clang"
)

type HeaderFile struct {
	name string

	enums     []*Enum
	functions []*Function
	structs   []*Struct

	lookupEnum        map[string]*Enum
	lookupNonTypedefs map[string]string
	lookupStruct      map[string]*Struct
}

func (h *HeaderFile) AddMethod(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	// Needs to be renamed manually since clang_getTranslationUnitCursor will conflict with clang_getCursor
	if f.CName == "clang_getTranslationUnitCursor" {
		fname = "TranslationUnitCursor"
	} else {
		fname = UpperFirstCharacter(fnamePrefix + UpperFirstCharacter(fname))
	}

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
		fname = TrimCommonFunctionName(fname, rt.Type)
		if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") {
			fname = "New" + fname
		}

		return h.AddMethod(f, fname, fnamePrefix, rt)
	} else if (fname[0] == 'i' && fname[1] == 's' && unicode.IsUpper(rune(fname[2]))) || (fname[0] == 'h' && fname[1] == 'a' && fname[2] == 's' && unicode.IsUpper(rune(fname[3]))) {
		f.ReturnType.GoName = "bool"

		return h.AddMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 1 && h.IsEnumOrStruct(f.Parameters[0].Type.GoName) && strings.HasPrefix(fname, "dispose") && f.ReturnType.GoName == "void" {
		fname = "Dispose"

		return h.AddMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 2 && strings.HasPrefix(fname, "equal") && h.IsEnumOrStruct(f.Parameters[0].Type.GoName) && f.Parameters[0].Type == f.Parameters[1].Type {
		fname = "Equal"
		f.Parameters[0].Name = receiverName(f.Parameters[0].Type.GoName)
		f.Parameters[1].Name = f.Parameters[0].Name + "2"

		f.ReturnType.GoName = "bool"

		return h.AddMethod(f, fname, fnamePrefix, rt)
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

func handleHeaderFile(HeaderFilename string, clangArguments []string) error {
	h := &HeaderFile{
		name: HeaderFilename,

		lookupEnum:        map[string]*Enum{},
		lookupNonTypedefs: map[string]string{},
		lookupStruct: map[string]*Struct{
			"cxstring": &Struct{
				Name:  "cxstring",
				CName: "CXString",
			},
		},
	}

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
		return cmdFatal("Cannot read Index.h", nil)
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
	err = ioutil.WriteFile(h.name, []byte(fs), 0700)
	if err != nil {
		return cmdFatal("Cannot write Index.h", nil)
	}

	// Parse clang-c's Index.h to analyse everything we need to know
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.Parse(h.name, clangArguments, nil, 0)
	defer tu.Dispose()

	if !tu.IsValid() {
		return cmdFatal("Cannot parse Index.h", nil)
	}

	for _, diag := range tu.Diagnostics() {
		switch diag.Severity() {
		case clang.Diagnostic_Error:
			return cmdFatal("Diagnostic error in Index.h", errors.New(diag.Spelling()))
		case clang.Diagnostic_Fatal:
			return cmdFatal("Diagnostic fatal in Index.h", errors.New(diag.Spelling()))
		}
	}

	cursor := tu.ToCursor()
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		// Only handle code of the current file
		sourceFile, _, _, _ := cursor.Location().GetFileLocation()
		if sourceFile.Name() != h.name {
			return clang.CVR_Continue
		}

		cname := cursor.Spelling()
		cnameIsTypeDef := false

		if parentCName := parent.Spelling(); parent.Kind() == clang.CK_TypedefDecl && parentCName != "" {
			cname = parentCName
			cnameIsTypeDef = true
		}

		switch cursor.Kind() {
		case clang.CK_EnumDecl:
			if cname == "" {
				break
			}

			e := handleEnumCursor(cursor, cname, cnameIsTypeDef)

			h.lookupEnum[e.Name] = e
			h.lookupNonTypedefs["enum "+e.CName] = e.Name
			h.lookupEnum[e.CName] = e

			h.enums = append(h.enums, e)
		case clang.CK_FunctionDecl:
			f := handleFunctionCursor(cursor)
			if f != nil {
				h.functions = append(h.functions, f)
			}
		case clang.CK_StructDecl:
			if cname == "" {
				break
			}

			s := handleStructCursor(cursor, cname, cnameIsTypeDef)

			h.lookupStruct[s.Name] = s
			h.lookupNonTypedefs["struct "+s.CName] = s.Name
			h.lookupStruct[s.CName] = s

			h.structs = append(h.structs, s)
		case clang.CK_TypedefDecl:
			underlyingType := cursor.TypedefDeclUnderlyingType().TypeSpelling()
			underlyingStructType := strings.TrimSuffix(strings.TrimPrefix(underlyingType, "struct "), " *")

			if s, ok := h.lookupStruct[underlyingStructType]; ok && !s.CNameIsTypeDef && strings.HasPrefix(underlyingType, "struct "+s.CName) {
				// Sometimes the typedef is not a parent of the struct but a sibling
				sn := handleStructCursor(cursor, cname, true)

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

				h.lookupStruct[s.Name] = s
				h.lookupNonTypedefs["struct "+s.CName] = s.Name
				h.lookupStruct[s.CName] = s

				h.structs = append(h.structs, s)
			}
		}

		return clang.CVR_Recurse
	})

	clangFile := newFile("clang")

	for _, f := range h.functions {
		// Some functions are not compiled in (TODO only 3.4?) the library see https://lists.launchpad.net/desktop-packages/msg75835.html for a never resolved bug report https://github.com/zimmski/go-clang-phoenix/issues/59
		if f.CName == "clang_CompileCommand_getMappedSourceContent" || f.CName == "clang_CompileCommand_getMappedSourcePath" || f.CName == "clang_CompileCommand_getNumMappedSources" {
			fmt.Printf("Ignore function %q because it is not compiled within libClang\n", f.CName)

			continue
		}
		// Some functions can not be handled automatically by us
		if f.CName == "clang_executeOnThread" || f.CName == "clang_getInclusions" {
			fmt.Printf("Ignore function %q because it cannot be handled automatically\n", f.CName)

			continue
		}
		// Some functions are simply manually implemented
		if f.CName == "clang_annotateTokens" || f.CName == "clang_getCursorPlatformAvailability" || f.CName == "clang_visitChildren" {
			fmt.Printf("Ignore function %q because it is manually implemented\n", f.CName)

			continue
		}

		/*
			TODO mark the enum https://github.com/zimmski/go-clang-phoenix/issues/40
				typedef enum CXChildVisitResult (*CXCursorVisitor)(CXCursor cursor, CXCursor parent, CXClientData client_data);
			as manually implemented
		*/
		// TODO report other enums like callbacks that they are not implemented https://github.com/zimmski/go-clang-phoenix/issues/51

		// Prepare the parameters
		for i := range f.Parameters {
			p := &f.Parameters[i]

			if n, ok := h.lookupNonTypedefs[p.Type.CGoName]; ok {
				p.Type.GoName = n
			}
			if e, ok := h.lookupEnum[p.Type.GoName]; ok {
				p.CName = e.Receiver.CName
				// TODO remove the receiver... and copy only names here to preserve the original pointers and so https://github.com/zimmski/go-clang-phoenix/issues/40
				p.Type.GoName = e.Receiver.Type.GoName
				p.Type.CGoName = e.Receiver.Type.CGoName
				p.Type.CGoName = e.Receiver.Type.CGoName
			} else if _, ok := h.lookupStruct[p.Type.GoName]; ok {
			}

			if f.CName == "clang_getRemappingsFromFileList" {
				switch p.CName {
				case "filePaths":
					p.Type.IsSlice = true
				case "numFiles":
					p.Type.LengthOfSlice = "filePaths"
				}

				continue
			}

			// TODO happy hack, whiteflag types that are return arguments https://github.com/zimmski/go-clang-phoenix/issues/40
			if p.Type.PointerLevel == 1 && (p.Type.GoName == "File" || p.Type.GoName == "FileUniqueID" || p.Type.GoName == "IdxClientFile" || p.Type.GoName == "cxstring" || p.Type.GoName == GoInt16 || p.Type.GoName == GoUInt16 || p.Type.GoName == "CompilationDatabase_Error" || p.Type.GoName == "PlatformAvailability" || p.Type.GoName == "SourceRange" || p.Type.GoName == "LoadDiag_Error") {
				p.Type.IsReturnArgument = true
			}
			if p.Type.PointerLevel == 2 && (p.Type.GoName == "Token" || p.Type.GoName == "Cursor") {
				p.Type.IsReturnArgument = true
			}

			if f.CName == "clang_disposeOverriddenCursors" && p.CName == "overridden" {
				p.Type.IsSlice = true
			}

			// TODO happy hack, if this is an array length parameter we need to find its partner https://github.com/zimmski/go-clang-phoenix/issues/40
			paCName := ArrayNameFromLength(p.CName)

			if paCName != "" {
				for j := range f.Parameters {
					pa := &f.Parameters[j]

					if strings.ToLower(pa.CName) == strings.ToLower(paCName) {
						if pa.Type.GoName == "struct CXUnsavedFile" || pa.Type.GoName == "UnsavedFile" {
							pa.Type.GoName = "UnsavedFile"
							pa.Type.CGoName = "struct_CXUnsavedFile"
						} else if pa.Type.CGoName == CSChar && pa.Type.PointerLevel == 2 {
						} else if pa.Type.GoName == "CompletionResult" {
						} else if pa.Type.GoName == "Token" {
						} else if pa.Type.GoName == "Cursor" {
						} else {
							break
						}

						p.Type.LengthOfSlice = pa.Name
						pa.Type.IsSlice = true

						if pa.Type.IsReturnArgument && p.Type.PointerLevel > 0 {
							p.Type.IsReturnArgument = true
						}

						break
					}
				}
			}
		}
		for i := range f.Parameters {
			p := &f.Parameters[i]

			if p.Type.CGoName == CSChar && p.Type.PointerLevel == 2 && !p.Type.IsSlice {
				p.Type.IsReturnArgument = true
			}
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
			rt.Name = receiverName(f.Parameters[0].Type.GoName)
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
			// These pointers are ok
			if p.Type.PointerLevel == 1 && (p.Type.CGoName == CSChar || p.Type.GoName == "UnsavedFile" || p.Type.GoName == "CodeCompleteResults" || p.Type.GoName == "CursorKind" || p.Type.GoName == "IdxContainerInfo" || p.Type.GoName == "IdxDeclInfo" || p.Type.GoName == "IndexerCallbacks" || p.Type.GoName == "TranslationUnit" || p.Type.GoName == "IdxEntityInfo" || p.Type.GoName == "IdxAttrInfo") {
				continue
			}
			// Return arguments are always ok since we mark them earlier
			if p.Type.IsReturnArgument {
				continue
			}
			// We whiteflag slices
			if p.Type.IsSlice {
				continue
			}

			if (!h.IsEnumOrStruct(p.Type.GoName) && !p.Type.IsPrimitive) || p.Type.PointerLevel != 0 {
				found = true

				fmt.Printf("Cannot handle parameter %s -> %#v\n", f.CName, p)

				break
			}
		}

		fname := f.Name

		// TODO happy hack we trim some whitelisted prefixes https://github.com/zimmski/go-clang-phoenix/issues/40
		if fn := strings.TrimPrefix(fname, "indexLoc_"); len(fn) != len(fname) {
			fname = fn
		} else if fn := strings.TrimPrefix(fname, "index_"); len(fn) != len(fname) {
			fname = fn
		} else if fn := strings.TrimPrefix(fname, "Location_"); len(fn) != len(fname) {
			fname = fn
		} else if fn := strings.TrimPrefix(fname, "Range_"); len(fn) != len(fname) {
			fname = fn
		} else if fn := strings.TrimPrefix(fname, "remap_"); len(fn) != len(fname) {
			fname = fn
		}

		if len(f.Parameters) > 0 && h.IsEnumOrStruct(f.Parameters[0].Type.GoName) {
			switch f.Parameters[0].Type.GoName {
			case "CodeCompleteResults":
				fname = strings.TrimPrefix(fname, "codeComplete")
			case "CompletionString":
				if f.CName == "clang_getNumCompletionChunks" {
					fname = "NumChunks"
				} else {
					fname = strings.TrimPrefix(fname, "getCompletion")
				}
			case "SourceRange":
				fname = strings.TrimPrefix(fname, "getRange")
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
						added = h.addBasicMethods(f, strings.Join(s[1:], ""), s[0]+"_", rt)
					}
				}
			}

			if !added {
				if len(f.Parameters) == 0 {
					f.Name = UpperFirstCharacter(f.Name)

					clangFile.Functions = append(clangFile.Functions, f.generate())

					added = true
				} else if h.IsEnumOrStruct(f.ReturnType.GoName) || f.ReturnType.IsPrimitive {
					fname = TrimCommonFunctionName(fname, rt.Type)
					if fn := strings.TrimPrefix(fname, f.ReturnType.GoName+"_"); len(fn) != len(fname) {
						fname = fn
					}

					added = h.AddMethod(f, fname, "", rt)

					if !added && h.IsEnumOrStruct(f.ReturnType.GoName) {
						fname = TrimCommonFunctionName(fname, rt.Type)
						if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") {
							fname = "New" + fname
						}

						rtc := rt
						rtc.Type = f.ReturnType

						added = h.AddMethod(f, fname, "", rtc)
					}
					if !added {
						f.Name = UpperFirstCharacter(f.Name)

						clangFile.Functions = append(clangFile.Functions, f.generate())

						added = true
					}
				}
			}
		}

		if !added {
			fmt.Println("Unused function:", f.Name)
		}
	}

	for _, e := range h.enums {
		if h.name != "./clang-c/Index.h" {
			e.HeaderFile = h.name
		}

		if err := e.addEnumStringMethods(); err != nil {
			return cmdFatal("Cannot generate enum string methods", err)
		}

		for i, m := range e.Methods {
			if r := h.handleMethod(e.Name, m); r != "" {
				e.Methods[i] = r
			}
			switch m := m.(type) {
			case *Function:
				if len(m.Parameters) > 0 && !m.Parameters[0].Type.IsSlice && m.Parameters[0].Type.GoName == e.Name {
					m.Receiver = Receiver{
						Name: receiverName(e.Name),
						Type: m.Parameters[0].Type,
					}

					h.setIsPointerComposition(&m.Receiver.Type)
				}
				for i := range m.Parameters {
					h.setIsPointerComposition(&m.Parameters[i].Type)
				}

				e.Methods[i] = m.generate()
			}
		}

		if err := e.generate(); err != nil {
			return cmdFatal("Cannot generate enum", err)
		}
	}

	for _, s := range h.structs {
		if h.name != "./clang-c/Index.h" {
			s.HeaderFile = h.name
		}

		if err := s.addMemberGetters(); err != nil {
			return cmdFatal("Cannot generate struct member getters", err)
		}

		for i, m := range s.Methods {
			if r := h.handleMethod(s.Name, m); r != "" {
				s.Methods[i] = r
			}
		}

		if err := s.generate(); err != nil {
			return cmdFatal("Cannot generate struct", err)
		}
	}

	if len(clangFile.Functions) > 0 {
		if h.name != "./clang-c/Index.h" {
			clangFile.HeaderFiles[h.name] = struct{}{}
		}

		if err := clangFile.generate(); err != nil {
			return cmdFatal("Cannot generate clang file", err)
		}
	}

	return nil
}

func (h *HeaderFile) handleMethod(rname string, m interface{}) string {
	switch m := m.(type) {
	case *Function:
		if len(m.Parameters) > 0 && !m.Parameters[0].Type.IsSlice && m.Parameters[0].Type.GoName == rname {
			m.Receiver = Receiver{
				Name: receiverName(rname),
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
