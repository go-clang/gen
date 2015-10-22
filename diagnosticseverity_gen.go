package phoenix

// #include "go-clang.h"
import "C"

// Describes the severity of a particular diagnostic.
type DiagnosticSeverity uint32

const (
	// A diagnostic that has been suppressed, e.g., by a command-line option.
	Diagnostic_Ignored DiagnosticSeverity = C.CXDiagnostic_Ignored
	// This diagnostic is a note that should be attached to the previous (non-note) diagnostic.
	Diagnostic_Note = C.CXDiagnostic_Note
	// This diagnostic indicates suspicious code that may not be wrong.
	Diagnostic_Warning = C.CXDiagnostic_Warning
	// This diagnostic indicates that the code is ill-formed.
	Diagnostic_Error = C.CXDiagnostic_Error
	// This diagnostic indicates that the code is ill-formed such that future parser recovery is unlikely to produce useful results.
	Diagnostic_Fatal = C.CXDiagnostic_Fatal
)
