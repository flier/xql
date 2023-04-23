package xql

import "fmt"

type MultiSetType struct {
	Type DataType
}

var _ DataType = &MultiSetType{}

func (t *MultiSetType) dataType() DataType          { return t }
func (t *MultiSetType) applyColumnDef(d *ColumnDef) { d.Type = t }
func (t *MultiSetType) String() string              { return fmt.Sprintf("%s MULTISET", t.Type) }
