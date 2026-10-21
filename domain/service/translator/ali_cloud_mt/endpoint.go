package ali_cloud_mt

// Public endpoints for the regions currently listed in the Aliyun Machine
// Translation endpoint documentation. The bundled SDK's endpoint table is
// outdated for me-east-1 and does not include me-central-1.
var publicEndpoints = map[string]string{
	"cn-hangzhou":           "mt.cn-hangzhou.aliyuncs.com",
	"cn-shanghai":           "mt.aliyuncs.com",
	"cn-beijing":            "mt.aliyuncs.com",
	"cn-shenzhen":           "mt.aliyuncs.com",
	"cn-qingdao":            "mt.aliyuncs.com",
	"cn-huhehaote":          "mt.aliyuncs.com",
	"cn-zhangjiakou":        "mt.aliyuncs.com",
	"cn-chengdu":            "mt.aliyuncs.com",
	"cn-hongkong":           "mt.aliyuncs.com",
	"ap-southeast-1":        "mt.ap-southeast-1.aliyuncs.com",
	"ap-southeast-3":        "mt.aliyuncs.com",
	"ap-northeast-1":        "mt.aliyuncs.com",
	"eu-central-1":          "mt.aliyuncs.com",
	"eu-west-1":             "mt.aliyuncs.com",
	"us-east-1":             "mt.aliyuncs.com",
	"us-west-1":             "mt.aliyuncs.com",
	"me-central-1":          "mt.me-central-1.aliyuncs.com",
	"me-east-1":             "alimt.me-east-1.aliyuncs.com",
	"cn-beijing-finance-1":  "mt.aliyuncs.com",
	"cn-hangzhou-finance":   "mt.aliyuncs.com",
	"cn-north-2-gov-1":      "mt.aliyuncs.com",
	"cn-shanghai-finance-1": "mt.aliyuncs.com",
	"cn-shenzhen-finance-1": "mt.aliyuncs.com",
}
