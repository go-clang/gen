package gen

import (
	"testing"
)

func TestUpperFirstCharacter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		data   string
		expect string
	}{
		{
			name:   "lowerToUpperOneChar",
			data:   "a",
			expect: "A",
		},
		{
			name:   "UpperToUpperOneChar",
			data:   "A",
			expect: "A",
		},
		{
			name:   "lowerCaseToUpperCase",
			data:   "abc",
			expect: "Abc",
		},
		{
			name:   "UpperCaseToUpperCase",
			data:   "Abc",
			expect: "Abc",
		},
		{
			name:   "UpperToUpper",
			data:   "ABC",
			expect: "ABC",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got, want := UpperFirstCharacter(tt.data), tt.expect; got != want {
				t.Fatalf("got %s but want %s", got, want)
			}
		})
	}
}
