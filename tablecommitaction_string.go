// Code generated by "stringer -type=TableCommitAction -linecomment"; DO NOT EDIT.

package xql

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OnCommitPreserveRows-0]
	_ = x[OnCommitDeleteRows-1]
	_ = x[OnCommitDrop-2]
}

const _TableCommitAction_name = "PRESERVE ROWSDELETE ROWSDROP"

var _TableCommitAction_index = [...]uint8{0, 13, 24, 28}

func (i TableCommitAction) String() string {
	if i < 0 || i >= TableCommitAction(len(_TableCommitAction_index)-1) {
		return "TableCommitAction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TableCommitAction_name[_TableCommitAction_index[i]:_TableCommitAction_index[i+1]]
}