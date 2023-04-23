package xql

import (
	"fmt"
	"strings"
)

type JoinedTable struct {
	*CrossJoin
	*QualifiedJoin
	*NaturalJoin
}

func (j *JoinedTable) String() string {
	if j.CrossJoin != nil {
		return j.CrossJoin.String()
	}

	if j.QualifiedJoin != nil {
		return j.QualifiedJoin.String()
	}

	if j.NaturalJoin != nil {
		return j.NaturalJoin.String()
	}

	panic("unreachable")
}

type CrossJoin struct {
	Left  TableRef
	Right TableFactor
}

func (j *CrossJoin) String() string { return fmt.Sprintf("%s CROSS JOIN %s", j.Left, j.Right) }

//go:generate stringer -type JoinType -linecomment

type JoinType int

const (
	JoinInner JoinType = iota // INNER
	JoinLeft                  // LEFT
	JoinRight                 // RIGHT
	JoinFull                  // FULL
)

func (t JoinType) Outer() bool { return t != JoinInner }

type QualifiedJoin struct {
	Left  Either[*TableRef, *PartitionedJoinedTable]
	Type  JoinType
	Right Either[*TableRef, *PartitionedJoinedTable]
	Spec  JoinSpec
}

func (j *QualifiedJoin) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s ", j.Left)
	if j.Type.Outer() {
		fmt.Fprintf(&b, "%s ", j.Type)
	}
	fmt.Fprintf(&b, "JOIN %s %s", j.Right, j.Spec)

	return b.String()
}

type JoinSpec struct {
	On    *JoinCond
	Using *NamedColumnsJoin
}

type JoinCond struct {
	Search SearchCond
}

func (j *JoinCond) String() string { return fmt.Sprintf("ON %s", j.Search) }

type NamedColumnsJoin struct {
	Columns []ColumnName
	As      string
}

func (j *NamedColumnsJoin) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "USING (%s)", strings.Join(j.Columns, ", "))

	if j.As != "" {
		fmt.Fprintf(&b, "AS %s", j.As)
	}

	return b.String()
}

type NaturalJoin struct {
	Left  TableRef
	Type  JoinType
	Right TableFactor
}

func (j *NaturalJoin) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s NATURAL ", j.Left)
	if j.Type.Outer() {
		fmt.Fprintf(&b, "%s ", j.Type)
	}
	fmt.Fprintf(&b, "JOIN %s", j.Right)

	return b.String()
}

type PartitionedJoinedTable struct {
	Table   TableFactor
	Columns []*ColumnRef
}

func (t *PartitionedJoinedTable) String() string {
	return fmt.Sprintf("%s PARTITION BY %s", t.Table, Join(t.Columns, ", "))
}
