package phoenix

// #include "go-clang.h"
import "C"
import "unsafe"

type IdxEntityInfo struct {
	c C.CXIdxEntityInfo
}

// For retrieving a custom CXIdxClientEntity attached to an entity.
func (iei *IdxEntityInfo) Index_getClientEntity() IdxClientEntity {
	return IdxClientEntity{C.clang_index_getClientEntity(&iei.c)}
}

// For setting a custom CXIdxClientEntity attached to an entity.
func (iei *IdxEntityInfo) Index_setClientEntity(ice IdxClientEntity) {
	C.clang_index_setClientEntity(&iei.c, ice.c)
}

func (iei IdxEntityInfo) Kind() IdxEntityKind {
	return IdxEntityKind(iei.c.kind)
}

func (iei IdxEntityInfo) TemplateKind() IdxEntityCXXTemplateKind {
	return IdxEntityCXXTemplateKind(iei.c.templateKind)
}

func (iei IdxEntityInfo) Lang() IdxEntityLanguage {
	return IdxEntityLanguage(iei.c.lang)
}

func (iei IdxEntityInfo) Name() string {
	return C.GoString(iei.c.name)
}

func (iei IdxEntityInfo) USR() string {
	return C.GoString(iei.c.USR)
}

func (iei IdxEntityInfo) Cursor() Cursor {
	return Cursor{iei.c.cursor}
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
	return uint16(iei.c.numAttributes)
}
