package xql

type rawExpr string

func Raw(expr string) rawExpr { return rawExpr(expr) }

func (e rawExpr) expr() Expr                           { return e }
func (e rawExpr) boolValueExpr() BoolValueExpr         { return e }
func (e rawExpr) numberValueExpr() NumberValueExpr     { return e }
func (e rawExpr) unsignedValueExpr() UnsignedValueExpr { return e }
func (e rawExpr) insertFrom() InsertFrom               { return e }
func (e rawExpr) setClause() SetClause                 { return e }
func (e rawExpr) String() string                       { return string(e) }
