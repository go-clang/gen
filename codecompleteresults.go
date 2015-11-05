package phoenix

func (ccr *CodeCompleteResults) Diagnostics() []Diagnostic { // TODO this can be generated
	s := make([]Diagnostic, ccr.CodeCompleteGetNumDiagnostics())

	for i := range s {
		s[i] = ccr.CodeCompleteGetDiagnostic(uint16(i))
	}

	return s
}
