// Code generated by "stringer -type=ControlTaskStatus"; DO NOT EDIT

package rest

import "fmt"

const _ControlTaskStatus_name = "ControlTaskStatusSuccessControlTaskStatusErrorControlTaskStatusCanceledControlTaskStatusIncomplete"

var _ControlTaskStatus_index = [...]uint8{0, 24, 46, 71, 98}

func (i ControlTaskStatus) String() string {
	if i < 0 || i >= ControlTaskStatus(len(_ControlTaskStatus_index)-1) {
		return fmt.Sprintf("ControlTaskStatus(%d)", i)
	}
	return _ControlTaskStatus_name[_ControlTaskStatus_index[i]:_ControlTaskStatus_index[i+1]]
}