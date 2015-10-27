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

	isEnumOrStruct := func(name string) bool {
		if _, ok := lookupEnum[name]; ok {
			return true
		} else if _, ok := lookupStruct[name]; ok {
			return true
		}

		return false
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

	trimCommonFName := func(fname string, rt Receiver) string {
		fname = strings.TrimPrefix(fname, "create")
		fname = strings.TrimPrefix(fname, "get")

		fname = trimClangPrefix(fname)

		if fn := strings.TrimPrefix(fname, rt.Type+"_"); len(fn) != len(fname) {
			fname = fn
		} else if fn := strings.TrimPrefix(fname, rt.Type); len(fn) != len(fname) {
			fname = fn
		} else if fn := strings.TrimSuffix(fname, rt.CName); len(fn) != len(fname) {
			fname = fn
		}

		fname = strings.TrimPrefix(fname, "create")
		fname = strings.TrimPrefix(fname, "get")

		fname = trimClangPrefix(fname)

		// If the function name is empty at this point, it is a constructor
		if fname == "" {
			fname = rt.Type
		}

		return fname
	}

	addMethod := func(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
		fname = upperFirstCharacter(fname)

		// TODO this is a big HACK. Figure out how we can trim receiver names while still not having two "Cursor" methods for TranslationUnit
		if f.CName == "clang_getTranslationUnitCursor" {
			fname = "TranslationUnitCursor"
		}

		if e, ok := lookupEnum[rt.Type]; ok {
			f.Name = fnamePrefix + fname
			f.Receiver = e.Receiver
			f.Receiver.Type = rt.Type

			e.Methods = append(e.Methods, generateASTFunction(f))

			return true
		} else if s, ok := lookupStruct[rt.Type]; ok {
			f.Name = fnamePrefix + fname
			f.Receiver = s.Receiver
			f.Receiver.Type = rt.Type

			s.Methods = append(s.Methods, generateASTFunction(f))

			return true
		}

		return false
	}

	addBasicMethods := func(f *Function, fname string, fnamePrefix string, rt Receiver) bool {
		if len(f.Parameters) == 0 && isEnumOrStruct(f.ReturnType) {
			fname = trimCommonFName(fname, rt)
			if strings.HasPrefix(f.CName, "clang_create") || strings.HasPrefix(f.CName, "clang_get") {
				fname = "New" + fname
			}

			return addMethod(f, fname, fnamePrefix, rt)
		} else if (fname[0] == 'i' && fname[1] == 's' && unicode.IsUpper(rune(fname[2]))) || (fname[0] == 'h' && fname[1] == 'a' && fname[2] == 's' && unicode.IsUpper(rune(fname[3]))) {
			f.ReturnType = "bool"

			return addMethod(f, fname, fnamePrefix, rt)
		} else if len(f.Parameters) == 1 && strings.HasPrefix(fname, "dispose") && f.ReturnType == "void" {
			fname = "Dispose"

			return addMethod(f, fname, fnamePrefix, rt)
		} else if len(f.Parameters) == 2 && strings.HasPrefix(fname, "equal") && isEnumOrStruct(f.Parameters[0].Type) && f.Parameters[0].Type == f.Parameters[1].Type {
			f.Parameters[0].Name = receiverName(f.Parameters[0].Type)
			f.Parameters[1].Name = f.Parameters[0].Name + "2"

			f.ReturnType = "bool"

			return addMethod(f, fname, fnamePrefix, rt)
		}

		return false
	}

	for _, f := range functions {
		fname := f.Name

		for i := range f.Parameters {
			p := &f.Parameters[i]

			if n, ok := lookupNonTypedefs[p.Type]; ok {
				p.Type = n
			}
			if e, ok := lookupEnum[p.Type]; ok {
				p.CName = e.Receiver.CName
				p.Type = e.Receiver.Type
				p.CType = e.Receiver.CType
				p.PrimitiveType = e.Receiver.PrimitiveType
			} else if _, ok := lookupStruct[p.Type]; ok {
			} else {
				if goType, primitiveType := goAndPrimitiveType(p.Type); goType != "" {
					p.Type = goType
					p.PrimitiveType = primitiveType
				}
			}
		}

		if n, ok := lookupNonTypedefs[f.ReturnType]; ok {
			f.ReturnType = n
		}
		if e, ok := lookupEnum[f.ReturnType]; ok {
			f.ReturnPrimitiveType = e.Receiver.PrimitiveType
		} else if _, ok := lookupStruct[f.ReturnType]; ok {
		}
		if goType, primitiveType := goAndPrimitiveType(f.ReturnType); goType != "" {
			f.ReturnType = goType
			f.ReturnPrimitiveType = primitiveType
		}

		var rt Receiver
		if len(f.Parameters) > 0 {
			rt.Name = receiverName(f.Parameters[0].Type)
			rt.CName = f.Parameters[0].CName
			rt.Type = f.Parameters[0].Type
			rt.CType = f.Parameters[0].CType
			rt.PrimitiveType = f.Parameters[0].PrimitiveType
		} else {
			if e, ok := lookupEnum[f.ReturnType]; ok {
				rt.Type = e.Receiver.Type
			} else if s, ok := lookupStruct[f.ReturnType]; ok {
				rt.Type = s.Name
			}
		}

		added := addBasicMethods(f, fname, "", rt)

		if !added {
			if s := strings.Split(f.Name, "_"); len(s) == 2 {
				if s[0] == rt.Type {
					rtc := rt
					rtc.Name = s[0]

					added = addBasicMethods(f, s[1], "", rtc)
				} else {
					added = addBasicMethods(f, strings.Join(s[1:], ""), s[0]+"_", rt)
				}
			}
		}
		if !added {
			if len(f.Parameters) == 1 && (f.ReturnType == "int" || f.ReturnType == "unsigned int" || f.ReturnType == "long long" || f.ReturnType == "unsigned long long") && isEnumOrStruct(f.Parameters[0].Type) {
				fname = trimCommonFName(fname, rt)

				f.ReturnPrimitiveType = f.ReturnType
				switch f.ReturnType { // TODO refactor to use getTypeConversion(...) somehow, maybe do the conversion during the create of a Function instance
				case "int":
					f.ReturnType = "uint16"
				case "unsigned int":
					f.ReturnType = "uint16"
				case "long long":
					f.ReturnType = "int64"
				case "unsigned long long":
					f.ReturnType = "uint64"
				}

				added = addMethod(f, fname, "", rt)
			}
		}

		if !added {
			if len(f.Parameters) > 0 && (isEnumOrStruct(f.ReturnType) || f.ReturnPrimitiveType != "") {
				found := false
				for _, p := range f.Parameters {
					if !isEnumOrStruct(p.Type) && p.PrimitiveType == "" {
						found = true

						break
					}
				}

				if !found {
					fname = trimCommonFName(fname, rt)

					added = addMethod(f, fname, "", rt)
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

	if _, _, err = execToBuffer("gofmt", "-w", "./"); err != nil { // TODO do this before saving the files using go/fmt
		exitWithFatal("Gofmt failed", err)
	}
}

func goAndPrimitiveType(typ string) (string, string) {
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
	case "String":
		return "string", "String"
	}

	return "", ""
}
