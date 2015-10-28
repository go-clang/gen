package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// Retrieve the name of a particular diagnostic category. This is now deprecated. Use clang_getDiagnosticCategoryText() instead. \param Category A diagnostic category number, as returned by \c clang_getDiagnosticCategory(). \returns The name of the given diagnostic category.
func getDiagnosticCategoryName(Category uint16) string {
	o := cxstring{C.clang_getDiagnosticCategoryName(C.uint(Category))}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C class.
func constructUSR_ObjCClass(class_name string) string {
	c_class_name := C.CString(class_name)
	defer C.free(unsafe.Pointer(c_class_name))

	o := cxstring{C.clang_constructUSR_ObjCClass(c_class_name)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C category.
func constructUSR_ObjCCategory(class_name string, category_name string) string {
	c_class_name := C.CString(class_name)
	defer C.free(unsafe.Pointer(c_class_name))
	c_category_name := C.CString(category_name)
	defer C.free(unsafe.Pointer(c_category_name))

	o := cxstring{C.clang_constructUSR_ObjCCategory(c_class_name, c_category_name)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C protocol.
func constructUSR_ObjCProtocol(protocol_name string) string {
	c_protocol_name := C.CString(protocol_name)
	defer C.free(unsafe.Pointer(c_protocol_name))

	o := cxstring{C.clang_constructUSR_ObjCProtocol(c_protocol_name)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C instance variable and the USR for its containing class.
func constructUSR_ObjCIvar(name string, classUSR cxstring) string {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	o := cxstring{C.clang_constructUSR_ObjCIvar(c_name, classUSR.c)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C method and the USR for its containing class.
func constructUSR_ObjCMethod(name string, isInstanceMethod uint16, classUSR cxstring) string {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	o := cxstring{C.clang_constructUSR_ObjCMethod(c_name, C.uint(isInstanceMethod), classUSR.c)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C property and the USR for its containing class.
func constructUSR_ObjCProperty(property string, classUSR cxstring) string {
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_property))

	o := cxstring{C.clang_constructUSR_ObjCProperty(c_property, classUSR.c)}
	defer o.Dispose()

	return o.String()
}

// Enable/disable crash recovery. \param isEnabled Flag to indicate if crash recovery is enabled. A non-zero value enables crash recovery, while 0 disables it.
func toggleCrashRecovery(isEnabled uint16) {
	C.clang_toggleCrashRecovery(C.uint(isEnabled))
}
