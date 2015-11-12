package generate

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type API struct {
	// ClangArguments hold command line arguments for Clang
	ClangArguments []string

	// PrepareFunctionName returns a prepared function name for further processing
	PrepareFunctionName func(h *HeaderFile, f *Function) string
	// PrepareFunction prepare a function for further processing
	PrepareFunction func(f *Function)
	// FilterFunction determines if a function is generateable
	FilterFunction func(f *Function) bool
	// FilterFunctionParameter determines if a function parameter is generateable
	FilterFunctionParameter func(p FunctionParameter) bool
	// FixedFunctionName returns an unempty string if a function needs to receive a specific name
	FixedFunctionName func(f *Function) string

	// PrepareStructMembers called before adding struct member getters
	PrepareStructMembers func(s *Struct)
	// FilterStructMemberGetter determines if a getter should be generated for a member
	FilterStructMemberGetter func(m *StructMember) bool
}

func (a *API) HandleDirectory(dir string) error {
	headers, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Cannot list clang-c directory: %v", err)
	}

	h := &HeaderFile{
		api: a,

		dir: dir,

		lookupEnum:        map[string]*Enum{},
		lookupNonTypedefs: map[string]string{},
		lookupStruct: map[string]*Struct{
			"cxstring": &Struct{
				Name:  "cxstring",
				CName: "CXString",
			},
		},
	}

	for _, hf := range headers {
		name := hf.Name()

		if hf.IsDir() || !strings.HasSuffix(name, ".h") {
			continue
		}

		h.name = dir + name

		if err := h.parse(a.ClangArguments); err != nil {
			return fmt.Errorf("Cannot handle header file %q: %v", name, err)
		}
	}

	if err := h.handle(); err != nil {
		return err
	}

	return nil
}
