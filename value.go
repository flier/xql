package xql

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

type SimpleValue interface {
	Expr
}

type ValueExpr interface {
	Expr
}

type BoolValueExpr interface {
	Expr
}

type NumberValue interface {
	Expr
}

type UnsignedValue interface {
	Expr
}

type TypedRowValueExpr interface {
	Expr
}

var (
	_ TypedRowValueExpr = &nullValue{}
	_ TypedRowValueExpr = boolValue(false)
	_ TypedRowValueExpr = strValue("")
	_ TypedRowValueExpr = binValue(nil)
	_ TypedRowValueExpr = int8Value(0)
	_ TypedRowValueExpr = int16Value(0)
	_ TypedRowValueExpr = int32Value(0)
	_ TypedRowValueExpr = int64Value(0)
	_ TypedRowValueExpr = intValue(0)
	_ TypedRowValueExpr = floatValue(0)
	_ TypedRowValueExpr = &anyValue{}
	_ TypedRowValueExpr = &DefaultSpec{}
)

func typedRowValueExpr(value any) TypedRowValueExpr {
	if value == nil {
		return Nil
	}

	switch v := value.(type) {
	case bool:
		return boolValue(v)

	case string:
		return strValue(v)

	case []byte:
		return binValue(v)

	case int8:
		return int8Value(v)

	case int16:
		return int16Value(v)

	case int32:
		return int32Value(v)

	case int64:
		return int64Value(v)

	case int:
		return intValue(v)

	case float32:
		return floatValue(v)

	case float64:
		return floatValue(v)

	case Row:
		row := make(rowValue, len(v))

		for i, vv := range v {
			row[i] = typedRowValueExpr(vv)
		}

		return rowValue(row)

	case Rows:
		rows := make(rowsValue, len(v))

		for i, vv := range v {
			rows[i] = typedRowValueExpr(Row(vv)).(rowValue)
		}

		return rowsValue(rows)

	case *DefaultSpec:
		return v

	default:
		return &anyValue{v}
	}
}

type nullValue struct{}

var Nil = &nullValue{}

func (v nullValue) Expr() Expr     { return v }
func (v nullValue) String() string { return "NULL" }

type boolValue bool

func (v boolValue) Expr() Expr     { return v }
func (v boolValue) String() string { return strconv.FormatBool(bool(v)) }

type strValue string

func (v strValue) Expr() Expr     { return v }
func (v strValue) String() string { return strconv.Quote(string(v)) }

type binValue []byte

func (v binValue) Expr() Expr     { return v }
func (v binValue) String() string { return fmt.Sprintf("X'%s'", hex.EncodeToString(v)) }

type int8Value int8

func (v int8Value) Expr() Expr     { return v }
func (v int8Value) String() string { return strconv.Itoa(int(v)) }

type int16Value int16

func (v int16Value) Expr() Expr     { return v }
func (v int16Value) String() string { return strconv.Itoa(int(v)) }

type int32Value int32

func (v int32Value) Expr() Expr     { return v }
func (v int32Value) String() string { return strconv.Itoa(int(v)) }

type int64Value int64

func (v int64Value) Expr() Expr     { return v }
func (v int64Value) String() string { return strconv.FormatInt(int64(v), 10) }

type intValue int

func (v intValue) Expr() Expr     { return v }
func (v intValue) String() string { return strconv.Itoa(int(v)) }

type floatValue float64

func (v floatValue) Expr() Expr     { return v }
func (v floatValue) String() string { return strconv.FormatFloat(float64(v), 'g', -1, 64) }

type anyValue struct{ any }

func (v anyValue) Expr() Expr     { return v }
func (v anyValue) String() string { return fmt.Sprintf("%v", v.any) }

type Row []any

type rowValue []TypedRowValueExpr

func (v rowValue) Expr() Expr     { return v }
func (v rowValue) String() string { return fmt.Sprintf("ROW(%s)", Join(v, ", ")) }

type Rows [][]any

type rowsValue []rowValue

func (v rowsValue) Expr() Expr     { return v }
func (v rowsValue) String() string { return "\n\t" + Join(v, ",\n\t") }
