package phoenix

// #include "go-clang.h"
import "C"

import (
	"unsafe"
)

// Retrieve the set of display options most similar to the default behavior of the clang compiler. \returns A set of display options suitable for use with \c clang_formatDiagnostic().
func DefaultDiagnosticDisplayOptions() uint16 {
	return uint16(C.clang_defaultDiagnosticDisplayOptions())
}

// Retrieve the name of a particular diagnostic category. This is now deprecated. Use clang_getDiagnosticCategoryText() instead. \param Category A diagnostic category number, as returned by \c clang_getDiagnosticCategory(). \returns The name of the given diagnostic category.
func GetDiagnosticCategoryName(Category uint16) string {
	o := cxstring{C.clang_getDiagnosticCategoryName(C.uint(Category))}
	defer o.Dispose()

	return o.String()
}

// Returns the set of flags that is suitable for parsing a translation unit that is being edited. The set of flags returned provide options for \c clang_parseTranslationUnit() to indicate that the translation unit is likely to be reparsed many times, either explicitly (via \c clang_reparseTranslationUnit()) or implicitly (e.g., by code completion (\c clang_codeCompletionAt())). The returned flag set contains an unspecified set of optimizations (e.g., the precompiled preamble) geared toward improving the performance of these routines. The set of optimizations enabled may change from one version to the next.
func DefaultEditingTranslationUnitOptions() uint16 {
	return uint16(C.clang_defaultEditingTranslationUnitOptions())
}

// Construct a USR for a specified Objective-C class.
func ConstructUSR_ObjCClass(class_name string) string {
	c_class_name := C.CString(class_name)
	defer C.free(unsafe.Pointer(c_class_name))

	o := cxstring{C.clang_constructUSR_ObjCClass(c_class_name)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C category.
func ConstructUSR_ObjCCategory(class_name string, category_name string) string {
	c_class_name := C.CString(class_name)
	defer C.free(unsafe.Pointer(c_class_name))
	c_category_name := C.CString(category_name)
	defer C.free(unsafe.Pointer(c_category_name))

	o := cxstring{C.clang_constructUSR_ObjCCategory(c_class_name, c_category_name)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C protocol.
func ConstructUSR_ObjCProtocol(protocol_name string) string {
	c_protocol_name := C.CString(protocol_name)
	defer C.free(unsafe.Pointer(c_protocol_name))

	o := cxstring{C.clang_constructUSR_ObjCProtocol(c_protocol_name)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C instance variable and the USR for its containing class.
func ConstructUSR_ObjCIvar(name string, classUSR cxstring) string {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	o := cxstring{C.clang_constructUSR_ObjCIvar(c_name, classUSR.c)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C method and the USR for its containing class.
func ConstructUSR_ObjCMethod(name string, isInstanceMethod uint16, classUSR cxstring) string {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	o := cxstring{C.clang_constructUSR_ObjCMethod(c_name, C.uint(isInstanceMethod), classUSR.c)}
	defer o.Dispose()

	return o.String()
}

// Construct a USR for a specified Objective-C property and the USR for its containing class.
func ConstructUSR_ObjCProperty(property string, classUSR cxstring) string {
	c_property := C.CString(property)
	defer C.free(unsafe.Pointer(c_property))

	o := cxstring{C.clang_constructUSR_ObjCProperty(c_property, classUSR.c)}
	defer o.Dispose()

	return o.String()
}

func EnableStackTraces() {
	C.clang_enableStackTraces()
}

// Returns a default set of code-completion options that can be passed to\c clang_codeCompleteAt().
func DefaultCodeCompleteOptions() uint16 {
	return uint16(C.clang_defaultCodeCompleteOptions())
}

// Return a version string, suitable for showing to a user, but not intended to be parsed (the format is not guaranteed to be stable).
func GetClangVersion() string {
	o := cxstring{C.clang_getClangVersion()}
	defer o.Dispose()

	return o.String()
}

// Enable/disable crash recovery. \param isEnabled Flag to indicate if crash recovery is enabled. A non-zero value enables crash recovery, while 0 disables it.
func ToggleCrashRecovery(isEnabled uint16) {
	C.clang_toggleCrashRecovery(C.uint(isEnabled))
}
