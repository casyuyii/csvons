import 'dart:io';

import 'package:file_selector/file_selector.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../core/csv_preview.dart';
import '../core/local_state_store.dart';

class WorkspacePage extends StatefulWidget {
  const WorkspacePage({super.key});

  @override
  State<WorkspacePage> createState() => _WorkspacePageState();
}

class _WorkspacePageState extends State<WorkspacePage> {
  final _stateStore = LocalStateStore();
  final _workspaceController = TextEditingController(text: Directory.current.path);
  bool _loading = false;
  String? _error;
  List<FileSystemEntity> _csvFiles = const <FileSystemEntity>[];
  List<String> _recentWorkspacePaths = const <String>[];
  String? _selectedCsvPath;
  CsvPreview? _preview;
  bool _previewLoading = false;
  int _previewLoadToken = 0;

  @override
  void initState() {
    super.initState();
    _loadState();
    _scanWorkspace();
  }

  Future<void> _loadState() async {
    final state = await _stateStore.load();
    if (!mounted) return;

    setState(() {
      _recentWorkspacePaths = state.recentWorkspacePaths;
      if (_recentWorkspacePaths.isNotEmpty) {
        _workspaceController.text = _recentWorkspacePaths.first;
      }
    });
  }

  @override
  void dispose() {
    _workspaceController.dispose();
    super.dispose();
  }

  Future<void> _pickWorkspaceDirectory() async {
    try {
      final selected = await getDirectoryPath();
      if (selected == null || selected.trim().isEmpty) return;
      setState(() {
        _workspaceController.text = selected;
        _error = null;
      });
      await _scanWorkspace();
    } on PlatformException catch (e) {
      setState(() {
        _error = 'Unable to open workspace picker: ${e.message ?? e.code}';
      });
    }
  }

  Future<void> _scanWorkspace() async {
    final workspacePath = _workspaceController.text.trim();
    if (workspacePath.isEmpty) {
      setState(() => _error = 'Workspace path is required.');
      return;
    }

    setState(() {
      _loading = true;
      _error = null;
      _preview = null;
      _previewLoading = false;
      _selectedCsvPath = null;
    });

    try {
      final dir = Directory(workspacePath);
      if (!await dir.exists()) {
        throw const FileSystemException('Directory does not exist');
      }

      final files = await dir
          .list(followLinks: false)
          .where((entity) {
            if (entity is! File) return false;
            return entity.path.toLowerCase().endsWith('.csv');
          })
          .toList();
      files.sort((a, b) => a.path.compareTo(b.path));
      await _stateStore.saveRecentWorkspace(workspacePath: workspacePath);
      final latestState = await _stateStore.load();

      if (!mounted) return;
      setState(() {
        _csvFiles = files;
        _recentWorkspacePaths = latestState.recentWorkspacePaths;
      });
    } on FileSystemException catch (e) {
      setState(() {
        _error = 'Cannot read workspace: ${e.message}';
        _csvFiles = const <FileSystemEntity>[];
      });
    } finally {
      if (mounted) {
        setState(() => _loading = false);
      }
    }
  }

