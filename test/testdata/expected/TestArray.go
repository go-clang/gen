package expected

var TestArrayStructGo string = `
type TestArray struct {
	c C.TestArray
}`

var TestStructsGetter string = `
func (ta TestArray) Structs() []EmptyStruct {
	var s []EmptyStruct
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(10)
	gos_s.Len = int(10)
	gos_s.Data = uintptr(unsafe.Pointer(ta.c.structs))

	return s
}`

var TestFixedSizedArrayGetter string = `
func (ta TestArray) FixedSizedArray() []uint32 {
	var s []uint32
	gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	gos_s.Cap = int(10)
	gos_s.Len = int(10)
	gos_s.Data = uintptr(unsafe.Pointer(ta.c.fixedSizedArray))

	return s
}`

/*var TestInitArrayGetter string = `
	func (TestArray ta) InitArray() []uint{
		var s []uint
		gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		gos_s.Cap = int(3)
		gos_s.Len = int(3)
		gos_s.Data = uintptr(unsafe.Pointer(ta.c.initArray))

		return s
	}
`*/

var TestFunctionWithStructArrayParam = `
func (ta TestArray) FunctionWithStructArrayParam(earr [10]EmptyStruct) {
	ca_earr := make([]C.EmptyStruct, len(earr))
	var cp_earr *C.EmptyStruct
	if len(earr) > 0 {
		cp_earr = &ca_earr[0]
	}
	for i := range earr {
		ca_earr[i] = earr[i].c
	}

	C.functionWithStructArrayParam(ta.c, cp_earr)
}`

var TestFunctionWithULongArrayParam = `
func (ta TestArray) FunctionWithULongArrayParam(larr [10]uint32) {
	ca_larr := make([]C.ulong, len(larr))
	var cp_larr *C.ulong
	if len(larr) > 0 {
		cp_larr = &ca_larr[0]
	}
	for i := range larr {
		ca_larr[i] = larr[i].c
	}

	C.functionWithULongArrayParam(ta.c, cp_larr)
}`

var TestFunctionWithStructArrayParamNoSize = `
func (ta TestArray) FunctionWithStructArrayParamNoSize(earr []EmptyStruct, size int16) {
	ca_earr := make([]C.EmptyStruct, len(earr))
	var cp_earr *C.EmptyStruct
	if len(earr) > 0 {
		cp_earr = &ca_earr[0]
	}
	for i := range earr {
		ca_earr[i] = earr[i].c
	}

	C.functionWithStructArrayParamNoSize(ta.c, cp_earr, C.int(size))
}`

var TestFunctionWithULongArrayParamNoSize = `
func (ta TestArray) FunctionWithULongArrayParamNoSize(larr []uint32, size int16) {
	ca_larr := make([]C.ulong, len(larr))
	var cp_larr *C.ulong
	if len(larr) > 0 {
		cp_larr = &ca_larr[0]
	}
	for i := range larr {
		ca_larr[i] = larr[i].c
	}

	C.functionWithULongArrayParamNoSize(ta.c, cp_larr, C.int(size))
}`
