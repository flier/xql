package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleDeleteFrom() {
	fmt.Println(DeleteFrom("products"))
	fmt.Println(DeleteFrom(Only("products")))
	fmt.Println(DeleteFrom("products").As("p"))
	fmt.Println(DeleteFrom("products").Where(Raw("price = 10")))
	fmt.Println(DeleteFrom("products").WhereCurrentOf("c_tasks"))
	// Output:
	// DELETE FROM products
	// DELETE FROM ONLY products
	// DELETE FROM products AS p
	// DELETE FROM products WHERE price = 10
	// DELETE FROM products WHERE CURRENT OF c_tasks
}
