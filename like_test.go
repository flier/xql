package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleLike() {
	fmt.Println(CreateTable("bar", Like("foo", IncludingAll, ExcludingDefaults)))
	// Output:
	// CREATE TABLE bar (
	// 	LIKE foo INCLUDING ALL EXCLUDING DEFAULTS
	// )
}
