package phoenix

// #include "go-clang.h"
import "C"

type IdxContainerInfo struct {
	c C.CXIdxContainerInfo
}

// For retrieving a custom CXIdxClientContainer attached to a container.
func (ici *IdxContainerInfo) Index_getClientContainer() IdxClientContainer {
	return IdxClientContainer{C.clang_index_getClientContainer(&ici.c)}
}

// For setting a custom CXIdxClientContainer attached to a container.
func (ici *IdxContainerInfo) Index_setClientContainer(icc IdxClientContainer) {
	C.clang_index_setClientContainer(&ici.c, icc.c)
}

func (ici IdxContainerInfo) Cursor() Cursor {
	return Cursor{ici.c.cursor}
}
