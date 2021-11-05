package clang

// #include "go-clang.h"
import "C"

func NewUnsavedFile(filename, contents string) UnsavedFile {
	return UnsavedFile{
		C.struct_CXUnsavedFile{
			Filename: C.CString(filename),
			Contents: C.CString(contents),
			Length:   C.ulong(len(contents)),
		},
	}
}
