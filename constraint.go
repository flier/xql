package xql

import (
	"fmt"
)

type ConstraintName = SchemaQualifiedName

type ConstraintNameDef struct {
	Name ConstraintName
}

func constraintName[T ~string | SchemaQualifiedName | *SchemaQualifiedName](name T) ConstraintName {
	switch v := any(name).(type) {
	case string:
		return *SchemaQName(v)
	case SchemaQualifiedName:
		return v
	case *SchemaQualifiedName:
		return *v
	default:
		panic("unreachable")
	}
}

func Constraint[T ~string | *SchemaQualifiedName](name T) *ConstraintNameDef {
	return &ConstraintNameDef{constraintName(name)}
}

const kConstraint = Keyword("CONSTRAINT")

func (d *ConstraintNameDef) Accept(v Visitor) Visitor {
	return v.Visit(kConstraint, WS, &d.Name)
}

func (d *ConstraintNameDef) String() string { return XQL(d) }

func (d *ConstraintNameDef) NotNull() *ColumnConstraintDef {
	return &ColumnConstraintDef{Name: d, Constraint: NotNull}
}

func (d *ConstraintNameDef) PrimaryKey(x ...ColumnName) *GenericConstraintDef {
	return &GenericConstraintDef{Name: d, Constraint: &UniqueConstraintDef{Spec: SpecPrimaryKey, Columns: x}}
}

func (d *ConstraintNameDef) Unique(x ...ColumnName) *GenericConstraintDef {
	return &GenericConstraintDef{Name: d, Constraint: &UniqueConstraintDef{Spec: SpecUnique, Columns: x}}
}

type GenericConstraint interface {
	fmt.Stringer

	ColumnConstraint
	TableConstraint
}

type GenericConstraintDef struct {
	Name       *ConstraintNameDef
	Constraint GenericConstraint
}

func (d *GenericConstraintDef) applyColumnDef(c *ColumnDef) {
	c.Constraints = append(c.Constraints, &ColumnConstraintDef{Name: d.Name, Constraint: d.Constraint})
}

func (d *GenericConstraintDef) applyTableDef(t *TableDef) {
	l, _ := t.Content.(TableElementList)
	t.Content = TableElementList(append(l, &TableConstraintDef{Name: d.Name, Constraint: d.Constraint}))
}

//go:generate stringer -type=ConstraintCheckTime -linecomment

type ConstraintCheckTime int

const (
	CheckTimeInitiallyDeferred  ConstraintCheckTime = iota // INITIALLY DEFERRED
	CheckTimeInitiallyImmediate                            // INITIALLY IMMEDIATE
)

func (c *ConstraintCheckTime) Accept(v Visitor) Visitor {
	return v.Visit(Keyword(c.String()))
}

type ConstraintCharacteristics struct {
	CheckTime   *ConstraintCheckTime
	Deferrable  *ConstraintDeferrable
	Enforcement *ConstraintEnforcement
}

func (c *ConstraintCharacteristics) Accept(v Visitor) Visitor {
	return v.Visit(c.CheckTime, WS).Visit(c.Deferrable, WS).Visit(c.Enforcement)
}

func (c *ConstraintCharacteristics) String() string { return XQL(c) }

type ConstraintDeferrable bool

const (
	kDeferrable    = Keyword("DEFERRABLE")
	kNotDeferrable = Keyword("NOT DEFERRABLE")
)

func (d *ConstraintDeferrable) Accept(v Visitor) Visitor {
	if *d {
		return v.Visit(kDeferrable)
	}

	return v.Visit(kNotDeferrable)
}

func (d *ConstraintDeferrable) String() string { return XQL(d) }

type ConstraintEnforcement bool

const (
	kEnforced    = Keyword("ENFORCED")
	kNotEnforced = Keyword("NOT ENFORCED")
)

func (c *ConstraintEnforcement) Accept(v Visitor) Visitor {
	if *c {
		return v.Visit(kEnforced)
	}

	return v.Visit(kNotEnforced)
}

func (c *ConstraintEnforcement) String() string { return XQL(c) }

type NotNullConstraint struct{}

var NotNull = &NotNullConstraint{}

func (c *NotNullConstraint) columnConstraint() ColumnConstraint { return c }

func (c *NotNullConstraint) applyColumnDef(d *ColumnDef) {
	d.Constraints = append(d.Constraints, &ColumnConstraintDef{Constraint: NotNull})
}

const kNotNull = Keyword("NOT NULL")

func (c *NotNullConstraint) Accept(v Visitor) Visitor { return v.Visit(kNotNull) }
func (c *NotNullConstraint) String() string           { return XQL(c) }

//go:generate stringer -type=UniqueSpec -linecomment

type UniqueSpec int

const (
	SpecUnique     UniqueSpec = iota // UNIQUE
	SpecPrimaryKey                   // PRIMARY KEY
)

func (s UniqueSpec) columnConstraint() ColumnConstraint { return s }
func (s UniqueSpec) applyColumnDef(d *ColumnDef) {
	d.Constraints = append(d.Constraints, &ColumnConstraintDef{Constraint: s})
}
func (s UniqueSpec) Accept(v Visitor) Visitor { return v.Visit(Keyword(s.String())) }

