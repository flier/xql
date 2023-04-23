package xql

type UserDefinedTypeName = SchemaQualifiedName

var _ DataType = &UserDefinedTypeName{}

func (n *UserDefinedTypeName) dataType() DataType          { return n }
func (n *UserDefinedTypeName) applyColumnDef(d *ColumnDef) { d.Type = n }
