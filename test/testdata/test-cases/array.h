#include <stdint.h>


/*
	TestCases:

	function parameter
	return argument -> cannot be done http://stackoverflow.com/questions/11656532/returning-an-array-using-c
	struct member

	fixed size arrays
	pre-initialized arrays

	each with struct or primitive

	missing?
*/

typedef struct {
} EmptyStruct;

typedef struct {
	EmptyStruct structs[10];
	unsigned long fixedSizedArray[10];

	// unsigned long initArray[] = {10, 10, 10}; TODO: ask zimmski
} TestArray;

void functionWithStructArrayParam(TestArray ta, EmptyStruct earr[10]);

void functionWithULongArrayParam(TestArray ta, unsigned long larr[10]);

void functionWithStructArrayParamNoSize(TestArray ta, EmptyStruct earr[], int size);

void functionWithULongArrayParamNoSize(TestArray ta, unsigned long larr[], int size);
