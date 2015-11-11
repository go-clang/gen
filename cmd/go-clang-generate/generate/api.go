package generate

type API struct {
	// ClangArguments holds the command line arguments for Clang
	ClangArguments []string

	// PrepareFunctionName returns a prepared function name for further processing
	PrepareFunctionName func(h *HeaderFile, f *Function) string
	// PrepareFunction prepares a function for further processing
	PrepareFunction func(f *Function)
	// FilterFunction determines if a function is generateable
	FilterFunction func(f *Function) bool
	// FilterFunctionParameter determines if a function parameter is generateable
	FilterFunctionParameter func(p FunctionParameter) bool
	// FixedFunctionName returns an unempty string if a function needs to receive a specific name
	FixedFunctionName func(f *Function) string

	// PrepareStructMembers is called before adding struct member getters
	PrepareStructMembers func(s *Struct)
	// FilterStructMemberGetter determines if a getter should be generated for a member
	FilterStructMemberGetter func(m *StructMember) bool
}

func (a *API) HandleHeaderFile(filename string) error {
	return handleHeaderFile(a, filename, a.ClangArguments)
}
