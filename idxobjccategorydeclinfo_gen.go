package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCCategoryDeclInfo struct {
	c C.CXIdxObjCCategoryDeclInfo
}

func (ioccdi IdxObjCCategoryDeclInfo) ContainerInfo() *IdxObjCContainerDeclInfo {
	value := IdxObjCContainerDeclInfo{*ioccdi.c.containerInfo}
	return &value
}

func (ioccdi IdxObjCCategoryDeclInfo) ObjcClass() *IdxEntityInfo {
	value := IdxEntityInfo{*ioccdi.c.objcClass}
	return &value
}

func (ioccdi IdxObjCCategoryDeclInfo) ClassCursor() Cursor {
	value := Cursor{ioccdi.c.classCursor}
	return value
}

func (ioccdi IdxObjCCategoryDeclInfo) ClassLoc() IdxLoc {
	value := IdxLoc{ioccdi.c.classLoc}
	return value
}

func (ioccdi IdxObjCCategoryDeclInfo) Protocols() *IdxObjCProtocolRefListInfo {
	value := IdxObjCProtocolRefListInfo{*ioccdi.c.protocols}
	return &value
}
