package xql

import (
	"fmt"
	"strings"
)

type ColumnNameList []ColumnName

func (l ColumnNameList) Accept(v Visitor) Visitor {
	if len(l) > 0 {
		v.Token('(')
		for i, n := range l {
			if i > 0 {
				v.Sep().WS()
			}
			v.Visit(QName(n))
		}
		v.Token(')')
	}
	return v
}

func (l ColumnNameList) String() string { return XQL(l) }

type (
	ColumnName = string
	ColumnRef  = string
)

// ColumnDef define a column of a base table.
//
//	<column definition> ::=
//		<column name>  [ <data type or domain name>  ]
//		    [ <default clause>  | <identity column specification>  | <generation clause>
//		    | <system time period start column specification>
//		    | <system time period end column specification>  ]
//		    [ <column constraint definition> ... ]
//		    [ <collate clause>  ]
//
// https://jakewheat.github.io/sql-overview/sql-2016-foundation-grammar.html#column-definition
type ColumnDef struct {
	Name        ColumnName
	Type        DataType
	Value       ColumnValue
	Constraints []*ColumnConstraintDef
	Collate     *CollateClause
}

type ColumnDefOption interface {
	applyColumnDef(*ColumnDef)
}

var (
	_ ColumnDefOption = CreateUniqueConstraintDefFunc(nil)
	C                 = Column
)

func Column(name string, x ...ColumnDefOption) *ColumnDef {
	d := &ColumnDef{
		Name: ColumnName(name),
	}

	for _, opt := range x {
		opt.applyColumnDef(d)
	}

	return d
}

type ColumnExpr struct {
	*ColumnDef
}

func (e *ColumnExpr) expr() Expr     { return e }
func (e *ColumnExpr) String() string { return e.Name }

func (d *ColumnDef) expr() Expr                 { return &ColumnExpr{d} }
func (d *ColumnDef) tableElement() TableElement { return d }
func (d *ColumnDef) applyTableDef(t *TableDef) {
	l, _ := t.Content.(TableElementList)
	t.Content = TableElementList(append(l, d))
}

type ColumnOption interface {
	applyColumnOptions(*ColumnOptions)
}

func (d *ColumnDef) WithOptions(x ...ColumnOption) *ColumnOptions {
	o := &ColumnOptions{Name: d.Name, Constraints: d.Constraints}

	if c, ok := d.Value.(*DefaultClause); ok && c != nil {
		o.Default = c
	}

	for _, opt := range x {
		opt.applyColumnOptions(o)
	}

	return o
}

func (d *ColumnDef) applyTypedTableDef(t *TableDef) {
	c := t.Content.(*TypedTableClause)
	c.Elements = append(c.Elements, d.WithOptions())
}

func (d *ColumnDef) Accept(v Visitor) Visitor {
	return v.Visit(QName(d.Name)).
		IfNotNil(d.Type, AcceptFunc(func(v Visitor) Visitor { return v.WS().DataType(d.Type) })).
		Visit(WS, d.Value).
		IfNotNil(d.Constraints, WS, Joins(d.Constraints, WS)).
		Visit(WS, d.Collate)
}

func (d *ColumnDef) String() string { return XQL(d) }

type DomainName struct {
	SchemaQualifiedName
}

var _ DataType = &DomainName{}

func (n *DomainName) dataType() DataType          { return n }
func (n *DomainName) applyColumnDef(d *ColumnDef) { d.Type = n }

type ToColumnValue interface {
	columnValue() ColumnValue
}

type ColumnValue interface {
	fmt.Stringer

	Accepter

	ToColumnValue
}

var (
	_ ColumnValue = &DefaultClause{}
	_ ColumnValue = &IdentityColumnSpec{}
	_ ColumnValue = &GenerationClause{}
	_ ColumnValue = &SystemTimePeriodStartColumnSpec{}
	_ ColumnValue = &SystemTimePeriodEndColumnSpec{}
)

//go:generate stringer -type=GeneratedAction -linecomment

type GeneratedAction int

const (
	GeneratedAlways    GeneratedAction = iota // ALWAYS
	GeneratedByDefault                        // BY DEFAULT
)

func (a GeneratedAction) Accept(v Visitor) Visitor { return v.Keyword(a) }

type GeneratedSpec struct {
	Action GeneratedAction
}

var Generated = &GeneratedSpec{}

func (g *GeneratedSpec) Always() *GeneratedSpec {
	g.Action = GeneratedAlways
	return g
}

func (g *GeneratedSpec) ByDefault() *GeneratedSpec {
	g.Action = GeneratedByDefault
	return g
}

func (g *GeneratedSpec) AsIdentity(x ...SequenceGeneratorOption) *IdentityColumnSpec {
	return &IdentityColumnSpec{
		Action:  g.Action,
		Options: x,
	}
}

func (g *GeneratedSpec) AsRowStart() *SystemTimePeriodStartColumnSpec {
	return &SystemTimePeriodStartColumnSpec{}
}

func (g *GeneratedSpec) AsRowEnd() *SystemTimePeriodEndColumnSpec {
	return &SystemTimePeriodEndColumnSpec{}
}

