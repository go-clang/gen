package phoenix

// #include "go-clang.h"
import "C"

// Describes a version number of the form major.minor.subminor.
type Version struct {
	c C.CXVersion
}

// The major version number, e.g., the '10' in '10.7.3'. A negative value indicates that there is no version number at all.
func (v Version) Major() int16 {
	return int16(v.c.Major)
}

// The minor version number, e.g., the '7' in '10.7.3'. This value will be negative if no minor version number was provided, e.g., for version '10'.
func (v Version) Minor() int16 {
	return int16(v.c.Minor)
}

// The subminor version number, e.g., the '3' in '10.7.3'. This value will be negative if no minor or subminor version number was provided, e.g., in version '10' or '10.7'.
func (v Version) Subminor() int16 {
	return int16(v.c.Subminor)
}
