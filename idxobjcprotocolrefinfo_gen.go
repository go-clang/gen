package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCProtocolRefInfo struct {
	c C.CXIdxObjCProtocolRefInfo
}

func (iocpri IdxObjCProtocolRefInfo) Protocol() *IdxEntityInfo {
	value := IdxEntityInfo{*iocpri.c.protocol}
	return &value
}

func (iocpri IdxObjCProtocolRefInfo) Cursor() Cursor {
	value := Cursor{iocpri.c.cursor}
	return value
}

func (iocpri IdxObjCProtocolRefInfo) Loc() IdxLoc {
	value := IdxLoc{iocpri.c.loc}
	return value
}