func (g *GeneratedSpec) AsValue(value ValueExpr) *GenerationClause {
	return &GenerationClause{Value: value}
}

type IdentityColumnSpec struct {
	Action  GeneratedAction
	Options []SequenceGeneratorOption
}

func (s *IdentityColumnSpec) columnValue() ColumnValue { return s }
func (s *IdentityColumnSpec) applyColumnDef(d *ColumnDef) {
	d.Value = s
}

const (
	kGenerated  = Keyword("GENERATED")
	kAsIdentity = Keyword("AS IDENTITY")
)

func (s *IdentityColumnSpec) Accept(v Visitor) Visitor {
	return v.Visit(kGenerated, WS, s.Action, WS, kAsIdentity).
		IfNotNil(s.Options, WS, Paren(Joins(s.Options, Sep)))
}

func (s *IdentityColumnSpec) String() string { return XQL(s) }

type GenerationClause struct {
	GenerationRule
	Value ValueExpr
}

const kAs = Keyword("AS")

func (c *GenerationClause) columnValue() ColumnValue    { return c }
func (c *GenerationClause) applyColumnDef(d *ColumnDef) { d.Value = c }
func (c *GenerationClause) Accept(v Visitor) Visitor {
	return v.Visit(&c.GenerationRule, kAs, Paren(Raw(c.Value.String())))
}
func (c *GenerationClause) String() string { return XQL(c) }

type GenerationRule struct{}

func (r *GenerationRule) Accept(v Visitor) Visitor { return v.Keyword(r) }
func (r *GenerationRule) String() string           { return "GENERATED ALWAYS" }

type TimestampGenerationRule = GenerationRule

type SystemTimePeriodStartColumnSpec struct {
	TimestampGenerationRule
}

const kAsRowStart = Keyword("AS ROW START")

func (s *SystemTimePeriodStartColumnSpec) columnValue() ColumnValue    { return s }
func (s *SystemTimePeriodStartColumnSpec) applyColumnDef(d *ColumnDef) { d.Value = s }
func (s *SystemTimePeriodStartColumnSpec) Accept(v Visitor) Visitor {
	return v.Visit(&s.TimestampGenerationRule, WS, kAsRowStart)
}
func (s *SystemTimePeriodStartColumnSpec) String() string { return XQL(s) }

type SystemTimePeriodEndColumnSpec struct {
	TimestampGenerationRule
}

const kAsRowEnd = Keyword("AS ROW END")

func (s *SystemTimePeriodEndColumnSpec) columnValue() ColumnValue    { return s }
func (s *SystemTimePeriodEndColumnSpec) applyColumnDef(d *ColumnDef) { d.Value = s }
func (s *SystemTimePeriodEndColumnSpec) Accept(v Visitor) Visitor {
	return v.Visit(&s.TimestampGenerationRule, WS, kAsRowEnd)
}
func (s *SystemTimePeriodEndColumnSpec) String() string { return XQL(s) }

type ColumnConstraintDef struct {
	Name            *ConstraintNameDef
	Constraint      ColumnConstraint
	Characteristics *ConstraintCharacteristics
}

func (d *ColumnConstraintDef) applyColumnDef(c *ColumnDef) {
	c.Constraints = append(c.Constraints, d)
}

func (d *ColumnConstraintDef) Accept(v Visitor) Visitor {
	return v.Visit(d.Name, WS).Visit(d.Constraint).Visit(WS, d.Characteristics)
}

func (d *ColumnConstraintDef) String() string { return XQL(d) }

type ToColumnConstraint interface {
	columnConstraint() ColumnConstraint
}

type ColumnConstraint interface {
	fmt.Stringer

	Accepter

	ToColumnConstraint
}

var (
	_ ColumnConstraint = NotNull
	_ ColumnConstraint = UniqueSpec(0)
	_ ColumnConstraint = &ReferencesSpec{}
	_ ColumnConstraint = &CheckConstraintDef{}
)

type ColumnOptions struct {
	Name        ColumnName
	Scope       *ScopeClause
	Default     *DefaultClause
	Constraints []*ColumnConstraintDef
}

func (o *ColumnOptions) typedTableElement() TypedTableElement { return o }
func (o *ColumnOptions) applyTypedTableDef(t *TableDef) {
	c := t.Content.(*TypedTableClause)
	c.Elements = append(c.Elements, o)
}

const kWithOptions = Keyword("WITH OPTIONS")

func (o *ColumnOptions) Accept(v Visitor) Visitor {
	return v.Visit(QName(o.Name), WS, kWithOptions).
		IfNotNil(o.Scope, WS, o.Scope).
		IfNotNil(o.Default, WS, o.Default).
		IfNotNil(o.Constraints, WS, Joins(o.Constraints, WS))
}

func (o *ColumnOptions) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s WITH OPTIONS", o.Name)

	if o.Scope != nil {
		fmt.Fprintf(&b, " %s", o.Scope)
	}

	if o.Default != nil {
		fmt.Fprintf(&b, " %s", o.Default)
	}

	if len(o.Constraints) > 0 {
		fmt.Fprintf(&b, " %s", Join(o.Constraints, " "))
	}

	return b.String()
}
