package xql

import (
	"fmt"
	"strings"
)

//go:generate stringer -type SampleMethod -linecomment

type SampleMethod int

const (
	SampleBernoulli SampleMethod = iota // BERNOULLI
	SampleSystem                        // SYSTEM
)

type SampleClause struct {
	Method     SampleMethod
	Percent    NumberValue
	Repeatable *Repeatable
}

func (s *SampleClause) String() string {
	var b strings.Builder

	b.WriteString("TABLESAMPLE (")
	b.WriteString(s.Percent.String())
	b.WriteString(")")

	if s.Repeatable != nil {
		b.WriteByte(' ')
		b.WriteString(s.Repeatable.String())
	}

	return b.String()
}

type Repeatable struct {
	Repeat NumberValue
}

func (r *Repeatable) String() string { return fmt.Sprintf("REPEATABLE ( %s )", r.Repeat) }
