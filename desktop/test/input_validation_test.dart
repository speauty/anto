import 'package:flutter_test/flutter_test.dart';

import 'package:anto_desktop/input_validation.dart';

void main() {
  test('parseSrtPaths filters non-SRT paths and duplicates', () {
    expect(
      parseSrtPaths(' a.srt; b.txt\nA.SRT\n ; c.srt '),
      ['a.srt', 'A.SRT', 'c.srt'],
    );
  });

  test('language pair rejects the same language case-insensitively', () {
    expect(isLanguagePairValid(' en ', 'EN'), isFalse);
    expect(isLanguagePairValid('en', 'zh-cn'), isTrue);
    expect(isLanguagePairValid('', ''), isTrue);
  });

  test('only explicitly valid translators can be selected', () {
    const translators = [
      {'id': 'ali_cloud_mt', 'valid': true},
    ];
    expect(isUsableTranslator(translators, 'unknown'), isFalse);
    expect(isUsableTranslator(translators, 'ali_cloud_mt'), isTrue);
    expect(isUsableTranslator(translators, null), isFalse);
  });
}
