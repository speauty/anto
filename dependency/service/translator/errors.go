package translator

import "errors"

var (
	ErrSrcAndTgtNotMatched = errors.New("翻译异常, 错误: 原文和译文数量不对等")
)
