package gen

import (
	"reflect"
	"testing"
)

func TestUpperFirstCharacter(t *testing.T) {
	for _, tc := range []struct {
		Data   string
		Expect string
	}{
		{
			Data:   "a",
			Expect: "A",
		},
		{
			Data:   "A",
			Expect: "A",
		},
		{
			Data:   "abc",
			Expect: "Abc",
		},
		{
			Data:   "Abc",
			Expect: "Abc",
		},
		{
			Data:   "ABC",
			Expect: "ABC",
		},
	} {
		reflect.DeepEqual(tc.Expect, UpperFirstCharacter(tc.Data))
	}
}
