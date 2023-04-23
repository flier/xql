package xql

import (
	"fmt"
	"reflect"
)

type Either[L fmt.Stringer, R fmt.Stringer] struct {
	Left  L
	Right R
}

func Left[L fmt.Stringer, R fmt.Stringer](left L) Either[L, R] {
	return Either[L, R]{Left: left}
}

func Right[L fmt.Stringer, R fmt.Stringer](right R) Either[L, R] {
	return Either[L, R]{Right: right}
}

func (e *Either[L, R]) String() string {
	if !reflect.ValueOf(e.Left).IsZero() {
		return e.Left.String()
	}

	return e.Right.String()
}
