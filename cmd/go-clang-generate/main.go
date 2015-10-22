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
		"-I", "/usr/local/lib/clang/3.4.2/include/", // Include clang headers TODO make this generic
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

	var enums []*Enum
	var functions []*Function
	var structs []*Struct

	lookupEnum := map[string]*Enum{}
	lookupNonTypedefs := map[string]string{}
	lookupStruct := map[string]*Struct{}

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

			switch cname {
			// TODO ignore declarations like "typedef struct CXTranslationUnitImpl *CXTranslationUnit" for now
			case "CXCursorSetImpl", "CXTranslationUnitImpl":
				return clang.CVR_Recurse
			}

			s := handleStructCursor(cursor, cname, cnameIsTypeDef)

			lookupStruct[s.Name] = s
			lookupNonTypedefs["struct "+s.CName] = s.Name
			lookupStruct[s.CName] = s

			structs = append(structs, s)
		case clang.CK_TypedefDecl:
			if cursor.TypedefDeclUnderlyingType().TypeSpelling() == "void *" {
				s := handleVoidStructCursor(cursor, cname, true)

				lookupStruct[s.Name] = s
				lookupNonTypedefs["struct "+s.CName] = s.Name
				lookupStruct[s.CName] = s

				structs = append(structs, s)
			}
		}

		return clang.CVR_Recurse
	})

	addMethod := func(f *Function, fname string, rt string, method func(f *Function) string) bool {
		fname = upperFirstCharacter(fname)

		if e, ok := lookupEnum[rt]; ok {
			f.Name = fname
			f.Receiver = e.Receiver
			f.ReceiverType = rt
			f.ReceiverPrimitiveType = e.ReceiverPrimitiveType

			e.Methods = append(e.Methods, method(f))

			return true
		} else if s, ok := lookupStruct[rt]; ok {
			f.Name = fname
			f.ReceiverType = rt
			f.Receiver = s.Receiver

			s.Methods = append(s.Methods, method(f))

			return true
		}

		return false
	}

	addBasicMethods := func(f *Function, fname string, rt string) bool {
		if len(f.Parameters) == 1 && f.ReturnType == "String" {
			fname = strings.TrimPrefix(fname, rt+"_")

			fname = strings.TrimPrefix(fname, "get")
			fname = strings.TrimPrefix(fname, rt)

			return addMethod(f, fname, rt, generateFunctionStringGetter)
		} else if len(f.Parameters) == 1 && fname[0] == 'i' && fname[1] == 's' && unicode.IsUpper(rune(fname[2])) && f.ReturnType == "unsigned int" {
			return addMethod(f, fname, rt, generateGenerateFunctionIs)
		} else if len(f.Parameters) == 1 && strings.HasPrefix(fname, "dispose") && f.ReturnType == "void" && (fname == "dispose" || fname[len("dispose"):] == rt) {
			fname = "Dispose"

			return addMethod(f, fname, rt, generateFunctionVoidMethod)
		} else if len(f.Parameters) == 2 && strings.HasPrefix(fname, "equal") && f.ReturnType == "unsigned int" && f.Parameters[0].Type == f.Parameters[1].Type {
			return addMethod(f, fname, rt, generateFunctionEqual)
		}

		return false
	}

	for _, f := range functions {
		fname := f.Name
		var rt string
		if len(f.Parameters) > 0 {
			rt = trimClangPrefix(f.Parameters[0].Type)
			if n, ok := lookupNonTypedefs[rt]; ok {
				rt = n
			}
		}

		added := addBasicMethods(f, fname, rt)

		if !added {
			if s := strings.SplitN(f.Name, "_", 2); len(s) == 2 {
				if s[0] == rt {
					added = addBasicMethods(f, s[1], s[0])
				}
			}
		}

		if !added {
			fmt.Println("Unused:", f.Name)
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

	if _, _, err = execToBuffer("gofmt", "-w", "./"); err != nil { // TODO do this before saving the files using go/fmt
		exitWithFatal("Gofmt failed", err)
	}
}
