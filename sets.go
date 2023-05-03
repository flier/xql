package xql

//go:generate stringer -type SetQuantifier -linecomment

type SetQuantifier int

const (
	SetAll      SetQuantifier = iota // ALL
	SetDistinct                      // DISTINCT
)

var (
	Distinct = SetDistinct
)

func (q SetQuantifier) applySelectStmt(s *SelectStmt) { s.Quantifier = q }

//go:generate stringer -type SetOperation -linecomment

type SetOperation int

const (
	SetUnion      SetOperation = iota // UNION
	SetExceptions                     // EXCEPT
	SetIntersect                      // INTERSECT
)
