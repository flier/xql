package xql

import "fmt"

type FromClause TableRefList

func From(x ...ToTableRef) FromClause {
	var refs []TableRef

	for _, r := range x {
		refs = append(refs, r.tableRef())
	}

	return FromClause(refs)
}

func (c FromClause) String() string { return fmt.Sprintf("FROM %s", TableRefList(c)) }
