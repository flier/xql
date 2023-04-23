package xql

import (
	"fmt"
	"strings"
)

type RefType struct {
	Name  UserDefinedTypeName
	Scope *ScopeClause
}

var _ DataType = &RefType{}

func (t *RefType) dataType() DataType          { return t }
func (t *RefType) applyColumnDef(d *ColumnDef) { d.Type = t }
func (t *RefType) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "REF(%s)", &t.Name)

	if t.Scope != nil {
		fmt.Fprintf(&b, " %s", t.Scope)
	}

	return b.String()
}
