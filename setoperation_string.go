// Code generated by "stringer -type SetOperation -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SetUnion-0]
	_ = x[SetExceptions-1]
	_ = x[SetIntersect-2]
}

const _SetOperation_name = "UNIONEXCEPTINTERSECT"

var _SetOperation_index = [...]uint8{0, 5, 11, 20}

func (i SetOperation) String() string {
	if i < 0 || i >= SetOperation(len(_SetOperation_index)-1) {
		return "SetOperation(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SetOperation_name[_SetOperation_index[i]:_SetOperation_index[i+1]]
}
