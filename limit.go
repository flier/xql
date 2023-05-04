package xql

import (
	"bytes"
	"fmt"
	"strings"
)

type LimitClause struct {
	RowCount SimpleValue
	Percent  bool
	WithTies bool
}

func Limit(rowCount SimpleValue) *LimitClause {
	return &LimitClause{RowCount: rowCount}
}

func (c *LimitClause) Offset(offset SimpleValue) *LimitsClause {
	return &LimitsClause{c, Offset(offset)}
}

func (c *LimitClause) limitsClause() *LimitsClause { return &LimitsClause{LimitClause: c} }

func (c *LimitClause) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "LIMIT %s", c.RowCount)

	if c.Percent {
		b.WriteString(" PERCENT")
	}

	if c.WithTies {
		b.WriteString(" WITH TIES")
	}

	return b.String()
}

type ToLimitsClause interface {
	limitsClause() *LimitsClause
}

var (
	_ ToLimitsClause = &LimitsClause{}
	_ ToLimitsClause = &LimitClause{}
	_ ToLimitsClause = &OffsetClause{}
)

type LimitsClause struct {
	*LimitClause
	*OffsetClause
}

func Limits(offset, rowCount SimpleValue) *LimitsClause {
	return &LimitsClause{
		&LimitClause{RowCount: rowCount},
		&OffsetClause{Offset: offset},
	}
}

func (c *LimitsClause) limitsClause() *LimitsClause { return c }
func (c *LimitsClause) String() string {
	var b bytes.Buffer

	if c.LimitClause != nil {
		b.WriteString(c.LimitClause.String())
	}

	if c.OffsetClause != nil {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}

		b.WriteString(c.OffsetClause.String())
	}

	return b.String()
}
