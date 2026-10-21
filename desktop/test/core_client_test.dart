import 'dart:io';

import 'package:flutter_test/flutter_test.dart';

import 'package:anto_desktop/core_client.dart';

void main() {
  test('parseCoreReadyLine parses the Core handshake', () {
    final ready =
        parseCoreReadyLine('ANTO_CORE_READY {"port":41234,"token":"one-time"}');
    expect(ready?.port, 41234);
    expect(ready?.token, 'one-time');
  });

  test('parseCoreReadyLine ignores unrelated output', () {
    expect(parseCoreReadyLine('log line'), isNull);
  });

  test('parseCoreReadyLine rejects malformed handshake fields', () {
    expect(
      () => parseCoreReadyLine('ANTO_CORE_READY {"port":"41234"}'),
      throwsFormatException,
    );
  });

  test('start reports a missing Core executable and clears running state',
      () async {
    final client = CoreClient();
    try {
      await expectLater(
        client.start(executable: r'C:\anto\missing-anto-core.exe'),
        throwsA(anyOf(isA<ProcessException>(), isA<StateError>())),
      );
      expect(client.isRunning, isFalse);
    } finally {
      await client.dispose();
    }
  });
}
