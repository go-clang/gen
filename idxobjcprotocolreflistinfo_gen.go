package phoenix

// #include "go-clang.h"
import "C"
import "unsafe"

type IdxObjCProtocolRefListInfo struct {
	c C.CXIdxObjCProtocolRefListInfo
}

func (iocprli IdxObjCProtocolRefListInfo) Protocols() []*IdxObjCProtocolRefInfo {
	sc := []*IdxObjCProtocolRefInfo{}

	length := int(iocprli.c.numProtocols)
	goslice := (*[1 << 30]*C.CXIdxObjCProtocolRefInfo)(unsafe.Pointer(&iocprli.c.protocols))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, &IdxObjCProtocolRefInfo{*goslice[is]})
	}

	return sc
}

func (iocprli IdxObjCProtocolRefListInfo) NumProtocols() uint16 {
	return uint16(iocprli.c.numProtocols)
}
