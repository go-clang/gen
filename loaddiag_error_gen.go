package phoenix

// #include "go-clang.h"
import "C"

// Describes the kind of error that occurred (if any) in a call to \c clang_loadDiagnostics.
type LoadDiag_Error int32

const (
	// Indicates that no error occurred.
	LoadDiag_None LoadDiag_Error = C.CXLoadDiag_None
	// Indicates that an unknown error occurred while attempting to deserialize diagnostics.
	LoadDiag_Unknown = C.CXLoadDiag_Unknown
	// Indicates that the file containing the serialized diagnostics could not be opened.
	LoadDiag_CannotLoad = C.CXLoadDiag_CannotLoad
	// Indicates that the serialized diagnostics file is invalid or corrupt.
	LoadDiag_InvalidFile = C.CXLoadDiag_InvalidFile
)
