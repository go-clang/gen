package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

/**
 * \param Module a module object.
 *
 * \returns the module file where the provided module object came from.
 */
func (m Module) ASTFile() File {
	return File{C.clang_Module_getASTFile(m.c)}
}

/**
 * \param Module a module object.
 *
 * \returns the parent of a sub-module or NULL if the given module is top-level,
 * e.g. for 'std.vector' it will return the 'std' module.
 */
func (m Module) Parent() Module {
	return Module{C.clang_Module_getParent(m.c)}
}
