package gen

// Receiver represents a generation receiver.
//
// TODO refactor https://github.com/go-clang/gen/issues/52
type Receiver struct {
	Name  string
	CName string
	Type  Type
}
