package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCContainerDeclInfo struct {
	c C.CXIdxObjCContainerDeclInfo
}

func (ioccdi IdxObjCContainerDeclInfo) DeclInfo() *IdxDeclInfo {
	value := IdxDeclInfo{*ioccdi.c.declInfo}
	return &value
}

func (ioccdi IdxObjCContainerDeclInfo) Kind() IdxObjCContainerKind {
	value := IdxObjCContainerKind(ioccdi.c.kind)
	return value
}
