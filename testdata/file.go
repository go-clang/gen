package clang

// #include <stdlib.h>
// #include "go-clang.h"
import "C"

import (
	"time"
)

// ModTime retrieves the last modification time of the given file.
func (c File) ModTime() time.Time {
	// time_t is in seconds since epoch
	sec := C.clang_getFileTime(c.c)
	const nsec = 0
	return time.Unix(int64(sec), nsec)
}
