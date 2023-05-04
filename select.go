package xql

import (
	"fmt"
	"strings"
	"time"
)

type SelectFinalStep struct {
	Stmt *SelectStmt
}

func (s *SelectFinalStep) stmt() *SelectStmt {
	if s.Stmt == nil {
		s.Stmt = &SelectStmt{}
	}
	return s.Stmt
}

func (s *SelectFinalStep) Query() *SelectStmt { return s.Stmt }
func (s *SelectFinalStep) String() string     { return s.Stmt.String() }

type SelectUnionStep struct {
	SelectFinalStep
}

type SelectOptionStep struct {
	SelectUnionStep
}

func (s *SelectOptionStep) Option(o string) *SelectUnionStep {
	s.stmt().expr().Option = o

	return &s.SelectUnionStep
}

type SelectForStep struct {
	SelectOptionStep
}

type SelectForUpdateWaitStep struct {
	SelectForStep
}

func (s *SelectForUpdateWaitStep) Wait(d time.Duration) *SelectForStep {
	s.stmt().expr().forLock().Wait(d)

	return &s.SelectForStep
}

func (s *SelectForUpdateWaitStep) NoWait() *SelectForStep {
	s.stmt().expr().forLock().NoWait()

	return &s.SelectForStep
}

func (s *SelectForUpdateWaitStep) SkipLocked() *SelectForStep {
	s.stmt().expr().forLock().SkipLocked()

	return &s.SelectForStep
}

type SelectForUpdateOfStep struct {
	SelectForUpdateWaitStep
}

func (s *SelectForUpdateOfStep) Of(x ...Field) *SelectForUpdateWaitStep {
	return &s.SelectForUpdateWaitStep
}

type SelectForUpdateStep struct {
	SelectForStep
}

func (s *SelectForUpdateStep) ForUpdate() *SelectForUpdateOfStep {
	s.stmt().expr().forLock().Mode = ForUpdate

	return &SelectForUpdateOfStep{SelectForUpdateWaitStep{s.SelectForStep}}
}

func (s *SelectForUpdateStep) ForNoKeyUpdate() *SelectForUpdateOfStep {
	s.stmt().expr().forLock().Mode = ForNoKeyUpdate

	return &SelectForUpdateOfStep{SelectForUpdateWaitStep{s.SelectForStep}}
}

func (s *SelectForUpdateStep) ForShare() *SelectForUpdateOfStep {
	s.stmt().expr().forLock().Mode = ForShare

	return &SelectForUpdateOfStep{SelectForUpdateWaitStep{s.SelectForStep}}
}

func (s *SelectForUpdateStep) ForKeyShare() *SelectForUpdateOfStep {
	s.stmt().expr().forLock().Mode = ForKeyShare

	return &SelectForUpdateOfStep{SelectForUpdateWaitStep{s.SelectForStep}}
}

func (s *SelectForUpdateStep) withCheckOption() *SelectFinalStep {
	s.stmt().expr().WithCheckOption = true

	return &s.SelectFinalStep
}

func (s *SelectForUpdateStep) withReadOnly() *SelectFinalStep {
	s.stmt().expr().WithReadOnly = true

	return &s.SelectFinalStep
}

type SelectLimitAfterOffsetStep struct {
	SelectForUpdateStep
}

func (s *SelectLimitAfterOffsetStep) Limit(n int) *SelectForUpdateStep {
	s.stmt().expr().Limits.LimitClause = Limit(intValue(n))

	return &s.SelectForUpdateStep
}

type SelectWithTiesAfterOffsetStep struct {
	SelectForUpdateStep
}

func (s *SelectWithTiesAfterOffsetStep) WithTies() *SelectForUpdateStep {
	s.stmt().expr().Limits.WithTies = true

	return &s.SelectForUpdateStep
}

type SelectOffsetStep struct {
	SelectForUpdateStep
}

func (s *SelectOffsetStep) Offset(n int) *SelectForUpdateStep {
	s.stmt().expr().Limits.OffsetClause = Offset(intValue(n))

	return &s.SelectForUpdateStep
}

type SelectWithTiesStep struct {
	SelectOffsetStep
}

func (s *SelectWithTiesStep) WithTies() *SelectOffsetStep {
	s.stmt().expr().Limits.WithTies = true

	return &s.SelectOffsetStep
}

type SelectLimitPercentStep struct {
	SelectWithTiesStep
}

