import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:csvons_gui/core/validation_runner.dart';
import 'package:test/test.dart';

class _FakeProcess implements Process {
  _FakeProcess({
    required this.stdoutText,
    required this.stderrText,
    required this.code,
  });

  final String stdoutText;
  final String stderrText;
  final int code;

  @override
  Future<int> get exitCode async => code;

  @override
  int get pid => 42;

  @override
  IOSink get stdin => throw UnimplementedError();

  @override
  Stream<List<int>> get stderr => Stream<List<int>>.value(utf8.encode(stderrText));

  @override
  Stream<List<int>> get stdout => Stream<List<int>>.value(utf8.encode(stdoutText));

  @override
  bool kill([ProcessSignal signal = ProcessSignal.sigterm]) => true;
}

void main() {
  test('builds JSON run args and parses report output', () async {
    late String capturedExecutable;
    late List<String> capturedArgs;

    final runner = ValidationRunner(
      binaryPath: '/tmp/csvons',
      processStarter: (executable, arguments) async {
        capturedExecutable = executable;
        capturedArgs = arguments;
        return _FakeProcess(
          code: 1,
          stdoutText:
              '{"schema_version":"1","summary":{"files_checked":1,"failed":1},"issues":[{"message":"x"}]}',
          stderrText: '',
        );
      },
    );

    final result = await runner.run(rulerPath: '/tmp/ruler.json');

    expect(capturedExecutable, '/tmp/csvons');
    expect(capturedArgs, ['--format', 'json', '/tmp/ruler.json']);
    expect(result.exitCode, 1);
    expect(result.report, isNotNull);
    expect(result.report!.summary.failed, 1);
  });

  test('keeps report null when stdout is not valid json', () async {
    final runner = ValidationRunner(
      binaryPath: '/tmp/csvons',
      processStarter: (_, __) async => _FakeProcess(
        code: 2,
        stdoutText: 'panic: bad config',
        stderrText: 'trace...',
      ),
    );

    final result = await runner.run(rulerPath: '/tmp/ruler.json');

    expect(result.report, isNull);
    expect(result.stderrText, contains('trace'));
  });

  test('does not force json argument for text runs', () async {
    late List<String> capturedArgs;
    final runner = ValidationRunner(
      binaryPath: '/tmp/csvons',
      processStarter: (_, arguments) async {
        capturedArgs = arguments;
        return _FakeProcess(code: 0, stdoutText: 'ok', stderrText: '');
      },
    );

    final result = await runner.run(
      rulerPath: '/tmp/ruler.json',
      jsonFormat: false,
    );

    expect(capturedArgs, ['/tmp/ruler.json']);
    expect(result.report, isNull);
  });
}
