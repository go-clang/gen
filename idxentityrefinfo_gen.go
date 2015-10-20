package phoenix

// #include "go-clang.h"
import "C"

// Data for IndexerCallbacks#indexEntityReference.
type IdxEntityRefInfo struct {
	c C.CXIdxEntityRefInfo
}
