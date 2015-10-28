package phoenix

// #include "go-clang.h"
import "C"

import (
	"time"
)

// A particular source file that is part of a translation unit.
type File struct {
	c C.CXFile
}

// Retrieve the complete file and path name of the given file.
func (f File) Name() string {
	o := cxstring{C.clang_getFileName(f.c)}
	defer o.Dispose()

	return o.String()
}

// Retrieve the last modification time of the given file.
func (f File) Time() time.Time {
	return time.Unix(int64(C.clang_getFileTime(f.c)), 0)
}
