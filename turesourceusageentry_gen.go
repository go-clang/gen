package phoenix

// #include "go-clang.h"
import "C"

type TUResourceUsageEntry struct {
	c C.CXTUResourceUsageEntry
}
