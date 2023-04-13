package msg

import "github.com/lxn/walk"

const msgTitle = "提示"

func Info(owner walk.Form, msg string) int {
	return walk.MsgBox(owner, msgTitle, msg, walk.MsgBoxOK|walk.MsgBoxIconInformation)
}

func Warn(owner walk.Form, msg string) int {
	return walk.MsgBox(owner, msgTitle, msg, walk.MsgBoxOK|walk.MsgBoxIconWarning)
}

func Err(owner walk.Form, err error) int {
	return walk.MsgBox(owner, msgTitle, err.Error(), walk.MsgBoxOK|walk.MsgBoxIconError)
}

/* 这个。。。似乎意思不大
func Help(owner walk.Form, msg string) int {
	return walk.MsgBox(owner, msgTitle, msg, walk.MsgBoxHelp|walk.MsgBoxIconExclamation)
}
*/

func Confirm(owner walk.Form, msg string) (bool, bool) {
	ret := walk.MsgBox(owner, msgTitle, msg, walk.MsgBoxYesNo|walk.MsgBoxIconExclamation)
	return ret == 6, ret == 7
}
