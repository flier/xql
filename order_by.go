package xql

import (
	"fmt"
	"strings"
)

// https://jakewheat.github.io/sql-overview/sql-2016-foundation-grammar.html#order-by-clause

type OrderByClause struct {
	Sort []*SortSpec
}

func (e *OrderByClause) String() string {
	return fmt.Sprintf("ORDER BY %s", Join(e.Sort, ", "))
}

type SortSpec struct {
	Key ValueExpr
	OrderingSpec
	NullOrdering
}

func (s *SortSpec) String() string {
	var b strings.Builder

	b.WriteString(s.Key.String())

	if s.OrderingSpec != OrderingAsc {
		b.WriteByte(' ')
		b.WriteString(s.OrderingSpec.String())
	}

	if s.NullOrdering != NullsFirst {
		b.WriteByte(' ')
		b.WriteString(s.NullOrdering.String())
	}

	return b.String()
}

//go:generate stringer -type=OrderingSpec -linecomment

type OrderingSpec int

const (
	OrderingAsc  OrderingSpec = iota // ASC
	OrderingDesc                     // DESC
)

//go:generate stringer -type=NullOrdering -linecomment

type NullOrdering int

const (
	NullsFirst NullOrdering = iota // NULLS FIRST
	NullsLast                      // NULLS LAST
)
