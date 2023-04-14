package ui

import "github.com/lxn/walk"

func (customG *Gui) eventCustomClose(canceled *bool, reason walk.CloseReason) {
	reason = walk.CloseReasonUser
	*canceled = false
	customG.ctxCancelFn()
}
