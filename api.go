package gen

import (
	"fmt"
	"os"
	"strings"
)

// API represents a Clang bindings generation.
type API struct {
	// ClangArguments holds the command line arguments for Clang
	ClangArguments []string

	// PrepareFunctionName returns a prepared function name for further processing
	PrepareFunctionName func(g *Generation, f *Function) string

	// PrepareFunction prepares a function for further processing
	PrepareFunction func(f *Function)

	// FilterFunction determines if a function is generateable
	FilterFunction func(f *Function) bool

	// FilterFunctionParameter determines if a function parameter is generateable
	FilterFunctionParameter func(p FunctionParameter) bool

	// FixedFunctionName returns an unempty string if a function needs to receive a specific name
	FixedFunctionName func(f *Function) string

	// PrepareStructFields is called before adding struct field getters
	PrepareStructFields func(s *Struct)

	// FilterStructFieldGetter determines if a getter should be generated for a field
	FilterStructFieldGetter func(m *StructField) bool
}

// HandleDirectory handles header files on dir and returns the *HeaderFile slice.
func (a *API) HandleDirectory(dir string) ([]*HeaderFile, error) {
	headers, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read clang-c directory: %w", err)
	}

	headerFiles := make([]*HeaderFile, 0, len(headers))
	for _, hf := range headers {
		if hf.IsDir() || !strings.HasSuffix(hf.Name(), ".h") {
			continue
		}

		h := NewHeaderFile(a, hf.Name(), dir)

		if err := h.Parse(a.ClangArguments); err != nil {
			return nil, fmt.Errorf("cannot handle header file %q: %w", h.FullPath(), err)
		}

		headerFiles = append(headerFiles, h)
	}

	return headerFiles, nil
}
