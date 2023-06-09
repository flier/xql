package xql

import "fmt"

type Assignment interface {
	fmt.Stringer

	SetClause
}

type AsSetTarget interface {
	~string | *MutatedSetClause | []SetTarget
}

func Assign[T AsSetTarget](target T, source any) Assignment {
	switch t := any(target).(type) {
	case string:
		return newColumnAssignment(t, newTypedRowValueExpr(source))
	case *MutatedSetClause:
		return newColumnAssignment(t, newTypedRowValueExpr(source))
	case []SetTarget:
		return newMultiColumnAssignment(t, newTypedRowValueExpr(source))
	default:
		panic("unreachable")
	}
}

type SetClauseList []SetClause

func (l SetClauseList) String() string { return Join(l, ", ") }

type ToSetClause interface {
	setClause() SetClause
}

type SetClause interface {
	fmt.Stringer

	ToSetClause
}

var (
	_ SetClause = Raw("")
	_ SetClause = &ColumnAssignment{}
	_ SetClause = &MultiColumnAssignment{}
)

type ColumnAssignment struct {
	Target SetTarget
	Source UpdateSource
}

func newColumnAssignment[T ~string | *MutatedSetClause](target T, source UpdateSource) *ColumnAssignment {
	switch t := any(target).(type) {
	case string:
		return &ColumnAssignment{Target: ObjectColumn(t), Source: source}
	case *MutatedSetClause:
		return &ColumnAssignment{Target: t, Source: source}
	default:
		panic("unreachable")
	}
}

func (c *ColumnAssignment) setClause() SetClause { return c }
func (c *ColumnAssignment) String() string       { return fmt.Sprintf("%s = %s", c.Target, c.Source) }

type AssignedRow = TypedRowValueExpr

type MultiColumnAssignment struct {
	Targets []SetTarget
	Source  AssignedRow
}

func newMultiColumnAssignment(target []SetTarget, source AssignedRow) *MultiColumnAssignment {
	return &MultiColumnAssignment{Targets: target, Source: source}
}

func (c *MultiColumnAssignment) setClause() SetClause { return c }

func (c *MultiColumnAssignment) String() string {
	return fmt.Sprintf("(%s) = %s", Join(c.Targets, ", "), c.Source)
}

type ToSetTarget interface {
	setTarget() SetTarget
}

type SetTarget interface {
	fmt.Stringer

	ToSetTarget
}

var (
	_ SetTarget = ObjectColumn("")
	_ SetTarget = &MutatedSetClause{}
)

type UpdateSource = ValueExpr

type ToUpdateTarget interface {
	updateTarget() UpdateTarget
}

type UpdateTarget interface {
	fmt.Stringer

	SetTarget

	ToUpdateTarget
}

var (
	_ UpdateTarget = ObjectColumn("")
)

type MutatedSetClause struct {
	Target MutatedTarget
	Method MethodName
}

func (c *MutatedSetClause) setTarget() SetTarget         { return c }
func (c *MutatedSetClause) mutatedTarget() MutatedTarget { return c }
func (c *MutatedSetClause) String() string {
	return fmt.Sprintf("%s.%s", c.Target, c.Method)
}

type ToMutatedTarget interface {
	mutatedTarget() MutatedTarget
}

type MutatedTarget interface {
	fmt.Stringer

	ToMutatedTarget
}

var (
	_ MutatedTarget = ObjectColumn("")
	_ MutatedTarget = &MutatedSetClause{}
)

type ObjectColumn ColumnName

func (c ObjectColumn) mutatedTarget() MutatedTarget { return c }
func (c ObjectColumn) setTarget() SetTarget         { return c }
func (c ObjectColumn) updateTarget() UpdateTarget   { return c }
func (c ObjectColumn) String() string               { return string(c) }
