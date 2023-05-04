package xql

import (
	"fmt"
	"strings"
)

//go:generate stringer -type OffsetSuffix -linecomment

type OffsetSuffix int

const (
	SuffixRows OffsetSuffix = iota // ROWS
	SuffixRow                      // ROW
)

func (s *OffsetSuffix) applyOffsetClause(c *OffsetClause) { c.Suffix = s }

type OffsetClause struct {
	Offset SimpleValue
	Suffix *OffsetSuffix
}

type OffsetOption interface {
	applyOffsetClause(*OffsetClause)
}

func Offset(offset SimpleValue, x ...OffsetOption) *OffsetClause {
	c := &OffsetClause{Offset: offset}

	for _, opt := range x {
		opt.applyOffsetClause(c)
	}

	return c
}

func (c *OffsetClause) Limit(n SimpleValue) *LimitsClause {
	return &LimitsClause{Limit(n), c}
}

func (c *OffsetClause) limitsClause() *LimitsClause { return &LimitsClause{OffsetClause: c} }

func (c *OffsetClause) Row() *OffsetClause {
	row := SuffixRow
	c.Suffix = &row
	return c
}

func (c *OffsetClause) Rows() *OffsetClause {
	rows := SuffixRows
	c.Suffix = &rows
	return c
}

func (c *OffsetClause) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "OFFSET %s", c.Offset)

	if c.Suffix != nil {
		fmt.Fprintf(&b, " %s", c.Suffix)
	}

	return c.String()
}
