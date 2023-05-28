package xql

import "time"

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
