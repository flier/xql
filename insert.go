package xql

import (
	"fmt"
	"strings"
)

type InsertStmt struct {
	Target *TableName
	From   InsertFrom
}

func InsertInto[T ToLocalOrSchemaQualifiedName](name T, from InsertFrom) *InsertStmt {
	return &InsertStmt{newTableName(name), from}
}

func (i *InsertStmt) String() string {
	return fmt.Sprintf("INSERT INTO %s %s", i.Target, i.From)
}

type InsertFrom interface {
	insertFrom() InsertFrom
}

var (
	_ InsertFrom = rawExpr("")
	_ InsertFrom = &FromSubQuery{}
	_ InsertFrom = &FromConstructor{}
	_ InsertFrom = &FromDefault{}
)

//go:generate stringer -type OverridingClause -linecomment

type OverridingClause int

const (
	OverridingUserValue OverridingClause = iota
	OverridingSystemValue
)

type FromSubQuery struct {
	Columns    ColumnNameList
	Overriding *OverridingClause
	SubQuery   QueryExpr
}

func (f *FromSubQuery) insertFrom() InsertFrom { return f }

func (f *FromSubQuery) String() string {
	var b strings.Builder

	if len(f.Columns) > 0 {
		fmt.Fprintf(&b, "(%s) ", f.Columns)
	}
	if f.Overriding != nil {
		fmt.Fprintf(&b, "%s ", f.Overriding)
	}
	fmt.Fprintf(&b, "%s", f.SubQuery)

	return b.String()
}

type ColumnsConstructor struct {
	Columns ColumnNameList
}

func Columns(x ...ColumnName) *ColumnsConstructor { return &ColumnsConstructor{x} }

func (c *ColumnsConstructor) Values(x ...any) *FromConstructor {
	var values ValueConstructor

	for _, v := range x {
		values = append(values, newTypedRowValueExpr(v))
	}

	return &FromConstructor{Columns: c.Columns, Values: values}
}

type FromConstructor struct {
	Columns    ColumnNameList
	Overriding *OverridingClause
	Values     ValueConstructor
}

func Values(x ...any) *FromConstructor {
	var values ValueConstructor

	for _, v := range x {
		values = append(values, newTypedRowValueExpr(v))
	}

	return &FromConstructor{Values: values}
}

func (f *FromConstructor) insertFrom() InsertFrom { return f }

func (f *FromConstructor) OverridingSystemValue() *FromConstructor {
	o := OverridingSystemValue
	f.Overriding = &o
	return f
}

func (f *FromConstructor) OverridingUserValue() *FromConstructor {
	o := OverridingUserValue
	f.Overriding = &o
	return f
}

func (f *FromConstructor) String() string {
	var b strings.Builder

	if len(f.Columns) > 0 {
		fmt.Fprintf(&b, "(%s) ", f.Columns)
	}
	if f.Overriding != nil {
		fmt.Fprintf(&b, "%s ", f.Overriding)
	}
	fmt.Fprintf(&b, "%s", f.Values)

	return b.String()
}

type ValueConstructor []TypedRowValueExpr

func (c ValueConstructor) String() string {
	if len(c) == 1 {
		if _, ok := c[0].(rowsValue); ok {
			return fmt.Sprintf("VALUES%s", Join(c, ", "))
		}
	} else if len(c) > 1 {
		if _, ok := c[1].(rowValue); ok {
			return fmt.Sprintf("VALUES\n\t%s", Join(c, ",\n\t"))
		}
	}

	return fmt.Sprintf("VALUES (%s)", Join(c, ", "))
}

type FromDefault struct{}

var DefaultValues = &FromDefault{}

func (f *FromDefault) insertFrom() InsertFrom { return f }

func (f *FromDefault) String() string { return "DEFAULT VALUES" }
