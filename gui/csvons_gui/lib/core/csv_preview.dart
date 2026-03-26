import 'dart:convert';
import 'dart:io';

class CsvPreview {
  CsvPreview({
    required this.header,
    required this.rows,
  });

  final List<String> header;
  final List<List<String>> rows;

  static Future<CsvPreview> load(
    String path, {
    int maxRows = 6,
  }) async {
    final file = File(path);
    if (!await file.exists()) {
      throw const FileSystemException('CSV file does not exist');
    }

    final records = <List<String>>[];
    final currentRecord = StringBuffer();
    var inQuotes = false;

    await for (final line in file
        .openRead()
        .transform(utf8.decoder)
        .transform(const LineSplitter())) {
      if (currentRecord.isNotEmpty) {
        currentRecord.write('\n');
      }
      currentRecord.write(line);

      inQuotes = _nextQuoteState(line, initialState: inQuotes);
      if (inQuotes) {
        continue;
      }

      records.add(parseCsvLine(currentRecord.toString()));
      currentRecord.clear();

      if (records.length >= maxRows + 1) {
        break;
      }
    }

    if (currentRecord.isNotEmpty) {
      records.add(parseCsvLine(currentRecord.toString()));
    }

    if (records.isEmpty) {
      return CsvPreview(
        header: const <String>[],
        rows: const <List<String>>[],
      );
    }

    return CsvPreview(
      header: records.first,
      rows: records.skip(1).take(maxRows).toList(growable: false),
    );
  }

  static List<String> parseCsvLine(String line) {
    final cells = <String>[];
    final current = StringBuffer();
    var inQuotes = false;

    for (var i = 0; i < line.length; i++) {
      final ch = line[i];
      if (ch == '"') {
        if (inQuotes && i + 1 < line.length && line[i + 1] == '"') {
          current.write('"');
          i++;
          continue;
        }
        inQuotes = !inQuotes;
        continue;
      }

      if (ch == ',' && !inQuotes) {
        cells.add(current.toString().trim());
        current.clear();
        continue;
      }

      current.write(ch);
    }

    cells.add(current.toString().trim());
    return cells;
  }

  static bool _nextQuoteState(
    String line, {
    required bool initialState,
  }) {
    var inQuotes = initialState;
    for (var i = 0; i < line.length; i++) {
      final ch = line[i];
      if (ch != '"') continue;
      if (inQuotes && i + 1 < line.length && line[i + 1] == '"') {
        i++;
        continue;
      }
      inQuotes = !inQuotes;
    }
    return inQuotes;
  }
}
