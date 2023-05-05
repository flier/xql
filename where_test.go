package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleWhere() {
	fmt.Println(Where(Raw("t1.name = t2.name")))
	// Output:
	// WHERE t1.name = t2.name
}
