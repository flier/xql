package xql

import (
	"fmt"
	"strings"
)

type CorrelationName = string

type UpdateStmt struct {
	Table  TargetTable
	As     CorrelationName
	Sets   []*SetClause
	Cursor CursorName
	Search SearchCond
}

func (s *UpdateStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "UPDATE %s", s.Table)

	if len(s.As) > 0 {
		fmt.Fprintf(&b, " AS %s", s.As)
	}

	fmt.Fprintf(&b, " SET %s", Join(s.Sets, ", "))

	if len(s.Cursor.Name()) > 0 {
		fmt.Fprintf(&b, " WHERE CURRENT OF %s", &s.Cursor)
	} else if s.Search != nil {
		fmt.Fprintf(&b, " WHERE %s", s.Search)
	}

	return b.String()
}

type TargetTable interface {
	fmt.Stringer

	targetTable() TargetTable
}
