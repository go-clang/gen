package phoenix

// #include "go-clang.h"
import "C"

// A particular source file that is part of a translation unit.
type File struct {
	c C.CXFile
}

// Retrieve the complete file and path name of the given file.
func (f File) Name() string {
	cstr := cxstring{C.clang_getFileName(f.c)}
	defer cstr.Dispose()

	return cstr.String()
}
