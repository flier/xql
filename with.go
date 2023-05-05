package xql

type WithClause struct {
	Recursive bool
}

var (
	With          = &WithClause{}
	WithRecursive = &WithClause{true}
)

const (
	kWith      = Keyword("WITH")
	kRecursive = Keyword("RECURSIVE")
)

func (w *WithClause) Accept(v Visitor) Visitor {
	return v.Visit(kWith).If(w.Recursive, WS, kRecursive)
}

func (w *WithClause) String() string { return XQL(w) }
