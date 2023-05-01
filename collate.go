package xql

import "fmt"

type CollationName = LocalOrSchemaQualifiedName

type CollateClause struct {
	Name *CollationName
}

func (c *CollateClause) String() string {
	return fmt.Sprintf("COLLATE %s", c.Name)
}

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
