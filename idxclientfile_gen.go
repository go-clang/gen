package phoenix

// #include "go-clang.h"
import "C"

// The client's data object that is associated with a CXFile.
type IdxClientFile struct {
	c C.CXIdxClientFile
}
