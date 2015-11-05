package phoenix

func (ccr *CodeCompleteResults) Diagnostics() []Diagnostic { // TODO this can be generated https://github.com/zimmski/go-clang-phoenix/issues/47
	s := make([]Diagnostic, ccr.NumDiagnostics())

	for i := range s {
		s[i] = ccr.Diagnostic(uint16(i))
	}

	return s
}
