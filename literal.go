package xql

type Literal string

func (l Literal) Accept(v Visitor) Visitor { return v.Raw(string(l)) }
func (l Literal) String() string           { return string(l) }
