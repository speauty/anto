package baidu_cloud_mt

import "anto/domain/service/translator"

// These are the common languages available to all API service tiers.
var langSupported = []translator.LangPair{
	{Key: "zh", Name: "中文"}, {Key: "cht", Name: "繁体中文"},
	{Key: "en", Name: "英语"}, {Key: "jp", Name: "日语"},
	{Key: "kor", Name: "韩语"}, {Key: "fra", Name: "法语"},
	{Key: "spa", Name: "西班牙语"}, {Key: "th", Name: "泰语"},
	{Key: "ara", Name: "阿拉伯语"}, {Key: "ru", Name: "俄语"},
	{Key: "pt", Name: "葡萄牙语"}, {Key: "de", Name: "德语"},
	{Key: "it", Name: "意大利语"}, {Key: "el", Name: "希腊语"},
	{Key: "nl", Name: "荷兰语"}, {Key: "pl", Name: "波兰语"},
	{Key: "bul", Name: "保加利亚语"}, {Key: "est", Name: "爱沙尼亚语"},
	{Key: "dan", Name: "丹麦语"}, {Key: "fin", Name: "芬兰语"},
	{Key: "cs", Name: "捷克语"}, {Key: "rom", Name: "罗马尼亚语"},
	{Key: "slo", Name: "斯洛文尼亚语"}, {Key: "swe", Name: "瑞典语"},
	{Key: "hu", Name: "匈牙利语"}, {Key: "vie", Name: "越南语"},
	{Key: "yue", Name: "粤语"}, {Key: "wyw", Name: "文言文"},
}
