import 'dart:io';

import 'package:csvons_gui/core/csv_preview.dart';
import 'package:test/test.dart';

void main() {
  test('parseCsvLine handles quoted commas and escaped quotes', () {
    final row = CsvPreview.parseCsvLine('1,"hello,world","a""b"');

    expect(row, ['1', 'hello,world', 'a"b']);
  });

  test('load reads header and sample rows', () async {
    final tempDir = await Directory.systemTemp.createTemp('csv_preview_test');
    final csvFile = File('${tempDir.path}${Platform.pathSeparator}sample.csv');
    await csvFile.writeAsString('id,name\n1,Alice\n2,Bob\n');

    final preview = await CsvPreview.load(csvFile.path);

    expect(preview.header, ['id', 'name']);
    expect(preview.rows, [
      ['1', 'Alice'],
      ['2', 'Bob'],
    ]);

    await tempDir.delete(recursive: true);
  });

  test('load keeps multiline quoted field as a single row', () async {
    final tempDir = await Directory.systemTemp.createTemp('csv_preview_test');
    final csvFile = File('${tempDir.path}${Platform.pathSeparator}sample_multiline.csv');
    await csvFile.writeAsString('id,notes\n1,"line one\nline two"\n2,done\n');

    final preview = await CsvPreview.load(csvFile.path);

    expect(preview.header, ['id', 'notes']);
    expect(preview.rows.first, ['1', 'line one\nline two']);
    expect(preview.rows[1], ['2', 'done']);

    await tempDir.delete(recursive: true);
  });

  test('load respects maxRows', () async {
    final tempDir = await Directory.systemTemp.createTemp('csv_preview_test');
    final csvFile = File('${tempDir.path}${Platform.pathSeparator}sample_many.csv');
    await csvFile.writeAsString('id\n1\n2\n3\n4\n');

    final preview = await CsvPreview.load(csvFile.path, maxRows: 2);

    expect(preview.header, ['id']);
    expect(preview.rows, [
      ['1'],
      ['2'],
    ]);

    await tempDir.delete(recursive: true);
  });
}
