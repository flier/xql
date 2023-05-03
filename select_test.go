package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleSelect() {
	fmt.Println(SelectAllFrom("table1"))
	fmt.Println(Select(All, From(QName("table1"))))
	fmt.Println(Select(Distinct, All, From(QName("table1"))))
	fmt.Println(Select(Raw("3 + 4").As("sum")))

	var random = Func("random")
	fmt.Println(Select(random()))

	fmt.Println(Select(Column("a"), Column("b"), Column("c")).From(QName("table1")))
	fmt.Println(Select(Column("a").As("value"), Raw("b + c").As("sum")).From(QName("table1")))

	var tbl1 = QName("tbl1")
	var tbl2 = QName("tbl2")
	fmt.Println(Select(tbl1.Join("a"), tbl2.Join("a"), tbl2.Join("b")).From(tbl1, tbl2))
	// Output:
	// SELECT * FROM table1
	// SELECT * FROM table1
	// SELECT DISTINCT * FROM table1
	// SELECT 3 + 4 AS sum
	// SELECT random()
	// SELECT a, b, c FROM table1
	// SELECT a AS value, b + c AS sum FROM table1
	// SELECT tbl1.a, tbl2.a, tbl2.b FROM tbl1, tbl2
}
