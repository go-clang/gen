package gen

type includeFiles map[string]struct{}

// NewIncludeFiles returns the new include files map.
func NewIncludeFiles() includeFiles {
	return includeFiles(map[string]struct{}{})
}

// AddIncludeFile adds include file to inf.
func (inf includeFiles) AddIncludeFile(includeFile string) {
	inf[includeFile] = struct{}{}
}

// unifyIncludeFiles unified include files.
func (inf includeFiles) unifyIncludeFiles(m includeFiles) {
	for i := range m {
		inf[i] = struct{}{}
	}
}
