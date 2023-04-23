package xql

import (
	"fmt"
	"strings"
)

type IDColumnRestart int

const (
	IDContinue IDColumnRestart = iota // CONTINUE IDENTITY
	IDRestart                         // RESTART IDENTITY
)

type TruncateStmt struct {
	Table   TargetTable
	Restart IDColumnRestart
}

func (s *TruncateStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "TRUNCATE TABLE %s", s.Table)
	if s.Restart != IDContinue {
		fmt.Fprintf(&b, " RESTART IDENTITY")
	}

	return b.String()
}
