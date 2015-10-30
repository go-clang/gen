package phoenix

// #include "go-clang.h"
import "C"

// Describes the availability of a given entity on a particular platform, e.g., a particular class might only be available on Mac OS 10.7 or newer.
type PlatformAvailability struct {
	c C.CXPlatformAvailability
}

// A string that describes the platform for which this structure provides availability information. Possible values are "ios" or "macosx".
func (pa PlatformAvailability) Platform() cxstring {
	value := cxstring{pa.c.Platform}
	return value
}

// The version number in which this entity was introduced.
func (pa PlatformAvailability) Introduced() Version {
	value := Version{pa.c.Introduced}
	return value
}

// The version number in which this entity was deprecated (but is still available).
func (pa PlatformAvailability) Deprecated() Version {
	value := Version{pa.c.Deprecated}
	return value
}

// The version number in which this entity was obsoleted, and therefore is no longer available.
func (pa PlatformAvailability) Obsoleted() Version {
	value := Version{pa.c.Obsoleted}
	return value
}

// Whether the entity is unconditionally unavailable on this platform.
func (pa PlatformAvailability) Unavailable() int32 {
	value := int32(pa.c.Unavailable)
	return value
}

// An optional message to provide to a user of this API, e.g., to suggest replacement APIs.
func (pa PlatformAvailability) Message() cxstring {
	value := cxstring{pa.c.Message}
	return value
}
