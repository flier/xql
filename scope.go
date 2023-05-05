package xql

type ScopeClause struct {
	Table TableName
}

const kScope = Keyword("SCOPE")

func (c *ScopeClause) Accept(v Visitor) Visitor { return v.Visit(kScope, WS, &c.Table) }
func (c *ScopeClause) String() string           { return XQL(c) }