  Future<void> _selectCsv(String path) async {
    final requestToken = ++_previewLoadToken;
    setState(() {
      _selectedCsvPath = path;
      _preview = null;
      _previewLoading = true;
      _error = null;
    });

    try {
      final preview = await CsvPreview.load(path);
      if (!mounted || requestToken != _previewLoadToken) return;
      setState(() {
        _preview = preview;
        _previewLoading = false;
      });
    } on IOException catch (e) {
      if (!mounted || requestToken != _previewLoadToken) return;
      setState(() {
        _error = 'Unable to preview CSV: $e';
        _previewLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Workspace')),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _workspaceController,
              decoration: InputDecoration(
                labelText: 'Workspace directory path',
                suffixIcon: IconButton(
                  tooltip: 'Browse workspace',
                  icon: const Icon(Icons.folder_open),
                  onPressed: _pickWorkspaceDirectory,
                ),
              ),
              onSubmitted: (_) => _scanWorkspace(),
            ),
            if (_recentWorkspacePaths.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(top: 6),
                child: Wrap(
                  spacing: 8,
                  runSpacing: 8,
                  children: [
                    const Text('Recent:'),
                    ..._recentWorkspacePaths.take(4).map(
                          (path) => ActionChip(
                            label: Text(path, overflow: TextOverflow.ellipsis),
                            onPressed: () {
                              _workspaceController.text = path;
                              _scanWorkspace();
                            },
                          ),
                        ),
                  ],
                ),
              ),
            const SizedBox(height: 12),
            FilledButton.icon(
              onPressed: _loading ? null : _scanWorkspace,
              icon: const Icon(Icons.refresh),
              label: Text(_loading ? 'Scanning...' : 'Scan CSV files'),
            ),
            const SizedBox(height: 12),
            if (_error != null)
              Text('Error: $_error', style: const TextStyle(color: Colors.red)),
            if (_error == null)
              Text('Detected ${_csvFiles.length} CSV file(s) in workspace.'),
            const SizedBox(height: 8),
            Expanded(
              child: Row(
                children: [
                  Expanded(
                    child: _csvFiles.isEmpty
                        ? const _EmptyWorkspaceState()
                        : ListView.builder(
                            itemCount: _csvFiles.length,
                            itemBuilder: (_, index) {
                              final path = _csvFiles[index].path;
                              final selected = path == _selectedCsvPath;
                              return ListTile(
                                selected: selected,
                                dense: true,
                                leading: const Icon(Icons.description_outlined),
                                title: Text(path.split(Platform.pathSeparator).last),
                                subtitle: Text(path),
                                onTap: () => _selectCsv(path),
                              );
                            },
                          ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: _CsvPreviewCard(
                      selectedCsvPath: _selectedCsvPath,
                      preview: _preview,
                      loading: _previewLoading,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _EmptyWorkspaceState extends StatelessWidget {
  const _EmptyWorkspaceState();

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.folder_copy_outlined, size: 42, color: Colors.grey.shade500),
          const SizedBox(height: 8),
          Text('No CSV files found', style: textTheme.titleMedium),
          const SizedBox(height: 4),
          Text(
            'Update the workspace path and scan again.',
            style: textTheme.bodyMedium?.copyWith(color: Colors.grey.shade700),
          ),
        ],
      ),
    );
  }
}

class _CsvPreviewCard extends StatelessWidget {
  const _CsvPreviewCard({
    required this.selectedCsvPath,
    required this.preview,
    required this.loading,
  });

  final String? selectedCsvPath;
  final CsvPreview? preview;
  final bool loading;

  @override
  Widget build(BuildContext context) {
    if (selectedCsvPath == null) {
      return const Card(
        child: Center(
          child: Padding(
            padding: EdgeInsets.all(12),
            child: Text('Select a CSV file to preview header and sample rows.'),
          ),
        ),
      );
    }

    if (loading) {
      return const Card(
        child: Center(child: CircularProgressIndicator()),
      );
    }

    if (preview == null) {
      return const Card(
        child: Center(
          child: Padding(
            padding: EdgeInsets.all(12),
            child: Text('Preview unavailable for selected file.'),
          ),
        ),
      );
    }

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              selectedCsvPath!.split(Platform.pathSeparator).last,
              style: Theme.of(context).textTheme.titleSmall,
            ),
            const SizedBox(height: 4),
            Text(
              'Columns: ${preview!.header.length} · Sample rows: ${preview!.rows.length}',
            ),
            const SizedBox(height: 8),
            Wrap(
              spacing: 6,
              runSpacing: 6,
              children: preview!.header
                  .map((h) => Chip(label: Text(h.isEmpty ? '(empty)' : h)))
                  .toList(growable: false),
            ),
            const Divider(height: 20),
            Expanded(
              child: preview!.rows.isEmpty
                  ? const Center(child: Text('No data rows in this file.'))
                  : ListView.builder(
                      itemCount: preview!.rows.length,
                      itemBuilder: (_, index) {
                        final row = preview!.rows[index];
                        return Padding(
                          padding: const EdgeInsets.only(bottom: 8),
                          child: Text('Row ${index + 1}: ${row.join(' | ')}'),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }
}
