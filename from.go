package xql

import "fmt"

type FromClause TableRefList

func (c FromClause) String() string { return fmt.Sprintf("FROM %s", TableRefList(c)) }
