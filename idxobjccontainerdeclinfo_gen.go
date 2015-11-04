package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCContainerDeclInfo struct {
	c C.CXIdxObjCContainerDeclInfo
}

func (ioccdi IdxObjCContainerDeclInfo) DeclInfo() *IdxDeclInfo {
	o := *ioccdi.c.declInfo

	return &IdxDeclInfo{o}
}

func (ioccdi IdxObjCContainerDeclInfo) Kind() IdxObjCContainerKind {
	return IdxObjCContainerKind(ioccdi.c.kind)
}
