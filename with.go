package xql

import "strings"

type WithClause struct {
	Recursive bool
}

func (w *WithClause) String() string {
	var b strings.Builder

	b.WriteString("WITH ")
	if w.Recursive {
		b.WriteString("RECURSIVE ")
	}

	return b.String()
}
