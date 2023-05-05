package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleWith() {
	fmt.Println(With)
	fmt.Println(WithRecursive)
	// Output:
	// WITH
	// WITH RECURSIVE
}
