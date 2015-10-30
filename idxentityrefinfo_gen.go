package phoenix

// #include "go-clang.h"
import "C"

// Data for IndexerCallbacks#indexEntityReference.
type IdxEntityRefInfo struct {
	c C.CXIdxEntityRefInfo
}

func (ieri IdxEntityRefInfo) Kind() IdxEntityRefKind {
	value := IdxEntityRefKind(ieri.c.kind)
	return value
}

// Reference cursor.
func (ieri IdxEntityRefInfo) Cursor() Cursor {
	value := Cursor{ieri.c.cursor}
	return value
}

func (ieri IdxEntityRefInfo) Loc() IdxLoc {
	value := IdxLoc{ieri.c.loc}
	return value
}

// The entity that gets referenced.
func (ieri IdxEntityRefInfo) ReferencedEntity() *IdxEntityInfo {
	value := IdxEntityInfo{*ieri.c.referencedEntity}
	return &value
}

/*
 * \brief Immediate "parent" of the reference. For example:
 *
 * \code
 * Foo *var;
 * \endcode
 *
 * The parent of reference of type 'Foo' is the variable 'var'.
 * For references inside statement bodies of functions/methods,
 * the parentEntity will be the function/method.
 */
func (ieri IdxEntityRefInfo) ParentEntity() *IdxEntityInfo {
	value := IdxEntityInfo{*ieri.c.parentEntity}
	return &value
}

// Lexical container context of the reference.
func (ieri IdxEntityRefInfo) Container() *IdxContainerInfo {
	value := IdxContainerInfo{*ieri.c.container}
	return &value
}
