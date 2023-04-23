package xql

import (
	"fmt"
)

const NameDataLen = 64

type Catalog string

func (c Catalog) Schema(schema string) *SchemaName {
	return &SchemaName{string(c), schema}
}

type SchemaName struct {
	Catalog string
	Schema  string
}

func Schema(name string) *SchemaName { return &SchemaName{Schema: name} }

func (n *SchemaName) QName(name string) *SchemaQualifiedName {
	return &SchemaQualifiedName{n, name}
}

func (n *SchemaName) String() string {
	if len(n.Catalog) == 0 {
		return n.Schema
	}

	return fmt.Sprintf("%s.%s", n.Catalog, n.Schema)
}

type SchemaQualifiedName struct {
	schema *SchemaName
	name   string
}

func (n *SchemaQualifiedName) Name() string                      { return n.name }
func (n *SchemaQualifiedName) IsLocal() bool                     { return false }
func (n *SchemaQualifiedName) SchemaName() *SchemaName           { return n.schema }
func (n *SchemaQualifiedName) LocalQualifier() *LocalQualifier   { return nil }
func (n *SchemaQualifiedName) SchemaQName() *SchemaQualifiedName { return n }
func (n *SchemaQualifiedName) LocalQName() *LocalQualifiedName   { return nil }

type ToSchemaQualifiedName interface {
	~string | *SchemaQualifiedName
}

func SchemaQName[T ToSchemaQualifiedName](name T) *SchemaQualifiedName {
	switch v := any(name).(type) {
	case string:
		return &SchemaQualifiedName{name: v}
	case *SchemaQualifiedName:
		return v
	default:
		panic("unreachable")
	}
}

func (n *SchemaQualifiedName) String() string {
	if n.schema != nil {
		return fmt.Sprintf("%s.%s", n.schema, n.name)
	}

	return n.name
}

type LocalOrSchemaQualifiedName interface {
	fmt.Stringer

	SchemaName() *SchemaName

	LocalQualifier() *LocalQualifier

	IsLocal() bool

	SchemaQName() *SchemaQualifiedName

	LocalQName() *LocalQualifiedName

	Name() string
}

var (
	_ LocalOrSchemaQualifiedName = &SchemaQualifiedName{}
	_ LocalOrSchemaQualifiedName = &LocalQualifiedName{}
)

type ToLocalOrSchemaQualifiedName interface {
	~string | *LocalQualifiedName | *SchemaQualifiedName
}

func LocalOrSchemaQName[T ToLocalOrSchemaQualifiedName](name T) LocalOrSchemaQualifiedName {
	switch v := any(name).(type) {
	case string:
		return SchemaQName(v)
	case *SchemaQualifiedName:
		return v
	case *LocalQualifiedName:
		return v
	default:
		panic("unreachable")
	}
}

type LocalQualifier struct{}

var Local = &LocalQualifier{}

func (l *LocalQualifier) Name(name string) *LocalQualifiedName { return &LocalQualifiedName{l, name} }

func (q LocalQualifier) String() string { return "MODULE" }

type LocalQualifiedName struct {
	qualifier *LocalQualifier
	name      string
}

type ToLocalQualifiedName interface {
	~string | *LocalQualifiedName
}

func LocalQName[T ToLocalQualifiedName](name T) *LocalQualifiedName {
	switch v := any(name).(type) {
	case string:
		return &LocalQualifiedName{nil, v}
	case *LocalQualifiedName:
		return v
	default:
		panic("unreachable")
	}
}

func (n *LocalQualifiedName) Name() string                      { return n.name }
func (n *LocalQualifiedName) IsLocal() bool                     { return n.qualifier != nil }
func (n *LocalQualifiedName) LocalQualifier() *LocalQualifier   { return n.qualifier }
func (n *LocalQualifiedName) SchemaName() *SchemaName           { return nil }
func (n *LocalQualifiedName) SchemaQName() *SchemaQualifiedName { return nil }
func (n *LocalQualifiedName) LocalQName() *LocalQualifiedName   { return n }
func (n *LocalQualifiedName) WithQualifier() *LocalQualifiedName {
	n.qualifier = &LocalQualifier{}
	return n
}

func (n *LocalQualifiedName) String() string {
	if n.qualifier != nil {
		return fmt.Sprintf("%s.%s", n.qualifier, n.name)
	}

	return n.name
}
