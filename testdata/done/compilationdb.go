package clang

// #include <stdlib.h>
// #include "go-clang.h"
// #include "clang-c/CXCompilationDatabase.h"
//
import "C"

import (
	"fmt"
	"unsafe"
)

/**
 * \brief Creates a compilation database from the database found in directory
 * buildDir. For example, CMake can output a compile_commands.json which can
 * be used to build the database.
 *
 * It must be freed by \c clang_CompilationDatabase_dispose.
 */
func NewCompilationDatabase(builddir string) (CompilationDatabase, error) {
	var db CompilationDatabase

	c_dir := C.CString(builddir)
	defer C.free(unsafe.Pointer(c_dir))
	var c_err C.CXCompilationDatabase_Error
	c_db := C.clang_CompilationDatabase_fromDirectory(c_dir, &c_err)
	if c_err == C.CXCompilationDatabase_NoError {
		return CompilationDatabase{c_db}, nil
	}
	return db, CompilationDatabaseError(c_err)
}
