package phoenix

// #include "go-clang.h"
import "C"
import "unsafe"

type IdxCXXClassDeclInfo struct {
	c C.CXIdxCXXClassDeclInfo
}

func (icxxcdi IdxCXXClassDeclInfo) DeclInfo() *IdxDeclInfo {
	value := IdxDeclInfo{*icxxcdi.c.declInfo}
	return &value
}

func (icxxcdi IdxCXXClassDeclInfo) Bases() []*IdxBaseClassInfo {
	sc := []*IdxBaseClassInfo{}

	length := int(icxxcdi.c.numBases)
	goslice := (*[1 << 30]*C.CXIdxBaseClassInfo)(unsafe.Pointer(&icxxcdi.c.bases))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, &IdxBaseClassInfo{*goslice[is]})
	}

	return sc
}

func (icxxcdi IdxCXXClassDeclInfo) NumBases() uint16 {
	value := uint16(icxxcdi.c.numBases)
	return value
}
