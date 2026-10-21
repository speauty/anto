List<String> parseSrtPaths(String text) => text
    .split(RegExp(r'[;\r\n]+'))
    .map((value) => value.trim())
    .where((value) => value.isNotEmpty && value.toLowerCase().endsWith('.srt'))
    .toSet()
    .toList();

bool isLanguagePairValid(String from, String to) {
  final source = from.trim().toLowerCase();
  final target = to.trim().toLowerCase();
  return source.isEmpty || target.isEmpty || source != target;
}

bool isUsableTranslator(Iterable<dynamic> translators, String? id) {
  if (id == null || id.isEmpty) return false;
  for (final item in translators) {
    if (item is Map && item['id'] == id) return item['valid'] == true;
  }
  return false;
}
