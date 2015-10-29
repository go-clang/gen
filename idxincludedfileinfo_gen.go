package phoenix

// #include "go-clang.h"
import "C"

// Data for ppIncludedFile callback.
type IdxIncludedFileInfo struct {
	c C.CXIdxIncludedFileInfo
}

// Location of '#' in the \#include/\#import directive.
func (iifi IdxIncludedFileInfo) HashLoc() IdxLoc {
	value := IdxLoc{iifi.c.hashLoc}
	return value
}

// Filename as written in the \#include/\#import directive.
func (iifi IdxIncludedFileInfo) Filename() *int8 {
	value := int8(*iifi.c.filename)
	return &value
}

// The actual file that the \#include/\#import directive resolved to.
func (iifi IdxIncludedFileInfo) File() File {
	value := File{iifi.c.file}
	return value
}

func (iifi IdxIncludedFileInfo) IsImport() bool {
	value := iifi.c.isImport
	return value != C.int(0)
}

func (iifi IdxIncludedFileInfo) IsAngled() bool {
	value := iifi.c.isAngled
	return value != C.int(0)
}

// Non-zero if the directive was automatically turned into a module import.
func (iifi IdxIncludedFileInfo) IsModuleImport() bool {
	value := iifi.c.isModuleImport
	return value != C.int(0)
}
