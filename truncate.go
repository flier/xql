package xql

import (
	"fmt"
	"strings"
)

//go:generate stringer -type IdentityColumnRestart -linecomment

type IdentityColumnRestart int

const (
	// Do not change the values of sequences. This is the default.
	ContinueIdentity IdentityColumnRestart = iota // CONTINUE IDENTITY
	// Automatically restart sequences owned by columns of the truncated table(s).
	RestartIdentity // RESTART IDENTITY
)

func (i IdentityColumnRestart) applyTruncateStmt(s *TruncateStmt) { s.Restart = &i }

//go:generate stringer -type=DropBehavior -linecomment

type DropBehavior int

const (
	// Refuse to truncate if any of the tables have foreign-key references from tables that are not listed in the command.
	// This is the default.
	DropRestrict DropBehavior = iota // RESTRICT
	// Automatically truncate all tables that have foreign-key references to any of the named tables,
	// or to any tables added to the group due to CASCADE.
	DropCascade // CASCADE
)

func (b DropBehavior) applyTruncateStmt(s *TruncateStmt) { s.Drop = &b }

// TruncateTable quickly removes all rows from a set of tables.
//
// It has the same effect as an unqualified DELETE on each table,
// but since it does not actually scan the tables it is faster.
// This is most useful on large tables.
type TruncateStmt struct {
	Targets []TargetTable
	Restart *IdentityColumnRestart
	Drop    *DropBehavior
}

// TruncateTable empty a table or set of tables.
func TruncateTable[T ToTargetTable](x ...T) *TruncateStmt {
	var targets []TargetTable

	for _, t := range x {
		targets = append(targets, toTargetTable(t))
	}

	s := &TruncateStmt{
		Targets: targets,
	}

	return s
}

// Do not change the values of sequences. This is the default.
func (s *TruncateStmt) ContinueIdentity() *TruncateStmt {
	ContinueIdentity.applyTruncateStmt(s)
	return s
}

// Automatically restart sequences owned by columns of the truncated table(s).
func (s *TruncateStmt) RestartIdentity() *TruncateStmt {
	RestartIdentity.applyTruncateStmt(s)
	return s
}

// Automatically truncate all tables that have foreign-key references to any of the named tables,
// or to any tables added to the group due to CASCADE.
func (s *TruncateStmt) Cascade() *TruncateStmt {
	DropCascade.applyTruncateStmt(s)
	return s
}

// Refuse to truncate if any of the tables have foreign-key references from tables that are not listed in the command.
// This is the default.
func (s *TruncateStmt) Restrict() *TruncateStmt {
	DropRestrict.applyTruncateStmt(s)
	return s
}

func (s *TruncateStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "TRUNCATE TABLE %s", Join(s.Targets, ", "))

	if s.Restart != nil {
		fmt.Fprintf(&b, " %s", s.Restart)
	}

	return b.String()
}
