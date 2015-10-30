package clang

/**
 * \brief Tokenize the source code described by the given range into raw
 * lexical tokens.
 *
 * \param TU the translation unit whose text is being tokenized.
 *
 * \param Range the source range in which text should be tokenized. All of the
 * tokens produced by tokenization will fall within this source range,
 *
 * \param Tokens this pointer will be set to point to the array of tokens
 * that occur within the given source range. The returned pointer must be
 * freed with clang_disposeTokens() before the translation unit is destroyed.
 *
 * \param NumTokens will be set to the number of tokens in the \c *Tokens
 * array.
 *
 */
func Tokenize(tu TranslationUnit, src SourceRange) Tokens {
	tokens := Tokens{}
	tokens.tu = tu.c
	C.clang_tokenize(tu.c, src.c, &tokens.c, &tokens.n)
	return tokens
}

// an array of tokens
type Tokens struct {
	tu C.CXTranslationUnit
	c  *C.CXToken
	n  C.uint
}

/**
 * \brief Annotate the given set of tokens by providing cursors for each token
 * that can be mapped to a specific entity within the abstract syntax tree.
 *
 * This token-annotation routine is equivalent to invoking
 * clang_getCursor() for the source locations of each of the
 * tokens. The cursors provided are filtered, so that only those
 * cursors that have a direct correspondence to the token are
 * accepted. For example, given a function call \c f(x),
 * clang_getCursor() would provide the following cursors:
 *
 *   * when the cursor is over the 'f', a DeclRefExpr cursor referring to 'f'.
 *   * when the cursor is over the '(' or the ')', a CallExpr referring to 'f'.
 *   * when the cursor is over the 'x', a DeclRefExpr cursor referring to 'x'.
 *
 * Only the first and last of these cursors will occur within the
 * annotate, since the tokens "f" and "x' directly refer to a function
 * and a variable, respectively, but the parentheses are just a small
 * part of the full syntax of the function call expression, which is
 * not provided as an annotation.
 *
 * \param TU the translation unit that owns the given tokens.
 *
 * \param Tokens the set of tokens to annotate.
 *
 * \param NumTokens the number of tokens in \p Tokens.
 *
 * \param Cursors an array of \p NumTokens cursors, whose contents will be
 * replaced with the cursors corresponding to each token.
 */
func (t Tokens) Annotate() []Cursor {
	cursors := make([]Cursor, int(t.n))
	if t.n <= 0 {
		return cursors
	}
	c_cursors := make([]C.CXCursor, int(t.n))
	C.clang_annotateTokens(t.tu, t.c, t.n, &c_cursors[0])
	for i, _ := range cursors {
		cursors[i] = Cursor{C._go_clang_ocursor_at(&c_cursors[0], C.int(i))}
	}
	return cursors
}

/**
 * \brief Free the given set of tokens.
 */
func (t Tokens) Dispose() {
	C.clang_disposeTokens(t.tu, t.c, t.n)
}

// EOF
