package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleArrayType() {
	fmt.Println(ArrayOf(Integer, Caps(100)))

	// Output:
	// INTEGER ARRAY[100]
}
