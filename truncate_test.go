package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleTruncateTable() {
	fmt.Println(TruncateTable("bigtable"))
	fmt.Println(TruncateTable("bigtable", "fattable"))
	fmt.Println(TruncateTable("bigtable", "fattable").RestartIdentity())
	fmt.Println(TruncateTable("othertable").Cascade())
	// Output:
	// TRUNCATE TABLE bigtable
	// TRUNCATE TABLE bigtable, fattable
	// TRUNCATE TABLE bigtable, fattable RESTART IDENTITY
	// TRUNCATE TABLE othertable
}
