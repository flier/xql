package xql

import "fmt"

type FromClause TableRefList

func From(x ...TableRef) FromClause {
	return FromClause(x)
}

func (c FromClause) applySelectStmt(s *SelectStmt) { s.expr().From = c }
func (c FromClause) String() string                { return fmt.Sprintf("FROM %s", TableRefList(c)) }
