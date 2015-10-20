package phoenix

// #include "go-clang.h"
import "C"

// \defgroup CINDEX_MODULE Module introspection The functions in this group provide access to information about modules. @{
type Module struct {
	c C.CXModule
}

// \param Module a module object. \returns the name of the module, e.g. for the 'std.vector' sub-module it will return "vector".
func (m Module) Name() string {
	cstr := cxstring{C.clang_Module_getName(m.c)}
	defer cstr.Dispose()

	return cstr.String()
}

// \param Module a module object. \returns the full name of the module, e.g. "std.vector".
func (m Module) FullName() string {
	cstr := cxstring{C.clang_Module_getFullName(m.c)}
	defer cstr.Dispose()

	return cstr.String()
}
