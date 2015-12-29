package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

var testDataFolder string = "../../test/testdata/"

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
		err := os.Remove(fmt.Sprintf("%s_gen.go", strings.ToLower(k)))
		if err != nil {
			panic(err)
		}
	}

}
