package xql

import (
	"fmt"
	"strings"
)

type MergeCorrelationName = CorrelationName

type MergeStmt struct {
	Target TargetTable
	Alias  MergeCorrelationName
	Source TableRef
	Join   SearchCond
	When   MergeWhenClause
}

type MergeIntoClause struct {
	s *MergeStmt
}

func MergeInto[T ToTargetTable](target T) *MergeIntoClause {
	return &MergeIntoClause{&MergeStmt{
		Target: toTargetTable(target),
	}}
}

func (c *MergeIntoClause) As(alias MergeCorrelationName) *MergeIntoClause {
	c.s.Alias = alias
	return c
}

func (c *MergeIntoClause) Using(source TableRef) *MergeIntoUsingClause {
	c.s.Source = source
	return &MergeIntoUsingClause{c}
}

type MergeIntoUsingClause struct {
	*MergeIntoClause
}

func (c *MergeIntoUsingClause) On(join SearchCond) *MergeStmt {
	c.s.Join = join
	return c.s
}

func (s *MergeStmt) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "MERGE INTO %s", s.Target)

	if len(s.Alias) > 0 {
		fmt.Fprintf(&b, " AS %s", s.Alias)
	}

	fmt.Fprintf(&b, " USING %s ON %s", s.Source, s.Join)

	return b.String()
}

type MergeWhenClause interface {
	mergeWhenClause() MergeWhenClause
}

var (
	_ MergeWhenClause = &MergeWhenMatchedClause{}
	_ MergeWhenClause = &MergeWhenNotMatchedClause{}
)

type MergeWhenMatchedClause struct {
	Cond           SearchCond
	UpdateOrDelete MergeUpdateOrDeleteSpec
}

var WhenMatched = &MergeWhenMatchedClause{}

func (c *MergeWhenMatchedClause) And(cond SearchCond) *MergeWhenMatchedClause {
	c.Cond = cond
	return c
}

func (c *MergeWhenMatchedClause) ThenUpdate(x ...SetClause) *MergeWhenMatchedClause {
	c.UpdateOrDelete = MergeUpdateSpec(x)
	return c
}

func (c *MergeWhenMatchedClause) ThenDelete() *MergeWhenMatchedClause {
	c.UpdateOrDelete = &MergeDeleteSpec{}
	return c
}

func (c *MergeWhenMatchedClause) ThenDoNothing() *MergeWhenMatchedClause {
	c.UpdateOrDelete = nil
	return c
}

func (c *MergeWhenMatchedClause) mergeWhenClause() MergeWhenClause { return c }
func (c *MergeWhenMatchedClause) String() string {
	var b strings.Builder

	b.WriteString("WHEN MATCHED")

	if c.Cond != nil {
		fmt.Fprintf(&b, " AND %s", c.Cond)
	}

	if c.UpdateOrDelete != nil {
		fmt.Fprintf(&b, " THEN %s", c.UpdateOrDelete)
	} else {
		b.WriteString(" THEN DO NOTHING")
	}

	return b.String()
}

type MergeWhenNotMatchedClause struct {
	Cond   SearchCond
	Insert *MergeInsertSpec
}

var WhenNotMatched = &MergeWhenNotMatchedClause{}

func (c *MergeWhenNotMatchedClause) And(cond SearchCond) *MergeWhenNotMatchedClause {
	c.Cond = cond
	return c
}

func (c *MergeWhenNotMatchedClause) ThenInsert(from *FromConstructor) *MergeWhenNotMatchedClause {
	c.Insert = (*MergeInsertSpec)(from)
	return c
}

func (c *MergeWhenNotMatchedClause) ThenDoNothing() *MergeWhenNotMatchedClause {
	c.Insert = nil
	return c
}

func (c *MergeWhenNotMatchedClause) mergeWhenClause() MergeWhenClause { return c }

func (c *MergeWhenNotMatchedClause) String() string {
	var b strings.Builder

	b.WriteString("WHEN NOT MATCHED")

	if c.Cond != nil {
		fmt.Fprintf(&b, " AND %s", c.Cond)
	}

	if c.Insert != nil {
		fmt.Fprintf(&b, " THEN %s", c.Insert)
	} else {
		b.WriteString(" THEN DO NOTHING")
	}

	return b.String()
}

type MergeUpdateOrDeleteSpec interface {
	fmt.Stringer

	mergeUpdateOrDeleteSpec() MergeUpdateOrDeleteSpec
}

var (
	_ MergeUpdateOrDeleteSpec = MergeUpdateSpec(nil)
	_ MergeUpdateOrDeleteSpec = &MergeDeleteSpec{}
)

type MergeUpdateSpec SetClauseList

func (s MergeUpdateSpec) mergeUpdateOrDeleteSpec() MergeUpdateOrDeleteSpec { return s }
func (s MergeUpdateSpec) String() string                                   { return fmt.Sprintf("UPDATE SET %s", SetClauseList(s)) }

type MergeDeleteSpec struct{}

func (s *MergeDeleteSpec) mergeUpdateOrDeleteSpec() MergeUpdateOrDeleteSpec { return s }
func (s *MergeDeleteSpec) String() string                                   { return "DELETE" }

type MergeInsertSpec FromConstructor

func (s *MergeInsertSpec) String() string { return fmt.Sprintf("INSERT %s", (*FromConstructor)(s)) }
