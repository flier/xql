package xql

type ArrayType struct {
	Type DataType
	Caps uint
}

var _ DataType = &ArrayType{}

type CreateArrayFunc func(uint) *ArrayType

func (f CreateArrayFunc) dataType() DataType          { return f }
func (f CreateArrayFunc) applyColumnDef(d *ColumnDef) { d.Type = f }
func (f CreateArrayFunc) Accept(v Visitor) Visitor    { return f(0).Accept(v) }
func (f CreateArrayFunc) String() string              { return XQL(f) }

func ArrayOf(t DataType) CreateArrayFunc {
	return func(caps uint) *ArrayType {
		return &ArrayType{t, caps}
	}
}

func (t *ArrayType) dataType() DataType          { return t }
func (t *ArrayType) applyColumnDef(d *ColumnDef) { d.Type = t }

const kArray = Keyword("ARRAY")

func (t *ArrayType) Accept(v Visitor) Visitor {
	return v.DataType(t.Type).WS().Visit(kArray).
		IfElse(t.Caps > 0, Bracket(Uint(t.Caps)), Bracket)
}

func (t *ArrayType) String() string { return XQL(t) }
