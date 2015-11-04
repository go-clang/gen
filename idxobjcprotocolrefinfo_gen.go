package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCProtocolRefInfo struct {
	c C.CXIdxObjCProtocolRefInfo
}

func (iocpri IdxObjCProtocolRefInfo) Protocol() *IdxEntityInfo {
	o := *iocpri.c.protocol

	return &IdxEntityInfo{o}
}

func (iocpri IdxObjCProtocolRefInfo) Cursor() Cursor {
	return Cursor{iocpri.c.cursor}
}

func (iocpri IdxObjCProtocolRefInfo) Loc() IdxLoc {
	return IdxLoc{iocpri.c.loc}
}
