package xql

type WhereClause struct {
	Search SearchCond
}

const kWhere = Keyword("WHERE")

func Where(x SearchCond) *WhereClause { return &WhereClause{x} }

func (w *WhereClause) Accept(v Visitor) Visitor { return v.Visit(kWhere, WS, Raw(w.Search.String())) }
func (w *WhereClause) String() string           { return XQL(w) }
