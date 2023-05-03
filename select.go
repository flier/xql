package xql

import (
	"fmt"
	"strings"
)

type SelectStmt struct {
	Quantifier SetQuantifier
	Select     SelectList
	Into       *TargetSpec
	Expr       *TableExpr
}

type SelectOption interface {
	applySelectStmt(*SelectStmt)
}

var (
	_ SelectOption = SetQuantifier(0)
	_ SelectOption = Raw("")
	_ SelectOption = &LocalOrSchemaQualifiedName{}
	_ SelectOption = &SchemaQualifiedName{}
	_ SelectOption = &ColumnDef{}
	_ SelectOption = FromClause(nil)
	_ SelectOption = &WhereClause{}
	_ SelectOption = &GroupByClause{}
	_ SelectOption = &HavingClause{}
	_ SelectOption = &WindowClause{}
	_ SelectOption = &OrderByClause{}
	_ SelectOption = &LimitClause{}
)

func SelectAllFrom[T ToTableName](name T, x ...SelectOption) *SelectStmt {
	return Select(append([]SelectOption{All, From(newTableName(name))}, x...)...)
}

func SelectDistinct(x ...SelectOption) *SelectStmt {
	return Select(append([]SelectOption{Distinct}, x...)...)
}

func Select(x ...SelectOption) *SelectStmt {
	s := &SelectStmt{}

	for _, opt := range x {
		opt.applySelectStmt(s)
	}

	return s
}

func (s *SelectStmt) expr() *TableExpr {
	if s.Expr == nil {
		s.Expr = &TableExpr{}
	}
	return s.Expr
}

func (s *SelectStmt) From(x ...ToTableRef) *SelectStmt {
	var refs []TableRef

	for _, v := range x {
		refs = append(refs, v.tableRef())
	}

	s.expr().From = refs

	return s
}

func (s *SelectStmt) Where(search SearchCond) *SelectStmt {
	s.expr().Where = &WhereClause{search}
	return s
}

func (s *SelectStmt) GroupBy(x ...ToGroupingElement) *SelectStmt {
	var elems []GroupingElement

	for _, e := range x {
		elems = append(elems, e.groupingElement())
	}

	s.expr().GroupBy = &GroupByClause{Elems: elems}
	return s
}

func (s *SelectStmt) Having(search SearchCond) *SelectStmt {
	s.expr().Having = &HavingClause{search}
	return s
}

func (s *SelectStmt) Window(x ...ToWindowDef) *SelectStmt {
	var defs []*WindowDef

	for _, e := range x {
		defs = append(defs, e.windowDef())
	}

	s.expr().Window = defs
	return s
}

func (s *SelectStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "SELECT")
	if s.Quantifier != SetAll {
		fmt.Fprintf(&b, " %s", s.Quantifier)
	}
	fmt.Fprintf(&b, " %s", s.Select)
	if s.Into != nil {
		fmt.Fprintf(&b, " INTO %s", s.Into)
	}
	if s.Expr != nil {
		fmt.Fprintf(&b, " %s", s.Expr)
	}

	return b.String()
}

type ToSelectList interface {
	selectList() SelectList
}

type SelectList interface {
	fmt.Stringer

	ToSelectList
}

var (
	_ SelectList = All
	_ SelectList = SelectSubLists(nil)
)

type allFields struct{}

var All = &allFields{}

func (a *allFields) selectList() SelectList        { return a }
func (a *allFields) applySelectStmt(s *SelectStmt) { s.Select = a }
func (a *allFields) String() string                { return "*" }

type SelectSubLists []*SelectSubList

func (l SelectSubLists) selectList() SelectList        { return l }
func (l SelectSubLists) applySelectStmt(s *SelectStmt) { s.Select = l }
func (l SelectSubLists) String() string                { return Join(l, ", ") }

type ToSelectSubList interface {
	selectSubList() *SelectSubList
}

var (
	_ ToSelectSubList = Raw("")
	_ ToSelectSubList = &LocalOrSchemaQualifiedName{}
	_ ToSelectSubList = &SchemaQualifiedName{}
	_ ToSelectSubList = &ColumnDef{}
	_ ToSelectSubList = &CallExpr{}
	_ ToSelectSubList = &SelectSubList{}
)

func (r Raw) As(name ColumnName) *SelectSubList { return &SelectSubList{r, AsClause(name)} }
func (r Raw) applySelectStmt(s *SelectStmt)     { r.selectSubList().applySelectStmt(s) }
func (r Raw) selectSubList() *SelectSubList     { return &SelectSubList{Value: r} }

func (n *LocalOrSchemaQualifiedName) applySelectStmt(s *SelectStmt) {
	n.selectSubList().applySelectStmt(s)
}
func (n *LocalOrSchemaQualifiedName) selectSubList() *SelectSubList { return &SelectSubList{Value: n} }

func (n *SchemaQualifiedName) applySelectStmt(s *SelectStmt) {
	n.selectSubList().applySelectStmt(s)
}
func (n *SchemaQualifiedName) selectSubList() *SelectSubList { return &SelectSubList{Value: n} }

func (e *CallExpr) As(name ColumnName) *SelectSubList { return &SelectSubList{e, AsClause(name)} }
func (e *CallExpr) applySelectStmt(s *SelectStmt)     { e.selectSubList().applySelectStmt(s) }
func (e *CallExpr) selectSubList() *SelectSubList     { return &SelectSubList{Value: e} }

func (d *ColumnDef) As(name ColumnName) *SelectSubList {
	return &SelectSubList{d.expr(), AsClause(name)}
}
func (d *ColumnDef) applySelectStmt(s *SelectStmt) { d.selectSubList().applySelectStmt(s) }
func (d *ColumnDef) selectSubList() *SelectSubList { return &SelectSubList{Value: d.expr()} }

type SelectSubList struct {
	Value ValueExpr
	As    AsClause
}

func (l *SelectSubList) applySelectStmt(s *SelectStmt) {
	switch v := s.Select.(type) {
	case SelectSubLists:
		s.Select = append(v, l)
	default:
		s.Select = SelectSubLists{l}
	}
}
func (l *SelectSubList) selectSubList() *SelectSubList { return l }

func (l *SelectSubList) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s", l.Value)

	if len(l.As) > 0 {
		fmt.Fprintf(&b, " %s", l.As)
	}

	return b.String()
}

type AsClause ColumnName

func (c AsClause) String() string { return fmt.Sprintf("AS %s", string(c)) }

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
	OrderBy OrderByClause
	Limit   *LimitClause
}

func (e *TableExpr) String() string {
	var b strings.Builder

	if e.From != nil {
		b.WriteString(e.From.String())
	}

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

	if e.OrderBy != nil {
		fmt.Fprintf(&b, " %s", e.OrderBy)
	}

	if e.Limit != nil {
		fmt.Fprintf(&b, " %s", e.Limit)
	}

	return b.String()
}
