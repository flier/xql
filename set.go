package xql

import "fmt"

type SetClause struct {
	Assignment
}

type Assignment struct {
	Column ColumnName
	Value  ValueExpr
}

func (a *Assignment) String() string {
	return fmt.Sprintf("%s = %s", a.Column, a.Value)
}
