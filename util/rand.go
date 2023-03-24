package util

import (
	"github.com/twharmon/gouid"
	_const "translator/const"
)

func Uid() string {
	return gouid.String(_const.GoUidLen, gouid.LowerCaseAlphaNum)
}
