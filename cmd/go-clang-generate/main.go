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

	addMethod := func(f *Function, method func(f *Function) string) {
		if e, ok := lookupEnum[f.ReceiverType]; ok {
			f.Receiver = e.Receiver
			f.ReceiverPrimitiveType = e.ReceiverPrimitiveType

			e.Methods = append(e.Methods, method(f))
		} else if s, ok := lookupStruct[f.ReceiverType]; ok {
			f.Receiver = s.Receiver

			s.Methods = append(s.Methods, method(f))
		}
	}

	for _, f := range functions {
		var rt string
		if len(f.Parameters) > 0 {
			rt = trimClangPrefix(f.Parameters[0].Type)
			if n, ok := lookupNonTypedefs[rt]; ok {
				rt = n
			}
		}

		if len(f.Parameters) == 1 && f.ReturnType == "CXString" {
			f.ReceiverType = rt

			f.Name = strings.TrimPrefix(f.Name, f.ReceiverType+"_")

			f.Name = strings.TrimPrefix(f.Name, "get")
			f.Name = strings.TrimPrefix(f.Name, f.ReceiverType)

			f.Name = upperFirstCharacter(f.Name)

			addMethod(f, generateFunctionStringGetter)
		} else if len(f.Parameters) == 1 && f.Name[0] == 'i' && f.Name[1] == 's' && unicode.IsUpper(rune(f.Name[2])) && f.ReturnType == "unsigned int" {
			f.ReceiverType = rt

			f.Name = upperFirstCharacter(f.Name)

			addMethod(f, generateGenerateFunctionIs)
		} else if len(f.Parameters) == 1 && strings.HasPrefix(f.Name, "dispose") && f.Name[len("dispose"):] == rt {
			f.ReceiverType = rt

			f.Name = "Dispose"

			addMethod(f, generateFunctionVoidMethod)
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
