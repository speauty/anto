/// Machine Translation public endpoints from the Aliyun service endpoint page.
class AliyunRegion {
  const AliyunRegion(this.id, this.name, this.publicEndpoint);

  final String id;
  final String name;
  final String publicEndpoint;
}

const aliyunRegions = <AliyunRegion>[
  AliyunRegion('cn-hangzhou', '华东1（杭州）', 'mt.cn-hangzhou.aliyuncs.com'),
  AliyunRegion('cn-shanghai', '华东2（上海）', 'mt.aliyuncs.com'),
  AliyunRegion('cn-beijing', '华北2（北京）', 'mt.aliyuncs.com'),
  AliyunRegion('cn-shenzhen', '华南1（深圳）', 'mt.aliyuncs.com'),
  AliyunRegion('cn-qingdao', '华北1（青岛）', 'mt.aliyuncs.com'),
  AliyunRegion('cn-huhehaote', '华北5（呼和浩特）', 'mt.aliyuncs.com'),
  AliyunRegion('cn-zhangjiakou', '华北3（张家口）', 'mt.aliyuncs.com'),
  AliyunRegion('cn-chengdu', '西南1（成都）', 'mt.aliyuncs.com'),
  AliyunRegion('cn-hongkong', '中国香港', 'mt.aliyuncs.com'),
  AliyunRegion('ap-southeast-1', '新加坡', 'mt.ap-southeast-1.aliyuncs.com'),
  AliyunRegion('ap-southeast-3', '马来西亚（吉隆坡）', 'mt.aliyuncs.com'),
  AliyunRegion('ap-northeast-1', '日本（东京）', 'mt.aliyuncs.com'),
  AliyunRegion('eu-central-1', '德国（法兰克福）', 'mt.aliyuncs.com'),
  AliyunRegion('eu-west-1', '英国（伦敦）', 'mt.aliyuncs.com'),
  AliyunRegion('us-east-1', '美国（弗吉尼亚）', 'mt.aliyuncs.com'),
  AliyunRegion('us-west-1', '美国（硅谷）', 'mt.aliyuncs.com'),
  AliyunRegion('me-central-1', '沙特（利雅得）', 'mt.me-central-1.aliyuncs.com'),
  AliyunRegion('me-east-1', '阿联酋（迪拜）', 'alimt.me-east-1.aliyuncs.com'),
  AliyunRegion('cn-beijing-finance-1', '华北2 金融云（邀测）', 'mt.aliyuncs.com'),
  AliyunRegion('cn-hangzhou-finance', '华东1 金融云', 'mt.aliyuncs.com'),
  AliyunRegion('cn-north-2-gov-1', '北京政务云', 'mt.aliyuncs.com'),
  AliyunRegion('cn-shanghai-finance-1', '华东2 金融云', 'mt.aliyuncs.com'),
  AliyunRegion('cn-shenzhen-finance-1', '华南1 金融云', 'mt.aliyuncs.com'),
];

AliyunRegion? aliyunRegionById(String id) {
  for (final region in aliyunRegions) {
    if (region.id == id) return region;
  }
  return null;
}
