package phoenix

// #include "go-clang.h"
import "C"

type IdxAttrInfo struct {
	c C.CXIdxAttrInfo
}

func (iai IdxAttrInfo) Kind() IdxAttrKind {
	value := IdxAttrKind(iai.c.kind)
	return value
}

func (iai IdxAttrInfo) Cursor() Cursor {
	value := Cursor{iai.c.cursor}
	return value
}

func (iai IdxAttrInfo) Loc() IdxLoc {
	value := IdxLoc{iai.c.loc}
	return value
}
