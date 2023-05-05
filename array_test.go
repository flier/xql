package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleArrayType() {
	IntArray := ArrayOf(Integer)
	IntArrayOfArray := ArrayOf(IntArray)

	fmt.Println(IntArray)
	fmt.Println(IntArray(100))
	fmt.Println(IntArrayOfArray)

	// Output:
	// INTEGER ARRAY[]
	// INTEGER ARRAY[100]
	// INTEGER ARRAY[] ARRAY[]
}
