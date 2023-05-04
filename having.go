package xql

import "fmt"

type HavingClause struct {
	Search SearchCond
}

func Having(cond SearchCond) *HavingClause { return &HavingClause{cond} }
func (h *HavingClause) String() string     { return fmt.Sprintf("HAVING %s", h.Search) }
