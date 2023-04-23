package xql_test

import (
	"fmt"

	. "github.com/flier/xql"
)

func ExampleEscapeName() {
	fmt.Println(EscapeName("foobar", '`'))
	fmt.Println(EscapeName("测试", '`'))
	fmt.Println(EscapeName("foo_bar", '`'))
	fmt.Println(EscapeName("$foobar", '`'))
	fmt.Println(EscapeName("foo+bar", '`'))
	fmt.Println(EscapeName("foo'bar", '`'))
	fmt.Println(EscapeName("foo`bar", '`'))
	fmt.Println(EscapeName("4foobar2", '`'))
	fmt.Println(EscapeName("42", '`'))
	// Output:
	// foobar
	// 测试
	// foo_bar
	// $foobar
	// `foo+bar`
	// `foo'bar`
	// `foo``bar`
	// 4foobar2
	// `42`
}
