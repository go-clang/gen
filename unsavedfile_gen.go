package phoenix

// #include "go-clang.h"
import "C"

/*
	Provides the contents of a file that has not yet been saved to disk.

	Each CXUnsavedFile instance provides the name of a file on the
	system along with the current contents of that file that have not
	yet been saved to disk.
*/
type UnsavedFile struct {
	c C.struct_CXUnsavedFile
}

/*
	The file whose contents have not yet been saved.

	This file must already exist in the file system.
*/
func (uf UnsavedFile) Filename() *int8 {
	value := int8(*uf.c.Filename)
	return &value
}

// A buffer containing the unsaved contents of this file.
func (uf UnsavedFile) Contents() *int8 {
	value := int8(*uf.c.Contents)
	return &value
}

// The length of the unsaved contents of this buffer.
func (uf UnsavedFile) Length() uint32 {
	value := uint32(uf.c.Length)
	return value
}
