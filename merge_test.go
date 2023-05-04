package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleMergeInto() {
	fmt.Println(MergeInto("customer_account").As("ca").
		Using(QName("recent_transactions").As("t")).
		On(Raw("t.customer_id = ca.customer_id")))
	// Output:
	// MERGE INTO customer_account AS ca USING recent_transactions AS t ON t.customer_id = ca.customer_id
}
