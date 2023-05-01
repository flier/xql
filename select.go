package xql

import (
	"fmt"
	"strings"
)

type SelectStmt struct {
	Set    SetQuantifier
	Select SelectList
	Into   *TargetSpec
	Expr   TableExpr
}

func (s *SelectStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "SELECT")
	if s.Set != SetAll {
		fmt.Fprintf(&b, " %s", s.Set)
	}
	fmt.Fprintf(&b, " %s", s.Select)
	if s.Into != nil {
		fmt.Fprintf(&b, " INTO %s", s.Into)
	}
	fmt.Fprintf(&b, " %s", s.Expr)

	return b.String()
}

type SelectList interface {
	fmt.Stringer

	selectList() SelectList
}

var (
	_ SelectList = &All{}
)

type All struct{}

func (a *All) selectList() SelectList { return a }
func (a *All) String() string         { return "*" }

type TargetSpec struct {
	Name *TableName
	Vars []VarName
}

type VarName = string

func (s *TargetSpec) String() string {
	if s.Name != nil {
		return fmt.Sprintf("TABLE %s", s.Name)
	}

	if len(s.Vars) > 0 {
		return strings.Join(s.Vars, ", ")
	}

	panic("unreachable")
}

type TableExpr struct {
	From    FromClause
	Where   *WhereClause
	GroupBy *GroupByClause
	Having  *HavingClause
	Window  WindowClause
}

func (e *TableExpr) String() string {
	var b strings.Builder

	b.WriteString(e.From.String())
	if e.Where != nil {
		fmt.Fprintf(&b, " %s", e.Where)
	}
	if e.GroupBy != nil {
		fmt.Fprintf(&b, " %s", e.GroupBy)
	}
	if e.Having != nil {
		fmt.Fprintf(&b, " %s", e.Having)
	}
	if e.Window != nil {
		fmt.Fprintf(&b, " %s", e.Window)
	}

	return b.String()
}
