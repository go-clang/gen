package gen_test

import (
	"testing"

	"github.com/go-clang/gen"
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

			if got, want := gen.UpperFirstCharacter(tt.data), tt.expect; got != want {
				t.Fatalf("got %s but want %s", got, want)
			}
		})
	}
}

func TestTrimCommonFunctionNamePrefix(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		name string
		want string
	}{
		"createCXCursorSet": {
			name: "createCXCursorSet",
			want: "CursorSet",
		},
		"getKind": {
			name: "getKind",
			want: "Kind",
		},
		"GetObjCSelector": {
			name: "GetObjCSelector",
			want: "Selector",
		},
		"getCXXManglings": {
			name: "getCXXManglings",
			want: "CXXManglings",
		},
		"getObjCManglings": {
			name: "getObjCManglings",
			want: "ObjCManglings",
		},
		"SourceLocation": {
			name: "SourceLocation",
			want: "SourceLocation",
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := gen.TrimCommonFunctionNamePrefix(tt.name); got != tt.want {
				t.Fatalf("TrimCommonFunctionNamePrefix(%v) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
