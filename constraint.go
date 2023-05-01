package xql

import (
	"fmt"
	"strings"
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

func (d *ConstraintNameDef) String() string {
	return fmt.Sprintf("CONSTRAINT %s", &d.Name)
}

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

type ConstraintCharacteristics struct {
	CheckTime   *ConstraintCheckTime
	Deferrable  *ConstraintDeferrable
	Enforcement *ConstraintEnforcement
}

func (c *ConstraintCharacteristics) String() string {
	var v []string

	if c.CheckTime != nil {
		v = append(v, c.CheckTime.String())
	}
	if c.Deferrable != nil {
		v = append(v, c.Deferrable.String())
	}
	if c.Enforcement != nil {
		v = append(v, c.Enforcement.String())
	}

	return strings.Join(v, " ")
}

type ConstraintDeferrable bool

func (d ConstraintDeferrable) String() string {
	if d {
		return "DEFERRABLE"
	}

	return "NOT DEFERRABLE"
}

type ConstraintEnforcement bool

func (c ConstraintEnforcement) String() string {
	if c {
		return "ENFORCED"
	}

	return "NOT ENFORCED"
}

type NotNullConstraint struct{}

var NotNull = &NotNullConstraint{}

func (c *NotNullConstraint) columnConstraint() ColumnConstraint { return c }

func (c *NotNullConstraint) applyColumnDef(d *ColumnDef) {
	d.Constraints = append(d.Constraints, &ColumnConstraintDef{Constraint: NotNull})
}

func (c *NotNullConstraint) String() string { return "NOT NULL" }

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
func (d *UniqueConstraintDef) String() string {
	var b strings.Builder

	b.WriteString(d.Spec.String())
	if len(d.Columns) > 0 {
		fmt.Fprintf(&b, " (%s)", d.Columns)
	}

	return b.String()
}

type ReferentialConstraintDef struct {
	Columns ColumnNameList
	Spec    ReferencesSpec
}

func (d *ReferentialConstraintDef) tableConstraint() TableConstraint { return d }

func (d *ReferentialConstraintDef) applyTypedTableDef(t *TableDef) {
	c := t.Content.(*TypedTableClause)
	c.Elements = append(c.Elements, &TableConstraintDef{Constraint: d})
}

func (d *ReferentialConstraintDef) String() string {
	return fmt.Sprintf("FOREIGN KEY (%s) %s", d.Columns, &d.Spec)
}

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

func (s *ReferencesSpec) columnConstraint() ColumnConstraint { return s }
func (s *ReferencesSpec) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "REFERENCES %s", s.Name)

	if len(s.Columns) > 0 {
		fmt.Fprintf(&b, " (%s)", s.Columns)
	}

	if s.Match != MatchSimple {
		fmt.Fprintf(&b, " MATCH %s", s.Match)
	}

	if s.Action != nil {
		fmt.Fprintf(&b, " %s", s.Action)
	}

	return b.String()
}

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

func (r *ReferentialTriggeredAction) String() string {
	var v []string

	if r.OnUpdate != NoAction {
		v = append(v, fmt.Sprintf("ON UPDATE %s", r.OnUpdate))
	}

	if r.OnDelete != NoAction {
		v = append(v, fmt.Sprintf("ON DELETE %s", r.OnDelete))
	}

	return strings.Join(v, " ")
}

type CheckConstraintDef struct {
	Cond SearchCond
}

func Check(cond SearchCond) *CheckConstraintDef {
	return &CheckConstraintDef{Cond: cond}
}

func (c *CheckConstraintDef) columnConstraint() ColumnConstraint { return c }
func (c *CheckConstraintDef) applyColumnDef(d *ColumnDef) {
	d.Constraints = append(d.Constraints, &ColumnConstraintDef{Constraint: c})
}
func (c *CheckConstraintDef) applyTypedTableDef(t *TableDef) {
	tc := t.Content.(*TypedTableClause)
	tc.Elements = append(tc.Elements, &TableConstraintDef{Constraint: c})
}
func (c *CheckConstraintDef) tableConstraint() TableConstraint { return c }
func (c *CheckConstraintDef) String() string                   { return fmt.Sprintf("CHECK (%s)", c.Cond) }
