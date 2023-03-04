package region

const (
	T_AP_BeiJing Region = iota
	T_AP_ChengDu
	T_AP_ChongQing
	T_AP_GuangZhou
	T_AP_ShangHai
	T_AP_ShangHai_FSI
	T_AP_ShenZhen_FSI
	T_AP_HongKong
	T_AP_Mumbai
	T_AP_Seoul
	T_AP_Bangkok
	T_AP_Singapore
	T_AP_Tokyo
	T_AP_EU_Frankfurt
	T_AP_NA_Ashburn
	T_AP_Siliconvalley
	T_AP_Torronto
)

var (
	tencentZHMaps = []string{
		"华北地区-北京", "西南地区-成都", "西南地区-重庆", "华南地区-广州", "华东地区-上海",
		"华东地区-上海金融", "华南地区-深圳金融", "港澳台地区-中国香港",
		"亚太南部-孟买", "亚太东北-首尔", "亚太东南-曼谷", "亚太东南-新加坡", "亚太东北-东京",
		"欧洲地区-法兰克福", "美国东部-弗吉尼亚", "美国西部-硅谷", "北美地区-多伦多",
	}

	tencentENMaps = []string{
		"ap-beijing", "ap-chengdu", "ap-chongqing", "ap-guangzhou", "ap-shanghai",
		"ap-shanghai-fsi", "ap-shenzhen-fsi", "ap-hongkong",
		"ap-mumbai", "ap-seoul", "ap-bangkok", "ap-singapore", "ap-tokyo",
		"eu-frankfurt", "na-ashburn", "na-siliconvalley", "na-toronto",
	}
)
