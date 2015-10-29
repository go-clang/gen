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
var lookupStruct = map[string]*Struct{}

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

	if e, ok := lookupEnum[rt.Type.Name]; ok {
		f.Name = fnamePrefix + fname

		e.Methods = append(e.Methods, generateASTFunction(f))

		return true
	} else if s, ok := lookupStruct[rt.Type.Name]; ok {
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
		if strings.Contains(mem, fName+"()") {
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

	if e, ok := lookupEnum[rt.Type.Name]; ok {
		f.Name = fnamePrefix + fname
		f.Receiver = e.Receiver
		f.Receiver.Type = rt.Type

		e.Methods = append(e.Methods, generateASTFunction(f))

		return true
	} else if s, ok := lookupStruct[rt.Type.Name]; ok {
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

	// Parse clang-c's Index.h to analyse everything we need to know
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	clangIndexHeaderFilepath := "./clang-c/Index.h"
	tu := idx.Parse(clangIndexHeaderFilepath, []string{
		"-I", ".", // Include current folder
		"-I", "/usr/include/clang/3.6.2/include/", // Include clang headers TODO make this generic
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

		for i := range f.Parameters {
			p := &f.Parameters[i]

			if n, ok := lookupNonTypedefs[p.Type.Name]; ok {
				p.Type.Name = n
			}
			if e, ok := lookupEnum[p.Type.Name]; ok {
				p.CName = e.Receiver.CName
				p.Type = e.Receiver.Type
			} else if _, ok := lookupStruct[p.Type.Name]; ok {
			} else {
				if goType, primitiveType := goAndTypePrimitive(p.Type.Name); goType != "" {
					p.Type.Name = goType
					p.Type.Primitive = primitiveType
				}
			}
		}

		if n, ok := lookupNonTypedefs[f.ReturnType.Name]; ok {
			f.ReturnType.Name = n
		}
		if e, ok := lookupEnum[f.ReturnType.Name]; ok {
			f.ReturnType.Primitive = e.Receiver.Type.Primitive
		} else if _, ok := lookupStruct[f.ReturnType.Name]; ok {
		}
		if goType, primitiveType := goAndTypePrimitive(f.ReturnType.Name); goType != "" {
			f.ReturnType.Name = goType
			f.ReturnType.Primitive = primitiveType
		}
		if f.ReturnType.Name == "cxstring" {
			f.ReturnType.Name = "string"
		}

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
			if len(f.Parameters) == 1 && (f.ReturnType.Name == "int" || f.ReturnType.Name == "unsigned int" || f.ReturnType.Name == "long long" || f.ReturnType.Name == "unsigned long long") && isEnumOrStruct(f.Parameters[0].Type.Name) {
				fname = trimCommonFName(fname, rt)

				f.ReturnType.Primitive = f.ReturnType.Name
				switch f.ReturnType.Name { // TODO refactor to use getTypeConversion(...) somehow, maybe do the conversion during the create of a Function instance
				case "int":
					f.ReturnType.Name = "uint16"
				case "unsigned int":
					f.ReturnType.Name = "uint16"
				case "long long":
					f.ReturnType.Name = "int64"
				case "unsigned long long":
					f.ReturnType.Name = "uint64"
				}

				added = addMethod(f, fname, "", rt)
			}
		}

		if !added {
			if len(f.Parameters) > 0 && (isEnumOrStruct(f.ReturnType.Name) || f.ReturnType.Primitive != "") {
				found := false
				for _, p := range f.Parameters {
					if !isEnumOrStruct(p.Type.Name) && p.Type.Primitive == "" {
						found = true

						break
					}
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

	if _, _, err = execToBuffer("gofmt", "-w", "./"); err != nil { // TODO do this before saving the files using go/fmt
		exitWithFatal("Gofmt failed", err)
	}
}

func goAndTypePrimitive(typ string) (string, string) {
	switch typ {
	case "int":
		return "uint16", "int"
	case "unsigned int":
		return "uint16", "uint"
	case "long long":
		return "int64", "longlong"
	case "unsigned long long":
		return "uint64", "ulonglong"
	case "void":
		return "void", "void"
	case "time_t":
		return "time.Time", "time_t"
	case "const char *":
		return "string", "const char *"
	case "String":
		return "cxstring", "cxstring"
	}

	return "", ""
}
