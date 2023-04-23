package xql

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=DateTimeValueKind -linecomment

type DateTimeValueKind int

const (
	KindCurrentDate      DateTimeValueKind = iota // CURRENT_DATE
	KindCurrentTime                               // CURRENT_TIME
	KindLocalTime                                 // LOCALTIME
	KindCurrentTimestamp                          // CURRENT_TIMESTAMP
	KindLocalTimestamp                            // LOCALTIMESTAMP
)

type DateTimeValueFunc struct {
	Kind      DateTimeValueKind
	Precision uint
}

var (
	CurrentDate      = &DateTimeValueFunc{KindCurrentDate, 0}
	CurrentTime      = dateTimeValueFunc(KindCurrentTime)
	LocalTime        = dateTimeValueFunc(KindLocalTime)
	CurrentTimestamp = dateTimeValueFunc(KindCurrentTimestamp)
	LocalTimestamp   = dateTimeValueFunc(KindLocalTimestamp)
)

type CreateDateTimeValueFunc func(precision uint) *DateTimeValueFunc

func (f CreateDateTimeValueFunc) AsDefault() *DefaultClause { return &DefaultClause{f(0)} }
func (f CreateDateTimeValueFunc) String() string            { return f(0).String() }

func dateTimeValueFunc(kind DateTimeValueKind) CreateDateTimeValueFunc {
	return func(precision uint) *DateTimeValueFunc {
		return &DateTimeValueFunc{kind, precision}
	}
}

func (f *DateTimeValueFunc) AsDefault() *DefaultClause { return &DefaultClause{f} }
func (f *DateTimeValueFunc) String() string {
	var b strings.Builder

	b.WriteString(f.Kind.String())

	if f.Precision > 0 {
		fmt.Fprintf(&b, "(%d)", f.Precision)
	}

	return b.String()
}
