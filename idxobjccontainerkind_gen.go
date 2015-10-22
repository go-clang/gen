package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCContainerKind uint32

const (
	IdxObjCContainer_ForwardRef     IdxObjCContainerKind = C.CXIdxObjCContainer_ForwardRef
	IdxObjCContainer_Interface                           = C.CXIdxObjCContainer_Interface
	IdxObjCContainer_Implementation                      = C.CXIdxObjCContainer_Implementation
)
