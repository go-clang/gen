package gen

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-clang/bootstrap/clang"
)

// HeaderFile represents a generation headerfile.
type HeaderFile struct {
	Lookup

	api *API

	Filename string
	Path     string

	Enums     []*Enum
	Functions []*Function
	Structs   []*Struct
}

// NewHeaderFile returns the new initialized HeaderFile.
func NewHeaderFile(a *API, name string, dir string) *HeaderFile {
	return &HeaderFile{
		Lookup:   NewLookup(),
		api:      a,
		Filename: name,
		Path:     dir,
	}
}

var (
	reFindStructs     = regexp.MustCompile(`(?s)struct[\s\w]+{.+?}`)
	reFindVoidPointer = regexp.MustCompile(`(?:const\s+)?void\s*\*\s*(\w+(\[\d+\])?;)`)
)

// PrepareFile prepares header files name.
func (h *HeaderFile) PrepareFile() error {
	/*
		Hide all "void *" fields of structs by replacing the type with "uintptr_t".

		To paraphrase the original go-clang source code:
		Not hiding these fields confuses the Go GC during garbage collection and
		pointer scanning, making it think the heap/stack has been somehow corrupted.

		I do not know how the original author debugged this, but one thing: Thank you!
	*/
	f, err := os.ReadFile(h.FullPath())
	if err != nil {
		return fmt.Errorf("cannot read Index.h: %w", err)
	}

	voidPointerReplacements := map[string]string{}
	for _, s := range reFindStructs.FindAll(f, -1) {
		s2 := reFindVoidPointer.ReplaceAll(s, []byte("uintptr_t $1"))
		if len(s) != len(s2) {
			voidPointerReplacements[string(s)] = string(s2)
		}
	}

	fs := string(f)
	for s, r := range voidPointerReplacements {
		fs = strings.ReplaceAll(fs, s, r)
	}

	if incl := "#include <stdint.h>"; !strings.HasPrefix(fs, incl) { // Include for uintptr_t
		fs = "#include <stdint.h>\n\n" + fs
	}

	if err = os.WriteFile(h.FullPath(), []byte(fs), 0o600); err != nil {
		return fmt.Errorf("cannot write %s: %w", h.FullPath(), err)
	}

	return nil
}

// HandleFile handles header file.
func (h *HeaderFile) HandleFile(cursor clang.Cursor) {
	// TODO(go-clang): mark the enum https://github.com/go-clang/gen/issues/40
	//  	typedef enum CXChildVisitResult (*CXCursorVisitor)(CXCursor cursor, CXCursor parent, CXClientData client_data);
	// as manually implemented

	// TODO(go-clang): report other enums like callbacks that they are not implemented
	// https://github.com/go-clang/gen/issues/51

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		// only handle code of the current file
		sourceFile, _, _, _ := cursor.Location().FileLocation()
		isCurrentFile := sourceFile.Name() == h.FullPath()

		if !strings.HasPrefix(sourceFile.Name(), h.Path) {
			return clang.ChildVisit_Continue
		}
		// TODO(zchee): Documentation.h header haven't correct cursor information
		if h.Filename == "Documentation.h" && filepath.Base(sourceFile.Name()) == "Index.h" {
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

			e := HandleEnumCursor(cursor, cname, cnameIsTypeDef)
			e.IncludeFiles.AddIncludeFile(sourceFile.Name())

			if _, ok := h.HasEnum(e.Name); !ok {
				h.Enums = append(h.Enums, e)
			}

		case clang.Cursor_FunctionDecl:
			if !isCurrentFile {
				return clang.ChildVisit_Continue
			}

			f := HandleFunctionCursor(cursor)
			if f != nil {
				f.IncludeFiles.AddIncludeFile(sourceFile.Name())
				h.Functions = append(h.Functions, f)
			}

		case clang.Cursor_StructDecl:
			if cname == "" {
				break
			}

			s := HandleStructCursor(cursor, cname, cnameIsTypeDef)
			s.api = h.api
			s.IncludeFiles.AddIncludeFile(sourceFile.Name())

			if _, ok := h.HasStruct(s.Name); !ok {
				h.RegisterStruct(s)
				h.Structs = append(h.Structs, s)
			}

		case clang.Cursor_TypedefDecl:
			underlyingType := cursor.TypedefDeclUnderlyingType().Spelling()
			underlyingStructType := strings.TrimSuffix(strings.TrimPrefix(underlyingType, "struct "), " *")

			if s, ok := h.HasStruct(underlyingStructType); ok && !s.CNameIsTypeDef && strings.HasPrefix(underlyingType, "struct "+s.CName) {
				// sometimes the typedef is not a parent of the struct but a sibling
				sn := HandleStructCursor(cursor, cname, true)
				sn.api = h.api
				sn.IncludeFiles.AddIncludeFile(sourceFile.Name())

				if sn.Comment == "" {
					sn.Comment = s.Comment
				}
				sn.Fields = s.Fields
				sn.Methods = s.Methods

				h.RemoveStruct(s)
				h.RegisterStruct(sn)

				for i, si := range h.Structs {
					if si == s {
						h.Structs[i] = sn

						break
					}
				}
			} else if underlyingType == "void *" {
				s := HandleStructCursor(cursor, cname, true)
				s.api = h.api
				s.IncludeFiles.AddIncludeFile(sourceFile.Name())

				if _, ok := h.HasStruct(s.Name); !ok {
					h.RegisterStruct(s)
					h.Structs = append(h.Structs, s)
				}
			}
		}

		return clang.ChildVisit_Recurse
	})
}

// Parse parses header file with clangArguments.
func (h *HeaderFile) Parse(clangArguments []string) error {
	if err := h.PrepareFile(); err != nil {
		return err
	}

	// parse the header file to analyse everything we need to know
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit(h.FullPath(), clangArguments, nil, 0)
	defer tu.Dispose()

	if !tu.IsValid() {
		return errors.New("cannot parse Index.h")
	}

	for _, diag := range tu.Diagnostics() {
		switch diag.Severity() {
		case clang.Diagnostic_Error:
			return fmt.Errorf("diagnostic error in Index.h: %s", diag.Spelling())

		case clang.Diagnostic_Fatal:
			return fmt.Errorf("diagnostic fatal in Index.h: %s", diag.Spelling())
		}
	}

	h.HandleFile(tu.TranslationUnitCursor())

	return nil
}

// FullPath returns the full path of h.
func (h *HeaderFile) FullPath() string {
	path := h.Path
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return path + h.Filename
}
