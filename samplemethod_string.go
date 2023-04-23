// Code generated by "stringer -type SampleMethod -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SampleBernoulli-0]
	_ = x[SampleSystem-1]
}

const _SampleMethod_name = "BERNOULLISYSTEM"

var _SampleMethod_index = [...]uint8{0, 9, 15}

func (i SampleMethod) String() string {
	if i < 0 || i >= SampleMethod(len(_SampleMethod_index)-1) {
		return "SampleMethod(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SampleMethod_name[_SampleMethod_index[i]:_SampleMethod_index[i+1]]
}
