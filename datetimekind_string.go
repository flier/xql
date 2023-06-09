// Code generated by "stringer -type=DateTimeKind -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[KindDate-0]
	_ = x[KindSmallDateTime-1]
	_ = x[KindDateTime-2]
	_ = x[KindDateTime2-3]
	_ = x[KindYear-4]
	_ = x[KindTime-5]
	_ = x[KindTimestamp-6]
}

const _DateTimeKind_name = "DATESMALLDATETIMEDATETIMEDATETIME2YEARTIMETIMESTAMP"

var _DateTimeKind_index = [...]uint8{0, 4, 17, 25, 34, 38, 42, 51}

func (i DateTimeKind) String() string {
	if i < 0 || i >= DateTimeKind(len(_DateTimeKind_index)-1) {
		return "DateTimeKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DateTimeKind_name[_DateTimeKind_index[i]:_DateTimeKind_index[i+1]]
}
