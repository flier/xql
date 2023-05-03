package xql

import (
	"fmt"
	"strings"
)

type WindowClause []*WindowDef

func Window(x ...*WindowDef) WindowClause            { return WindowClause(x) }
func (w WindowClause) applySelectStmt(s *SelectStmt) { s.expr().Window = w }
func (w WindowClause) String() string                { return fmt.Sprintf("WINDOW %s", Join(w, ", ")) }

type ToWindowDef interface {
	windowDef() *WindowDef
}

type WindowDef struct {
	Name string
	Spec WindowSpec
}

func (w *WindowDef) windowDef() *WindowDef { return w }
func (w *WindowDef) String() string        { return fmt.Sprintf("%s AS %s", w.Name, w.Spec) }

type WindowName = string
type WindowSpec struct {
	Name        WindowName
	PartitionBy WindowPartitionClause
	OrderBy     []*SortSpec
	Frame       *WindowFrameClause
}

func (w *WindowSpec) String() string {
	var details []string

	if len(w.Name) > 0 {
		details = append(details, w.Name)
	}

	if len(w.PartitionBy) > 0 {
		details = append(details, w.PartitionBy.String())
	}

	if len(w.OrderBy) > 0 {
		details = append(details, Join(w.OrderBy, ", "))
	}

	if w.Frame != nil {
		details = append(details, w.Frame.String())
	}

	return fmt.Sprintf("(%s)", strings.Join(details, ", "))
}

type (
	WindowPartitionClause        WindowPartitionColumnRefList
	WindowPartitionColumnRefList []WindowPartitionColumnRef
	WindowPartitionColumnRef     = ColumnRef
)

func (p WindowPartitionClause) String() string {
	return fmt.Sprintf("PARTITION BY %s", strings.Join(p, ", "))
}

type WindowOrder []*SortSpec

func (o WindowOrder) String() string { return fmt.Sprintf("ORDER BY %s", Join(o, ", ")) }

//go:generate stringer -type WindowFrameUnits -linecomment

type WindowFrameUnits int

const (
	UnitRows   WindowFrameUnits = iota // ROWS
	UnitRange                          // RANGE
	UnitGroups                         // GROUPS
)

type WindowFrameClause struct {
	Units     WindowFrameUnits
	Extent    WindowFrameExtent
	Exclusion WindowFrameExclusion
}

func (f *WindowFrameClause) String() string {
	return fmt.Sprintf("%s", f.Units)
}

type WindowFrameExtent struct {
	Start   *WindowFrameStart
	Between *WindowFrameBetween
}

type WindowFrameStart struct {
	UnboundedPreceding bool
	Preceding          UnsignedValueExpr
	CurrentRow         bool
}

func (s *WindowFrameStart) String() string {
	if s.UnboundedPreceding {
		return "UNBOUNDED PRECEDING"
	}

	if s.Preceding != nil {
		return fmt.Sprintf("%s PRECEDING", s.Preceding)
	}

	if s.CurrentRow {
		return "CURRENT ROW"
	}

	return ""
}

type WindowFrameBetween struct {
	Lower WindowFrameBound
	Upper WindowFrameBound
}

func (b *WindowFrameBetween) String() string {
	return fmt.Sprintf("BETWEEN %s AND %s", &b.Lower, &b.Upper)
}

type WindowFrameBound struct {
	Start              *WindowFrameStart
	UnboundedFollowing bool
	Following          UnsignedValueExpr
}

func (b *WindowFrameBound) String() string {
	if b.Start == nil {
		return b.Start.String()
	}

	if b.UnboundedFollowing {
		return "UNBOUNDED FOLLOWING"
	}

	if b.Following != nil {
		return fmt.Sprintf("%d FOLLOWING", b.Following)
	}

	return ""
}

//go:generate stringer -type WindowFrameExclusion -linecomment

type WindowFrameExclusion int

const (
	ExcludeCurrentRow WindowFrameExclusion = iota // EXCLUDE CURRENT ROW
	ExcludeGroup                                  // EXCLUDE GROUP
	ExcludeTies                                   // EXCLUDE TIES
	ExcludeNoOthers                               // EXCLUDE NO OTHERS
)
