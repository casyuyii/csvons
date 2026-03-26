import 'package:flutter/material.dart';

import '../core/local_state_store.dart';
import '../core/validation_runner.dart';
import '../models/validation_report.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _store = LocalStateStore();
  final _binaryController = TextEditingController(
    text: ValidationRunner.defaultBinaryPath(),
  );
  final _rulerController = TextEditingController();
  final _filterController = TextEditingController();

  List<String> _recentBinaryPaths = const <String>[];
  List<String> _recentRulerPaths = const <String>[];

  bool _running = false;
  ValidationResult? _result;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadState();
  }

  Future<void> _loadState() async {
    final state = await _store.load();
    if (!mounted) return;

    setState(() {
      _recentBinaryPaths = state.recentBinaryPaths;
      _recentRulerPaths = state.recentRulerPaths;
      if (_rulerController.text.trim().isEmpty && state.recentRulerPaths.isNotEmpty) {
        _rulerController.text = state.recentRulerPaths.first;
      }
      if (_binaryController.text.trim().isEmpty && state.recentBinaryPaths.isNotEmpty) {
        _binaryController.text = state.recentBinaryPaths.first;
      }
    });
  }

  @override
  void dispose() {
    _binaryController.dispose();
    _rulerController.dispose();
    _filterController.dispose();
    super.dispose();
  }

  Future<void> _run() async {
    setState(() {
      _running = true;
      _error = null;
    });

    try {
      final runner = ValidationRunner(binaryPath: _binaryController.text.trim());
      final result = await runner.run(rulerPath: _rulerController.text.trim());
      await _store.saveRecentPaths(
        binaryPath: _binaryController.text,
        rulerPath: _rulerController.text,
      );
      await _loadState();

      setState(() {
        _result = result;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
      });
    } finally {
      if (mounted) {
        setState(() {
          _running = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('csvons GUI starter')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _PathInputRow(
              label: 'csvons binary path',
              controller: _binaryController,
              recents: _recentBinaryPaths,
            ),
            const SizedBox(height: 12),
            _PathInputRow(
              label: 'ruler.json absolute path',
              controller: _rulerController,
              recents: _recentRulerPaths,
            ),
            const SizedBox(height: 12),
            FilledButton.icon(
              onPressed: _running ? null : _run,
              icon: const Icon(Icons.play_arrow),
              label: Text(_running ? 'Running...' : 'Run Validation'),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _filterController,
              decoration: const InputDecoration(
                labelText: 'Filter issues (message/file/rule/field)',
              ),
              onChanged: (_) => setState(() {}),
            ),
            const SizedBox(height: 8),
            if (_error != null)
              Text('Error: $_error', style: const TextStyle(color: Colors.red)),
            if (_result != null)
              _ResultView(result: _result!, filterText: _filterController.text),
          ],
        ),
      ),
    );
  }
}

class _PathInputRow extends StatelessWidget {
  const _PathInputRow({
    required this.label,
    required this.controller,
    required this.recents,
  });

  final String label;
  final TextEditingController controller;
  final List<String> recents;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        TextField(
          controller: controller,
          decoration: InputDecoration(labelText: label),
        ),
        if (recents.isNotEmpty) ...[
          const SizedBox(height: 6),
          Wrap(
            spacing: 8,
            runSpacing: 6,
            children: recents
                .take(5)
                .map(
                  (path) => ActionChip(
                    label: Text(path, overflow: TextOverflow.ellipsis),
                    onPressed: () => controller.text = path,
                  ),
                )
                .toList(growable: false),
          ),
        ],
      ],
    );
  }
}

class _ResultView extends StatelessWidget {
  const _ResultView({required this.result, required this.filterText});

  final ValidationResult result;
  final String filterText;

  @override
  Widget build(BuildContext context) {
    final report = result.report;
    return Expanded(
      child: ListView(
        children: [
          Text('Exit code: ${result.exitCode}'),
          if (report != null) ...[
            Text('Files checked: ${report.summary.filesChecked}'),
            Text('Passed: ${report.summary.passed}'),
            Text('Failed: ${report.summary.failed}'),
            const Divider(),
            _IssueTable(
              issues: _filteredIssues(report.issues, filterText),
            ),
          ] else ...[
            const Text('Raw stdout:'),
            SelectableText(result.stdoutText),
            if (result.stderrText.trim().isNotEmpty) ...[
              const Divider(),
              const Text('Raw stderr:'),
              SelectableText(result.stderrText),
            ],
          ],
        ],
      ),
    );
  }

  List<ValidationIssue> _filteredIssues(List<ValidationIssue> issues, String filterText) {
    final q = filterText.trim().toLowerCase();
    if (q.isEmpty) return issues;

    return issues.where((i) {
      final joined = '${i.message} ${i.file} ${i.rule} ${i.field}'.toLowerCase();
      return joined.contains(q);
    }).toList(growable: false);
  }
}

class _IssueTable extends StatefulWidget {
  const _IssueTable({required this.issues});

  final List<ValidationIssue> issues;

  @override
  State<_IssueTable> createState() => _IssueTableState();
}

class _IssueTableState extends State<_IssueTable> {
  int _sortColumnIndex = 0;
  bool _sortAscending = true;

  @override
  Widget build(BuildContext context) {
    final rows = widget.issues.toList(growable: false);
    _sortRows(rows);

    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: DataTable(
        sortColumnIndex: _sortColumnIndex,
        sortAscending: _sortAscending,
        columns: [
          _column('Severity', 0),
          _column('Message', 1),
          _column('File', 2),
          _column('Rule', 3),
          _column('Field', 4),
          _column('Row', 5),
        ],
        rows: rows
            .map(
              (i) => DataRow(
                cells: [
                  DataCell(Text(i.severity)),
                  DataCell(SizedBox(width: 360, child: Text(i.message))),
                  DataCell(Text(i.file)),
                  DataCell(Text(i.rule)),
                  DataCell(Text(i.field)),
                  DataCell(Text(i.row?.toString() ?? '-')),
                ],
              ),
            )
            .toList(growable: false),
      ),
    );
  }

  DataColumn _column(String label, int index) {
    return DataColumn(
      label: Text(label),
      onSort: (_, ascending) {
        setState(() {
          _sortColumnIndex = index;
          _sortAscending = ascending;
        });
      },
    );
  }

  void _sortRows(List<ValidationIssue> rows) {
    String asString(ValidationIssue i) {
      switch (_sortColumnIndex) {
        case 0:
          return i.severity;
        case 1:
          return i.message;
        case 2:
          return i.file;
        case 3:
          return i.rule;
        case 4:
          return i.field;
        case 5:
          return i.row?.toString() ?? '';
        default:
          return i.message;
      }
    }

    rows.sort((a, b) {
      final cmp = asString(a).compareTo(asString(b));
      return _sortAscending ? cmp : -cmp;
    });
  }
}
