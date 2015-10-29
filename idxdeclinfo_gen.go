package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

type IdxDeclInfo struct {
	c C.CXIdxDeclInfo
}

func (idi IdxDeclInfo) EntityInfo() *IdxEntityInfo {
	value := IdxEntityInfo{*idi.c.entityInfo}
	return &value
}

func (idi IdxDeclInfo) Cursor() Cursor {
	value := Cursor{idi.c.cursor}
	return value
}

func (idi IdxDeclInfo) Loc() IdxLoc {
	value := IdxLoc{idi.c.loc}
	return value
}

func (idi IdxDeclInfo) SemanticContainer() *IdxContainerInfo {
	value := IdxContainerInfo{*idi.c.semanticContainer}
	return &value
}

// Generally same as #semanticContainer but can be different in cases like out-of-line C++ member functions.
func (idi IdxDeclInfo) LexicalContainer() *IdxContainerInfo {
	value := IdxContainerInfo{*idi.c.lexicalContainer}
	return &value
}

func (idi IdxDeclInfo) IsRedeclaration() bool {
	value := idi.c.isRedeclaration
	return value != C.int(0)
}

func (idi IdxDeclInfo) IsDefinition() bool {
	value := idi.c.isDefinition
	return value != C.int(0)
}

func (idi IdxDeclInfo) IsContainer() bool {
	value := idi.c.isContainer
	return value != C.int(0)
}

func (idi IdxDeclInfo) DeclAsContainer() *IdxContainerInfo {
	value := IdxContainerInfo{*idi.c.declAsContainer}
	return &value
}

// Whether the declaration exists in code or was created implicitly by the compiler, e.g. implicit objc methods for properties.
func (idi IdxDeclInfo) IsImplicit() bool {
	value := idi.c.isImplicit
	return value != C.int(0)
}

func (idi IdxDeclInfo) Attributes() []*IdxAttrInfo {
	sc := []*IdxAttrInfo{}

	length := int(idi.c.numAttributes)
	goslice := (*[1 << 30]*C.CXIdxAttrInfo)(unsafe.Pointer(&idi.c.attributes))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, &IdxAttrInfo{*goslice[is]})
	}

	return sc
}

func (idi IdxDeclInfo) NumAttributes() uint16 {
	value := uint16(idi.c.numAttributes)
	return value
}

func (idi IdxDeclInfo) Flags() uint16 {
	value := uint16(idi.c.flags)
	return value
}
