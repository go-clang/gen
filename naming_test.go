package gen

import (
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
		if tc.Expect != UpperFirstCharacter(tc.Data) {
			t.Fatalf("got %s but want %s", UpperFirstCharacter(tc.Data), tc.Expect)
		}
	}
}
