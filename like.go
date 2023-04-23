package xql

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=LikeAction -linecomment

type LikeAction int

const (
	LikeExcluding LikeAction = iota // EXCLUDING
	LikeIncluding                   // INCLUDING
)

//go:generate stringer -type=LikeProperty -linecomment

type LikeProperty int

const (
	LikeComments    LikeProperty = iota // COMMENTS
	LikeCompression                     // COMPRESSION
	LikeConstraints                     // CONSTRAINTS
	LikeDefaults                        // DEFAULTS
	LikeGenerated                       // GENERATED
	LikeIdentity                        // IDENTITY
	LikeIndexes                         // INDEXES
	LikeStatistics                      // STATISTICS
	LikeStorage                         // STORAGE
	LikeAll                             // ALL
)

type LikeOption struct {
	Action   LikeAction
	Property LikeProperty
}

var (
	IncludingComments    = &LikeOption{LikeIncluding, LikeComments}
	ExcludingComments    = &LikeOption{LikeExcluding, LikeComments}
	IncludingCompression = &LikeOption{LikeIncluding, LikeCompression}
	ExcludingCompression = &LikeOption{LikeExcluding, LikeCompression}
	IncludingConstraints = &LikeOption{LikeIncluding, LikeConstraints}
	ExcludingConstraints = &LikeOption{LikeExcluding, LikeConstraints}
	IncludingDefaults    = &LikeOption{LikeIncluding, LikeDefaults}
	ExcludingDefaults    = &LikeOption{LikeExcluding, LikeDefaults}
	IncludingGenerated   = &LikeOption{LikeIncluding, LikeGenerated}
	ExcludingGenerated   = &LikeOption{LikeExcluding, LikeGenerated}
	IncludingIdentity    = &LikeOption{LikeIncluding, LikeIdentity}
	ExcludingIdentity    = &LikeOption{LikeExcluding, LikeIdentity}
	IncludingIndexes     = &LikeOption{LikeIncluding, LikeIndexes}
	ExcludingIndexes     = &LikeOption{LikeExcluding, LikeIndexes}
	IncludingStatistics  = &LikeOption{LikeIncluding, LikeStatistics}
	ExcludingStatistics  = &LikeOption{LikeExcluding, LikeStatistics}
	IncludingStorage     = &LikeOption{LikeIncluding, LikeStorage}
	ExcludingStorage     = &LikeOption{LikeExcluding, LikeStorage}
	IncludingAll         = &LikeOption{LikeIncluding, LikeAll}
	ExcludingAll         = &LikeOption{LikeExcluding, LikeAll}
)

func Excluding(p LikeProperty) *LikeOption { return &LikeOption{LikeExcluding, p} }
func Including(p LikeProperty) *LikeOption { return &LikeOption{LikeIncluding, p} }

func (o *LikeOption) String() string { return fmt.Sprintf("%s %s", o.Action, o.Property) }

type LikeClause struct {
	Name    TableName
	Options []*LikeOption
}

// Like returns a LIKE clause to create an empty table based on the definition of another table,
// including any column attributes and indexes defined in the original table:
func Like[T ToTableName](name T, x ...*LikeOption) *LikeClause {
	return &LikeClause{tableName(name), x}
}

func (c *LikeClause) tableElement() TableElement { return c }

func (c *LikeClause) applyTableDef(t *TableDef) {
	t.Content = TableElementList{c}
}

func (c *LikeClause) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "LIKE %s", c.Name)

	if len(c.Options) > 0 {
		fmt.Fprintf(&b, " %s", Join(c.Options, " "))
	}

	return b.String()
}
