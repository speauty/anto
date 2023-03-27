package task

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"translator/tst/tt_translator/ling_va"
	_type "translator/type"
)

func defaultTranslateObj() *Translate {
	obj := new(Translate)
	obj.SetTranslator(ling_va.GetInstance()).
		SetFromLang(ling_va.GetInstance().GetLangSupported()[1].Key).
		SetToLang(ling_va.GetInstance().GetLangSupported()[0].Key).
		SetTranslateMode(_type.ModeFull).
		SetMainTrackReport(_type.LangDirectionFrom).
		SetSrtFile("E:\\工作空间\\ts\\阅读和理解源码的最佳实践.srt")
	return obj
}

func TestRun(t *testing.T) {
	assertObj := assert.New(t)
	obj := defaultTranslateObj()
	msgs, err := obj.Run()
	fmt.Println(msgs)
	assertObj.Nil(err)
}
