/// Regions accepted by Tencent Cloud Machine Translation's API 3.0.
class TencentRegion {
  const TencentRegion(this.id, this.name);

  final String id;
  final String name;
}

const tencentRegions = <TencentRegion>[
  TencentRegion('ap-guangzhou', '广州'),
  TencentRegion('ap-shanghai', '上海'),
  TencentRegion('ap-beijing', '北京'),
  TencentRegion('ap-chengdu', '成都'),
  TencentRegion('ap-chongqing', '重庆'),
  TencentRegion('ap-hongkong', '中国香港'),
  TencentRegion('ap-singapore', '新加坡'),
  TencentRegion('ap-bangkok', '曼谷'),
  TencentRegion('ap-seoul', '首尔'),
  TencentRegion('ap-tokyo', '东京'),
  TencentRegion('eu-frankfurt', '法兰克福'),
  TencentRegion('na-ashburn', '弗吉尼亚'),
  TencentRegion('na-siliconvalley', '硅谷'),
];

TencentRegion? tencentRegionById(String id) {
  for (final region in tencentRegions) {
    if (region.id == id) return region;
  }
  return null;
}
