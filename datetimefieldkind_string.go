// Code generated by "stringer -type=DateTimeFieldKind -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FieldYear-0]
	_ = x[FieldMonth-1]
	_ = x[FieldWeek-2]
	_ = x[FieldDay-3]
	_ = x[FieldHour-4]
	_ = x[FieldMinute-5]
	_ = x[FieldSecond-6]
	_ = x[FieldMillisecond-7]
	_ = x[FieldMicrosecond-8]
}

const _DateTimeFieldKind_name = "YEARMONTHWEEKDAYHOURMINUTESECONDMILLISECONDMICROSECOND"

var _DateTimeFieldKind_index = [...]uint8{0, 4, 9, 13, 16, 20, 26, 32, 43, 54}

func (i DateTimeFieldKind) String() string {
	if i < 0 || i >= DateTimeFieldKind(len(_DateTimeFieldKind_index)-1) {
		return "DateTimeFieldKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DateTimeFieldKind_name[_DateTimeFieldKind_index[i]:_DateTimeFieldKind_index[i+1]]
}
