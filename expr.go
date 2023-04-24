package xql

import "fmt"

type Expr interface {
	fmt.Stringer

	expr() Expr
}
