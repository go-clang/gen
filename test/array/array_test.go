package array

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-clang/gen"
	"github.com/go-clang/gen/test/expectations"
)

var testAPI *gen.API = &gen.API{
	PrepareFunctionName:     prepareFunctionName,
	PrepareFunction:         prepareFunction,
	FilterFunction:          filterFunction,
	FilterFunctionParameter: filterFunctionParameter,
	FixedFunctionName:       fixedFunctionName,

	PrepareStructMembers:     prepareStructMembers,
	FilterStructMemberGetter: filterStructMemberGetter,
}

var testDataFolder string = "../testdata/"

func TestDriver(t *testing.T) {
	// parse and generate test cases
	h, err := testAPI.HandleDirectory(testDataFolder + "test-cases/")
	if err != nil {
		panic(err)
	}

	expctPerConstruct := expectations.BuildUpExpectations(testDataFolder + "expected/")

	success := make(map[string]bool)
	for _, st := range h.GetStructs() {
		dat, err := ioutil.ReadFile(fmt.Sprintf("%s_gen.go", strings.ToLower(st.Name)))
		if err != nil {
			panic(err)
		}
		success[st.Name] = expectations.VerifyExpectations(string(dat), expctPerConstruct[st.Name]) //trick the branch optimization
	}

	for k, v := range success {
		assert.True(t, v, fmt.Sprintf("%s did not fulfill expectations", k))
	}

}

func prepareFunctionName(h *gen.HeaderFile, f *gen.Function) string {
	return f.Name
}

func fixedFunctionName(f *gen.Function) string {
	return ""
}

func prepareFunction(f *gen.Function) {
	for i := range f.Parameters {
		p := &f.Parameters[i]

		if p.Type.CGoName == gen.CSChar && p.Type.PointerLevel == 2 && !p.Type.IsSlice {
			p.Type.IsReturnArgument = true
		}
	}
}

func filterFunction(f *gen.Function) bool {
	return true
}

func filterFunctionParameter(p gen.FunctionParameter) bool {
	return true
}

func prepareStructMembers(s *gen.Struct) {
}

func filterStructMemberGetter(m *gen.StructMember) bool {
	// We do not want getters to *int_data member
	return true
}
