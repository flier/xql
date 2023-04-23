package xql

import "fmt"

type LimitClause struct {
	RowCount SimpleValue
}

func (o *LimitClause) String() string { return fmt.Sprintf("LIMIT %s", o.RowCount) }
