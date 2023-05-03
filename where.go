package xql

import "fmt"

type SearchCond BoolValueExpr

type WhereClause struct {
	Search SearchCond
}

func Where(cond SearchCond) *WhereClause             { return &WhereClause{cond} }
func (w *WhereClause) applySelectStmt(s *SelectStmt) { s.expr().Where = w }
func (w *WhereClause) String() string                { return fmt.Sprintf("WHERE %s", w.Search) }
