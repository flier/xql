package xql

import "fmt"

type SearchCond BoolValueExpr

type WhereClause struct {
	Search SearchCond
}

func (w *WhereClause) String() string { return fmt.Sprintf("WHERE %s", w.Search) }
