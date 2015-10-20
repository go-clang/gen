package phoenix

// #include "go-clang.h"
import "C"

// Data for IndexerCallbacks#importedASTFile.
type IdxImportedASTFileInfo struct {
	c C.CXIdxImportedASTFileInfo
}
