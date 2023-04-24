package xql

import "fmt"

type TargetTable interface {
	fmt.Stringer

	targetTable() TargetTable
}

var (
	_ TargetTable = &TableName{}
	_ TargetTable = &OnlyClause{}
)

type ToTargetTable interface {
	~string | *LocalQualifiedName | *SchemaQualifiedName | *TableName | *OnlyClause
}

func toTargetTable[T ToTargetTable](name T) TargetTable {
	switch v := any(name).(type) {
	case string:
		return tableName(v)
	case *LocalQualifiedName:
		return tableName(v)
	case *SchemaQualifiedName:
		return tableName(v)
	case *TableName:
		return v
	case *OnlyClause:
		return v
	default:
		panic("unreachable")
	}
}
