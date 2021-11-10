package clang

// TODO this can be generated https://github.com/go-clang/gen/issues/47

// Diagnostics determine the number of diagnostics produced prior to the
// location where code completion was performed.
func (ccr *CodeCompleteResults) Diagnostics() []Diagnostic {
	s := make([]Diagnostic, ccr.NumDiagnostics())

	for i := range s {
		s[i] = ccr.Diagnostic(uint32(i))
	}

	return s
}
