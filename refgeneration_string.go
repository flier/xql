// Code generated by "stringer -type=RefGeneration -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[RefSystemGenerated-0]
	_ = x[RefUserGenerated-1]
	_ = x[RefDerived-2]
}

const _RefGeneration_name = "SYSTEM GENERATEDUSER GENERATEDDERIVED"

var _RefGeneration_index = [...]uint8{0, 16, 30, 37}

func (i RefGeneration) String() string {
	if i < 0 || i >= RefGeneration(len(_RefGeneration_index)-1) {
		return "RefGeneration(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RefGeneration_name[_RefGeneration_index[i]:_RefGeneration_index[i+1]]
}
