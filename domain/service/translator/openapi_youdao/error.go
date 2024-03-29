package openapi_youdao

var errorMap = map[string]string{
	"101": "缺少必填的参数,首先确保必填参数齐全, 然后确认参数书写是否正确",
	"102": "不支持的语言类型",
	"103": "翻译文本过长",
	"104": "不支持的API类型",
	"105": "不支持的签名类型",
	"106": "不支持的响应类型",
	"107": "不支持的传输加密类型",
	"108": "应用ID无效, 注册账号, 登录后台创建应用和实例并完成绑定, 可获得应用ID和应用密钥等信息",
	"109": "batchLog格式不正确",
	"110": "无相关服务的有效实例,应用没有绑定服务, 可以新建服务, 绑定服务",
	"111": "开发者账号无效",
	"112": "请求服务无效",
	"113": "翻译字符串(q)不能为空",
	"118": "detectLevel取值错误",
	"201": "解密失败, 可能为DES,BASE64,URLDecode的错误",
	"202": "签名检验失败",
	"203": "访问IP地址不在可访问IP列表",
	"205": "请求的接口与应用的平台类型不一致",
	"206": "因为时间戳无效导致签名校验失败",
	"207": "重放请求",
	"301": "辞典查询失败",
	"302": "翻译查询失败",
	"303": "服务端的其它异常",
	"304": "会话闲置太久超时",
	"401": "账户已经欠费, 请进行账户充值",
	"402": "offlinesdk不可用",
	"411": "访问频率受限,请稍后访问",
	"412": "长请求过于频繁, 请稍后访问",
}
