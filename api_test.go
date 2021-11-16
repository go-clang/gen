package gen_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-clang/gen"
	"github.com/go-clang/gen/cmd/go-clang-gen/runtime"
)

func TestAPIPrepareFunctionName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		PrepareFunctionName func(g *gen.Generation, f *gen.Function) string
		name                string
		cname               string
		want                string
	}{
		"clang_indexLoc_getFileLocation": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_indexLoc_getFileLocation",
			cname:               "clang_indexLoc_getFileLocation",
			want:                "getFileLocation",
		},
		"clang_index_isEntityObjCContainerKind": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_index_isEntityObjCContainerKind",
			cname:               "clang_index_isEntityObjCContainerKind",
			want:                "isEntityObjCContainerKind",
		},
		"clang_Location_isInSystemHeader": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_Location_isInSystemHeader",
			cname:               "clang_Location_isInSystemHeader",
			want:                "isInSystemHeader",
		},
		"clang_Range_isNull": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_Range_isNull",
			cname:               "clang_Range_isNull",
			want:                "isNull",
		},
		"clang_remap_getNumFiles": {
			PrepareFunctionName: runtime.PrepareFunctionName,
			name:                "clang_remap_getNumFiles",
			cname:               "clang_remap_getNumFiles",
			want:                "getNumFiles",
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api := &gen.API{
				PrepareFunctionName: tt.PrepareFunctionName,
			}
			g := gen.NewGeneration(api)
			f := &gen.Function{
				Name:       tt.name,
				CName:      tt.cname,
				Parameters: []gen.FunctionParameter{},
			}

			got := g.API().PrepareFunctionName(g, f)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("API.PrepareFunctionName(): (-want +got):\n%s", diff)
			}
		})
	}
}
