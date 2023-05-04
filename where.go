package xql

import "fmt"

type WhereClause struct {
	Search SearchCond
}

func Where(x SearchCond) *WhereClause { return &WhereClause{x} }
func (w *WhereClause) String() string { return fmt.Sprintf("WHERE %s", w.Search) }
