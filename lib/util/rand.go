package util

import (
	_const "anto/common"
	"github.com/twharmon/gouid"
)

func Uid() string {
	return gouid.String(_const.GoUidLen, gouid.LowerCaseAlphaNum)
}
