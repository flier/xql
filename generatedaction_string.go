// Code generated by "stringer -type=GeneratedAction -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[GeneratedAlways-0]
	_ = x[GeneratedByDefault-1]
}

const _GeneratedAction_name = "ALWAYSBY DEFAULT"

var _GeneratedAction_index = [...]uint8{0, 6, 16}

func (i GeneratedAction) String() string {
	if i < 0 || i >= GeneratedAction(len(_GeneratedAction_index)-1) {
		return "GeneratedAction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _GeneratedAction_name[_GeneratedAction_index[i]:_GeneratedAction_index[i+1]]
}