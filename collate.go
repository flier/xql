package xql

import "fmt"

type CollationName = SchemaQualifiedName

type CollateClause struct {
	Name CollationName
}

func (c *CollateClause) String() string {
	return fmt.Sprintf("COLLATE %s", &c.Name)
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

func Collate(name CollationName) CollateOption {
	return &collateOption{
		&CollateClause{Name: name},
	}
}
