package xql

type Keyword string

func (k Keyword) Accept(v Visitor) Visitor { return v.Raw(string(k)) }
func (k Keyword) String() string           { return string(k) }
