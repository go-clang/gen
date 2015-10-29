package phoenix

// #include "go-clang.h"
import "C"

// Data for IndexerCallbacks#importedASTFile.
type IdxImportedASTFileInfo struct {
	c C.CXIdxImportedASTFileInfo
}

// Top level AST file containing the imported PCH, module or submodule.
func (iiastfi IdxImportedASTFileInfo) File() File {
	value := File{iiastfi.c.file}
	return value
}

// The imported module or NULL if the AST file is a PCH.
func (iiastfi IdxImportedASTFileInfo) Module() Module {
	value := Module{iiastfi.c.module}
	return value
}

// Location where the file is imported. Applicable only for modules.
func (iiastfi IdxImportedASTFileInfo) Loc() IdxLoc {
	value := IdxLoc{iiastfi.c.loc}
	return value
}

// Non-zero if an inclusion directive was automatically turned into a module import. Applicable only for modules.
func (iiastfi IdxImportedASTFileInfo) IsImplicit() bool {
	value := iiastfi.c.isImplicit
	return value != C.int(0)
}
