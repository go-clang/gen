package phoenix

// #include "go-clang.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type IdxDeclInfo struct {
	c C.CXIdxDeclInfo
}

func (idi *IdxDeclInfo) Index_getObjCContainerDeclInfo() *IdxObjCContainerDeclInfo {
	o := *C.clang_index_getObjCContainerDeclInfo(&idi.c)

	return &IdxObjCContainerDeclInfo{o}
}

func (idi *IdxDeclInfo) Index_getObjCInterfaceDeclInfo() *IdxObjCInterfaceDeclInfo {
	o := *C.clang_index_getObjCInterfaceDeclInfo(&idi.c)

	return &IdxObjCInterfaceDeclInfo{o}
}

func (idi *IdxDeclInfo) Index_getObjCCategoryDeclInfo() *IdxObjCCategoryDeclInfo {
	o := *C.clang_index_getObjCCategoryDeclInfo(&idi.c)

	return &IdxObjCCategoryDeclInfo{o}
}

func (idi *IdxDeclInfo) Index_getObjCProtocolRefListInfo() *IdxObjCProtocolRefListInfo {
	o := *C.clang_index_getObjCProtocolRefListInfo(&idi.c)

	return &IdxObjCProtocolRefListInfo{o}
}

func (idi *IdxDeclInfo) Index_getObjCPropertyDeclInfo() *IdxObjCPropertyDeclInfo {
	o := *C.clang_index_getObjCPropertyDeclInfo(&idi.c)

	return &IdxObjCPropertyDeclInfo{o}
}

func (idi *IdxDeclInfo) Index_getCXXClassDeclInfo() *IdxCXXClassDeclInfo {
	o := *C.clang_index_getCXXClassDeclInfo(&idi.c)

	return &IdxCXXClassDeclInfo{o}
}

func (idi IdxDeclInfo) EntityInfo() *IdxEntityInfo {
	o := *idi.c.entityInfo

	return &IdxEntityInfo{o}
}

func (idi IdxDeclInfo) Cursor() Cursor {
	return Cursor{idi.c.cursor}
}

func (idi IdxDeclInfo) Loc() IdxLoc {
	return IdxLoc{idi.c.loc}
}

func (idi IdxDeclInfo) SemanticContainer() *IdxContainerInfo {
	o := *idi.c.semanticContainer

	return &IdxContainerInfo{o}
}

// Generally same as #semanticContainer but can be different in cases like out-of-line C++ member functions.
func (idi IdxDeclInfo) LexicalContainer() *IdxContainerInfo {
	o := *idi.c.lexicalContainer

	return &IdxContainerInfo{o}
}

func (idi IdxDeclInfo) IsRedeclaration() bool {
	o := idi.c.isRedeclaration

	return o != C.int(0)
}

func (idi IdxDeclInfo) IsDefinition() bool {
	o := idi.c.isDefinition

	return o != C.int(0)
}

func (idi IdxDeclInfo) IsContainer() bool {
	o := idi.c.isContainer

	return o != C.int(0)
}

func (idi IdxDeclInfo) DeclAsContainer() *IdxContainerInfo {
	o := *idi.c.declAsContainer

	return &IdxContainerInfo{o}
}

// Whether the declaration exists in code or was created implicitly by the compiler, e.g. implicit objc methods for properties.
func (idi IdxDeclInfo) IsImplicit() bool {
	o := idi.c.isImplicit

	return o != C.int(0)
}

func (idi IdxDeclInfo) Attributes() []*IdxAttrInfo {
	var s []*IdxAttrInfo
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(idi.c.numAttributes)
	gos_s.Len = int(idi.c.numAttributes)
	gos_s.Data = uintptr(unsafe.Pointer(idi.c.attributes))

	return s
}

func (idi IdxDeclInfo) NumAttributes() uint16 {
	return uint16(idi.c.numAttributes)
}

func (idi IdxDeclInfo) Flags() uint16 {
	return uint16(idi.c.flags)
}
