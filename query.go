package xql

import (
	"fmt"
	"strings"
)

type QueryExpr struct {
	With    *WithClause
	Body    QueryExprBody
	OrderBy *OrderByClause
	Limit   *LimitClause
	Offset  *OffsetClause
	Fetch   *FetchClause
}

var _ Expr = &QueryExpr{}

func (q *QueryExpr) String() string {
	var b strings.Builder

	if q.With != nil {
		fmt.Fprintf(&b, "%s ", q.With)
	}

	b.WriteString(q.Body.String())

	if q.OrderBy != nil {
		fmt.Fprintf(&b, " %s", q.OrderBy)
	}

	if q.Limit != nil {
		fmt.Fprintf(&b, " %s", q.Limit)
	}

	if q.Offset != nil {
		fmt.Fprintf(&b, " %s", q.Offset)
	}

	if q.Fetch != nil {
		fmt.Fprintf(&b, " %s", q.Fetch)
	}

	return b.String()
}

type QueryExprBody interface {
	fmt.Stringer
}

type QueryTerm interface {
	QueryExprBody
}

type QuerySet struct {
	Left  QueryExprBody
	Op    SetOperation
	Set   SetQuantifier
	Right QueryTerm
}

func (s *QuerySet) String() string {
	return fmt.Sprintf("%s %s %s %s", s.Left, s.Op, s.Set, s.Right)
}

type QueryPrimary interface{}
