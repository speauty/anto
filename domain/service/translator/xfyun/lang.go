package xfyun

import (
	"anto/domain/service/translator"
)

// @link https://www.xfyun.cn/doc/nlp/xftrans_new/API.html#%E8%AF%AD%E7%A7%8D%E5%88%97%E8%A1%A8
var langSupported = []translator.LangPair{
	{"cn", "中文"}, {"en", "英语"}, {"cs", "捷克语"}, {"ha", "豪萨语"},
	{"ja", "日语"}, {"ro", "罗马尼亚语"}, {"hu", "匈牙利语"}, {"ko", "韩语"},
	{"sv", "瑞典语"}, {"sw", "斯瓦希里语"}, {"th", "泰语"}, {"nl", "荷兰语"},
	{"uz", "乌兹别克语"}, {"ru", "俄语"}, {"pl", "波兰语"}, {"zu", "祖鲁语"},
	{"bg", "保加利亚语"}, {"ar", "阿拉伯语"}, {"el", "希腊语"}, {"uk", "乌克兰语"},
	{"fa", "波斯语"}, {"he", "希伯来语"}, {"vi", "越南语"}, {"ps", "普什图语"},
	{"hy", "亚美尼亚语"}, {"ms", "马来语"}, {"ur", "乌尔都语"}, {"hy", "亚美尼亚语"},
	{"ms", "马来语"}, {"ur", "乌尔都语"}, {"ka", "格鲁吉亚语"}, {"id", "印尼语"},
	{"yue", "广东话"}, {"tl", "菲律宾语"}, {"bn", "孟加拉语"}, {"ii", "彝语"},
	{"de", "德语"}, {"nm", "外蒙语"}, {"zua", "壮语"}, {"es", "西班牙语"},
	{"kk", "外哈语"}, {"mn", "内蒙语"}, {"fr", "法语"}, {"tr", "土耳其语"},
	{"kka", "内哈萨克语"},
}
