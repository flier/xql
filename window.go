package xql

import (
	"fmt"
	"strings"
)

type WindowClause []*WindowDef

func (w WindowClause) String() string { return fmt.Sprintf("WINDOW %s", Join(w, ", ")) }

type WindowDef struct {
	Name string
	Spec WindowSpec
}

func (w *WindowDef) String() string { return fmt.Sprintf("%s AS ( %s )", w.Name, w.Spec) }

type WindowSpec struct {
	Name        string
	PartitionBy WindowPartition
	OrderBy     []*SortSpec
	Frame       *WindowFrameClause
}

func (w *WindowSpec) String() string {
	var b strings.Builder

	if len(w.Name) > 0 {
		b.WriteString(w.Name)
	}

	if len(w.PartitionBy) > 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(w.PartitionBy.String())
	}

	if len(w.OrderBy) > 0 {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(Join(w.OrderBy, ", "))
	}

	if w.Frame != nil {
		b.WriteByte(' ')
		b.WriteString(w.Frame.String())
	}

	return b.String()
}

type WindowPartition []*ColumnRef

func (p WindowPartition) String() string {
	return fmt.Sprintf("PARTITION BY %s", Join(p, ", "))
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
