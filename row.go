package xql

import "fmt"

type RowType struct {
	Fields []*FieldDef
}

func RowTy(fields ...*FieldDef) *RowType { return &RowType{fields} }

var _ DataType = &RowType{}

func (t *RowType) dataType() DataType          { return t }
func (t *RowType) applyColumnDef(d *ColumnDef) { d.Type = t }
func (t *RowType) String() string              { return fmt.Sprintf("ROW (%s)", Join(t.Fields, ", ")) }

type FieldDef struct {
	Name string
	Type DataType
}

func Field(name string, typ DataType) *FieldDef { return &FieldDef{name, typ} }

func (f *FieldDef) String() string {
	return fmt.Sprintf("%s %s", f.Name, f.Type)
}
