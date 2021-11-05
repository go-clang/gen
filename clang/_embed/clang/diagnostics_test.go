package clang

import (
	"strings"
	"testing"
)

func TestDiagnostics(t *testing.T) {
	idx := NewIndex(0, 0)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("cursor.c", nil, nil, 0)
	if !tu.IsValid() {
		t.Fatal("tu is invalid")
	}
	defer tu.Dispose()

	diags := tu.Diagnostics()
	defer func() {
		for _, d := range diags {
			d.Dispose()
		}
	}()

	ok := false
	for _, d := range diags {
		if strings.Contains(d.Spelling(), "_cgo_export.h") {
			ok = true
		}
		t.Log(d)
		t.Log(d.Severity(), d.Spelling())
		t.Log(d.FormatDiagnostic(uint32(Diagnostic_DisplayCategoryName | Diagnostic_DisplaySourceLocation)))
	}
	if !ok {
		t.Fatal("not found _cgo_export.h")
	}
}
