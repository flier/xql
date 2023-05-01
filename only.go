package xql

import "fmt"

type OnlyClause struct {
	Table *TableName
}

func Only[T ToTableName](name T) *OnlyClause {
	return &OnlyClause{Table: newTableName(name)}
}

func (c *OnlyClause) targetTable() TargetTable { return c }
func (c *OnlyClause) applyDeleteStmt(s *DeleteStmt) {
	switch t := s.Target.(type) {
	case *TableName:
		s.Target = &OnlyClause{t}
	case *OnlyClause:
		if c.Table != nil {
			s.Target = c
		}
	}
}
func (c *OnlyClause) String() string {
	return fmt.Sprintf("ONLY %s", c.Table)
}
