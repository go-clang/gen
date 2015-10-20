package phoenix

// #include "go-clang.h"
import "C"

// Describes a version number of the form major.minor.subminor.
type Version struct {
	c C.CXVersion
}
