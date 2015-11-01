package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// An indexing action/session, to be applied to one or multiple translation units.
type IndexAction struct {
	c C.CXIndexAction
}

// Destroy the given index action. The index action must not be destroyed until all of the translation units created within that index action have been destroyed.
func (ia IndexAction) Dispose() {
	C.clang_IndexAction_dispose(ia.c)
}

// Index the given source file and the translation unit corresponding to that file via callbacks implemented through #IndexerCallbacks. \param client_data pointer data supplied by the client, which will be passed to the invoked callbacks. \param index_callbacks Pointer to indexing callbacks that the client implements. \param index_callbacks_size Size of #IndexerCallbacks structure that gets passed in index_callbacks. \param index_options A bitmask of options that affects how indexing is performed. This should be a bitwise OR of the CXIndexOpt_XXX flags. \param out_TU [out] pointer to store a CXTranslationUnit that can be reused after indexing is finished. Set to NULL if you do not require it. \returns If there is a failure from which the there is no recovery, returns non-zero, otherwise returns 0. The rest of the parameters are the same as #clang_parseTranslationUnit.
func (ia IndexAction) IndexSourceFile(client_data ClientData, index_callbacks *IndexerCallbacks, index_callbacks_size uint16, index_options uint16, source_filename string, command_line_args []string, unsaved_files []UnsavedFile, out_TU *TranslationUnit, TU_options uint16) int16 {
	ca_command_line_args := make([]*C.char, len(command_line_args))
	var cp_command_line_args **C.char
	if len(command_line_args) > 0 {
		cp_command_line_args = &ca_command_line_args[0]
	}
	for i := range command_line_args {
		ci_str := C.CString(command_line_args[i])
		defer C.free(unsafe.Pointer(ci_str))
		ca_command_line_args[i] = ci_str
	}
	ca_unsaved_files := make([]C.struct_CXUnsavedFile, len(unsaved_files))
	var cp_unsaved_files *C.struct_CXUnsavedFile
	if len(unsaved_files) > 0 {
		cp_unsaved_files = &ca_unsaved_files[0]
	}
	for i := range unsaved_files {
		ca_unsaved_files[i] = unsaved_files[i].c
	}

	c_source_filename := C.CString(source_filename)
	defer C.free(unsafe.Pointer(c_source_filename))

	return int16(C.clang_indexSourceFile(ia.c, client_data.c, &index_callbacks.c, C.uint(index_callbacks_size), C.uint(index_options), c_source_filename, cp_command_line_args, C.int(len(command_line_args)), cp_unsaved_files, C.uint(len(unsaved_files)), &out_TU.c, C.uint(TU_options)))
}

// Index the given translation unit via callbacks implemented through #IndexerCallbacks. The order of callback invocations is not guaranteed to be the same as when indexing a source file. The high level order will be: -Preprocessor callbacks invocations -Declaration/reference callbacks invocations -Diagnostic callback invocations The parameters are the same as #clang_indexSourceFile. \returns If there is a failure from which the there is no recovery, returns non-zero, otherwise returns 0.
func (ia IndexAction) IndexTranslationUnit(client_data ClientData, index_callbacks *IndexerCallbacks, index_callbacks_size uint16, index_options uint16, tu TranslationUnit) int16 {
	return int16(C.clang_indexTranslationUnit(ia.c, client_data.c, &index_callbacks.c, C.uint(index_callbacks_size), C.uint(index_options), tu.c))
}
