import 'dart:io';

import 'package:csvons_gui/core/local_state_store.dart';
import 'package:test/test.dart';

void main() {
  test('saveRecentPaths persists and orders export paths', () async {
    final tempDir = await Directory.systemTemp.createTemp('local_state_store_test');
    final statePath = '${tempDir.path}${Platform.pathSeparator}state.json';
    final store = LocalStateStore(fileName: statePath);

    await store.saveRecentPaths(
      binaryPath: '/bin/a',
      rulerPath: '/tmp/ruler.json',
      exportPath: '/tmp/report',
    );
    await store.saveRecentExport(exportPath: '/tmp/report2');

    final state = await store.load();

    expect(state.recentBinaryPaths.first, '/bin/a');
    expect(state.recentRulerPaths.first, '/tmp/ruler.json');
    expect(state.recentExportPaths, ['/tmp/report2', '/tmp/report']);

    await tempDir.delete(recursive: true);
  });

  test('clear removes all recents', () async {
    final tempDir = await Directory.systemTemp.createTemp('local_state_store_test');
    final statePath = '${tempDir.path}${Platform.pathSeparator}state.json';
    final store = LocalStateStore(fileName: statePath);

    await store.saveRecentPaths(
      binaryPath: '/bin/a',
      rulerPath: '/tmp/ruler.json',
      workspacePath: '/tmp/workspace',
      exportPath: '/tmp/report',
    );

    await store.clear();
    final state = await store.load();

    expect(state.recentBinaryPaths, isEmpty);
    expect(state.recentRulerPaths, isEmpty);
    expect(state.recentWorkspacePaths, isEmpty);
    expect(state.recentExportPaths, isEmpty);

    await tempDir.delete(recursive: true);
  });
}

