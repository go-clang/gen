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

	if fn := strings.TrimPrefix(fname, rt.Type.Name+"_"); len(fn) != len(fname) {
		fname = fn
	} else if fn := strings.TrimPrefix(fname, rt.Type.Name); len(fn) != len(fname) {
		fname = fn
	} else if fn := strings.TrimSuffix(fname, rt.CName); len(fn) != len(fname) {
		fname = fn
	}

	fname = strings.TrimPrefix(fname, "create")
	fname = strings.TrimPrefix(fname, "get")

	fname = trimClangPrefix(fname)

	// If the function name is empty at this point, it is a constructor
	if fname == "" {
		fname = rt.Type.Name
	}

	return fname
}

func addFunction(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
	fname = upperFirstCharacter(fname)

	if !hasHandleablePointers(f.Parameters) {
		return false
	}

	if e, ok := lookupEnum[rt.Type.Name]; ok {
		f.Name = fnamePrefix + fname

		e.Methods = append(e.Methods, generateASTFunction(f))

		return true
	} else if s, ok := lookupStruct[rt.Type.Name]; ok && s.CName != "CXString" {
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

	if !hasHandleablePointers(f.Parameters) {
		return false
	}

	if e, ok := lookupEnum[rt.Type.Name]; ok {
		f.Name = fnamePrefix + fname
		f.Receiver = e.Receiver
		f.Receiver.Type = rt.Type

		e.Methods = append(e.Methods, generateASTFunction(f))

		return true
	} else if s, ok := lookupStruct[rt.Type.Name]; ok && s.CName != "CXString" {
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
	if len(f.Parameters) == 0 && isEnumOrStruct(f.ReturnType.Name) {
		fname = trimCommonFName(fname, rt)
		if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") {
			fname = "New" + fname
		}

		return addMethod(f, fname, fnamePrefix, rt)
	} else if (fname[0] == 'i' && fname[1] == 's' && unicode.IsUpper(rune(fname[2]))) || (fname[0] == 'h' && fname[1] == 'a' && fname[2] == 's' && unicode.IsUpper(rune(fname[3]))) {
		f.ReturnType.Name = "bool"

		return addMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 1 && isEnumOrStruct(f.Parameters[0].Type.Name) && strings.HasPrefix(fname, "dispose") && f.ReturnType.Name == "void" {
		fname = "Dispose"

		return addMethod(f, fname, fnamePrefix, rt)
	} else if len(f.Parameters) == 2 && strings.HasPrefix(fname, "equal") && isEnumOrStruct(f.Parameters[0].Type.Name) && f.Parameters[0].Type == f.Parameters[1].Type {
		f.Parameters[0].Name = receiverName(f.Parameters[0].Type.Name)
		f.Parameters[1].Name = f.Parameters[0].Name + "2"

		f.ReturnType.Name = "bool"

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
	fs = "#include <stdint.h>\n\n" + fs
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

			if n, ok := lookupNonTypedefs[p.Type.CName]; ok {
				p.Type.Name = n
			}
			if e, ok := lookupEnum[p.Type.Name]; ok {
				p.CName = e.Receiver.CName
				p.Type = e.Receiver.Type
			} else if _, ok := lookupStruct[p.Type.Name]; ok {
			}

			// TODO happy hack, whiteflag types that are return arguments
			if p.Type.PointerLevel == 1 && (p.Type.Name == "File" || p.Type.Name == "FileUniqueID" || p.Type.Name == "IdxClientFile" || p.Type.Name == "cxstring" || p.Type.Name == GoUInt16) {
				p.Type.IsReturnArgument = true
			}

			// TODO happy hack, if this is an array length parameter we need to find its partner
			if paName := strings.TrimPrefix(p.Name, "num_"); len(paName) != len(p.Name) {
				for j := range f.Parameters {
					pa := &f.Parameters[j]

					// TODO we handle only incoming array pointers for now
					if pa.Type.PointerLevel != 1 && pa.Type.CName != "const char *const *" {
						continue
					}

					if pa.Name == paName {
						p.Type.LengthOfSlice = pa.Name
						pa.Type.IsSlice = true

						// TODO remove this when getType cane handle this kind of conversion
						switch pa.Type.CName {
						case "const char *const *":
							pa.Type.Name = GoUInt8
							pa.Type.Primitive = "char"
						case "struct CXUnsavedFile *":
							pa.Type.Name = "UnsavedFile"
							pa.Type.Primitive = "struct_CXUnsavedFile"
						default:
							panic(pa.Type.CName)
						}

						break
					}
				}
			}
		}

		// Prepare the return argument
		if n, ok := lookupNonTypedefs[f.ReturnType.CName]; ok {
			f.ReturnType.Name = n
		}
		if e, ok := lookupEnum[f.ReturnType.Name]; ok {
			f.ReturnType.Primitive = e.Receiver.Type.Primitive
		} else if _, ok := lookupStruct[f.ReturnType.Name]; ok {
		}

		// Prepare the receiver
		var rt Receiver
		if len(f.Parameters) > 0 {
			rt.Name = receiverName(f.Parameters[0].Type.Name)
			rt.CName = f.Parameters[0].CName
			rt.Type = f.Parameters[0].Type
		} else {
			if e, ok := lookupEnum[f.ReturnType.Name]; ok {
				rt.Type = e.Receiver.Type
			} else if s, ok := lookupStruct[f.ReturnType.Name]; ok {
				rt.Type.Name = s.Name
			}
		}

		added := addBasicMethods(f, fname, "", rt)

		if !added {
			if s := strings.Split(f.Name, "_"); len(s) == 2 {
				if s[0] == rt.Type.Name {
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
				if !added {
					clangFile.Functions = append(clangFile.Functions, generateASTFunction(f))

					added = true
				}
			} else if isEnumOrStruct(f.ReturnType.Name) || f.ReturnType.Primitive != "" {
				found := false

				for _, p := range f.Parameters {
					if !isEnumOrStruct(p.Type.Name) && p.Type.Primitive == "" {
						found = true

						break
					}
				}

				if f.ReturnType.PointerLevel > 0 { // TODO implement to return slices
					found = true
				}

				if !found {
					fname = trimCommonFName(fname, rt)

					added = addMethod(f, fname, "", rt)

					if !added && isEnumOrStruct(f.ReturnType.Name) {
						fname = trimCommonFName(fname, rt)
						if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") {
							fname = "New" + fname
						}

						rtc := rt
						rtc.Type = f.ReturnType

						added = addFunction(f, fname, "", rtc)
					}
					if !added {
						if hasHandleablePointers(f.Parameters) {
							clangFile.Functions = append(clangFile.Functions, generateASTFunction(f))

							added = true
						}
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

func hasHandleablePointers(params []FunctionParameter) bool {
	for _, p := range params {
		if p.Type.IsSlice && (p.Type.PointerLevel == 1 || (p.Type.PointerLevel == 2 && p.Type.CName == "const char *const *")) { // TODO we can handle currently only ingoing array pointers
			continue
		}

		if p.Type.PointerLevel > 0 && !p.Type.IsReturnArgument {
			return false
		}
	}

	return true
}

func printFunctionDetails(f *Function) {
	fmt.Printf("@@ %s %#v %#v\n", f.CName, f.ReturnType, f.Parameters)
}
