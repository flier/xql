package xql

type CollationName = LocalOrSchemaQualifiedName

type CollateClause struct {
	Name *CollationName
}

const kCollate = Keyword("COLLATE")

func (c *CollateClause) Accept(v Visitor) Visitor {
	return v.Visit(kCollate, WS, c.Name)
}

func (c *CollateClause) String() string { return XQL(c) }

type CollateOption interface {
	StringTypeOption
}

type collateOption struct {
	Collate *CollateClause
}

func (o *collateOption) applyStringType(t *StringType) {
	t.Collate = o.Collate
}

func Collate[T ToLocalOrSchemaQualifiedName](name T) CollateOption {
	return &collateOption{
		&CollateClause{Name: LocalOrSchemaQName(name)},
	}
}
