package xql

import (
	"fmt"
	"strings"
)

type FetchClause struct {
	Next     int
	Percent  bool
	WithTies bool
}

func (f *FetchClause) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "FETCH NEXT %d", f.Next)
	if f.Percent {
		b.WriteString(" PERCENT")
	}

	if f.Next > 1 {
		b.WriteString(" ROWS")
	} else {
		b.WriteString(" ROW")
	}

	if f.WithTies {
		b.WriteString(" WITH TIES")
	} else {
		b.WriteString(" ONLY")
	}

	return b.String()
}
