package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/sbinet/go-clang"
	// "github.com/termie/go-shutil"
)

var enums []*Enum
var functions []*Function
var structs []*Struct

var lookupEnum = map[string]*Enum{}
var lookupNonTypedefs = map[string]string{}
var lookupStruct = map[string]*Struct{
	"cxstring": &Struct{
		Name:  "cxstring",
		CName: "CXString",
	},
}

func trimCommonFName(fname string, rt Receiver) string {
	fname = strings.TrimPrefix(fname, "create")
	fname = strings.TrimPrefix(fname, "get")

	fname = trimClangPrefix(fname)

	if fn := strings.TrimPrefix(fname, rt.Type.GoName+"_"); len(fn) != len(fname) {
		fname = fn
	} else if fn := strings.TrimPrefix(fname, rt.Type.GoName); len(fn) != len(fname) {
		fname = fn
	} else if fn := strings.TrimSuffix(fname, rt.CName); len(fn) != len(fname) {
		fname = fn
	}

	fname = strings.TrimPrefix(fname, "create")
	fname = strings.TrimPrefix(fname, "get")

	fname = trimClangPrefix(fname)

	// If the function name is empty at this point, it is a constructor
	if fname == "" {
		fname = rt.Type.GoName
	}

	return fname
}

func addFunction(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	fname = upperFirstCharacter(fname)

	if e, ok := lookupEnum[rt.Type.GoName]; ok {
		f.Name = fnamePrefix + fname

		e.Methods = append(e.Methods, generateASTFunction(f))

		return true
	} else if s, ok := lookupStruct[rt.Type.GoName]; ok && s.CName != "CXString" {
		f.Name = fnamePrefix + fname

		fStr := generateASTFunction(f)
		s.Methods = deleteMethod(s.Methods, fname)
		s.Methods = append(s.Methods, fStr)

		return true
	}

	return false
}

func deleteMethod(methods []string, fName string) []string {
	idx := -1
	for i, mem := range methods {
		if strings.Contains(mem, ") "+fName+"()") {
			idx = i
		}
	}

	if idx != -1 {
		methods = append(methods[:idx], methods[idx+1:]...)
	}

	return methods
}

func addMethod(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	fname = upperFirstCharacter(fname)

	// TODO this is a big HACK. Figure out how we can trim receiver names while still not having two "Cursor" methods for TranslationUnit
	if f.CName == "clang_getTranslationUnitCursor" {
		fname = "TranslationUnitCursor"
	}

	if e, ok := lookupEnum[rt.Type.GoName]; ok {
		f.Name = fnamePrefix + fname
		f.Receiver = e.Receiver
		f.Receiver.Type = rt.Type

		e.Methods = append(e.Methods, generateASTFunction(f))

		return true
	} else if s, ok := lookupStruct[rt.Type.GoName]; ok && s.CName != "CXString" {
		f.Name = fnamePrefix + fname
		f.Receiver = s.Receiver
		f.Receiver.Type = rt.Type

		fStr := generateASTFunction(f)
		s.Methods = deleteMethod(s.Methods, fname)
		s.Methods = append(s.Methods, fStr)

		return true
	}

	return false
}

