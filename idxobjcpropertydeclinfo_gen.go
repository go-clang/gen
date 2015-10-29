package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCPropertyDeclInfo struct {
	c C.CXIdxObjCPropertyDeclInfo
}

func (iocpdi IdxObjCPropertyDeclInfo) DeclInfo() *IdxDeclInfo {
	value := IdxDeclInfo{*iocpdi.c.declInfo}
	return &value
}

func (iocpdi IdxObjCPropertyDeclInfo) Getter() *IdxEntityInfo {
	value := IdxEntityInfo{*iocpdi.c.getter}
	return &value
}

func (iocpdi IdxObjCPropertyDeclInfo) Setter() *IdxEntityInfo {
	value := IdxEntityInfo{*iocpdi.c.setter}
	return &value
}
