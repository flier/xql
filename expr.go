package xql

import "fmt"

type ToExpr interface {
	expr() Expr
}

type Expr interface {
	fmt.Stringer

	ToExpr
}
