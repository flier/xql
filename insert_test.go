package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleInsertInto_default() {
	fmt.Println(InsertInto("products", DefaultValues))
	fmt.Println(InsertInto("products", Columns("product_no", "name", "price").Values(1, "Cheese", Default)))
	// Output:
	// INSERT INTO products DEFAULT VALUES
	// INSERT INTO products (product_no, name, price) VALUES (1, "Cheese", DEFAULT)
}

func ExampleInsertInto_values() {
	fmt.Println(InsertInto("products", Values(1, "Cheese", 9.99)))
	fmt.Println(InsertInto("products", Columns("product_no", "name", "price").Values(1, "Cheese", 9.99)))
	// Output:
	// INSERT INTO products VALUES (1, "Cheese", 9.99)
	// INSERT INTO products (product_no, name, price) VALUES (1, "Cheese", 9.99)
}

func ExampleInsertInto_rows() {
	fmt.Println(InsertInto("products", Columns("product_no", "name", "price").Values(Rows{
		{1, "Cheese", 9.99},
		{2, "Bread", 1.99},
		{3, "Milk", 2.99},
	})))
	fmt.Println(InsertInto("products", Columns("product_no", "name", "price").Values(
		Row{1, "Cheese", 9.99},
		Row{2, "Bread", 1.99},
		Row{3, "Milk", 2.99},
	)))
	// Output:
	// INSERT INTO products (product_no, name, price) VALUES
	// 	ROW(1, "Cheese", 9.99),
	// 	ROW(2, "Bread", 1.99),
	// 	ROW(3, "Milk", 2.99)
	// INSERT INTO products (product_no, name, price) VALUES
	// 	ROW(1, "Cheese", 9.99),
	// 	ROW(2, "Bread", 1.99),
	// 	ROW(3, "Milk", 2.99)
}
