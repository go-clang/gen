package gen

// IncludeFiles represents a generation include files.
type IncludeFiles map[string]struct{}

// NewIncludeFiles returns the new include files map.
func NewIncludeFiles() IncludeFiles {
	return IncludeFiles(map[string]struct{}{})
}

// AddIncludeFile adds include file to inf.
func (inf IncludeFiles) AddIncludeFile(includeFile string) {
	inf[includeFile] = struct{}{}
}

// unifyIncludeFiles unified include files.
func (inf IncludeFiles) unifyIncludeFiles(m IncludeFiles) {
	for i := range m {
		inf[i] = struct{}{}
	}
}
