package util

import (
	_const "anto/const"
	"github.com/twharmon/gouid"
)

func Uid() string {
	return gouid.String(_const.GoUidLen, gouid.LowerCaseAlphaNum)
}
