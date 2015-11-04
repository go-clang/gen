package phoenix

// #include "go-clang.h"
import "C"

type IdxObjCInterfaceDeclInfo struct {
	c C.CXIdxObjCInterfaceDeclInfo
}

func (iocidi IdxObjCInterfaceDeclInfo) ContainerInfo() *IdxObjCContainerDeclInfo {
	o := *iocidi.c.containerInfo

	return &IdxObjCContainerDeclInfo{o}
}

func (iocidi IdxObjCInterfaceDeclInfo) SuperInfo() *IdxBaseClassInfo {
	o := *iocidi.c.superInfo

	return &IdxBaseClassInfo{o}
}

func (iocidi IdxObjCInterfaceDeclInfo) Protocols() *IdxObjCProtocolRefListInfo {
	o := *iocidi.c.protocols

	return &IdxObjCProtocolRefListInfo{o}
}
