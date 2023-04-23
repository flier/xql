package xql

import (
	"fmt"
	"strings"
)

type DeleteStmt struct {
	Target TargetTable
	Alias  CorrelationName
	Cursor *CursorName
	Search SearchCond
}

type DeleteStmtOption interface {
	applyDeleteStmt(*DeleteStmt)
}

type ToTargetTable interface {
	~string | *LocalQualifiedName | *SchemaQualifiedName | *OnlyClause
}

func DeleteFrom[T ToTargetTable](name T, x ...DeleteStmtOption) *DeleteStmt {
	var target TargetTable

	switch v := any(name).(type) {
	case *OnlyClause:
		target = v
	case string:
		target = tableName(v)
	case *LocalQualifiedName:
		target = tableName(v)
	case *SchemaQualifiedName:
		target = tableName(v)
	}

	s := &DeleteStmt{
		Target: target,
		Search: nil,
	}

	for _, opt := range x {
		opt.applyDeleteStmt(s)
	}

	return s
}

func (s *DeleteStmt) As(name string) *DeleteStmt {
	s.Alias = name
	return s
}

func (s *DeleteStmt) Where(search SearchCond) *DeleteStmt {
	s.Search = search
	return s
}

func (s *DeleteStmt) WhereCurrentOf(cursor string) *DeleteStmt {
	s.Cursor = LocalQName(cursor)
	return s
}

func (s *DeleteStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "DELETE FROM %s", s.Target)

	if len(s.Alias) > 0 {
		fmt.Fprintf(&b, " AS %s", s.Alias)
	}

	if s.Cursor != nil {
		fmt.Fprintf(&b, " WHERE CURRENT OF %s", s.Cursor)
	} else if s.Search != nil {
		fmt.Fprintf(&b, " WHERE %s", s.Search)
	}

	return b.String()
}
