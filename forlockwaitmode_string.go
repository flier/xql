// Code generated by "stringer -type=ForLockWaitMode -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ForWait-0]
	_ = x[ForNoWait-1]
	_ = x[ForSkipLocked-2]
}

const _ForLockWaitMode_name = "WAITNO WAITSKIP LOCKED"

var _ForLockWaitMode_index = [...]uint8{0, 4, 11, 22}

func (i ForLockWaitMode) String() string {
	if i < 0 || i >= ForLockWaitMode(len(_ForLockWaitMode_index)-1) {
		return "ForLockWaitMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ForLockWaitMode_name[_ForLockWaitMode_index[i]:_ForLockWaitMode_index[i+1]]
}
