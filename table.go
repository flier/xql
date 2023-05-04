package xql

import (
	"fmt"
	"strings"
)

type TableLike interface{}

type Table interface{}

type TableName LocalOrSchemaQualifiedName

type ToTableName = ToLocalOrSchemaQualifiedName

func newTableName[T ToTableName](name T) *TableName {
	return (*TableName)(LocalOrSchemaQName(name))
}

func (n *TableName) tablePrimary() TablePrimary { return n }
func (n *TableName) targetTable() TargetTable   { return n }
func (n *TableName) tableRef() TableRef         { return n }
func (n *TableName) String() string             { return ((*LocalOrSchemaQualifiedName)(n)).String() }

type TableRefList []TableRef

func (l TableRefList) String() string { return Join(l, ", ") }

type ToTableRef interface {
	tableRef() TableRef
}

type TableRef interface {
	fmt.Stringer

	ToTableRef
}

var (
	_ TableRef = &LocalOrSchemaQualifiedName{}
	_ TableRef = &SchemaQualifiedName{}
	_ TableRef = &TableName{}
	_ TableRef = &TableFactor{}
	_ TableRef = JoinedTable(nil)
)

type TableFactor struct {
	Primary TablePrimary
	Sample  *SampleClause
}

func (f *TableFactor) tableRef() TableRef { return f }

func (f *TableFactor) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s", f.Primary)

	if f.Sample == nil {
		fmt.Fprintf(&b, " %s", f.Sample)
	}

	return b.String()
}

type TablePrimary interface {
	fmt.Stringer

	tablePrimary() TablePrimary
}

var (
	_ TablePrimary = &DataSource{}
	_ TablePrimary = QueryName("")
)

type DataSource struct {
	Table       *TableName
	Correlation *CorrelationClause
}

func (n *TableName) As(alias CorrelationName) *DataSource {
	return &DataSource{
		Table:       n,
		Correlation: &CorrelationClause{Name: alias},
	}
}

func (n *SchemaQualifiedName) As(alias CorrelationName) *DataSource {
	return &DataSource{
		Table:       (*TableName)(n.LocalOrSchemaQName()),
		Correlation: &CorrelationClause{Name: alias},
	}
}

func (s *DataSource) tableRef() TableRef         { return &TableFactor{Primary: s} }
func (s *DataSource) tablePrimary() TablePrimary { return s }
func (s *DataSource) String() string {
	var b strings.Builder

	b.WriteString(s.Table.String())

	if s.Correlation != nil {
		fmt.Fprintf(&b, " %s", s.Correlation)
	}

	return b.String()
}

type CorrelationClause struct {
	Name    CorrelationName
	Columns ColumnNameList
}

func (c *CorrelationClause) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "AS %s", c.Name)

	if len(c.Columns) > 0 {
		fmt.Fprintf(&b, " (%s)", c.Columns)
	}

	return b.String()
}

type TableDef struct {
	Scope            *TableScope
	Name             *TableName
	Content          TableContentSource
	SystemVersioning *SystemVersioningClause
	OnCommit         *TableCommitAction
}

type TableDefOption interface {
	applyTableDef(*TableDef)
}

type applyTableDefFunc func(*TableDef)

func (f applyTableDefFunc) applyTableDef(t *TableDef) { f(t) }

func CreateTable[T ToTableName](name T, x ...TableDefOption) *TableDef {
	t := &TableDef{
		Name: newTableName(name),
	}

	for _, opt := range x {
		opt.applyTableDef(t)
	}

	return t
}

func CreateTempTable(name string, x ...TableDefOption) *TableDef {
	return CreateTable(name, append(x, Temporary)...)
}

type TypedTableDefOption interface {
	applyTypedTableDef(*TableDef)
}

var (
	_ TypedTableDefOption = &ColumnDef{}
	_ TypedTableDefOption = &ColumnOptions{}
	_ TypedTableDefOption = &UniqueConstraintDef{}
)

func (t *TableDef) Of(typeName string, x ...TypedTableDefOption) *TableDef {
	t.Content = &TypedTableClause{
		Name: *SchemaQName(typeName),
	}

	for _, opt := range x {
		opt.applyTypedTableDef(t)
	}

	return t
}

func (t *TableDef) String() string {
	var b strings.Builder

	b.WriteString("CREATE")

	if t.Scope != nil {
		fmt.Fprintf(&b, " %s", t.Scope)
	}

	fmt.Fprintf(&b, " TABLE %s", t.Name)

	if t.Content != nil {
		fmt.Fprintf(&b, " %s", t.Content)
	}

	if t.SystemVersioning != nil {
		fmt.Fprintf(&b, " WITH %s", t.SystemVersioning)
	}

	if t.OnCommit != nil {
		fmt.Fprintf(&b, " ON COMMIT %s", t.OnCommit)
	}

	return b.String()
}

