package phoenix

// #include "go-clang.h"
import "C"
import "unsafe"

// Uniquely identifies a CXFile, that refers to the same underlying file, across an indexing session.
type FileUniqueID struct {
	c C.CXFileUniqueID
}

func (fuid FileUniqueID) Data() []uint64 {
	sc := []uint64{}

	length := 3
	goslice := (*[1 << 30]C.ulonglong)(unsafe.Pointer(&fuid.c.data))[:length:length]

	for is := 0; is < length; is++ {
		sc = append(sc, uint64(goslice[is]))
	}

	return sc
}
