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
	*SchemaName
	Name string
}

func (n *SchemaQualifiedName) LocalOrSchemaQName() *LocalOrSchemaQualifiedName {
	e := LocalOrSchemaQualifiedName(Right[*LocalQualifiedName](n))
	return &e
}

type ToSchemaQualifiedName interface {
	~string | *SchemaQualifiedName
}

func SchemaQName[T ToSchemaQualifiedName](name T) *SchemaQualifiedName {
	switch v := any(name).(type) {
	case string:
		return &SchemaQualifiedName{Name: v}
	case *SchemaQualifiedName:
		return v
	default:
		panic("unreachable")
	}
}

func (n *SchemaQualifiedName) String() string {
	if n.SchemaName != nil {
		return fmt.Sprintf("%s.%s", n.SchemaName, n.Name)
	}

	return n.Name
}

type LocalOrSchemaQualifiedName Either[*LocalQualifiedName, *SchemaQualifiedName]

func (n *LocalOrSchemaQualifiedName) String() string {
	return (*Either[*LocalQualifiedName, *SchemaQualifiedName])(n).String()
}

type ToLocalOrSchemaQualifiedName interface {
	~string | *LocalQualifiedName | *SchemaQualifiedName
}

func LocalOrSchemaQName[T ToLocalOrSchemaQualifiedName](name T) *LocalOrSchemaQualifiedName {
	switch v := any(name).(type) {
	case string:
		return SchemaQName(v).LocalOrSchemaQName()
	case *SchemaQualifiedName:
		return v.LocalOrSchemaQName()
	case *LocalQualifiedName:
		return v.LocalOrSchemaQName()
	default:
		panic("unreachable")
	}
}

type LocalQualifier struct{}

var Local = &LocalQualifier{}

func (l *LocalQualifier) Name(name string) *LocalQualifiedName { return &LocalQualifiedName{l, name} }
func (q LocalQualifier) String() string                        { return "MODULE" }

type LocalQualifiedName struct {
	*LocalQualifier
	Name string
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

func (n *LocalQualifiedName) WithQualifier() *LocalQualifiedName {
	n.LocalQualifier = &LocalQualifier{}
	return n
}

func (l *LocalQualifiedName) LocalOrSchemaQName() *LocalOrSchemaQualifiedName {
	n := LocalOrSchemaQualifiedName(Left[*LocalQualifiedName, *SchemaQualifiedName](l))
	return &n
}

func (n *LocalQualifiedName) String() string {
	if n.LocalQualifier != nil {
		return fmt.Sprintf("%s.%s", n.LocalQualifier, n.Name)
	}

	return n.Name
}
