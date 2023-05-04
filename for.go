package xql

import (
	"fmt"
	"strings"
	"time"
)

//go:generate stringer -type=ForLockMode -linecomment

type ForLockMode int

const (
	ForUpdate      ForLockMode = iota // UPDATE
	ForNoKeyUpdate                    // NO KEY UPDATE
	ForShare                          // SHARE
	ForKeyShare                       // KEY SHARE
)

//go:generate stringer -type=ForLockWaitMode -linecomment

type ForLockWaitMode int

const (
	ForWait       ForLockWaitMode = iota // WAIT
	ForNoWait                            // NO WAIT
	ForSkipLocked                        // SKIP LOCKED
)

type ForLockClause struct {
	Mode     ForLockMode
	WaitMode *ForLockWaitMode
	WaitTime time.Duration
}

func (c *ForLockClause) Wait(d time.Duration) *ForLockClause {
	m := ForWait
	c.WaitMode = &m
	c.WaitTime = d
	return c
}

func (c *ForLockClause) NoWait() *ForLockClause {
	m := ForNoWait
	c.WaitMode = &m
	return c
}

func (c *ForLockClause) SkipLocked() *ForLockClause {
	m := ForSkipLocked
	c.WaitMode = &m
	return c
}

func (c *ForLockClause) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "FOR %s", c.Mode)

	if c.WaitMode != nil {
		fmt.Fprintf(&b, " %s", c.WaitMode)
		if *c.WaitMode == ForWait {
			fmt.Fprintf(&b, " %s", c.WaitTime)
		}
	}

	return b.String()
}
