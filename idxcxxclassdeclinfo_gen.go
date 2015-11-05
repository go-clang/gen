package phoenix

// #include "go-clang.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type IdxCXXClassDeclInfo struct {
	c C.CXIdxCXXClassDeclInfo
}

func (icxxcdi IdxCXXClassDeclInfo) DeclInfo() *IdxDeclInfo {
	o := icxxcdi.c.declInfo

	var gop_o *IdxDeclInfo
	if o != nil {
		gop_o = &IdxDeclInfo{o}
	}

	return gop_o
}

func (icxxcdi IdxCXXClassDeclInfo) Bases() []*IdxBaseClassInfo {
	var s []*IdxBaseClassInfo
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(icxxcdi.c.numBases)
	gos_s.Len = int(icxxcdi.c.numBases)
	gos_s.Data = uintptr(unsafe.Pointer(icxxcdi.c.bases))

	return s
}

func (icxxcdi IdxCXXClassDeclInfo) NumBases() uint16 {
	return uint16(icxxcdi.c.numBases)
}
