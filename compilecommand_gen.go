package phoenix

// #include "./clang-c/CXCompilationDatabase.h"
// #include "go-clang.h"
import "C"

// Represents the command line invocation to compile a specific file.
type CompileCommand struct {
	c C.CXCompileCommand
}

// Get the working directory where the CompileCommand was executed from
func (cc CompileCommand) Directory() string {
	o := cxstring{C.clang_CompileCommand_getDirectory(cc.c)}
	defer o.Dispose()

	return o.String()
}

// Get the number of arguments in the compiler invocation.
func (cc CompileCommand) NumArgs() uint16 {
	return uint16(C.clang_CompileCommand_getNumArgs(cc.c))
}

// Get the I'th argument value in the compiler invocations Invariant : - argument 0 is the compiler executable
func (cc CompileCommand) Arg(I uint16) string {
	o := cxstring{C.clang_CompileCommand_getArg(cc.c, C.uint(I))}
	defer o.Dispose()

	return o.String()
}
