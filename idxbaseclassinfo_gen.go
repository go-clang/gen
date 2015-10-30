package phoenix

// #include "go-clang.h"
import "C"

type IdxBaseClassInfo struct {
	c C.CXIdxBaseClassInfo
}

func (ibci IdxBaseClassInfo) Base() *IdxEntityInfo {
	value := IdxEntityInfo{*ibci.c.base}
	return &value
}

func (ibci IdxBaseClassInfo) Cursor() Cursor {
	value := Cursor{ibci.c.cursor}
	return value
}

func (ibci IdxBaseClassInfo) Loc() IdxLoc {
	value := IdxLoc{ibci.c.loc}
	return value
}
