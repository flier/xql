// Code generated by "stringer -type=DropBehavior -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DropRestrict-0]
	_ = x[DropCascade-1]
}

const _DropBehavior_name = "RESTRICTCASCADE"

var _DropBehavior_index = [...]uint8{0, 8, 15}

func (i DropBehavior) String() string {
	if i < 0 || i >= DropBehavior(len(_DropBehavior_index)-1) {
		return "DropBehavior(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DropBehavior_name[_DropBehavior_index[i]:_DropBehavior_index[i+1]]
}
