package phoenix

// #include "go-clang.h"
import "C"

type IdxIBOutletCollectionAttrInfo struct {
	c C.CXIdxIBOutletCollectionAttrInfo
}

func (iibocai IdxIBOutletCollectionAttrInfo) AttrInfo() *IdxAttrInfo {
	value := IdxAttrInfo{*iibocai.c.attrInfo}
	return &value
}

func (iibocai IdxIBOutletCollectionAttrInfo) ObjcClass() *IdxEntityInfo {
	value := IdxEntityInfo{*iibocai.c.objcClass}
	return &value
}

func (iibocai IdxIBOutletCollectionAttrInfo) ClassCursor() Cursor {
	value := Cursor{iibocai.c.classCursor}
	return value
}

func (iibocai IdxIBOutletCollectionAttrInfo) ClassLoc() IdxLoc {
	value := IdxLoc{iibocai.c.classLoc}
	return value
}
