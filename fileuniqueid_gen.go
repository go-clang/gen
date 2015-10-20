package phoenix

// #include "go-clang.h"
import "C"

// Uniquely identifies a CXFile, that refers to the same underlying file, across an indexing session.
type FileUniqueID struct {
	c C.CXFileUniqueID
}