func (s *SelectLimitPercentStep) Percent() *SelectWithTiesStep {
	s.stmt().expr().Limits.Percent = true

	return &s.SelectWithTiesStep
}

type SelectLimitStep struct {
	SelectForUpdateStep
}

func (s *SelectLimitStep) Limit(n int) *SelectLimitPercentStep {
	s.stmt().expr().Limits = Limit(intValue(n)).limitsClause()

	return &SelectLimitPercentStep{SelectWithTiesStep{SelectOffsetStep{s.SelectForUpdateStep}}}
}

func (s *SelectLimitStep) Limits(off, count int) *SelectWithTiesAfterOffsetStep {
	s.stmt().expr().Limits = Limits(intValue(off), intValue(count))

	return &SelectWithTiesAfterOffsetStep{s.SelectForUpdateStep}
}

func (s *SelectLimitStep) Offset(off int) *SelectLimitAfterOffsetStep {
	s.stmt().expr().Limits = Offset(intValue(off)).limitsClause()

	return &SelectLimitAfterOffsetStep{s.SelectForUpdateStep}
}

type SelectOrderByStep struct {
	SelectLimitStep
}

func (s *SelectOrderByStep) OrderBy(x ...ToSortSpec) *SelectLimitStep {
	s.stmt().expr().OrderBy = OrderBy(x...)

	return &s.SelectLimitStep
}

type SelectWindowStep struct {
	SelectOrderByStep
}

func (s *SelectWindowStep) Window(x ...ToWindowDef) *SelectOrderByStep {
	s.stmt().expr().Window = Window(x...)
	return &s.SelectOrderByStep
}

type SelectHavingStep struct {
	SelectWindowStep
}

func (s *SelectHavingStep) Having(x SearchCond) *SelectWindowStep {
	s.stmt().expr().Having = Having(x)
	return &s.SelectWindowStep
}

type SelectGroupByStep struct {
	SelectHavingStep
}

func (s *SelectGroupByStep) GroupBy(x ...ToGroupingElement) *SelectHavingStep {
	s.stmt().expr().GroupBy = GroupBy(x...)
	return &s.SelectHavingStep
}

func (s *SelectGroupByStep) GroupByDistinct(x ...ToGroupingElement) *SelectHavingStep {
	s.stmt().expr().GroupBy = GroupBy(x...).Distinct()
	return &s.SelectHavingStep
}

type SelectWhereStep struct {
	SelectGroupByStep
}

func (s *SelectWhereStep) Where(x SearchCond) *SelectGroupByStep {
	s.stmt().expr().Where = Where(x)
	return &s.SelectGroupByStep
}

type SelectFromStep struct {
	SelectWhereStep
}

func (s *SelectFromStep) From(x ...ToTableRef) *SelectWhereStep {
	s.stmt().expr().From = From(x...)
	return &s.SelectWhereStep
}

type SelectIntoStep struct {
	SelectFromStep
}

func (s *SelectIntoStep) Into(table Table) *SelectFromStep {
	return &s.SelectFromStep
}

type SelectDistinctOnStep struct {
	SelectIntoStep
}

func (s *SelectDistinctOnStep) On(x ...SelectFieldOrAsterisk) *SelectIntoStep {
	return &s.SelectIntoStep
}

func (s *SelectDistinctOnStep) DistinctOn(x ...SelectFieldOrAsterisk) *SelectIntoStep {
	return &s.SelectIntoStep
}

type SelectSelectStep struct {
	SelectDistinctOnStep
}

func (s *SelectSelectStep) Select(x ...SelectFieldOrAsterisk) *SelectSelectStep {
	for _, f := range x {
		s.stmt().Select = f.applySelectList(s.stmt().Select)
	}

	return s
}

type SelectFieldOrAsterisk interface {
	applySelectList(selected SelectList) SelectList
}

var (
	_ SelectFieldOrAsterisk = Asterisk
	_ SelectFieldOrAsterisk = Raw("")
	_ SelectFieldOrAsterisk = &SchemaQualifiedName{}
	_ SelectFieldOrAsterisk = &LocalOrSchemaQualifiedName{}
	_ SelectFieldOrAsterisk = &CallExpr{}
	_ SelectFieldOrAsterisk = &ColumnDef{}
	_ SelectFieldOrAsterisk = &SelectSubList{}
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
func (a *asteriskClause) String() string                                 { return "*" }

type SelectSubLists []*SelectSubList

func (l SelectSubLists) selectList() SelectList { return l }
func (l SelectSubLists) String() string         { return Join(l, ", ") }

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
