package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCCategoryDeclInfo struct {
	c C.CXIdxObjCCategoryDeclInfo
}

func (ioccdi IdxObjCCategoryDeclInfo) ContainerInfo() *IdxObjCContainerDeclInfo {
	o := *ioccdi.c.containerInfo

	return &IdxObjCContainerDeclInfo{o}
}

func (ioccdi IdxObjCCategoryDeclInfo) ObjcClass() *IdxEntityInfo {
	o := *ioccdi.c.objcClass

	return &IdxEntityInfo{o}
}

func (ioccdi IdxObjCCategoryDeclInfo) ClassCursor() Cursor {
	return Cursor{ioccdi.c.classCursor}
}

func (ioccdi IdxObjCCategoryDeclInfo) ClassLoc() IdxLoc {
	return IdxLoc{ioccdi.c.classLoc}
}

func (ioccdi IdxObjCCategoryDeclInfo) Protocols() *IdxObjCProtocolRefListInfo {
	o := *ioccdi.c.protocols

	return &IdxObjCProtocolRefListInfo{o}
}
