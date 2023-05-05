package xql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type AcceptFactoryFunc func(x ...Accepter) AcceptFunc

func (f AcceptFactoryFunc) Accept(v Visitor) Visitor {
	return f()(v)
}

var (
	Bracket = AcceptFactoryFunc(func(x ...Accepter) AcceptFunc {
		return func(v Visitor) Visitor {
			v.Token('[')
			if len(x) > 0 {
				v.Visit(x[0], x[1:]...)
			}
			return v.Token(']')
		}
	})

	Paren = AcceptFactoryFunc(func(x ...Accepter) AcceptFunc {
		return func(v Visitor) Visitor {
			v.Token('(')
			if len(x) > 0 {
				v.Visit(x[0], x[1:]...)
			}
			return v.Token(')')
		}
	})

	Sep = AcceptFunc(func(v Visitor) Visitor { return v.Sep().WS() })
	WS  = AcceptFunc(func(v Visitor) Visitor { return v.WS() })
)

func Token(t rune) AcceptFunc {
	return func(v Visitor) Visitor {
		return v.Token(t)
	}
}

func Ident(s fmt.Stringer) AcceptFunc {
	return func(v Visitor) Visitor {
		return v.Ident(s)
	}
}

func Int(n int) AcceptFunc {
	return func(v Visitor) Visitor {
		return v.Int(n)
	}
}

func Uint(n uint) AcceptFunc {
	return func(v Visitor) Visitor {
		return v.Uint(n)
	}
}

func Joins[T Accepter](s []T, sep AcceptFunc) Accepter {
	return AcceptFunc(func(v Visitor) Visitor {
		for i, a := range s {
			if i > 0 {
				v = sep.Accept(v)
			}
			v = a.Accept(v)
		}
		return v
	})
}

type Accepter interface {
	Accept(Visitor) Visitor
}

type AcceptFunc func(Visitor) Visitor

func (f AcceptFunc) Accept(v Visitor) Visitor {
	if f == nil {
		return v
	}

	return f(v)
}

type Visitor interface {
	WS() Visitor

	Sep() Visitor

	Token(tok rune) Visitor

	Int(n int) Visitor

	Uint(n uint) Visitor

	Float(n float64) Visitor

	Str(s string) Visitor

	Ident(s fmt.Stringer) Visitor

	Keyword(s fmt.Stringer) Visitor

	Raw(s string) Visitor

	DataType(dt DataType) Visitor

	Visit(a Accepter, x ...Accepter) Visitor

	If(cond bool, a Accepter, x ...Accepter) Visitor

	IfElse(cond bool, then Accepter, or Accepter) Visitor

	IfNotNil(cond any, a Accepter, x ...Accepter) Visitor
}

var _ Visitor = &Builder{}

type Builder struct {
	strings.Builder
	WhiteSpace rune
	Separator  rune
	Quote      rune
}

func XQL(a Accepter) string {
	b := NewBuilder()
	a.Accept(b)
	return b.String()
}

func NewBuilder() *Builder {
	return &Builder{
		WhiteSpace: ' ',
		Separator:  ',',
		Quote:      '`',
	}
}

func (b *Builder) WS() Visitor {
	b.WriteRune(b.WhiteSpace)
	return b
}

func (b *Builder) Sep() Visitor {
	b.WriteRune(b.Separator)
	return b
}

func (b *Builder) Token(tok rune) Visitor {
	b.WriteRune(tok)
	return b
}

func (b *Builder) Int(n int) Visitor {
	b.WriteString(strconv.Itoa(n))
	return b
}

func (b *Builder) Uint(n uint) Visitor {
	b.WriteString(strconv.FormatUint(uint64(n), 10))
	return b
}

func (b *Builder) Float(n float64) Visitor {
	b.WriteString(strconv.FormatFloat(n, 'f', -1, 64))
	return b
}

func (b *Builder) Str(s string) Visitor {
	b.WriteString(Quote(s, b.Quote, true))
	return b
}

func (b *Builder) Ident(s fmt.Stringer) Visitor {
	b.WriteString(EscapeName(s.String(), b.Quote))
	return b
}

func (b *Builder) Keyword(s fmt.Stringer) Visitor {
	return b.Visit(Keyword(s.String()))
}

func (b *Builder) Raw(s string) Visitor {
	b.WriteString(s)
	return b
}

func (b *Builder) DataType(dt DataType) Visitor {
	b.WriteString(dt.String())
	return b
}

func (b *Builder) Visit(a Accepter, x ...Accepter) Visitor {
	x = append([]Accepter{a}, x...)

	for _, a := range x {
		if isNil(a) {
			return b
		}
	}

	for _, a := range x {
		a.Accept(b)
	}

	return b
}

func isNil(i any) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)

	switch v.Kind() {
	case reflect.Func, reflect.Interface, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func (b *Builder) If(cond bool, a Accepter, x ...Accepter) Visitor {
	if cond {
		return b.Visit(a, x...)
	}

	return b
}

func (b *Builder) IfElse(cond bool, then Accepter, or Accepter) Visitor {
	if cond {
		return b.Visit(then)
	}

	return b.Visit(or)
}

func (b *Builder) IfNotNil(cond any, a Accepter, x ...Accepter) Visitor {
	if isNil(cond) {
		return b
	}

	return b.Visit(a, x...)
}
