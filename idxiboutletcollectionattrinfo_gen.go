package phoenix

// #include "go-clang.h"
import "C"

type IdxIBOutletCollectionAttrInfo struct {
	c C.CXIdxIBOutletCollectionAttrInfo
}

func (iibocai IdxIBOutletCollectionAttrInfo) AttrInfo() *IdxAttrInfo {
	o := *iibocai.c.attrInfo

	return &IdxAttrInfo{o}
}

func (iibocai IdxIBOutletCollectionAttrInfo) ObjcClass() *IdxEntityInfo {
	o := *iibocai.c.objcClass

	return &IdxEntityInfo{o}
}

func (iibocai IdxIBOutletCollectionAttrInfo) ClassCursor() Cursor {
	return Cursor{iibocai.c.classCursor}
}

func (iibocai IdxIBOutletCollectionAttrInfo) ClassLoc() IdxLoc {
	return IdxLoc{iibocai.c.classLoc}
}
