package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCInterfaceDeclInfo struct {
	c C.CXIdxObjCInterfaceDeclInfo
}

func (iocidi IdxObjCInterfaceDeclInfo) ContainerInfo() *IdxObjCContainerDeclInfo {
	value := IdxObjCContainerDeclInfo{*iocidi.c.containerInfo}
	return &value
}

func (iocidi IdxObjCInterfaceDeclInfo) SuperInfo() *IdxBaseClassInfo {
	value := IdxBaseClassInfo{*iocidi.c.superInfo}
	return &value
}

func (iocidi IdxObjCInterfaceDeclInfo) Protocols() *IdxObjCProtocolRefListInfo {
	value := IdxObjCProtocolRefListInfo{*iocidi.c.protocols}
	return &value
}
