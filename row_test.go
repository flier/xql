package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleRowType() {
	fmt.Println(RowTy(
		Field("r_deptNo", VarChar(3)),
		Field("r_reportNo", VarChar(3)),
		Field("r_depTName", VarChar(29)),
		Field("r_mgrNo", VarChar(8)),
		Field("r_location", VarChar(128)),
	))

	// Output:
	// ROW (r_deptNo VARCHAR(3), r_reportNo VARCHAR(3), r_depTName VARCHAR(29), r_mgrNo VARCHAR(8), r_location VARCHAR(128))
}
