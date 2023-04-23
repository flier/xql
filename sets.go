package xql

//go:generate stringer -type SetQuantifier -linecomment

type SetQuantifier int

const (
	SetAll      SetQuantifier = iota // ALL
	SetDistinct                      // DISTINCT
)

//go:generate stringer -type SetOperation -linecomment

type SetOperation int

const (
	SetUnion      SetOperation = iota // UNION
	SetExceptions                     // EXCEPT
	SetIntersect                      // INTERSECT
)
