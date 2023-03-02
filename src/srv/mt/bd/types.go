package bd

import (
	"strings"
)

// GenericTextApiVersion 接口版本-通用文本翻译
type GenericTextApiVersion int

// FromInt 从int转化
func (gtAV GenericTextApiVersion) FromInt(num int) GenericTextApiVersion {
	return GenericTextApiVersion(num)
}

// ToInt 转化成int输出
func (gtAV GenericTextApiVersion) ToInt() int {
	return int(gtAV)
}

// GetZHArrays 获取所有中文
func (gtAV GenericTextApiVersion) GetZHArrays() []string {
	return gtZHMaps
}

// GetZH 获取对应中文
func (gtAV GenericTextApiVersion) GetZH() string {
	return gtZHMaps[gtAV]
}

// GetQPS 获取对应QPS
func (gtAV GenericTextApiVersion) GetQPS() int {
	return gtQPSMaps[gtAV]
}

// GetLenLimited 获取对应长度限制
func (gtAV GenericTextApiVersion) GetLenLimited() int {
	return gtLenLimitedPQMaps[gtAV]
}

// GetLenFreeMonth 获取每月免费长度
func (gtAV GenericTextApiVersion) GetLenFreeMonth() int {
	return gtLenFreeMonthMaps[gtAV]
}

// GetRetryLimited 获取重试限制
func (gtAV GenericTextApiVersion) GetRetryLimited() int {
	return gtRetryLimitedMaps[gtAV]
}

const (
	// GTApiStandard 标准版接口
	GTApiStandard GenericTextApiVersion = iota
	// GTApiHigh 高级版接口
	GTApiHigh
	// GTApiEnjoy 尊享版接口
	GTApiEnjoy
)

// gtQPSMaps  QPS限定
var gtQPSMaps = map[GenericTextApiVersion]int{
	GTApiStandard: 1, GTApiHigh: 10, GTApiEnjoy: 100,
}

// gtLenLimitedPQMaps  单次查询长度限定
var gtLenLimitedPQMaps = map[GenericTextApiVersion]int{
	GTApiStandard: 1000, GTApiHigh: 6000, GTApiEnjoy: 6000,
}

// gtLenFreeMonthMaps  每月免费额度
var gtLenFreeMonthMaps = map[GenericTextApiVersion]int{
	GTApiStandard: 5e4, GTApiHigh: 1e6, GTApiEnjoy: 2e6,
}

// gtRetryLimitedMaps 重试限制
var gtRetryLimitedMaps = map[GenericTextApiVersion]int{
	GTApiStandard: 1, GTApiHigh: 2, GTApiEnjoy: 4,
}

// gtZHMaps  中文映射
var gtZHMaps = []string{"标准版", "高级版", "尊享版"}

type ErrCode string

func (ec ErrCode) FromString(str string) ErrCode {
	return ErrCode(str)
}

func (ec ErrCode) ToString() string {
	return string(ec)
}

func (ec ErrCode) GetZH() string {
	return errCodeZHMaps[ec]
}

func (ec ErrCode) IsExit(err error) bool {
	var exitFlags = []ErrCode{
		ErrSystemBroken, ErrUnauthorized, ErrSign, ErrAccountNotEnough, ErrClientIPInvalid, ErrToLanguageNotSupported,
		ErrSrvClosed, ErrAuthenticationInvalid,
	}

	for _, flag := range exitFlags {
		if strings.Contains(err.Error(), flag.ToString()) {
			return true
		}
	}

	return false
}

const (
	ErrRequestTimeout         ErrCode = "52001"
	ErrSystemBroken           ErrCode = "52002"
	ErrUnauthorized           ErrCode = "52003"
	ErrArgRequired            ErrCode = "54000"
	ErrSign                   ErrCode = "54001"
	ErrRequestTooMany         ErrCode = "54003"
	ErrAccountNotEnough       ErrCode = "54004"
	ErrLongQueryTooMany       ErrCode = "54005"
	ErrClientIPInvalid        ErrCode = "58000"
	ErrToLanguageNotSupported ErrCode = "58001"
	ErrSrvClosed              ErrCode = "58002"
	ErrAuthenticationInvalid  ErrCode = "90107"
)

var errCodeZHMaps = map[ErrCode]string{
	ErrRequestTimeout:         "请求超时, 请稍后重试",
	ErrSystemBroken:           "系统错误, 请稍后重试",
	ErrUnauthorized:           "当前用户暂未授权, 请检查您的访问身份(AppId)是否正确, 或服务是否开通",
	ErrArgRequired:            "当前请求必填参数为空, 请检查是否少传参数",
	ErrSign:                   "当前请求签名错误, 请检查您的签名生成方法",
	ErrRequestTooMany:         "当前访问频率受限, 请降低您的调用频率, 或进行身份认证后切换为高级版/尊享版",
	ErrAccountNotEnough:       "您的账户余额不足, 请前往管理控制台为当前账户充值",
	ErrLongQueryTooMany:       "当前长query请求频繁, 请降低长query的发送频率, 3s后再试",
	ErrClientIPInvalid:        "当前客户端IP非法, 检查个人资料里填写的IP地址是否正确, 可前往开发者信息-基本信息修改",
	ErrToLanguageNotSupported: "当前目标语言方向不支持, 检查译文语言是否在语言列表里",
	ErrSrvClosed:              "当前服务已关闭, 请前往管理控制台开启服务",
	ErrAuthenticationInvalid:  "当前认证未通过或未生效, 请前往我的认证查看认证进度",
}
