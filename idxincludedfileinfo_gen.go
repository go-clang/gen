package phoenix

// #include "go-clang.h"
import "C"

// Data for ppIncludedFile callback.
type IdxIncludedFileInfo struct {
	c C.CXIdxIncludedFileInfo
}