type TableScope struct {
	Global    *bool
	Temporary *bool
}

var (
	Temporary = applyTableDefFunc(func(t *TableDef) {
		temp := true

		if t.Scope != nil {
			t.Scope.Temporary = &temp
		} else {
			t.Scope = &TableScope{Temporary: &temp}
		}
	})
	Temp = Temporary
)

func (s *TableScope) String() string {
	var elems []string

	if s.Global != nil {
		if *s.Global {
			elems = append(elems, "GLOBAL")
		} else {
			elems = append(elems, "LOCAL")
		}
	}

	if s.Temporary != nil {
		if *s.Temporary {
			elems = append(elems, "TEMPORARY")
		}
	}

	return strings.Join(elems, " ")
}

type SystemVersioningClause struct {
	On bool
}

var (
	WithSystemVersioning   = &SystemVersioningClause{}
	WithSystemVersioningOn = &SystemVersioningClause{On: true}
)

func (c *SystemVersioningClause) applyTableDef(t *TableDef) { t.SystemVersioning = c }
func (c *SystemVersioningClause) String() string {
	if c.On {
		return "(SYSTEM_VERSIONING = ON)"
	}
	return "SYSTEM VERSIONING"
}

func (t *TableDef) WithSystemVersioning() *TableDef {
	WithSystemVersioning.applyTableDef(t)
	return t
}

func (t *TableDef) WithSystemVersioningOn() *TableDef {
	WithSystemVersioningOn.applyTableDef(t)
	return t
}

//go:generate stringer -type=TableCommitAction -linecomment

type TableCommitAction int

const (
	OnCommitPreserveRows TableCommitAction = iota // PRESERVE ROWS
	OnCommitDeleteRows                            // DELETE ROWS
	OnCommitDrop                                  // DROP
)

func (a TableCommitAction) applyTableDef(t *TableDef) {
	t.OnCommit = &a
}

func (t *TableDef) OnCommitPreserveRows() *TableDef {
	OnCommitPreserveRows.applyTableDef(t)
	return t
}

func (t *TableDef) OnCommitDeleteRows() *TableDef {
	OnCommitDeleteRows.applyTableDef(t)
	return t
}

func (t *TableDef) OnCommitDrop() *TableDef {
	OnCommitDrop.applyTableDef(t)
	return t
}

type TableContentSource interface {
	fmt.Stringer

	TableDefOption

	tableContentSource() TableContentSource
}

var (
	_ TableContentSource = TableElementList{}
	_ TableContentSource = &TypedTableClause{}
	_ TableContentSource = &AsSubQueryClause{}
)

type TableElementList []TableElement

func (l TableElementList) tableContentSource() TableContentSource { return l }
func (l TableElementList) applyTableDef(t *TableDef)              { t.Content = l }
func (l TableElementList) String() string                         { return fmt.Sprintf("(\n\t%s\n)", Join(l, ",\n\t")) }

type ToTableElement interface {
	tableElement() TableElement
}

type TableElement interface {
	fmt.Stringer

	TableDefOption

	ToTableElement
}

var (
	_ TableElement = &ColumnDef{}
	_ TableElement = &TablePeriodDef{}
	_ TableElement = &TableConstraintDef{}
	_ TableElement = &LikeClause{}
)

type TablePeriodDef struct {
	Period Either[*SystemTimePeriodSpec, *ApplicationTimePeriodSpec]
	Begin  ColumnName
	End    ColumnName
}

func PeriodForSystemTime(begin, end ColumnName) *TablePeriodDef {
	return &TablePeriodDef{
		Period: Left[*SystemTimePeriodSpec, *ApplicationTimePeriodSpec](&SystemTimePeriodSpec{}),
		Begin:  begin,
		End:    end,
	}
}

func PeriodFor(name string) func(begin, end ColumnName) *TablePeriodDef {
	return func(begin, end ColumnName) *TablePeriodDef {
		return &TablePeriodDef{
			Period: Right[*SystemTimePeriodSpec](&ApplicationTimePeriodSpec{Name: name}),
			Begin:  begin,
			End:    end,
		}
	}
}

func (d *TablePeriodDef) tableElement() TableElement { return d }
func (d *TablePeriodDef) applyTableDef(t *TableDef) {
	l, _ := t.Content.(TableElementList)
	t.Content = TableElementList(append(l, d))
}

