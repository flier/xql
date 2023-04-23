package xql

import "fmt"

type DefaultClause struct {
	Option DefaultOption
}

func (c *DefaultClause) columnValue() ColumnValue            { return c }
func (c *DefaultClause) applyColumnDef(d *ColumnDef)         { d.Value = c }
func (c *DefaultClause) applyColumnOptions(o *ColumnOptions) { o.Default = c }
func (c *DefaultClause) String() string                      { return fmt.Sprintf("DEFAULT %s", c.Option) }

type AsDefault interface {
	AsDefault() *DefaultClause
}

type DefaultOption interface {
	fmt.Stringer

	AsDefault
}

var (
	_ DefaultOption = Literal("")
	_ DefaultOption = rawExpr("")
	_ DefaultOption = DefaultKind(0)
	_ DefaultOption = &DateTimeValueFunc{}
	_ DefaultOption = CreateDateTimeValueFunc(nil)
	_ DefaultOption = Null
	_ DefaultOption = EmptyArray
	_ DefaultOption = EmptyMultiSet
)

func (l Literal) AsDefault() *DefaultClause { return &DefaultClause{l} }
func (e rawExpr) AsDefault() *DefaultClause { return &DefaultClause{e} }

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

func (s *NullSpec) AsDefault() *DefaultClause { return &DefaultClause{s} }
func (s *NullSpec) String() string            { return "NULL" }

type EmptyArraySpec struct{}

var EmptyArray = &EmptyArraySpec{}

func (s *EmptyArraySpec) AsDefault() *DefaultClause { return &DefaultClause{s} }
func (s *EmptyArraySpec) String() string            { return "ARRAY[]" }

type EmptyMultiSetSpec struct{}

var EmptyMultiSet = &EmptyMultiSetSpec{}

func (s *EmptyMultiSetSpec) AsDefault() *DefaultClause { return &DefaultClause{s} }
func (s *EmptyMultiSetSpec) String() string            { return "MULTISET[]" }

type DefaultSpec struct{}

var Default = &DefaultSpec{}

func (s *DefaultSpec) String() string { return "DEFAULT" }
