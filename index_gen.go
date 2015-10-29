package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// An "index" that consists of a set of translation units that would typically be linked together into an executable or library.
type Index struct {
	c C.CXIndex
}

/*
 * \brief Provides a shared context for creating translation units.
 *
 * It provides two options:
 *
 * - excludeDeclarationsFromPCH: When non-zero, allows enumeration of "local"
 * declarations (when loading any new translation units). A "local" declaration
 * is one that belongs in the translation unit itself and not in a precompiled
 * header that was used by the translation unit. If zero, all declarations
 * will be enumerated.
 *
 * Here is an example:
 *
 * \code
 *   // excludeDeclsFromPCH = 1, displayDiagnostics=1
 *   Idx = clang_createIndex(1, 1);
 *
 *   // IndexTest.pch was produced with the following command:
 *   // "clang -x c IndexTest.h -emit-ast -o IndexTest.pch"
 *   TU = clang_createTranslationUnit(Idx, "IndexTest.pch");
 *
 *   // This will load all the symbols from 'IndexTest.pch'
 *   clang_visitChildren(clang_getTranslationUnitCursor(TU),
 *                       TranslationUnitVisitor, 0);
 *   clang_disposeTranslationUnit(TU);
 *
 *   // This will load all the symbols from 'IndexTest.c', excluding symbols
 *   // from 'IndexTest.pch'.
 *   char *args[] = { "-Xclang", "-include-pch=IndexTest.pch" };
 *   TU = clang_createTranslationUnitFromSourceFile(Idx, "IndexTest.c", 2, args,
 *                                                  0, 0);
 *   clang_visitChildren(clang_getTranslationUnitCursor(TU),
 *                       TranslationUnitVisitor, 0);
 *   clang_disposeTranslationUnit(TU);
 * \endcode
 *
 * This process of creating the 'pch', loading it separately, and using it (via
 * -include-pch) allows 'excludeDeclsFromPCH' to remove redundant callbacks
 * (which gives the indexer the same performance benefit as the compiler).
 */
func NewIndex(excludeDeclarationsFromPCH uint16, displayDiagnostics uint16) Index {
	return Index{C.clang_createIndex(C.int(excludeDeclarationsFromPCH), C.int(displayDiagnostics))}
}

// Destroy the given index. The index must not be destroyed until all of the translation units created within that index have been destroyed.
func (i Index) Dispose() {
	C.clang_disposeIndex(i.c)
}

/*
 * \brief Sets general options associated with a CXIndex.
 *
 * For example:
 * \code
 * CXIndex idx = ...;
 * clang_CXIndex_setGlobalOptions(idx,
 *     clang_CXIndex_getGlobalOptions(idx) |
 *     CXGlobalOpt_ThreadBackgroundPriorityForIndexing);
 * \endcode
 *
 * \param options A bitmask of options, a bitwise OR of CXGlobalOpt_XXX flags.
 */
func (i Index) SetGlobalOptions(options uint16) {
	C.clang_CXIndex_setGlobalOptions(i.c, C.uint(options))
}

// Gets the general options associated with a CXIndex. \returns A bitmask of options, a bitwise OR of CXGlobalOpt_XXX flags that are associated with the given CXIndex object.
func (i Index) GlobalOptions() uint16 {
	return uint16(C.clang_CXIndex_getGlobalOptions(i.c))
}

// Create a translation unit from an AST file (-emit-ast).
func (i Index) TranslationUnit(ast_filename string) TranslationUnit {
	c_ast_filename := C.CString(ast_filename)
	defer C.free(unsafe.Pointer(c_ast_filename))

	return TranslationUnit{C.clang_createTranslationUnit(i.c, c_ast_filename)}
}

// An indexing action/session, to be applied to one or multiple translation units. \param CIdx The index object with which the index action will be associated.
func (i Index) Action_create() IndexAction {
	return IndexAction{C.clang_IndexAction_create(i.c)}
}
