package xql

import (
	"fmt"
	"strings"
)

type CorrelationName = string

type UpdateStmt struct {
	Target TargetTable
	Alias  CorrelationName
	Sets   []SetClause
	Cursor *CursorName
	Search SearchCond
}

type UpdateOption interface {
	applyUpdateStmt(*UpdateStmt)
}

func Update[T ToTargetTable](target T, x ...UpdateOption) *UpdateStmt {
	s := &UpdateStmt{
		Target: toTargetTable(target),
	}

	for _, opt := range x {
		opt.applyUpdateStmt(s)
	}

	return s
}

func (s *UpdateStmt) As(name string) *UpdateStmt {
	s.Alias = name
	return s
}

func (s *UpdateStmt) Set(x ...SetClause) *UpdateStmt {
	s.Sets = append(s.Sets, x...)
	return s
}

func (s *UpdateStmt) Where(search SearchCond) *UpdateStmt {
	s.Search = search
	return s
}

func (s *UpdateStmt) WhereCurrentOf(cursor string) *UpdateStmt {
	s.Cursor = LocalQName(cursor)
	return s
}

func (s *UpdateStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "UPDATE %s", s.Target)

	if len(s.Alias) > 0 {
		fmt.Fprintf(&b, " AS %s", s.Alias)
	}

	fmt.Fprintf(&b, " SET %s", Join(s.Sets, ", "))

	if s.Cursor != nil {
		fmt.Fprintf(&b, " WHERE CURRENT OF %s", s.Cursor)
	} else if s.Search != nil {
		fmt.Fprintf(&b, " WHERE %s", s.Search)
	}

	return b.String()
}
