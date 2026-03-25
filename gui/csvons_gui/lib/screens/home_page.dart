import 'package:flutter/material.dart';

import '../core/validation_runner.dart';
import '../models/validation_report.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _binaryController = TextEditingController(
    text: ValidationRunner.defaultBinaryPath(),
  );
  final _rulerController = TextEditingController();

  bool _running = false;
  ValidationResult? _result;
  String? _error;

  @override
  void dispose() {
    _binaryController.dispose();
    _rulerController.dispose();
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
            TextField(
              controller: _binaryController,
              decoration: const InputDecoration(
                labelText: 'csvons binary path',
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _rulerController,
              decoration: const InputDecoration(
                labelText: 'ruler.json absolute path',
              ),
            ),
            const SizedBox(height: 12),
            FilledButton.icon(
              onPressed: _running ? null : _run,
              icon: const Icon(Icons.play_arrow),
              label: Text(_running ? 'Running...' : 'Run Validation'),
            ),
            const SizedBox(height: 16),
            if (_error != null)
              Text('Error: $_error', style: const TextStyle(color: Colors.red)),
            if (_result != null) _ResultView(result: _result!),
          ],
        ),
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
    return Expanded(
      child: ListView(
        children: [
          Text('Exit code: ${result.exitCode}'),
          if (report != null) ...[
            Text('Files checked: ${report.summary.filesChecked}'),
            Text('Passed: ${report.summary.passed}'),
            Text('Failed: ${report.summary.failed}'),
            const Divider(),
            ...report.issues.map(
              (i) => ListTile(
                dense: true,
                title: Text('[${i.severity}] ${i.message}'),
                subtitle: Text('${i.file} · ${i.rule} · ${i.field}'),
                trailing: i.row != null ? Text('row ${i.row}') : null,
              ),
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
}
