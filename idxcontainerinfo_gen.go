package phoenix

// #include "go-clang.h"
import "C"

type IdxContainerInfo struct {
	c C.CXIdxContainerInfo
}

func (ici IdxContainerInfo) Cursor() Cursor {
	value := Cursor{ici.c.cursor}
	return value
}
