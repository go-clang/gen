package gen_test

import (
	"testing"

	"github.com/go-clang/gen"
)

func TestCleanDoxygenComment(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		comment string
		want    string
	}{
		"empty": {
			comment: ``,
			want:    ``,
		},
		"Create a CXVirtualFileOverlay ...": {
			comment: `/**
 * \brief Create a \c CXVirtualFileOverlay object.
 * Must be disposed with \c clang_VirtualFileOverlay_dispose().
 *
 * \param options is reserved, always pass 0.
 */`,
			want: `// Create a CXVirtualFileOverlay object.
// Must be disposed with clang_VirtualFileOverlay_dispose().
//
// Parameter options is reserved, always pass 0.`,
		},
		"Return the timestamp ...": {
			comment: `/**
 * \brief Return the timestamp for use with Clang's
 * \c -fbuild-session-timestamp= option.
 */`,
			want: `// Return the timestamp for use with Clang's -fbuild-session-timestamp= option.`,
		},
		"Object encapsulating information ...": {
			comment: `/**
 * \brief Object encapsulating information about overlaying virtual
 * file/directories over the real file system.
 */`,
			want: `// Object encapsulating information about overlaying virtual file/directories over the real file system.`,
		},
		"Map an absolute ...": {
			comment: `/**
 * \brief Map an absolute virtual file path to an absolute real one.
 * The virtual path must be canonicalized (not contain "."/"..").
 * \returns 0 for success, non-zero to indicate an error.
 */`,
			want: `// Map an absolute virtual file path to an absolute real one. The virtual path must be canonicalized (not contain "."/".."). Returns 0 for success, non-zero to indicate an error.`,
		},
		"Write out the ...": {
			comment: `/**
 * \brief Write out the \c CXVirtualFileOverlay object to a char buffer.
 *
 * \param options is reserved, always pass 0.
 * \param out_buffer_ptr pointer to receive the buffer pointer, which should be
 * disposed using \c clang_free().
 * \param out_buffer_size pointer to receive the buffer size.
 * \returns 0 for success, non-zero to indicate an error.
 */`,
			want: `// Write out the CXVirtualFileOverlay object to a char buffer.
//
// Parameter options is reserved, always pass 0.
// Parameter out_buffer_ptr pointer to receive the buffer pointer, which should be
// disposed using clang_free().
// Parameter out_buffer_size pointer to receive the buffer size.
// Returns 0 for success, non-zero to indicate an error.`,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := gen.CleanDoxygenComment("", tt.comment); got != tt.want {
				t.Fatalf("CleanDoxygenComment(\n%v\n) = \n%v\nwant \n%v\n", tt.comment, got, tt.want)
			}
		})
	}
}
