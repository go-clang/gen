package generate

type API struct {
	ClangArguments []string
}

func (a *API) HandleHeaderFile(filename string) error {
	return handleHeaderFile(filename, a.ClangArguments)
}
