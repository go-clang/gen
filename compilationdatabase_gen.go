package phoenix

// #include "./clang-c/CXCompilationDatabase.h"
// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// A compilation database holds all information used to compile files in a project. For each file in the database, it can be queried for the working directory or the command line used for the compiler invocation. Must be freed by \c clang_CompilationDatabase_dispose
type CompilationDatabase struct {
	c C.CXCompilationDatabase
}

// Creates a compilation database from the database found in directory buildDir. For example, CMake can output a compile_commands.json which can be used to build the database. It must be freed by \c clang_CompilationDatabase_dispose.
func CompilationDatabase_fromDirectory(BuildDir string) (CompilationDatabase_Error, CompilationDatabase) {
	var ErrorCode C.CXCompilationDatabase_Error

	c_BuildDir := C.CString(BuildDir)
	defer C.free(unsafe.Pointer(c_BuildDir))

	o := CompilationDatabase{C.clang_CompilationDatabase_fromDirectory(c_BuildDir, &ErrorCode)}

	return CompilationDatabase_Error(ErrorCode), o
}

// Free the given compilation database
func (cd CompilationDatabase) Dispose() {
	C.clang_CompilationDatabase_dispose(cd.c)
}

// Find the compile commands used for a file. The compile commands must be freed by \c clang_CompileCommands_dispose.
func (cd CompilationDatabase) CompileCommands(CompleteFileName string) CompileCommands {
	c_CompleteFileName := C.CString(CompleteFileName)
	defer C.free(unsafe.Pointer(c_CompleteFileName))

	return CompileCommands{C.clang_CompilationDatabase_getCompileCommands(cd.c, c_CompleteFileName)}
}

// Get all the compile commands in the given compilation database.
func (cd CompilationDatabase) AllCompileCommands() CompileCommands {
	return CompileCommands{C.clang_CompilationDatabase_getAllCompileCommands(cd.c)}
}
