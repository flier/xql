package xql

import "fmt"

type ScopeClause struct {
	Table TableName
}

func (c *ScopeClause) String() string { return fmt.Sprintf("SCOPE %s", c.Table) }
