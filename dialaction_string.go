// Code generated by "stringer -type DialAction"; DO NOT EDIT.

package hilib

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DialDisconnect-0]
	_ = x[DialConnect-1]
}

const _DialAction_name = "DialDisconnectDialConnect"

var _DialAction_index = [...]uint8{0, 14, 25}

func (i DialAction) String() string {
	if i < 0 || i >= DialAction(len(_DialAction_index)-1) {
		return "DialAction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DialAction_name[_DialAction_index[i]:_DialAction_index[i+1]]
}
