package xql

type rawExpr string

func Raw(expr string) rawExpr { return rawExpr(expr) }

func (e rawExpr) Expr() Expr     { return e }
func (e rawExpr) String() string { return string(e) }
