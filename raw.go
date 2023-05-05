package xql

type (
	Raw string
	R   = Raw
)

func (r Raw) expr() Expr                           { return r }
func (r Raw) boolValueExpr() BoolValueExpr         { return r }
func (r Raw) numberValueExpr() NumberValueExpr     { return r }
func (r Raw) unsignedValueExpr() UnsignedValueExpr { return r }
func (r Raw) insertFrom() InsertFrom               { return r }
func (r Raw) setClause() SetClause                 { return r }
func (r Raw) Accept(v Visitor) Visitor             { return v.Raw(string(r)) }
func (r Raw) String() string                       { return string(r) }