func (d *TablePeriodDef) String() string {
	return fmt.Sprintf("%s (%s, %s)", &d.Period, d.Begin, d.End)
}

type SystemTimePeriodSpec struct{}

func (s *SystemTimePeriodSpec) String() string { return "PERIOD FOR SYSTEM_TIME" }

type ApplicationTimePeriodSpec struct {
	Name string
}

func (s *ApplicationTimePeriodSpec) String() string { return fmt.Sprintf("PERIOD FOR %s", s.Name) }

type TableConstraintDef struct {
	Name            *ConstraintNameDef
	Constraint      TableConstraint
	Characteristics *ConstraintCharacteristics
}

func (d *TableConstraintDef) tableElement() TableElement           { return d }
func (d *TableConstraintDef) typedTableElement() TypedTableElement { return d }
func (d *TableConstraintDef) applyTypedTableDef(t *TableDef) {
	c := t.Content.(*TypedTableClause)
	c.Elements = append(c.Elements, d)
}

func (d *TableConstraintDef) applyTableDef(t *TableDef) {
	l, _ := t.Content.(TableElementList)
	t.Content = TableElementList(append(l, d))
}

func (d *TableConstraintDef) String() string {
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

type ToTableConstraint interface {
	tableConstraint() TableConstraint
}

type TableConstraint interface {
	fmt.Stringer

	TypedTableDefOption

	ToTableConstraint
}

var (
	_ TableConstraint = &UniqueConstraintDef{}
	_ TableConstraint = &ReferentialConstraintDef{}
	_ TableConstraint = &CheckConstraintDef{}
)

type TypedTableClause struct {
	Name     UserDefinedTypeName
	SubTable *SubTableClause
	Elements TypedTableElementList
}

func (c *TypedTableClause) tableContentSource() TableContentSource { return c }
func (c *TypedTableClause) applyTableDef(t *TableDef)              { t.Content = c }
func (c *TypedTableClause) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "OF %s", &c.Name)
	if c.SubTable != nil {
		fmt.Fprintf(&b, " %s", c.SubTable)
	}
	if len(c.Elements) > 0 {
		fmt.Fprintf(&b, " (\n\t%s\n)", Join(c.Elements, ",\n\t"))
	}

	return b.String()
}

type SubTableClause struct {
}

func (c *SubTableClause) String() string { return "" }

type TypedTableElementList []TypedTableElement

type TypedTableElement interface {
	fmt.Stringer

	typedTableElement() TypedTableElement

	applyTypedTableDef(*TableDef)
}

var (
	_ TypedTableElement = &ColumnOptions{}
	_ TypedTableElement = &TableConstraintDef{}
	_ TypedTableElement = &SelfRefColumnSpec{}
)

//go:generate stringer -type=RefGeneration -linecomment

type RefGeneration int

const (
	RefSystemGenerated RefGeneration = iota // SYSTEM GENERATED
	RefUserGenerated                        // USER GENERATED
	RefDerived                              // DERIVED
)

type SelfRefColumnSpec struct {
	Name       ColumnName
	Generation *RefGeneration
}

func SelfRefColumn(name string) *SelfRefColumnSpec {
	return &SelfRefColumnSpec{ColumnName(name), nil}
}

func SystemGenerated(name string) *SelfRefColumnSpec {
	g := RefSystemGenerated
	return &SelfRefColumnSpec{ColumnName(name), &g}
}

func UserGenerated(name string) *SelfRefColumnSpec {
	g := RefUserGenerated
	return &SelfRefColumnSpec{ColumnName(name), &g}
}

func Derived(name string) *SelfRefColumnSpec {
	g := RefDerived
	return &SelfRefColumnSpec{ColumnName(name), &g}
}

func (s *SelfRefColumnSpec) typedTableElement() TypedTableElement { return s }

func (s *SelfRefColumnSpec) applyTypedTableDef(t *TableDef) {
	c := t.Content.(*TypedTableClause)
	c.Elements = append(c.Elements, s)
}

func (s *SelfRefColumnSpec) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "REF IS %s", s.Name)

	if s.Generation != nil {
		fmt.Fprintf(&b, " %s", s.Generation)
	}

	return b.String()
}

type AsSubQueryClause struct{}

func (c *AsSubQueryClause) tableContentSource() TableContentSource { return c }
func (c *AsSubQueryClause) applyTableDef(t *TableDef)              { t.Content = c }
func (c *AsSubQueryClause) String() string                         { return "AS" }
