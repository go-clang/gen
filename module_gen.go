package phoenix

// #include "go-clang.h"
import "C"

// \defgroup CINDEX_MODULE Module introspection The functions in this group provide access to information about modules. @{
type Module struct {
	c C.CXModule
}
