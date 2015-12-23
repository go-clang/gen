package expected

var TestArrayStructGo string = `
	type TestArray struct {
		EmptyStruct [10]structs
		uint 		[10]fixedSizedArray

		uint 		[]initArray
	}
`

var TestStructsGetter string = `
	func (TestArray ta) Structs() []EmptyStruct{
		var s []EmptyStruct
		gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		gos_s.Cap = int(10)
		gos_s.Len = int(10)
		gos_s.Data = uintptr(unsafe.Pointer(ta.c.structs))

		return s
	}
`

var TestFixedSizedArrayGetter string = `
	func (TestArray ta) FixedSizedArray() []uint{
		var s []uint
		gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		gos_s.Cap = int(10)
		gos_s.Len = int(10)
		gos_s.Data = uintptr(unsafe.Pointer(ta.c.fixedSizedArray))

		return s
	}
`

var TestInitArrayGetter string = `
	func (TestArray ta) InitArray() []uint{
		var s []uint
		gos_s := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		gos_s.Cap = int(3)
		gos_s.Len = int(3)
		gos_s.Data = uintptr(unsafe.Pointer(ta.c.initArray))

		return s
	}
`

var TestFunctionWithStructArrayParam = `
	func (TestArray ta) FunctionWithStructArrayParam(EmptyStruct arr[10]){
		var ca_emtysstructs [10]EmptyStruct;
		for i := range arr {
			ca_emtysstructs[i] = arr[i].c
		}

		C.functionWithStructArrayParam(ta.c, ca_emptystructs)
	}

`

var TestFunctionWithULongArrayParam = `
	func (TestArray ta) FunctionWithULongArrayParam(uint32 arr[10]){
		var ca_arr [10]uint32;
		for i := range arr {
			ca_arr[i] = arr[i].c
		}

		C.functionWithStructArrayParam(ta.c, ca_arr)
	}
`

var TestFunctionWithStructArrayParamNoSize = `
	func (TestArray ta) FunctionWithStructArrayParamNoSize(EmptyStruct arr[], int size){
		var ca_arr [10]EmptyStruct;
		for i := range arr {
			ca_arr[i] = arr[i].c
		}


		C.functionWithStructArrayParamNoSize(ta.c, ca_arr, size)
	}
`

var TestFunctionWithULongArrayParamNoSize = `
	func (TestArray ta) FunctionWithULongArrayParamNoSize(uint arr[], int size){
		var ca_arr [10]uint32;
		for i := range arr {
			ca_arr[i] = arr[i].c
		}


		C.functionWithULongArrayParamNoSize(arr, size)
	}
`
