package xql

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

type ValueExpr interface {
	Expr
}

type SimpleValue interface {
	ValueExpr
}

type BoolValueExpr interface {
	ValueExpr

	boolValueExpr() BoolValueExpr
}

type NumberValueExpr interface {
	ValueExpr

	numberValueExpr() NumberValueExpr
}

type UnsignedValueExpr interface {
	NumberValueExpr

	unsignedValueExpr() UnsignedValueExpr
}

type TypedRowValueExpr interface {
	ValueExpr
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

func newTypedRowValueExpr(value any) TypedRowValueExpr {
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

	case uint8:
		return uint8Value(v)

	case uint16:
		return uint16Value(v)

	case uint32:
		return uint32Value(v)

	case uint64:
		return uint64Value(v)

	case uint:
		return uintValue(v)

	case float32:
		return floatValue(v)

	case float64:
		return floatValue(v)

	case Row:
		row := make(rowValue, len(v))

		for i, vv := range v {
			row[i] = newTypedRowValueExpr(vv)
		}

		return rowValue(row)

	case Rows:
		rows := make(rowsValue, len(v))

		for i, vv := range v {
			rows[i] = newTypedRowValueExpr(Row(vv)).(rowValue)
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

func (v nullValue) expr() Expr     { return v }
func (v nullValue) String() string { return "NULL" }

type boolValue bool

func (v boolValue) expr() Expr                   { return v }
func (v boolValue) boolValueExpr() BoolValueExpr { return v }
func (v boolValue) String() string               { return strconv.FormatBool(bool(v)) }

type strValue string

func (v strValue) expr() Expr     { return v }
func (v strValue) String() string { return strconv.Quote(string(v)) }

type binValue []byte

func (v binValue) expr() Expr     { return v }
func (v binValue) String() string { return fmt.Sprintf("X'%s'", hex.EncodeToString(v)) }

type int8Value int8

func (v int8Value) expr() Expr                       { return v }
func (v int8Value) numberValueExpr() NumberValueExpr { return v }
func (v int8Value) String() string                   { return strconv.Itoa(int(v)) }

type int16Value int16

func (v int16Value) expr() Expr                       { return v }
func (v int16Value) numberValueExpr() NumberValueExpr { return v }
func (v int16Value) String() string                   { return strconv.Itoa(int(v)) }

type int32Value int32

func (v int32Value) expr() Expr                       { return v }
func (v int32Value) numberValueExpr() NumberValueExpr { return v }
func (v int32Value) String() string                   { return strconv.Itoa(int(v)) }

type int64Value int64

func (v int64Value) expr() Expr                       { return v }
func (v int64Value) numberValueExpr() NumberValueExpr { return v }
func (v int64Value) String() string                   { return strconv.FormatInt(int64(v), 10) }

type intValue int

func (v intValue) expr() Expr                       { return v }
func (v intValue) numberValueExpr() NumberValueExpr { return v }
func (v intValue) String() string                   { return strconv.Itoa(int(v)) }

type uint8Value uint8

func (v uint8Value) expr() Expr                           { return v }
func (v uint8Value) numberValueExpr() NumberValueExpr     { return v }
func (v uint8Value) unsignedValueExpr() UnsignedValueExpr { return v }
func (v uint8Value) String() string                       { return strconv.FormatUint(uint64(v), 10) }

type uint16Value uint16

func (v uint16Value) expr() Expr                           { return v }
func (v uint16Value) numberValueExpr() NumberValueExpr     { return v }
func (v uint16Value) unsignedValueExpr() UnsignedValueExpr { return v }
func (v uint16Value) String() string                       { return strconv.FormatUint(uint64(v), 10) }

type uint32Value uint32

func (v uint32Value) expr() Expr                           { return v }
func (v uint32Value) numberValueExpr() NumberValueExpr     { return v }
func (v uint32Value) unsignedValueExpr() UnsignedValueExpr { return v }
func (v uint32Value) String() string                       { return strconv.FormatUint(uint64(v), 10) }

type uint64Value uint64

func (v uint64Value) expr() Expr                           { return v }
func (v uint64Value) numberValueExpr() NumberValueExpr     { return v }
func (v uint64Value) unsignedValueExpr() UnsignedValueExpr { return v }
func (v uint64Value) String() string                       { return strconv.FormatUint(uint64(v), 10) }

type uintValue uint

func (v uintValue) expr() Expr                           { return v }
func (v uintValue) numberValueExpr() NumberValueExpr     { return v }
func (v uintValue) unsignedValueExpr() UnsignedValueExpr { return v }
func (v uintValue) String() string                       { return strconv.FormatUint(uint64(v), 10) }

type floatValue float64

func (v floatValue) expr() Expr                       { return v }
func (v floatValue) numberValueExpr() NumberValueExpr { return v }
func (v floatValue) String() string                   { return strconv.FormatFloat(float64(v), 'g', -1, 64) }

type anyValue struct{ any }

func (v anyValue) expr() Expr     { return v }
func (v anyValue) String() string { return fmt.Sprintf("%v", v.any) }

type Row []any

type rowValue []TypedRowValueExpr

func (v rowValue) expr() Expr     { return v }
func (v rowValue) String() string { return fmt.Sprintf("ROW(%s)", Join(v, ", ")) }

type Rows [][]any

type rowsValue []rowValue

func (v rowsValue) expr() Expr     { return v }
func (v rowsValue) String() string { return "\n\t" + Join(v, ",\n\t") }
