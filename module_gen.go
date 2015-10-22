package phoenix

// #include "go-clang.h"
import "C"

// \defgroup CINDEX_MODULE Module introspection The functions in this group provide access to information about modules. @{
type Module struct {
	c C.CXModule
}

// \param Module a module object. \returns the name of the module, e.g. for the 'std.vector' sub-module it will return "vector".
func (m Module) Name() string {
	o := cxstring{C.clang_Module_getName(m.c)}
	defer o.Dispose()

	return o.String()
}

// \param Module a module object. \returns the full name of the module, e.g. "std.vector".
func (m Module) FullName() string {
	o := cxstring{C.clang_Module_getFullName(m.c)}
	defer o.Dispose()

	return o.String()
}
