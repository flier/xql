package xql

import (
	"fmt"
	"strings"
)

func Join[T fmt.Stringer](elems []T, sep string) string {
	var b strings.Builder

	for i, e := range elems {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(e.String())
	}

	return b.String()
}
