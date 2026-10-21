package tencent_cloud_mt

const publicEndpoint = "tmt.tencentcloudapi.com"

var regions = map[string]struct{}{
	"ap-bangkok": {}, "ap-beijing": {}, "ap-chengdu": {}, "ap-chongqing": {},
	"ap-guangzhou": {}, "ap-hongkong": {}, "ap-seoul": {}, "ap-shanghai": {},
	"ap-singapore": {},
	"ap-tokyo": {}, "eu-frankfurt": {}, "na-ashburn": {}, "na-siliconvalley": {},
}

func validRegion(region string) bool {
	_, ok := regions[region]
	return ok
}
