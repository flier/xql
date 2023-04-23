package xql

import "fmt"

type OffsetClause struct {
	RowCount SimpleValue
}

func (o *OffsetClause) String() string { return fmt.Sprintf("OFFSET %s ROWS", o.RowCount) }
