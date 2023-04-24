package xql

import "fmt"

func ExampleUpdate() {
	fmt.Println(Update("products").Set(Assign("price", 10)).Where(Raw("price = 5")))
	fmt.Println(Update("mytable").Set(Raw("a = 5, b = 3, c = 1")).Where(Raw("a > 0")))
	// Output:
	// UPDATE products SET price = 10 WHERE price = 5
	// UPDATE mytable SET a = 5, b = 3, c = 1 WHERE a > 0
}
