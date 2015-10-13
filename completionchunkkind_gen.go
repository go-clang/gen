package phoenix

// #include "go-clang.h"
import "C"

/**
 * \brief Describes a single piece of text within a code-completion string.
 *
 * Each "chunk" within a code-completion string (\c CXCompletionString) is
 * either a piece of text with a specific "kind" that describes how that text
 * should be interpreted by the client or is another completion string.
 */
type CompletionChunkKind int

const (
	/**
	 * \brief A code-completion string that describes "optional" text that
	 * could be a part of the template (but is not required).
	 *
	 * The Optional chunk is the only kind of chunk that has a code-completion
	 * string for its representation, which is accessible via
	 * \c clang_getCompletionChunkCompletionString(). The code-completion string
	 * describes an additional part of the template that is completely optional.
	 * For example, optional chunks can be used to describe the placeholders for
	 * arguments that match up with defaulted function parameters, e.g. given:
	 *
	 * \code
	 * void f(int x, float y = 3.14, double z = 2.71828);
	 * \endcode
	 *
	 * The code-completion string for this function would contain:
	 *   - a TypedText chunk for "f".
	 *   - a LeftParen chunk for "(".
	 *   - a Placeholder chunk for "int x"
	 *   - an Optional chunk containing the remaining defaulted arguments, e.g.,
	 *       - a Comma chunk for ","
	 *       - a Placeholder chunk for "float y"
	 *       - an Optional chunk containing the last defaulted argument:
	 *           - a Comma chunk for ","
	 *           - a Placeholder chunk for "double z"
	 *   - a RightParen chunk for ")"
	 *
	 * There are many ways to handle Optional chunks. Two simple approaches are:
	 *   - Completely ignore optional chunks, in which case the template for the
	 *     function "f" would only include the first parameter ("int x").
	 *   - Fully expand all optional chunks, in which case the template for the
	 *     function "f" would have all of the parameters.
	 */
	CompletionChunk_Optional CompletionChunkKind = C.CXCompletionChunk_Optional
	/**
	 * \brief Text that a user would be expected to type to get this
	 * code-completion result.
	 *
	 * There will be exactly one "typed text" chunk in a semantic string, which
	 * will typically provide the spelling of a keyword or the name of a
	 * declaration that could be used at the current code point. Clients are
	 * expected to filter the code-completion results based on the text in this
	 * chunk.
	 */
	CompletionChunk_TypedText CompletionChunkKind = C.CXCompletionChunk_TypedText
	/**
	 * \brief Text that should be inserted as part of a code-completion result.
	 *
	 * A "text" chunk represents text that is part of the template to be
	 * inserted into user code should this particular code-completion result
	 * be selected.
	 */
	CompletionChunk_Text CompletionChunkKind = C.CXCompletionChunk_Text
	/**
	 * \brief Placeholder text that should be replaced by the user.
	 *
	 * A "placeholder" chunk marks a place where the user should insert text
	 * into the code-completion template. For example, placeholders might mark
	 * the function parameters for a function declaration, to indicate that the
	 * user should provide arguments for each of those parameters. The actual
	 * text in a placeholder is a suggestion for the text to display before
	 * the user replaces the placeholder with real code.
	 */
	CompletionChunk_Placeholder CompletionChunkKind = C.CXCompletionChunk_Placeholder
	/**
	 * \brief Informative text that should be displayed but never inserted as
	 * part of the template.
	 *
	 * An "informative" chunk contains annotations that can be displayed to
	 * help the user decide whether a particular code-completion result is the
	 * right option, but which is not part of the actual template to be inserted
	 * by code completion.
	 */
	CompletionChunk_Informative CompletionChunkKind = C.CXCompletionChunk_Informative
	/**
	 * \brief Text that describes the current parameter when code-completion is
	 * referring to function call, message send, or template specialization.
	 *
	 * A "current parameter" chunk occurs when code-completion is providing
	 * information about a parameter corresponding to the argument at the
	 * code-completion point. For example, given a function
	 *
	 * \code
	 * int add(int x, int y);
	 * \endcode
	 *
	 * and the source code \c add(, where the code-completion point is after the
	 * "(", the code-completion string will contain a "current parameter" chunk
	 * for "int x", indicating that the current argument will initialize that
	 * parameter. After typing further, to \c add(17, (where the code-completion
	 * point is after the ","), the code-completion string will contain a
	 * "current paremeter" chunk to "int y".
	 */
	CompletionChunk_CurrentParameter CompletionChunkKind = C.CXCompletionChunk_CurrentParameter
	/**
	 * \brief A left parenthesis ('('), used to initiate a function call or
	 * signal the beginning of a function parameter list.
	 */
	CompletionChunk_LeftParen CompletionChunkKind = C.CXCompletionChunk_LeftParen
	/**
	 * \brief A right parenthesis (')'), used to finish a function call or
	 * signal the end of a function parameter list.
	 */
	CompletionChunk_RightParen CompletionChunkKind = C.CXCompletionChunk_RightParen
	/**
	 * \brief A left bracket ('[').
	 */
	CompletionChunk_LeftBracket CompletionChunkKind = C.CXCompletionChunk_LeftBracket
	/**
	 * \brief A right bracket (']').
	 */
	CompletionChunk_RightBracket CompletionChunkKind = C.CXCompletionChunk_RightBracket
	/**
	 * \brief A left brace ('{').
	 */
	CompletionChunk_LeftBrace CompletionChunkKind = C.CXCompletionChunk_LeftBrace
	/**
	 * \brief A right brace ('}').
	 */
	CompletionChunk_RightBrace CompletionChunkKind = C.CXCompletionChunk_RightBrace
	/**
	 * \brief A left angle bracket ('<').
	 */
	CompletionChunk_LeftAngle CompletionChunkKind = C.CXCompletionChunk_LeftAngle
	/**
	 * \brief A right angle bracket ('>').
	 */
	CompletionChunk_RightAngle CompletionChunkKind = C.CXCompletionChunk_RightAngle
	/**
	 * \brief A comma separator (',').
	 */
	CompletionChunk_Comma CompletionChunkKind = C.CXCompletionChunk_Comma
	/**
	 * \brief Text that specifies the result type of a given result.
	 *
	 * This special kind of informative chunk is not meant to be inserted into
	 * the text buffer. Rather, it is meant to illustrate the type that an
	 * expression using the given completion string would have.
	 */
	CompletionChunk_ResultType CompletionChunkKind = C.CXCompletionChunk_ResultType
	/**
	 * \brief A colon (':').
	 */
	CompletionChunk_Colon CompletionChunkKind = C.CXCompletionChunk_Colon
	/**
	 * \brief A semicolon (';').
	 */
	CompletionChunk_SemiColon CompletionChunkKind = C.CXCompletionChunk_SemiColon
	/**
	 * \brief An '=' sign.
	 */
	CompletionChunk_Equal CompletionChunkKind = C.CXCompletionChunk_Equal
	/**
	 * Horizontal space (' ').
	 */
	CompletionChunk_HorizontalSpace CompletionChunkKind = C.CXCompletionChunk_HorizontalSpace
	/**
	 * Vertical space ('\n'), after which it is generally a good idea to
	 * perform indentation.
	 */
	CompletionChunk_VerticalSpace CompletionChunkKind = C.CXCompletionChunk_VerticalSpace
)