type UniqueConstraintDef struct {
	Spec    UniqueSpec
	Columns ColumnNameList
}

var (
	Unique     = uniqueSpec(SpecUnique)
	PrimaryKey = uniqueSpec(SpecPrimaryKey)
)

type CreateUniqueConstraintDefFunc func(columns ...ColumnName) *UniqueConstraintDef

func (f CreateUniqueConstraintDefFunc) applyColumnDef(d *ColumnDef) {
	d.Constraints = append(d.Constraints, &ColumnConstraintDef{Constraint: f()})
}

func uniqueSpec(spec UniqueSpec) CreateUniqueConstraintDefFunc {
	return func(columns ...ColumnName) *UniqueConstraintDef {
		return &UniqueConstraintDef{spec, columns}
	}
}

func (d *UniqueConstraintDef) tableConstraint() TableConstraint { return d }
func (d *UniqueConstraintDef) applyTypedTableDef(t *TableDef) {
	c := t.Content.(*TypedTableClause)
	c.Elements = append(c.Elements, &TableConstraintDef{Constraint: d})
}

func (d *UniqueConstraintDef) columnConstraint() ColumnConstraint { return d.Spec }
func (d *UniqueConstraintDef) Accept(v Visitor) Visitor {
	return v.Keyword(d.Spec).Visit(WS, d.Columns)
}
func (d *UniqueConstraintDef) String() string { return XQL(d) }

type ReferentialConstraintDef struct {
	Columns ColumnNameList
	Spec    ReferencesSpec
}

func (d *ReferentialConstraintDef) tableConstraint() TableConstraint { return d }

func (d *ReferentialConstraintDef) applyTypedTableDef(t *TableDef) {
	c := t.Content.(*TypedTableClause)
	c.Elements = append(c.Elements, &TableConstraintDef{Constraint: d})
}

const kForeignKey = Keyword("FOREIGN KEY")

func (d *ReferentialConstraintDef) Accept(v Visitor) Visitor {
	return v.Visit(kForeignKey, WS, d.Columns, WS, &d.Spec)
}

func (d *ReferentialConstraintDef) String() string { return XQL(d) }

//go:generate stringer -type=MatchType -linecomment

type MatchType int

const (
	MatchSimple  MatchType = iota // SIMPLE
	MatchFull                     // FULL
	MatchPartial                  // PARTIAL
)

type ReferencesSpec struct {
	Name    TableName
	Columns ColumnNameList
	Match   MatchType
	Action  *ReferentialTriggeredAction
}

const (
	kReferences = Keyword("REFERENCES")
	kMatch      = Keyword("MATCH")
)

func (s *ReferencesSpec) columnConstraint() ColumnConstraint { return s }
func (s *ReferencesSpec) Accept(v Visitor) Visitor {
	return v.Visit(kReferences, WS, &s.Name).
		IfNotNil(s.Columns, WS, Bracket(s.Columns)).
		If(s.Match != MatchSimple, WS, kMatch, WS, Keyword(s.Match.String())).
		IfNotNil(s.Action, WS, s.Action)
}

func (s *ReferencesSpec) String() string { return XQL(s) }

//go:generate stringer -type=ReferentialAction -linecomment

type ReferentialAction int

const (
	NoAction   ReferentialAction = iota // NO ACTION
	Cascade                             // CASCADE
	SetNull                             // SET NULL
	SetDefault                          // SET DEFAULT
	Restrict                            // RESTRICT
)

type ReferentialTriggeredAction struct {
	OnUpdate ReferentialAction
	OnDelete ReferentialAction
}

func (r *ReferentialTriggeredAction) Accept(v Visitor) Visitor {
	return v.If(r.OnUpdate != NoAction, Keyword("ON UPDATE"), WS, Keyword(r.OnUpdate.String())).
		If(r.OnDelete != NoAction, Keyword("ON DELETE"), WS, Keyword(r.OnDelete.String()))
}

func (r *ReferentialTriggeredAction) String() string { return XQL(r) }

type CheckConstraintDef struct {
	Cond SearchCond
}

func Check(cond SearchCond) *CheckConstraintDef {
	return &CheckConstraintDef{Cond: cond}
}

const kCheck = Keyword("CHECK")

func (c *CheckConstraintDef) columnConstraint() ColumnConstraint { return c }
func (c *CheckConstraintDef) applyColumnDef(d *ColumnDef) {
	d.Constraints = append(d.Constraints, &ColumnConstraintDef{Constraint: c})
}
func (c *CheckConstraintDef) applyTypedTableDef(t *TableDef) {
	tc := t.Content.(*TypedTableClause)
	tc.Elements = append(tc.Elements, &TableConstraintDef{Constraint: c})
}
func (c *CheckConstraintDef) tableConstraint() TableConstraint { return c }
func (c *CheckConstraintDef) Accept(v Visitor) Visitor {
	return v.Visit(kCheck, WS, Paren(Raw(c.Cond.String())))
}
func (c *CheckConstraintDef) String() string { return XQL(c) }
