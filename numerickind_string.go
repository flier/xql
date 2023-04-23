// Code generated by "stringer -type=NumericKind -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[KindNumeric-0]
	_ = x[KindDecimal-1]
	_ = x[KindDec-2]
	_ = x[KindFloat-3]
	_ = x[KindReal-4]
	_ = x[KindDoublePrecision-5]
	_ = x[KindDecFloat-6]
	_ = x[KindSmallSerial-7]
	_ = x[KindSerial-8]
	_ = x[KindBigSerial-9]
}

const _NumericKind_name = "NUMERICDECIMALDECFLOATREALDOUBLE PRECISIONDECFLOATSMALLSERIALSERIALBIGSERIAL"

var _NumericKind_index = [...]uint8{0, 7, 14, 17, 22, 26, 42, 50, 61, 67, 76}

func (i NumericKind) String() string {
	if i < 0 || i >= NumericKind(len(_NumericKind_index)-1) {
		return "NumericKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NumericKind_name[_NumericKind_index[i]:_NumericKind_index[i+1]]
}