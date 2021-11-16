package gen_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/go-clang/gen"
	"github.com/go-clang/gen/cmd/go-clang-gen/runtime"
)

func TestAPIPrepareFunctionName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		PrepareFunctionName func(g *gen.Generation, f *gen.Function) string
		name                string
		cname               string
		want                string
	}{
		"clang_indexLoc_getFileLocation": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_indexLoc_getFileLocation",
			cname:               "clang_indexLoc_getFileLocation",
			want:                "getFileLocation",
		},
		"clang_index_isEntityObjCContainerKind": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_index_isEntityObjCContainerKind",
			cname:               "clang_index_isEntityObjCContainerKind",
			want:                "isEntityObjCContainerKind",
		},
		"clang_Location_isInSystemHeader": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_Location_isInSystemHeader",
			cname:               "clang_Location_isInSystemHeader",
			want:                "isInSystemHeader",
		},
		"clang_Range_isNull": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_Range_isNull",
			cname:               "clang_Range_isNull",
			want:                "isNull",
		},
		"clang_remap_getNumFiles": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_remap_getNumFiles",
			cname:               "clang_remap_getNumFiles",
			want:                "getNumFiles",
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api := &gen.API{
				PrepareFunctionName: tt.PrepareFunctionName,
			}
			g := gen.NewGeneration(api)
			f := &gen.Function{
				Name:       tt.name,
				CName:      tt.cname,
				Parameters: []gen.FunctionParameter{},
			}

			got := g.API().PrepareFunctionName(g, f)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("API.PrepareFunctionName(): (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAPIPrepareFunction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		PrepareFunction func(f *gen.Function)
		f               *gen.Function
		want            *gen.Function
	}{
		"clang_getRemappingsFromFileList": {
			PrepareFunction: runtime.PrepareFunction,
			f: &gen.Function{
				Name:  "getRemappingsFromFileList",
				CName: "clang_getRemappingsFromFileList",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "filePaths",
						CName: "filePaths",
						Type: gen.Type{
							CName: "const char **", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numFiles",
						CName: "numFiles",
						Type: gen.Type{
							CName: "unsigned int", CGoName: "uint", GoName: "uint32",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
			want: &gen.Function{
				Name:  "getRemappingsFromFileList",
				CName: "clang_getRemappingsFromFileList",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "filePaths",
						CName: "filePaths",
						Type: gen.Type{
							CName: "const char **", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              true,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numFiles",
						CName: "numFiles",
						Type: gen.Type{
							CName: "unsigned int", CGoName: "uint", GoName: "uint32",
							LengthOfSlice: "filePaths", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
		},
		"PointerLevelOneAndGoNameFile": {
			PrepareFunction: runtime.PrepareFunction,
			f: &gen.Function{
				Name:  "getExpansionLocation",
				CName: "clang_getExpansionLocation",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "file",
						CName: "file",
						Type: gen.Type{
							CName: "CXFile *", CGoName: "CXFile", GoName: "File",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
			want: &gen.Function{
				Name:  "getExpansionLocation",
				CName: "clang_getExpansionLocation",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "file",
						CName: "file",
						Type: gen.Type{
							CName: "CXFile *", CGoName: "CXFile", GoName: "File",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     true,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
		},
		"PointerLevelOneAndGoNameFileUniqueID": {
			PrepareFunction: runtime.PrepareFunction,
			f: &gen.Function{
				Name:  "getFileUniqueID",
				CName: "clang_getFileUniqueID",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "outID",
						CName: "outID",
						Type: gen.Type{
							CName: "CXFileUniqueID *", CGoName: "CXFileUniqueID", GoName: "FileUniqueID",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
			want: &gen.Function{
				Name:  "getFileUniqueID",
				CName: "clang_getFileUniqueID",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "outID",
						CName: "outID",
						Type: gen.Type{
							CName: "CXFileUniqueID *", CGoName: "CXFileUniqueID", GoName: "FileUniqueID",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     true,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
		},
		"PointerLevelTwoAndGoNameToken": {
			PrepareFunction: runtime.PrepareFunction,
			f: &gen.Function{
				Name:  "tokenize",
				CName: "clang_tokenize",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "tokens",
						CName: "Tokens",
						Type: gen.Type{
							CName: "CXToken **", CGoName: "CXToken", GoName: "Token",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
			want: &gen.Function{
				Name:  "tokenize",
				CName: "clang_tokenize",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "tokens",
						CName: "Tokens",
						Type: gen.Type{
							CName: "CXToken **", CGoName: "CXToken", GoName: "Token",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     true,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
		},
		"PointerLevelTwoAndGoNameCursor": {
			PrepareFunction: runtime.PrepareFunction,
			f: &gen.Function{
				Name:  "getOverriddenCursors",
				CName: "clang_getOverriddenCursors",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "overridden",
						CName: "overridden",
						Type: gen.Type{
							CName: "CXCursor **", CGoName: "CXCursor", GoName: "Cursor",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
			want: &gen.Function{
				Name:  "getOverriddenCursors",
				CName: "clang_getOverriddenCursors",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "overridden",
						CName: "overridden",
						Type: gen.Type{
							CName: "CXCursor **", CGoName: "CXCursor", GoName: "Cursor",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     true,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
		},
		"clang_disposeOverriddenCursors_overridden": {
			PrepareFunction: runtime.PrepareFunction,
			f: &gen.Function{
				Name:  "disposeOverriddenCursors",
				CName: "clang_disposeOverriddenCursors",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "overridden",
						CName: "overridden",
						Type: gen.Type{
							CName: "CXCursor *", CGoName: "CXCursor", GoName: "Cursor",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
			want: &gen.Function{
				Name:  "disposeOverriddenCursors",
				CName: "clang_disposeOverriddenCursors",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "overridden",
						CName: "overridden",
						Type: gen.Type{
							CName: "CXCursor *", CGoName: "CXCursor", GoName: "Cursor",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              true,
							IsPointerComposition: false,
						},
					},
				},
			},
		},
		"UnsavedFile": {
			PrepareFunction: runtime.PrepareFunction,
			f: &gen.Function{
				Name:  "createTranslationUnitFromSourceFile",
				CName: "clang_createTranslationUnitFromSourceFile",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "cIdx",
						CName: "CIdx",
						Type: gen.Type{
							CName: "CXIndex", CGoName: "CXIndex", GoName: "Index",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "sourceFilename",
						CName: "source_filename",
						Type: gen.Type{
							CName: "const char *", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numClangCommandLineArgs",
						CName: "num_clang_command_line_args",
						Type: gen.Type{
							CName: "int", CGoName: "int", GoName: "int32",
							LengthOfSlice: "clangCommandLineArgs", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "clangCommandLineArgs",
						CName: "clang_command_line_args",
						Type: gen.Type{
							CName: "const char *const *", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              true,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numUnsavedFiles",
						CName: "num_unsaved_files",
						Type: gen.Type{
							CName: "unsigned int", CGoName: "uint", GoName: "uint32",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "unsavedFiles",
						CName: "unsaved_files",
						Type: gen.Type{
							CName: "struct CXUnsavedFile *", CGoName: "struct CXUnsavedFile", GoName: "UnsavedFile",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
			want: &gen.Function{
				Name:  "createTranslationUnitFromSourceFile",
				CName: "clang_createTranslationUnitFromSourceFile",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "cIdx",
						CName: "CIdx",
						Type: gen.Type{
							CName: "CXIndex", CGoName: "CXIndex", GoName: "Index",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "sourceFilename",
						CName: "source_filename",
						Type: gen.Type{
							CName: "const char *", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numClangCommandLineArgs",
						CName: "num_clang_command_line_args",
						Type: gen.Type{
							CName: "int", CGoName: "int", GoName: "int32",
							LengthOfSlice: "clangCommandLineArgs", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "clangCommandLineArgs",
						CName: "clang_command_line_args",
						Type: gen.Type{
							CName: "const char *const *", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              true,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numUnsavedFiles",
						CName: "num_unsaved_files",
						Type: gen.Type{
							CName: "unsigned int", CGoName: "uint", GoName: "uint32",
							LengthOfSlice: "unsavedFiles", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "unsavedFiles",
						CName: "unsaved_files",
						Type: gen.Type{
							CName: "struct CXUnsavedFile *", CGoName: "struct_CXUnsavedFile", GoName: "UnsavedFile",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              true,
							IsPointerComposition: false,
						},
					},
				},
			},
		},
		"CGoNameCSCharAndPointerLevelTwo": {
			PrepareFunction: runtime.PrepareFunction,
			f: &gen.Function{
				Name:  "parseTranslationUnit",
				CName: "clang_parseTranslationUnit",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "cIdx",
						CName: "CIdx",
						Type: gen.Type{
							CName: "CXIndex", CGoName: "CXIndex", GoName: "Index",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "commandLineArgs",
						CName: "command_line_args",
						Type: gen.Type{
							CName: "const char *const *", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numCommandLineArgs",
						CName: "num_command_line_args",
						Type: gen.Type{
							CName: "int", CGoName: "int", GoName: "int32",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numUnsavedFiles",
						CName: "num_unsaved_files",
						Type: gen.Type{
							CName: "unsigned int", CGoName: "uint", GoName: "uint32",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "options",
						CName: "options",
						Type: gen.Type{
							CName: "unsigned int", CGoName: "uint", GoName: "uint32",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "sourceFilename",
						CName: "source_filename",
						Type: gen.Type{
							CName: "const char *", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "unsavedFiles",
						CName: "unsaved_files",
						Type: gen.Type{
							CName: "struct CXUnsavedFile *", CGoName: "struct CXUnsavedFile", GoName: "UnsavedFile",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
				},
			},
			want: &gen.Function{
				Name:  "parseTranslationUnit",
				CName: "clang_parseTranslationUnit",
				Parameters: []gen.FunctionParameter{
					{
						Name:  "cIdx",
						CName: "CIdx",
						Type: gen.Type{
							CName: "CXIndex", CGoName: "CXIndex", GoName: "Index",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "commandLineArgs",
						CName: "command_line_args",
						Type: gen.Type{
							CName: "const char *const *", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 2,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              true,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numCommandLineArgs",
						CName: "num_command_line_args",
						Type: gen.Type{
							CName: "int", CGoName: "int", GoName: "int32",
							LengthOfSlice: "commandLineArgs", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "numUnsavedFiles",
						CName: "num_unsaved_files",
						Type: gen.Type{
							CName: "unsigned int", CGoName: "uint", GoName: "uint32",
							LengthOfSlice: "unsavedFiles", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "options",
						CName: "options",
						Type: gen.Type{
							CName: "unsigned int", CGoName: "uint", GoName: "uint32",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 0,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "sourceFilename",
						CName: "source_filename",
						Type: gen.Type{
							CName: "const char *", CGoName: "schar", GoName: "int8",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          true,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              false,
							IsPointerComposition: false,
						},
					},
					{
						Name:  "unsavedFiles",
						CName: "unsaved_files",
						Type: gen.Type{
							CName: "struct CXUnsavedFile *", CGoName: "struct_CXUnsavedFile", GoName: "UnsavedFile",
							LengthOfSlice: "", ArraySize: -1, PointerLevel: 1,
							IsPrimitive:          false,
							IsArray:              false,
							IsEnumLiteral:        false,
							IsFunctionPointer:    false,
							IsReturnArgument:     false,
							IsSlice:              true,
							IsPointerComposition: false,
						},
					},
				},
			},
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api := &gen.API{
				PrepareFunction: tt.PrepareFunction,
			}
			g := gen.NewGeneration(api)
			f := tt.f
			g.API().PrepareFunction(f)

			if diff := cmp.Diff(tt.want, f,
				cmpopts.SortSlices(
					func(x, y gen.FunctionParameter) bool {
						return x.Name < y.Name
					},
				),
			); diff != "" {
				t.Fatalf("API.PrepareFunction(): (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAPIFilterFunction(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		FilterFunction func(f *gen.Function) bool
		function       *gen.Function
		want           bool
	}{
		"clang_CompileCommand_getMappedSourceContent": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_CompileCommand_getMappedSourceContent",
			},
			want: false,
		},
		"clang_CompileCommand_getMappedSourcePath": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_CompileCommand_getMappedSourcePath",
			},
			want: false,
		},
		"clang_CompileCommand_getNumMappedSources": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_CompileCommand_getNumMappedSources",
			},
			want: false,
		},
		"clang_executeOnThread": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_executeOnThread",
			},
			want: false,
		},
		"clang_getInclusions": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_getInclusions",
			},
			want: false,
		},
		"clang_annotateTokens": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_annotateTokens",
			},
			want: false,
		},
		"clang_getCursorPlatformAvailability": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_getCursorPlatformAvailability",
			},
			want: false,
		},
		"clang_visitChildren": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_visitChildren",
			},
			want: false,
		},
		"clang_getCString": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_getCString",
				IncludeFiles: gen.IncludeFiles{
					"./clang/clang-c/CXString.h": struct{}{},
				},
			},
			want: false,
		},
		"clang_getNumDiagnosticsInSet": {
			FilterFunction: runtime.FilterFunction,
			function: &gen.Function{
				CName: "clang_getNumDiagnosticsInSet",
				IncludeFiles: gen.IncludeFiles{
					"./clang/clang-c/Index.h": struct{}{},
				},
			},
			want: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api := &gen.API{
				FilterFunction: tt.FilterFunction,
			}
			g := gen.NewGeneration(api)

			got := g.API().FilterFunction(tt.function)
			if tt.want != got {
				t.Fatalf("API.FilterFunction(): want: %t but got %t", tt.want, got)
			}
		})
	}
}

func TestAPIFilterFunctionParameter(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		FilterFunctionParameter func(p gen.FunctionParameter) bool
		parameter               gen.FunctionParameter
		want                    bool
	}{
		"CSChar": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CGoName:      "schar",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"UnsavedFile": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "struct CXUnsavedFile *",
					CGoName:      "struct_CXUnsavedFile",
					GoName:       "UnsavedFile",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"CodeCompleteResults": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "CXCodeCompleteResults *",
					CGoName:      "CXCodeCompleteResults",
					GoName:       "CodeCompleteResults",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"CursorKind": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "enum CXCursorKind *",
					CGoName:      "enum_CXCursorKind",
					GoName:       "CursorKind",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"IdxContainerInfo": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "const CXIdxContainerInfo *",
					CGoName:      "CXIdxContainerInfo",
					GoName:       "IdxContainerInfo",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"IdxDeclInfo": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "const CXIdxDeclInfo *",
					CGoName:      "CXIdxDeclInfo",
					GoName:       "IdxDeclInfo",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"IndexerCallbacks": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "IndexerCallbacks *",
					CGoName:      "IndexerCallbacks",
					GoName:       "IndexerCallbacks",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"TranslationUnit": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "CXTranslationUnit *",
					CGoName:      "CXTranslationUnit",
					GoName:       "TranslationUnit",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"IdxEntityInfo": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "const CXIdxEntityInfo *",
					CGoName:      "CXIdxEntityInfo",
					GoName:       "IdxEntityInfo",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"IdxAttrInfo": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "const CXIdxAttrInfo *",
					CGoName:      "CXIdxAttrInfo",
					GoName:       "IdxAttrInfo",
					PointerLevel: 1,
				},
			},
			want: false,
		},
		"uint32": {
			FilterFunctionParameter: runtime.FilterFunctionParameter,
			parameter: gen.FunctionParameter{
				Type: gen.Type{
					CName:        "unsigned int",
					CGoName:      "uint",
					GoName:       "uint32",
					PointerLevel: 0,
				},
			},
			want: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api := &gen.API{
				FilterFunctionParameter: tt.FilterFunctionParameter,
			}
			g := gen.NewGeneration(api)

			got := g.API().FilterFunctionParameter(tt.parameter)
			if tt.want != got {
				t.Fatalf("API.FilterFunction(): want: %t but got %t", tt.want, got)
			}
		})
	}
}

func TestAPIFixFunctionName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		FixFunctionName func(f *gen.Function) string
		function        *gen.Function
		want            string
	}{
		"clang_getTranslationUnitCursor": {
			FixFunctionName: runtime.FixFunctionName,
			function: &gen.Function{
				CName: "clang_getTranslationUnitCursor",
			},
			want: "TranslationUnitCursor",
		},
		"struct CXUnsavedFile *": {
			FixFunctionName: runtime.FixFunctionName,
			function: &gen.Function{
				CName: "struct CXUnsavedFile *",
			},
			want: "",
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api := &gen.API{
				FixFunctionName: tt.FixFunctionName,
			}
			g := gen.NewGeneration(api)

			got := g.API().FixFunctionName(tt.function)
			if tt.want != got {
				t.Fatalf("API.FilterFunction(): want: %s but got %s", tt.want, got)
			}
		})
	}
}