func addBasicMethods(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	if len(f.Parameters) == 0 && isEnumOrStruct(f.ReturnType.GoName) {
		fname = trimCommonFName(fname, rt)
		if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") {
			fname = "New" + fname
		}

		return addMethod(f, fname, fnamePrefix, rt)
	} else if (fname[0] == 'i' && fname[1] == 's' && unicode.IsUpper(rune(fname[2]))) || (fname[0] == 'h' && fname[1] == 'a' && fname[2] == 's' && unicode.IsUpper(rune(fname[3]))) {
		f.ReturnType.GoName = "bool"

		return addMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 1 && isEnumOrStruct(f.Parameters[0].Type.GoName) && strings.HasPrefix(fname, "dispose") && f.ReturnType.GoName == "void" {
		fname = "Dispose"

		return addMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 2 && strings.HasPrefix(fname, "equal") && isEnumOrStruct(f.Parameters[0].Type.GoName) && f.Parameters[0].Type == f.Parameters[1].Type {
		f.Parameters[0].Name = receiverName(f.Parameters[0].Type.GoName)
		f.Parameters[1].Name = f.Parameters[0].Name + "2"

		f.ReturnType.GoName = "bool"

		return addMethod(f, fname, fnamePrefix, rt)
	}

	return false
}

func isEnumOrStruct(name string) bool {
	if _, ok := lookupEnum[name]; ok {
		return true
	} else if _, ok := lookupStruct[name]; ok {
		return true
	}

	return false
}

func main() {
	rawLLVMVersion, _, err := execToBuffer("llvm-config", "--version")
	if err != nil {
		exitWithFatal("Cannot determine LLVM version", err)
	}

	matchLLVMVersion := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`).FindSubmatch(rawLLVMVersion)
	if matchLLVMVersion == nil {
		exitWithFatal("Cannot parse LLVM version", nil)
	}

	var llvmVersion struct {
		Major int
		Minor int
		Patch int // TODO rename to Subminor
	}

	llvmVersion.Major, _ = strconv.Atoi(string(matchLLVMVersion[1]))
	llvmVersion.Minor, _ = strconv.Atoi(string(matchLLVMVersion[2]))
	llvmVersion.Patch, _ = strconv.Atoi(string(matchLLVMVersion[3]))

	fmt.Println("Found LLVM version", string(matchLLVMVersion[0]))

	rawLLVMIncludeDir, _, err := execToBuffer("llvm-config", "--includedir")
	if err != nil {
		exitWithFatal("Cannot determine LLVM include directory", err)
	}

	clangCIncludeDir := strings.TrimSpace(string(rawLLVMIncludeDir)) + "/clang-c/"
	if err := dirExists(clangCIncludeDir); err != nil {
		exitWithFatal(fmt.Sprintf("Cannot find Clang-C include directory %q", clangCIncludeDir), err)
	}

	fmt.Println("Clang-C include directory", clangCIncludeDir)

	fmt.Printf("Will generate go-clang for LLVM version %d.%d in current directory\n", llvmVersion.Major, llvmVersion.Minor)

	/*// Copy the Clang-C include directory into the current directory
	_ = os.RemoveAll("./clang-c/")
	if err := shutil.CopyTree(clangCIncludeDir, "./clang-c/", nil); err != nil {
		exitWithFatal(fmt.Sprintf("Cannot copy Clang-C include directory %q into current directory", clangCIncludeDir), err)
	}*/

	// Remove all generated .go files
	if files, err := ioutil.ReadDir("./"); err != nil {
		exitWithFatal("Cannot read current directory", err)
	} else {
		for _, f := range files {
			fn := f.Name()

			if !f.IsDir() && strings.HasSuffix(fn, "_gen.go") {
				if err := os.Remove(fn); err != nil {
					exitWithFatal(fmt.Sprintf("Cannot remove generated file %q", fn), err)
				}
			}
		}
	}

	clangIndexHeaderFilepath := "./clang-c/Index.h"

	/*
		Hide all "void *" fields of structs by replacing the type with "uintptr_t".

		To paraphrase the original go-clang source code:
			Not hiding these fields confuses the Go GC during garbage collection and
			pointer scanning, making it think the heap/stack has been somehow corrupted.

		I do not know how the original author debugged this, but one thing: Thank you!
	*/
	findStructsRe := regexp.MustCompile(`(?s)struct[\s\w]+{.+?}`)
	f, err := ioutil.ReadFile(clangIndexHeaderFilepath)
	if err != nil {
		exitWithFatal("Cannot read Index.h", nil)
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
	if incl := "#include <stdint.h>"; !strings.HasPrefix(fs, incl) {
		fs = "#include <stdint.h>\n\n" + fs
	}
	err = ioutil.WriteFile(clangIndexHeaderFilepath, []byte(fs), 0700)
	if err != nil {
		exitWithFatal("Cannot write Index.h", nil)
	}

	// Parse clang-c's Index.h to analyse everything we need to know
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.Parse(clangIndexHeaderFilepath, []string{
		"-I", ".", // Include current folder
		"-I", "/usr/local/lib/clang/3.4.2/include/",
		"-I", "/usr/include/clang/3.6.2/include/",
	}, nil, 0)
	defer tu.Dispose()

	if !tu.IsValid() {
		exitWithFatal("Cannot parse Index.h", nil)
	}

	for _, diag := range tu.Diagnostics() {
		switch diag.Severity() {
		case clang.Diagnostic_Error:
			exitWithFatal("Diagnostic error in Index.h", errors.New(diag.Spelling()))
		case clang.Diagnostic_Fatal:
			exitWithFatal("Diagnostic fatal in Index.h", errors.New(diag.Spelling()))
		}
	}

	/*
		TODO mark the enum
			typedef enum CXChildVisitResult (*CXCursorVisitor)(CXCursor cursor, CXCursor parent, CXClientData client_data);
		as manually implemented
	*/

	/*
		TODO mark the function
			unsigned clang_visitChildren(CXCursor parent, CXCursorVisitor visitor, CXClientData client_data);
		as manually implemented
	*/

	cursor := tu.ToCursor()
	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		// Only handle code of the current file
		sourceFile, _, _, _ := cursor.Location().GetFileLocation()
		if sourceFile.Name() != clangIndexHeaderFilepath {
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

			lookupEnum[e.Name] = e
			lookupNonTypedefs["enum "+e.CName] = e.Name
			lookupEnum[e.CName] = e

			enums = append(enums, e)
		case clang.CK_FunctionDecl:
			f := handleFunctionCursor(cursor)
			if f != nil {
				functions = append(functions, f)
			}
		case clang.CK_StructDecl:
			if cname == "" {
				break
			}

			s := handleStructCursor(cursor, cname, cnameIsTypeDef)

			lookupStruct[s.Name] = s
			lookupNonTypedefs["struct "+s.CName] = s.Name
			lookupStruct[s.CName] = s

			structs = append(structs, s)
		case clang.CK_TypedefDecl:
			underlyingType := cursor.TypedefDeclUnderlyingType().TypeSpelling()
			underlyingStructType := strings.TrimSuffix(strings.TrimPrefix(underlyingType, "struct "), " *")

			if s, ok := lookupStruct[underlyingStructType]; ok && !s.CNameIsTypeDef && strings.HasPrefix(underlyingType, "struct "+s.CName) {
				// Sometimes the typedef is not a parent of the struct but a sibling TODO find out if this is a bug?

				sn := handleVoidStructCursor(cursor, cname, true)

				lookupStruct[sn.Name] = sn
				lookupNonTypedefs["struct "+sn.CName] = sn.Name
				lookupStruct[sn.CName] = sn

				// Update the lookups for the old struct
				lookupStruct[s.Name] = sn
				lookupStruct[s.CName] = sn

				for i, si := range structs {
					if si == s {
						structs[i] = sn

						break
					}
				}
			} else if underlyingType == "void *" {
				s := handleVoidStructCursor(cursor, cname, true)

				lookupStruct[s.Name] = s
				lookupNonTypedefs["struct "+s.CName] = s.Name
				lookupStruct[s.CName] = s

				structs = append(structs, s)
			}
		}

		return clang.CVR_Recurse
	})

	clangFile := &File{
		Name: "clang",

		Imports: map[string]struct{}{},
	}

	for _, f := range functions {
		fname := f.Name

		// Prepare the parameters
		for i := range f.Parameters {
			p := &f.Parameters[i]

			if n, ok := lookupNonTypedefs[p.Type.CGoName]; ok {
				p.Type.GoName = n
			}
			if e, ok := lookupEnum[p.Type.GoName]; ok {
				p.CName = e.Receiver.CName
				// TODO remove the receiver... and copy only names here to preserve the original pointers and so
				p.Type.GoName = e.Receiver.Type.GoName
				p.Type.CGoName = e.Receiver.Type.CGoName
				p.Type.CGoName = e.Receiver.Type.CGoName
			} else if _, ok := lookupStruct[p.Type.GoName]; ok {
			}

			// TODO happy hack, whiteflag types that are return arguments
			if p.Type.PointerLevel == 1 && (p.Type.GoName == "File" || p.Type.GoName == "FileUniqueID" || p.Type.GoName == "IdxClientFile" || p.Type.GoName == "cxstring" || p.Type.GoName == GoUInt16) {
				p.Type.IsReturnArgument = true
			}

			// TODO happy hack, if this is an array length parameter we need to find its partner
			if paName := strings.TrimPrefix(p.Name, "num_"); len(paName) != len(p.Name) {
				for j := range f.Parameters {
					pa := &f.Parameters[j]

					if pa.Name == paName {

						// TODO remove this when getType cane handle this kind of conversion
						if pa.Type.GoName == "struct CXUnsavedFile" || pa.Type.GoName == "UnsavedFile" {
							pa.Type.GoName = "UnsavedFile"
							pa.Type.CGoName = "struct_CXUnsavedFile"
						} else if pa.Type.CGoName == CSChar && pa.Type.PointerLevel == 2 {
						} else if pa.Type.GoName == "CompletionResult" || pa.Type.GoName == "Token" {
							pa.Type.CGoName = "struct_CX" + pa.Type.GoName
						} else {
							break
						}

						p.Type.LengthOfSlice = pa.Name
						pa.Type.IsSlice = true

						break
					}
				}
			}
		}

		// Prepare the return argument
		if n, ok := lookupNonTypedefs[f.ReturnType.CGoName]; ok {
			f.ReturnType.GoName = n
		}
		if e, ok := lookupEnum[f.ReturnType.GoName]; ok {
			f.ReturnType.CGoName = e.Receiver.Type.CGoName
		} else if _, ok := lookupStruct[f.ReturnType.GoName]; ok {
		}

		// Prepare the receiver
		var rt Receiver
		if len(f.Parameters) > 0 {
			rt.Name = receiverName(f.Parameters[0].Type.GoName)
			rt.CName = f.Parameters[0].CName
			rt.Type = f.Parameters[0].Type
		} else {
			if e, ok := lookupEnum[f.ReturnType.GoName]; ok {
				rt.Type = e.Receiver.Type
			} else if s, ok := lookupStruct[f.ReturnType.GoName]; ok {
				rt.Type.GoName = s.Name
			}
		}

		// Check upfront if we can handle a function
		found := false

		for _, p := range f.Parameters {
			// These pointers are ok
			if p.Type.PointerLevel == 1 && (p.Type.CGoName == CSChar || p.Type.GoName == "UnsavedFile") {
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

			if (!isEnumOrStruct(p.Type.GoName) && !p.Type.IsPrimitive) || p.Type.PointerLevel != 0 {
				found = true

				break
			}
		}

		if f.ReturnType.PointerLevel > 0 && !(f.ReturnType.PointerLevel == 1 && f.ReturnType.CGoName == CSChar) { // TODO implement to return slices
			found = true
		}

		// If we find a heuristic to add the function, add it!
		added := false

		if !found {
			added = addBasicMethods(f, fname, "", rt)

			if !added {
				if s := strings.Split(f.Name, "_"); len(s) == 2 {
					if s[0] == rt.Type.GoName {
						rtc := rt
						rtc.Name = s[0]

						added = addBasicMethods(f, s[1], "", rtc)
					} else {
						added = addBasicMethods(f, strings.Join(s[1:], ""), s[0]+"_", rt)
					}
				}
			}

			if !added {
				if len(f.Parameters) == 0 {
					clangFile.Functions = append(clangFile.Functions, generateASTFunction(f))

					added = true
				} else if isEnumOrStruct(f.ReturnType.GoName) || f.ReturnType.IsPrimitive {
					fname = trimCommonFName(fname, rt)

					added = addMethod(f, fname, "", rt)

					if !added && isEnumOrStruct(f.ReturnType.GoName) {
						fname = trimCommonFName(fname, rt)
						if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") {
							fname = "New" + fname
						}

						rtc := rt
						rtc.Type = f.ReturnType

						added = addFunction(f, fname, "", rtc)
					}
					if !added {
						clangFile.Functions = append(clangFile.Functions, generateASTFunction(f))

						added = true
					}
				}
			}
		}

		if !added {
			fmt.Println("Unused function:", f.Name)
		}
	}

	for _, e := range enums {
		if err := generateEnum(e); err != nil {
			exitWithFatal("Cannot generate enum", err)
		}
	}

	for _, s := range structs {
		if err := generateStruct(s); err != nil {
			exitWithFatal("Cannot generate struct", err)
		}
	}

	if len(clangFile.Functions) > 0 {
		if err := generateFile(clangFile); err != nil {
			exitWithFatal("Cannot generate clang file", err)
		}
	}

	if out, _, err := execToBuffer("gofmt", "-w", "./"); err != nil { // TODO do this before saving the files using go/fmt
		fmt.Printf("gofmt:\n%s\n", out)

		exitWithFatal("Gofmt failed", err)
	}
}

func printFunctionDetails(f *Function) {
	fmt.Printf("@@ %s %#v %#v\n", f.CName, f.ReturnType, f.Parameters)
}
