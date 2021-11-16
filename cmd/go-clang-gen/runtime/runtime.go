package runtime

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-clang/gen"
)

// PrepareFunctionName prepares C function naming to Go function name.
func PrepareFunctionName(g *gen.Generation, f *gen.Function) string {
	fname := strings.TrimPrefix(f.Name, "clang_")

	// trim some allowlisted prefixes by their function name
	switch {
	case strings.HasPrefix(fname, "indexLoc_"):
		fname = strings.TrimPrefix(fname, "indexLoc_")

	case strings.HasPrefix(fname, "index_"):
		fname = strings.TrimPrefix(fname, "index_")

	case strings.HasPrefix(fname, "Location_"):
		fname = strings.TrimPrefix(fname, "Location_")

	case strings.HasPrefix(fname, "Range_"):
		fname = strings.TrimPrefix(fname, "Range_")

	case strings.HasPrefix(fname, "remap_"):
		fname = strings.TrimPrefix(fname, "remap_")
	}

	// trim some allowlisted prefixes by their types
	if len(f.Parameters) > 0 && g.IsEnumOrStruct(f.Parameters[0].Type.GoName) {
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

	return fname
}

// PrepareFunction prepares C function to Go function.
func PrepareFunction(f *gen.Function) {
	for i := range f.Parameters {
		p := &f.Parameters[i]

		if f.CName == "clang_getRemappingsFromFileList" {
			switch p.CName {
			case "filePaths":
				p.Type.IsSlice = true

			case "numFiles":
				p.Type.LengthOfSlice = "filePaths"
			}

			continue
		}

		// allowflag types that are return arguments
		switch p.Type.PointerLevel {
		case 1:
			switch p.Type.GoName {
			case
				gen.GoInt32,
				gen.GoUInt32,
				"File",
				"FileUniqueID",
				"IdxClientFile",
				"cxstring",
				"CompilationDatabase_Error",
				"PlatformAvailability",
				"SourceRange",
				"LoadDiag_Error":

				p.Type.IsReturnArgument = true
			}

		case 2:
			switch p.Type.GoName {
			case
				"Token",
				"Cursor":

				p.Type.IsReturnArgument = true
			}
		}

		if f.CName == "clang_disposeOverriddenCursors" && p.CName == "overridden" {
			p.Type.IsSlice = true
		}

		// if this is an array length parameter we need to find its partner
		paCName := gen.ArrayNameFromLength(p.CName)

		if paCName != "" {
			for j := range f.Parameters {
				pa := &f.Parameters[j]

				if strings.EqualFold(pa.CName, paCName) {
					switch pa.Type.GoName {
					case "UnsavedFile":
						pa.Type.GoName = "UnsavedFile"
						pa.Type.CGoName = "struct_CXUnsavedFile"

					case "CompletionResult", "Token", "Cursor":
						// nothing to do

					default:
						if pa.Type.CGoName != gen.CSChar && pa.Type.PointerLevel != 2 {
							break
						}
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

		if p.Type.CGoName == gen.CSChar && p.Type.PointerLevel == 2 && !p.Type.IsSlice {
			p.Type.IsReturnArgument = true
		}
	}
}

// FilterFunction reports whether the f function filtered to a particular condition.
func FilterFunction(f *gen.Function) bool {
	switch f.CName {
	case "clang_CompileCommand_getMappedSourceContent", "clang_CompileCommand_getMappedSourcePath", "clang_CompileCommand_getNumMappedSources":
		// some functions are not compiled in the library see https://lists.launchpad.net/desktop-packages/msg75835.html for a never resolved bug report
		fmt.Fprintf(os.Stderr, "Ignore function %q because it is not compiled within libClang\n", f.CName)

		return false

	case "clang_executeOnThread", "clang_getInclusions":
		// some functions can not be handled automatically by us
		fmt.Fprintf(os.Stderr, "Ignore function %q because it cannot be handled automatically\n", f.CName)

		return false

	case "clang_annotateTokens", "clang_getCursorPlatformAvailability", "clang_visitChildren":
		// some functions are simply manually implemented
		fmt.Fprintf(os.Stderr, "Ignore function %q because it is manually implemented\n", f.CName)

		return false
	}

	// TODO(go-clang): if this function is from CXString.h we ignore it https://github.com/go-clang/gen/issues/25
	for i := range f.IncludeFiles {
		if strings.HasSuffix(i, "CXString.h") {
			return false
		}
	}

	return true
}

// FilterFunctionParameter reports whether the p function parameter filtered to a particular condition.
func FilterFunctionParameter(p gen.FunctionParameter) bool {
	// these pointers are ok
	if p.Type.PointerLevel == 1 {
		if p.Type.CGoName == gen.CSChar {
			return false
		}

		switch p.Type.GoName {
		case "UnsavedFile",
			"CodeCompleteResults",
			"CursorKind",
			"IdxContainerInfo",
			"IdxDeclInfo",
			"IndexerCallbacks",
			"TranslationUnit",
			"IdxEntityInfo",
			"IdxAttrInfo":

			return false
		}
	}

	return true
}

// FixFunctionName fixes the function name under certain conditions.
func FixFunctionName(f *gen.Function) string {
	// needs to be renamed manually since clang_getTranslationUnitCursor will conflict with clang_getCursor
	if f.CName == "clang_getTranslationUnitCursor" {
		return "TranslationUnitCursor"
	}

	return ""
}

// PrepareStructFields prepares struct fields names.
func PrepareStructFields(s *gen.Struct) {
	for _, f := range s.Fields {
		if (strings.HasPrefix(f.CName, "has") || strings.HasPrefix(f.CName, "is")) && f.Type.GoName == gen.GoInt32 {
			f.Type.GoName = gen.GoBool
		}

		// if this is an array length parameter we need to find its partner
		faCName := gen.ArrayNameFromLength(f.CName)

		if faCName != "" {
			for _, fa := range s.Fields {
				if strings.EqualFold(fa.CName, faCName) {
					f.Type.LengthOfSlice = fa.CName
					fa.Type.IsSlice = true
					// TODO(go-clang): wrong usage but needed for the getter generation...
					// maybe refactor this LengthOfSlice all together?
					// https://github.com/go-clang/gen/issues/49
					fa.Type.LengthOfSlice = f.CName

					break
				}
			}
		}
	}

	prepareStructFieldsArrayStruct(s)
}

// prepareStructFieldsArrayStruct checks if the struct has two field variables, one is an array and the other a plain
// int/uint with size/length/count/len is its name because then this should be an array struct, and we connect them to handle a slice.
func prepareStructFieldsArrayStruct(s *gen.Struct) {
	if len(s.Fields) != 2 {
		return
	}

	if !arrayLengthCombination(&s.Fields[0].Type, &s.Fields[1].Type) && !arrayLengthCombination(&s.Fields[1].Type, &s.Fields[0].Type) {
		return
	}

	// if one of the fields is already marked as array/slice another heuristic has already covered both fields
	switch {
	case s.Fields[0].Type.IsArray,
		s.Fields[1].Type.IsArray,
		s.Fields[0].Type.IsSlice,
		s.Fields[1].Type.IsSlice:

		return
	}

	var a *gen.StructField
	var c *gen.StructField

	if s.Fields[0].Type.PointerLevel == 1 {
		a = s.Fields[0]
		c = s.Fields[1]
	} else {
		c = s.Fields[0]
		a = s.Fields[1]
	}

	lengthName := strings.ToLower(c.CName)
	if lengthName != "count" &&
		lengthName != "len" && lengthName != "length" && lengthName != "size" {

		return
	}

	c.Type.LengthOfSlice = a.CName
	a.Type.IsSlice = true
	// TODO(go-clang): wrong usage but needed for the getter generation...
	// maybe refactor this LengthOfSlice all together?
	// https://github.com/go-clang/gen/issues/49
	a.Type.LengthOfSlice = c.CName
}

// arrayLengthCombination reports whether the x and y to correct combination.
func arrayLengthCombination(x *gen.Type, y *gen.Type) bool {
	return x.PointerLevel == 1 && y.PointerLevel == 0 &&
		!gen.IsInteger(x) && gen.IsInteger(y)
}

// FilterStructFieldGetter reports whether the m struct field filtered to a particular condition.
func FilterStructFieldGetter(f *gen.StructField) bool {
	// we do not want getters to *int_data fields
	return !strings.HasSuffix(f.CName, "int_data")
}
