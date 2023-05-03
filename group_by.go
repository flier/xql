package xql

import (
	"fmt"
	"strings"
)

type GroupByClause struct {
	Set   SetQuantifier
	Elems []GroupingElement
}

func GroupBy(x ...GroupingElement) *GroupByClause      { return &GroupByClause{Elems: x} }
func (g *GroupByClause) Distinct() *GroupByClause      { g.Set = SetDistinct; return g }
func (g *GroupByClause) applySelectStmt(s *SelectStmt) { s.expr().GroupBy = g }

func (g *GroupByClause) String() string {
	var b strings.Builder

	b.WriteString("GROUP BY ")
	if g.Set != SetAll {
		fmt.Fprintf(&b, "%s ", g.Set)
	}
	b.WriteString(Join(g.Elems, ", "))

	return b.String()
}

type ToGroupingElement interface {
	groupingElement() GroupingElement
}

type GroupingElement interface {
	fmt.Stringer

	ToGroupingElement
}

var (
	_ GroupingElement = OrdinaryGroupingSet(nil)
	_ GroupingElement = RollUpClause(nil)
	_ GroupingElement = CubeClause(nil)
	_ GroupingElement = &GroupingSetsSpec{}
)

type OrdinaryGroupingSet []*GroupingColumnRef

func (s OrdinaryGroupingSet) groupingElement() GroupingElement { return s }
func (s OrdinaryGroupingSet) groupingSet() GroupingSet         { return s }
func (s OrdinaryGroupingSet) String() string {
	if len(s) == 1 {
		return s[0].String()
	}

	return fmt.Sprintf("(%s)", Join(s, ", "))
}

type GroupingColumnRef struct {
	Column  ColumnRef
	Collate *CollateClause
}

func (r *GroupingColumnRef) String() string {
	var b strings.Builder

	b.WriteString(r.Column)

	if r.Collate != nil {
		fmt.Fprintf(&b, " %s", r.Collate)
	}

	return b.String()
}

type RollUpClause []*OrdinaryGroupingSet

func (r RollUpClause) groupingElement() GroupingElement { return r }
func (r RollUpClause) groupingSet() GroupingSet         { return r }
func (r RollUpClause) String() string                   { return fmt.Sprintf("ROLLUP (%s)", Join(r, ", ")) }

type CubeClause []*OrdinaryGroupingSet

func (r CubeClause) groupingElement() GroupingElement { return r }
func (r CubeClause) groupingSet() GroupingSet         { return r }
func (r CubeClause) String() string                   { return fmt.Sprintf("CUBE (%s)", Join(r, ", ")) }

type GroupingSetsSpec []GroupingSet

func (s GroupingSetsSpec) groupingElement() GroupingElement { return s }
func (s GroupingSetsSpec) groupingSet() GroupingSet         { return s }
func (s GroupingSetsSpec) String() string                   { return fmt.Sprintf("GROUPING SETS (%s)", Join(s, ", ")) }

type ToGroupingSet interface {
	groupingSet() GroupingSet
}

type GroupingSet interface {
	fmt.Stringer

	ToGroupingSet
}

var (
	_ GroupingSet = OrdinaryGroupingSet(nil)
	_ GroupingSet = RollUpClause(nil)
	_ GroupingSet = CubeClause(nil)
	_ GroupingSet = GroupingSetsSpec(nil)
)
