package xql

import (
	"fmt"
	"strings"
)

type GroupByClause struct {
	Set   SetQuantifier
	Elems []*GroupingElement
}

func (g *GroupByClause) String() string {
	var b strings.Builder

	b.WriteString("GROUP BY ")
	if g.Set != SetAll {
		fmt.Fprintf(&b, "%s ", g.Set)
	}
	b.WriteString(Join(g.Elems, ", "))

	return b.String()
}

type GroupingElement struct {
	*OrdinaryGroupingSet
	*Rollup
	*Cube
	*GroupingSetsSpec
}

func (e *GroupingElement) String() string {
	if e.OrdinaryGroupingSet != nil {
		return e.OrdinaryGroupingSet.String()
	}

	if e.Rollup != nil {
		return e.Rollup.String()
	}

	if e.Cube != nil {
		return e.Cube.String()
	}

	if e.GroupingSetsSpec != nil {
		return e.GroupingSetsSpec.String()
	}

	return "()"
}

type OrdinaryGroupingSet struct {
	Columns []*GroupingColumnRef
}

func (s *OrdinaryGroupingSet) String() string {
	if len(s.Columns) == 1 {
		return s.Columns[0].String()
	}

	return fmt.Sprintf("(%s)", Join(s.Columns, ", "))
}

type GroupingColumnRef struct {
	Column  ColumnRef
	Collate *CollateClause
}

func (r *GroupingColumnRef) String() string {
	var b strings.Builder

	b.WriteString(r.Column.String())

	if r.Collate != nil {
		fmt.Fprintf(&b, " %s", r.Collate)
	}

	return b.String()
}

type Rollup struct {
	Sets []*OrdinaryGroupingSet
}

func (r Rollup) String() string {
	return fmt.Sprintf("ROLLUP (%s)", Join(r.Sets, ", "))
}

type Cube struct {
	Sets []*OrdinaryGroupingSet
}

func (r Cube) String() string {
	return fmt.Sprintf("CUBE (%s)", Join(r.Sets, ", "))
}

type GroupingSetsSpec struct {
	Sets []*GroupingSet
}

func (s *GroupingSetsSpec) String() string {
	return fmt.Sprintf("GROUPING SETS (%s)", Join(s.Sets, ", "))
}

type GroupingSet struct {
	*OrdinaryGroupingSet
	*Rollup
	*Cube
	*GroupingSetsSpec
}

func (s *GroupingSet) String() string {
	if s.OrdinaryGroupingSet != nil {
		return s.OrdinaryGroupingSet.String()
	}

	if s.Rollup != nil {
		return s.Rollup.String()
	}

	if s.Cube != nil {
		return s.Cube.String()
	}

	if s.GroupingSetsSpec != nil {
		return s.GroupingSetsSpec.String()
	}

	return "()"
}
