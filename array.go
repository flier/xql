package xql

import (
	"fmt"
	"strings"
)

type ArrayType struct {
	Type DataType
	Caps uint
}

var _ DataType = &ArrayType{}

type ArrayOption func(*ArrayType)

func Caps(n uint) ArrayOption { return func(t *ArrayType) { t.Caps = n } }

func ArrayOf(t DataType, x ...ArrayOption) *ArrayType {
	a := &ArrayType{t, 0}

	for _, opt := range x {
		opt(a)
	}

	return a
}

func (t *ArrayType) dataType() DataType          { return t }
func (t *ArrayType) applyColumnDef(d *ColumnDef) { d.Type = t }
func (t *ArrayType) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "%s ARRAY", t.Type)

	if t.Caps > 0 {
		fmt.Fprintf(&b, "[%d]", t.Caps)
	}

	return b.String()
}
