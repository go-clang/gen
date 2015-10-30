package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

type IdxEntityInfo struct {
	c C.CXIdxEntityInfo
}

func (iei IdxEntityInfo) Kind() IdxEntityKind {
	value := IdxEntityKind(iei.c.kind)
	return value
}

func (iei IdxEntityInfo) TemplateKind() IdxEntityCXXTemplateKind {
	value := IdxEntityCXXTemplateKind(iei.c.templateKind)
	return value
}

func (iei IdxEntityInfo) Lang() IdxEntityLanguage {
	value := IdxEntityLanguage(iei.c.lang)
	return value
}

func (iei IdxEntityInfo) Name() *int8 {
	value := int8(*iei.c.name)
	return &value
}

func (iei IdxEntityInfo) USR() *int8 {
	value := int8(*iei.c.USR)
	return &value
}

func (iei IdxEntityInfo) Cursor() Cursor {
	value := Cursor{iei.c.cursor}
	return value
}

func (iei IdxEntityInfo) Attributes() []*IdxAttrInfo {
	sc := []*IdxAttrInfo{}

	length := int(iei.c.numAttributes)
	goslice := (*[1 << 30]*C.CXIdxAttrInfo)(unsafe.Pointer(&iei.c.attributes))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, &IdxAttrInfo{*goslice[is]})
	}

	return sc
}

func (iei IdxEntityInfo) NumAttributes() uint16 {
	value := uint16(iei.c.numAttributes)
	return value
}
