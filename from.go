package xql

import "fmt"

type FromClause struct {
	Tables []*TableRef
}

func (f *FromClause) String() string { return fmt.Sprintf("FROM %s", Join(f.Tables, ", ")) }
