package xql

import (
	"fmt"
	"strings"
)

type SelectStmt struct {
	*TableExpr
	Quantifier *SetQuantifier
	Select     SelectList
	Into       *TargetSpec
}

func SelectAllFrom(x ...ToTableRef) *SelectWhereStep {
	return Select(Asterisk).From(x...)
}

func Select(x ...SelectFieldOrAsterisk) *SelectSelectStep {
	s := &SelectSelectStep{}
	s.Select(x...)
	return s
}

func SelectDistinct(x ...SelectFieldOrAsterisk) *SelectSelectStep {
	s := Select(x...)
	s.stmt().Quantifier = &Distinct
	return s
}

func (s *SelectStmt) expr() *TableExpr {
	if s.TableExpr == nil {
		s.TableExpr = &TableExpr{}
	}
	return s.TableExpr
}

var kSelect = Keyword("SELECT")

func (s *SelectStmt) Accept(v Visitor) Visitor {
	return v.Keyword(kSelect).
		IfNotNil(s.Quantifier, WS, Stringer(s.Quantifier)).
		Visit(WS, s.Select)
}

func (s *SelectStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "SELECT")

	if s.Quantifier != nil {
		fmt.Fprintf(&b, " %s", s.Quantifier)
	}

	fmt.Fprintf(&b, " %s", s.Select)

	if s.TableExpr != nil {
		fmt.Fprintf(&b, " %s", s.TableExpr)
	}

	if s.Into != nil {
		fmt.Fprintf(&b, " INTO %s", s.Into)
	}

	return b.String()
}

type ToSelectList interface {
	selectList() SelectList
}

type SelectList interface {
	fmt.Stringer

	Accepter

	ToSelectList
}

var (
	_ SelectList = Asterisk
	_ SelectList = SelectSubLists(nil)
)

type asteriskClause struct{}

var Asterisk = &asteriskClause{}

func (a *asteriskClause) selectList() SelectList                         { return a }
func (a *asteriskClause) applySelectList(selected SelectList) SelectList { return a }
func (a *asteriskClause) Accept(v Visitor) Visitor                       { return v.Raw("*") }
func (a *asteriskClause) String() string                                 { return "*" }

type SelectSubLists []*SelectSubList

func (l SelectSubLists) selectList() SelectList { return l }
func (l SelectSubLists) String() string         { return Join(l, ", ") }
func (l SelectSubLists) Accept(v Visitor) Visitor {
	for i, s := range l {
		if i > 0 {
			v.Sep()
		}

		s.Accept(v)
	}

	return v
}

type ToSelectSubList interface {
	selectSubList() *SelectSubList
}

func appendSelectList(l SelectList, v ToSelectSubList) SelectList {
	if l != nil {
		if l, ok := l.(SelectSubLists); ok {
			return append(l, v.selectSubList())
		}
	}

	return SelectSubLists{v.selectSubList()}
}

var (
	_ ToSelectSubList = Raw("")
	_ ToSelectSubList = &LocalOrSchemaQualifiedName{}
	_ ToSelectSubList = &SchemaQualifiedName{}
	_ ToSelectSubList = &ColumnDef{}
	_ ToSelectSubList = &CallExpr{}
	_ ToSelectSubList = &SelectSubList{}
)

func (r Raw) As(name ColumnName) *SelectSubList       { return &SelectSubList{r, AsClause(name)} }
func (r Raw) selectSubList() *SelectSubList           { return &SelectSubList{Value: r} }
func (r Raw) applySelectList(l SelectList) SelectList { return appendSelectList(l, r) }

func (n *LocalOrSchemaQualifiedName) selectSubList() *SelectSubList { return &SelectSubList{Value: n} }
func (n *LocalOrSchemaQualifiedName) applySelectList(l SelectList) SelectList {
	return appendSelectList(l, n)
}

func (n *SchemaQualifiedName) selectSubList() *SelectSubList           { return &SelectSubList{Value: n} }
func (n *SchemaQualifiedName) applySelectList(l SelectList) SelectList { return appendSelectList(l, n) }

func (e *CallExpr) As(name ColumnName) *SelectSubList       { return &SelectSubList{e, AsClause(name)} }
func (e *CallExpr) selectSubList() *SelectSubList           { return &SelectSubList{Value: e} }
func (e *CallExpr) applySelectList(l SelectList) SelectList { return appendSelectList(l, e) }

func (d *ColumnDef) As(name ColumnName) *SelectSubList {
	return &SelectSubList{d.expr(), AsClause(name)}
}
func (d *ColumnDef) selectSubList() *SelectSubList           { return &SelectSubList{Value: d.expr()} }
func (d *ColumnDef) applySelectList(l SelectList) SelectList { return appendSelectList(l, d) }

type SelectSubList struct {
	Value ValueExpr
	As    AsClause
}

func (l *SelectSubList) selectSubList() *SelectSubList           { return l }
func (l *SelectSubList) applySelectList(s SelectList) SelectList { return appendSelectList(s, l) }

func (l *SelectSubList) Accept(v Visitor) Visitor {
	return v.Raw(l.Value.String()).If(len(l.As) > 0, WS, Stringer(l.As))
}

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
	From            FromClause
	Where           *WhereClause
	GroupBy         *GroupByClause
	Having          *HavingClause
	Window          WindowClause
	OrderBy         OrderByClause
	Limits          *LimitsClause
	ForLock         *ForLockClause
	WithCheckOption bool
	WithReadOnly    bool
	Option          string
}

func (e *TableExpr) forLock() *ForLockClause {
	if e.ForLock == nil {
		e.ForLock = &ForLockClause{}
	}
	return e.ForLock
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

	if e.Limits != nil {
		fmt.Fprintf(&b, " %s", e.Limits)
	}

	if e.ForLock != nil {
		fmt.Fprintf(&b, " %s", e.ForLock)
	}

	if e.WithCheckOption {
		fmt.Fprintf(&b, " WITH CHECK OPTION")
	}

	if e.WithReadOnly {
		fmt.Fprintf(&b, " WITH READ ONLY")
	}

	if len(e.Option) > 0 {
		fmt.Fprintf(&b, " %s", e.Option)
	}

	return b.String()
}
