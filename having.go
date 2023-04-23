package xql

import "fmt"

type HavingClause struct {
	Search SearchCond
}

func (h *HavingClause) String() string { return fmt.Sprintf("HAVING %s", h.Search) }
