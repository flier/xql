// Code generated by "stringer -type OffsetSuffix -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SuffixRows-0]
	_ = x[SuffixRow-1]
}

const _OffsetSuffix_name = "ROWSROW"

var _OffsetSuffix_index = [...]uint8{0, 4, 7}

func (i OffsetSuffix) String() string {
	if i < 0 || i >= OffsetSuffix(len(_OffsetSuffix_index)-1) {
		return "OffsetSuffix(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OffsetSuffix_name[_OffsetSuffix_index[i]:_OffsetSuffix_index[i+1]]
}
