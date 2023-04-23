package xql

import (
	"fmt"
	"strings"
)

type ColumnName = string

type ColumnRef string

func (r ColumnRef) String() string { return string(r) }

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

func (d *ColumnDef) String() string {
	var b strings.Builder

	b.WriteString(d.Name)

	if d.Type != nil {
		fmt.Fprintf(&b, " %s", d.Type)
	}
	if d.Value != nil {
		fmt.Fprintf(&b, " %s", d.Value)
	}
	if len(d.Constraints) > 0 {
		fmt.Fprintf(&b, " %s", Join(d.Constraints, " "))
	}
	if d.Collate != nil {
		fmt.Fprintf(&b, " %s", d.Collate)
	}

	return b.String()
}

type DomainName struct {
	SchemaQualifiedName
}

var _ DataType = &DomainName{}

func (n *DomainName) dataType() DataType          { return n }
func (n *DomainName) applyColumnDef(d *ColumnDef) { d.Type = n }

type ColumnValue interface {
	fmt.Stringer

	columnValue() ColumnValue
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

func (s *IdentityColumnSpec) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "GENERATED %s AS IDENTITY", s.Action)
	if len(s.Options) > 0 {
		fmt.Fprintf(&b, " (%s)", Join(s.Options, ", "))
	}

	return b.String()
}

type GenerationClause struct {
	GenerationRule
	Value ValueExpr
}

func (c *GenerationClause) columnValue() ColumnValue    { return c }
func (c *GenerationClause) applyColumnDef(d *ColumnDef) { d.Value = c }
func (c *GenerationClause) String() string {
	return fmt.Sprintf("%s AS (%s)", &c.GenerationRule, c.Value)
}

type GenerationRule struct{}

func (r *GenerationRule) String() string { return "GENERATED ALWAYS" }

type TimestampGenerationRule = GenerationRule

type SystemTimePeriodStartColumnSpec struct {
	TimestampGenerationRule
}

func (s *SystemTimePeriodStartColumnSpec) columnValue() ColumnValue    { return s }
func (s *SystemTimePeriodStartColumnSpec) applyColumnDef(d *ColumnDef) { d.Value = s }
func (s *SystemTimePeriodStartColumnSpec) String() string {
	return fmt.Sprintf("%s AS ROW START", &s.TimestampGenerationRule)
}

type SystemTimePeriodEndColumnSpec struct {
	TimestampGenerationRule
}

func (s *SystemTimePeriodEndColumnSpec) columnValue() ColumnValue    { return s }
func (s *SystemTimePeriodEndColumnSpec) applyColumnDef(d *ColumnDef) { d.Value = s }
func (s *SystemTimePeriodEndColumnSpec) String() string {
	return fmt.Sprintf("%s AS ROW END", &s.TimestampGenerationRule)
}

type ColumnConstraintDef struct {
	Name            *ConstraintNameDef
	Constraint      ColumnConstraint
	Characteristics *ConstraintCharacteristics
}

func (d *ColumnConstraintDef) applyColumnDef(c *ColumnDef) {
	c.Constraints = append(c.Constraints, d)
}

func (d *ColumnConstraintDef) String() string {
	var b strings.Builder

	if d.Name != nil {
		fmt.Fprintf(&b, "%s ", d.Name)
	}

	b.WriteString(d.Constraint.String())

	if d.Characteristics != nil {
		fmt.Fprintf(&b, " %s", d.Characteristics)
	}

	return b.String()
}

type ColumnConstraint interface {
	fmt.Stringer

	columnConstraint() ColumnConstraint
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
