package phoenix

/* TODO
import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	idx := NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("cursor.c", nil, nil, 0)
	assert.True(t, tu.IsValid())
	defer tu.Dispose()

	res := tu.CodeCompleteAt("cursor.c", 10, 16, nil, 0)
	assert.True(t, res.IsValid())
	defer res.Dispose()

	if n := len(res.Results()); n < 10 {
		t.Errorf("Expected more results than %d", n)
	}

	t.Logf("%+v", res)
	for _, r := range res.Results() {
		t.Logf("%+v", r)
		for _, c := range r.CompletionString.Chunks() {
			t.Logf("\t%+v", c)
		}
	}

	diags := res.Diagnostics()
	defer diags.Dispose()

	ok := false
	for _, d := range diags {
		if strings.Contains(d.Spelling(), "_cgo_export.h") {
			ok = true
		}
		t.Log(d.Severity(), d.Spelling())
	}
	assert.True(t, ok)
}
*/
