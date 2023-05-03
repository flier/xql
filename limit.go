package xql

import (
	"fmt"
	"strings"
)

type LimitClause struct {
	RowCount SimpleValue
	Offset   *OffsetClause
}

func Limit(rowCount SimpleValue) *LimitClause {
	return &LimitClause{RowCount: rowCount}
}

func Limits(offset, rowCount SimpleValue) *LimitClause {
	return &LimitClause{RowCount: rowCount, Offset: Offset(offset)}
}

func (c *LimitClause) applySelectStmt(s *SelectStmt) { s.expr().Limit = c }
func (c *LimitClause) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "LIMIT %s", c.RowCount)

	if c.Offset != nil {
		fmt.Fprintf(&b, " %s", c.Offset)
	}

	return b.String()
}
