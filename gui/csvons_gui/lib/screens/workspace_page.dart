import 'dart:io';

import 'package:flutter/material.dart';

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

  Future<void> _scanWorkspace() async {
    final workspacePath = _workspaceController.text.trim();
    if (workspacePath.isEmpty) {
      setState(() => _error = 'Workspace path is required.');
      return;
    }

    setState(() {
      _loading = true;
      _error = null;
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
              decoration: const InputDecoration(
                labelText: 'Workspace directory path',
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
              child: _csvFiles.isEmpty
                  ? const _EmptyWorkspaceState()
                  : ListView.builder(
                      itemCount: _csvFiles.length,
                      itemBuilder: (_, index) {
                        final path = _csvFiles[index].path;
                        return ListTile(
                          dense: true,
                          leading: const Icon(Icons.description_outlined),
                          title: Text(path.split(Platform.pathSeparator).last),
                          subtitle: Text(path),
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
