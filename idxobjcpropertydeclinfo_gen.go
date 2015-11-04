package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCPropertyDeclInfo struct {
	c C.CXIdxObjCPropertyDeclInfo
}

func (iocpdi IdxObjCPropertyDeclInfo) DeclInfo() *IdxDeclInfo {
	o := *iocpdi.c.declInfo

	return &IdxDeclInfo{o}
}

func (iocpdi IdxObjCPropertyDeclInfo) Getter() *IdxEntityInfo {
	o := *iocpdi.c.getter

	return &IdxEntityInfo{o}
}

func (iocpdi IdxObjCPropertyDeclInfo) Setter() *IdxEntityInfo {
	o := *iocpdi.c.setter

	return &IdxEntityInfo{o}
}
