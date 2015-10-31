package phoenix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicParsing(t *testing.T) {
	idx := NewIndex(0, 1)
	defer idx.Dispose()

	tu := idx.ParseTranslationUnit("testdata/basicparsing.c", nil, nil, 0)
	defer tu.Dispose()

	cursor := tu.TranslationUnitCursor()

	found := 0

	cursor.Visit(func(cursor, parent Cursor) ChildVisitResult {
		if cursor.IsNull() {
			return ChildVisit_Continue
		}

		switch cursor.Kind() {
		case Cursor_FunctionDecl:
			assert.Equal(t, "foo", cursor.Spelling())

			found++
		case Cursor_ParmDecl:
			assert.Equal(t, "bar", cursor.Spelling())

			found++
		}

		return ChildVisit_Recurse
	})

	assert.Equal(t, 2, found, "Did not find all nodes")
}
