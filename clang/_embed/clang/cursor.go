package clang

// #include "go-clang.h"
import "C"
import (
	"sync"
	"unsafe"
)

// PlatformAvailability determine the availability of the entity that this cursor refers to on any platforms for which availability information is known.
//
// availabilitySize is the number of elements available in the availability array.
//
// alwaysDeprecated is if non-NULL, will be set to indicate whether the entity is deprecated on all platforms.
// deprecatedMessage is if non-NULL, will be set to the message text provided along with the unconditional deprecation of this entity. The client is responsible for deallocating this string.
// alwaysUnavailable is if non-NULL, will be set to indicate whether the entity is unavailable on all platforms.
// unavailableMessage is if non-NULL, will be set to the message text provided along with the unconditional unavailability of this entity. The client is responsible for deallocating this string.
// availability is if non-NULL, an array of PlatformAvailability instances that will be populated with platform availability information, up to either the number of platforms for which availability information is available (as returned by this function) or availabilitySize, whichever is smaller.
//
// Returns The number of platforms (N) for which availability information is available (which is unrelated to availabilitySize).
// Note that the client is responsible for calling Dispose to free each of the PlatformAvailability structures returned. There are min(N, availabilitySize) such structures.
func (c Cursor) PlatformAvailability(availabilitySize int) (alwaysDeprecated bool, deprecated_msg string, always_unavailable bool, unavailable_msg string, availability []PlatformAvailability) {
	var cAlwaysDeprecated C.int

	var cDeprecatedMessage cxstring
	defer cDeprecatedMessage.Dispose()

	var cAlwaysUnavailable C.int

	var cUnavailableMessage cxstring
	defer cUnavailableMessage.Dispose()

	cpAvailability := make([]C.CXPlatformAvailability, availabilitySize)

	nn := int(C.clang_getCursorPlatformAvailability(c.c, &cAlwaysDeprecated, &cDeprecatedMessage.c, &cAlwaysUnavailable, &cUnavailableMessage.c, &cpAvailability[0], C.int(len(cpAvailability))))

	if nn > availabilitySize {
		nn = availabilitySize
	}

	availability = make([]PlatformAvailability, nn)
	for i := 0; i < nn; i++ {
		availability[i] = PlatformAvailability{&cpAvailability[i]}
	}

	return cAlwaysDeprecated != 0, cDeprecatedMessage.String(), cAlwaysUnavailable != 0, cUnavailableMessage.String(), availability
}

// CursorVisitor invoked for each cursor found by a traversal.
//
// This visitor function will be invoked for each cursor found by
// clang_visitCursorChildren(). Its first argument is the cursor being
// visited, its second argument is the parent visitor for that cursor,
// and its third argument is the client data provided to
// clang_visitCursorChildren().
//
// The visitor should return one of the ChildVisitResult values
// to direct clang_visitCursorChildren().
type CursorVisitor func(cursor, parent Cursor) (status ChildVisitResult)

type funcRegistry struct {
	sync.RWMutex

	index int
	funcs map[int]*CursorVisitor
}

func (fm *funcRegistry) register(f *CursorVisitor) int {
	fm.Lock()
	defer fm.Unlock()

	fm.index++
	for fm.funcs[fm.index] != nil {
		fm.index++
	}

	fm.funcs[fm.index] = f

	return fm.index
}

func (fm *funcRegistry) lookup(index int) *CursorVisitor {
	fm.RLock()
	defer fm.RUnlock()

	return fm.funcs[index]
}

func (fm *funcRegistry) unregister(index int) {
	fm.Lock()

	delete(fm.funcs, index)

	fm.Unlock()
}

var visitors = &funcRegistry{
	funcs: map[int]*CursorVisitor{},
}

// GoClangCursorVisitor calls the cursor visitor.
//export GoClangCursorVisitor
func GoClangCursorVisitor(cursor, parent C.CXCursor, cfct unsafe.Pointer) (status ChildVisitResult) {
	i := *(*C.int)(cfct)
	f := visitors.lookup(int(i))

	return (*f)(Cursor{cursor}, Cursor{parent})
}

// Visit the children of a particular cursor.
//
// This function visits all the direct children of the given cursor,
// invoking the given visitor function with the cursors of each
// visited child. The traversal may be recursive, if the visitor returns
// ChildVisit_Recurse. The traversal may also be ended prematurely, if
// the visitor returns ChildVisit_Break.
//
// parent the cursor whose child may be visited. All kinds of
// cursors can be visited, including invalid cursors (which, by
// definition, have no children).
//
// visitor the visitor function that will be invoked for each
// child of parent.
//
// clientData pointer data supplied by the client, which will
// be passed to the visitor each time it is invoked.
//
// Returns a non-zero value if the traversal was terminated.
func (c Cursor) Visit(visitor CursorVisitor) bool {
	i := visitors.register(&visitor)
	defer visitors.unregister(i)

	// we need a pointer to the index because clang_visitChildren data parameter is a void pointer.
	ci := C.int(i)

	o := C.go_clang_visit_children(c.c, unsafe.Pointer(&ci))

	return o == C.uint(0)
}
