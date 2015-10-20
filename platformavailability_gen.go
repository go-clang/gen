package phoenix

// #include "go-clang.h"
import "C"

// Describes the availability of a given entity on a particular platform, e.g., a particular class might only be available on Mac OS 10.7 or newer.
type PlatformAvailability struct {
	c C.CXPlatformAvailability
}
