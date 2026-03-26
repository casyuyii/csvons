import 'dart:io';

import 'package:flutter/material.dart';

import '../core/local_state_store.dart';
import '../core/validation_runner.dart';
import '../models/validation_report.dart';
import '../widgets/issues_table.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _stateStore = LocalStateStore();
  final _binaryController = TextEditingController(
    text: ValidationRunner.defaultBinaryPath(),
  );
  final _rulerController = TextEditingController();

  bool _running = false;
  ValidationResult? _result;
  String? _error;
  List<String> _recentBinaryPaths = const <String>[];
  List<String> _recentRulerPaths = const <String>[];

  @override
  void initState() {
    super.initState();
    _loadState();
  }

  Future<void> _loadState() async {
    final state = await _stateStore.load();
    if (!mounted) return;

    setState(() {
      _recentBinaryPaths = state.recentBinaryPaths;
      _recentRulerPaths = state.recentRulerPaths;
      if (_binaryController.text.trim().isEmpty &&
          _recentBinaryPaths.isNotEmpty) {
        _binaryController.text = _recentBinaryPaths.first;
      }
      if (_rulerController.text.trim().isEmpty && _recentRulerPaths.isNotEmpty) {
        _rulerController.text = _recentRulerPaths.first;
      }
    });
  }

  @override
  void dispose() {
    _binaryController.dispose();
    _rulerController.dispose();
    super.dispose();
  }

  Future<void> _run() async {
    final binary = _binaryController.text.trim();
    final ruler = _rulerController.text.trim();
    final preflightError = await _validateInputs(binary: binary, ruler: ruler);
    if (preflightError != null) {
      setState(() {
        _error = preflightError;
      });
      return;
    }

    setState(() {
      _running = true;
      _error = null;
    });

    try {
      await _stateStore.saveRecentPaths(binaryPath: binary, rulerPath: ruler);
      final runner = ValidationRunner(binaryPath: binary);
      final result = await runner.run(rulerPath: ruler);
      final latestState = await _stateStore.load();
      setState(() {
        _result = result;
        _recentBinaryPaths = latestState.recentBinaryPaths;
        _recentRulerPaths = latestState.recentRulerPaths;
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

  Future<String?> _validateInputs({
    required String binary,
    required String ruler,
  }) async {
    if (binary.isEmpty || ruler.isEmpty) {
      return 'Both binary path and ruler path are required.';
    }

    try {
      if (!await File(binary).exists()) {
        return 'Binary not found: $binary';
      }
      if (!await File(ruler).exists()) {
        return 'Ruler file not found: $ruler';
      }
    } on FileSystemException catch (e) {
      return 'Unable to access selected paths: ${e.message}';
    }
    return null;
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
            TextField(
              controller: _binaryController,
              decoration: const InputDecoration(
                labelText: 'csvons binary path',
              ),
            ),
            _RecentPathRow(
              label: 'Recent binary paths',
              items: _recentBinaryPaths,
              onTap: (value) => _binaryController.text = value,
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _rulerController,
              decoration: const InputDecoration(
                labelText: 'ruler.json absolute path',
              ),
            ),
            _RecentPathRow(
              label: 'Recent ruler paths',
              items: _recentRulerPaths,
              onTap: (value) => _rulerController.text = value,
            ),
            const SizedBox(height: 12),
            FilledButton.icon(
              onPressed: _running ? null : _run,
              icon: const Icon(Icons.play_arrow),
              label: Text(_running ? 'Running...' : 'Run Validation'),
            ),
            const SizedBox(height: 16),
            if (_result != null) _RunStatus(result: _result!),
            if (_error != null)
              Text('Error: $_error', style: const TextStyle(color: Colors.red)),
            Expanded(
              child: _result == null
                  ? const _EmptyRunState()
                  : _ResultView(result: _result!),
            ),
          ],
        ),
      ),
    );
  }
}

class _EmptyRunState extends StatelessWidget {
  const _EmptyRunState();

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.analytics_outlined, size: 42, color: Colors.grey.shade500),
          const SizedBox(height: 8),
          Text('No validation run yet', style: textTheme.titleMedium),
          const SizedBox(height: 4),
          Text(
            'Pick a binary and ruler file, then click "Run Validation".',
            style: textTheme.bodyMedium?.copyWith(color: Colors.grey.shade700),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }
}

class _ResultView extends StatelessWidget {
  const _ResultView({required this.result});

  final ValidationResult result;

  @override
  Widget build(BuildContext context) {
    final report = result.report;
    return ListView(
      children: [
        Text('Exit code: ${result.exitCode}'),
        if (report != null) ...[
          if (report.schemaVersion.isNotEmpty)
            Text('Schema: ${report.schemaVersion}'),
          Text('Files checked: ${report.summary.filesChecked}'),
          Text('Passed: ${report.summary.passed}'),
          Text('Failed: ${report.summary.failed}'),
          const Divider(),
          IssuesTable(report: report),
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
    );
  }
}

class _RunStatus extends StatelessWidget {
  const _RunStatus({required this.result});

  final ValidationResult result;

  @override
  Widget build(BuildContext context) {
    final (label, color) = switch (result.exitCode) {
      0 => ('Passed', Colors.green),
      1 => ('Validation Issues', Colors.orange),
      _ => ('Runtime Error', Colors.red),
    };

    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Row(
        children: [
          Icon(Icons.circle, size: 10, color: color),
          const SizedBox(width: 8),
          Text(
            '$label (exit ${result.exitCode})',
            style: TextStyle(fontWeight: FontWeight.w600, color: color),
          ),
        ],
      ),
    );
  }
}

class _RecentPathRow extends StatelessWidget {
  const _RecentPathRow({
    required this.label,
    required this.items,
    required this.onTap,
  });

  final String label;
  final List<String> items;
  final ValueChanged<String> onTap;

  @override
  Widget build(BuildContext context) {
    if (items.isEmpty) return const SizedBox.shrink();

    return Padding(
      padding: const EdgeInsets.only(top: 6),
      child: Wrap(
        spacing: 8,
        runSpacing: 8,
        crossAxisAlignment: WrapCrossAlignment.center,
        children: [
          Text('$label:'),
          ...items.take(3).map(
                (item) => ActionChip(
                  label: Text(
                    item,
                    overflow: TextOverflow.ellipsis,
                  ),
                  onPressed: () => onTap(item),
                ),
              ),
        ],
      ),
    );
  }
}
