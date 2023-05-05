package xql

import "fmt"

type DefaultClause struct {
	Option DefaultOption
}

func (c *DefaultClause) columnValue() ColumnValue            { return c }
func (c *DefaultClause) applyColumnDef(d *ColumnDef)         { d.Value = c }
func (c *DefaultClause) applyColumnOptions(o *ColumnOptions) { o.Default = c }
func (c *DefaultClause) Accept(v Visitor) Visitor            { return v.Visit(Keyword("DEFAULT"), WS, c.Option) }
func (c *DefaultClause) String() string                      { return XQL(c) }

type AsDefault interface {
	AsDefault() *DefaultClause
}

type DefaultOption interface {
	fmt.Stringer

	Accepter

	AsDefault
}

var (
	_ DefaultOption = Literal("")
	_ DefaultOption = Raw("")
	_ DefaultOption = DefaultKind(0)
	_ DefaultOption = &DateTimeValueFunc{}
	_ DefaultOption = CreateDateTimeValueFunc(nil)
	_ DefaultOption = Null
	_ DefaultOption = EmptyArray
	_ DefaultOption = EmptyMultiSet
)

func (l Literal) AsDefault() *DefaultClause { return &DefaultClause{l} }
func (e Raw) AsDefault() *DefaultClause     { return &DefaultClause{e} }

//go:generate stringer -type=DefaultKind -linecomment

type DefaultKind int

const (
	DefaultUser           DefaultKind = iota // USER
	DefaultCurrentUser                       // CURRENT_USER
	DefaultCurrentRole                       // CURRENT_ROLE
	DefaultSessionUser                       // SESSION_USER
	DefaultSystemUser                        // SYSTEM_USER
	DefaultCurrentCatalog                    // CURRENT_CATALOG
	DefaultCurrentSchema                     // CURRENT_SCHEMA
	DefaultCurrentPath                       // CURRENT_PATH
)

func (k DefaultKind) AsDefault() *DefaultClause { return &DefaultClause{k} }
func (k DefaultKind) Accept(v Visitor) Visitor  { return v.Visit(Keyword(k.String())) }

var (
	User           = &DefaultClause{DefaultUser}
	CurrentUser    = &DefaultClause{DefaultCurrentUser}
	CurrentRole    = &DefaultClause{DefaultCurrentRole}
	SessionUser    = &DefaultClause{DefaultSessionUser}
	SystemUser     = &DefaultClause{DefaultSystemUser}
	CurrentCatalog = &DefaultClause{DefaultCurrentCatalog}
	CurrentSchema  = &DefaultClause{DefaultCurrentSchema}
	CurrentPath    = &DefaultClause{DefaultCurrentPath}
)

type NullSpec struct{}

var Null = &NullSpec{}

const kNull = Keyword("NULL")

func (s *NullSpec) AsDefault() *DefaultClause { return &DefaultClause{s} }
func (s *NullSpec) Accept(v Visitor) Visitor  { return v.Visit(kNull) }
func (s *NullSpec) String() string            { return XQL(s) }

type EmptyArraySpec struct{}

var EmptyArray = &EmptyArraySpec{}

func (s *EmptyArraySpec) AsDefault() *DefaultClause { return &DefaultClause{s} }
func (s *EmptyArraySpec) Accept(v Visitor) Visitor  { return v.Visit(kArray, Bracket) }
func (s *EmptyArraySpec) String() string            { return XQL(s) }

type EmptyMultiSetSpec struct{}

var EmptyMultiSet = &EmptyMultiSetSpec{}

const kMultiSet = Keyword("MULTISET")

func (s *EmptyMultiSetSpec) AsDefault() *DefaultClause { return &DefaultClause{s} }
func (s *EmptyMultiSetSpec) Accept(v Visitor) Visitor  { return v.Visit(kMultiSet, Bracket) }
func (s *EmptyMultiSetSpec) String() string            { return XQL(s) }

type DefaultSpec struct{}

var Default = &DefaultSpec{}

const kDefault = Keyword("DEFAULT")

func (s *DefaultSpec) expr() Expr               { return s }
func (s *DefaultSpec) Accept(v Visitor) Visitor { return v.Visit(kDefault) }
func (s *DefaultSpec) String() string           { return XQL(s) }
