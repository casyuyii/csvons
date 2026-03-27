import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../core/issue_filters.dart';
import '../models/validation_report.dart';

class IssuesTable extends StatefulWidget {
  const IssuesTable({required this.report, super.key});

  final ValidationReport report;

  @override
  State<IssuesTable> createState() => _IssuesTableState();
}

class _IssuesTableState extends State<IssuesTable> {
  final _searchController = TextEditingController();
  String _severityFilter = 'all';
  String _fileFilter = 'all';
  String _ruleFilter = 'all';
  int _sortColumnIndex = 0;
  bool _sortAscending = true;
  IssueSortField _sortField = IssueSortField.severity;

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (widget.report.issues.isEmpty) {
      return const Padding(
        padding: EdgeInsets.symmetric(vertical: 8),
        child: Text('No issues reported for this run.'),
      );
    }

    final scopedIssues = filterAndSortIssues(
      issues: widget.report.issues,
      query: _searchController.text,
      severityFilter: 'all',
      fileFilter: _fileFilter,
      ruleFilter: _ruleFilter,
      sortField: _sortField,
      ascending: _sortAscending,
    );
    final severityCounts = <String, int>{};
    for (final issue in scopedIssues) {
      final normalized = issue.severity.toLowerCase();
      severityCounts[normalized] = (severityCounts[normalized] ?? 0) + 1;
    }
    final severityOptions = severityCounts.keys.toList();
    severityOptions.sort((a, b) {
      final rankDelta = _severityChipRank(a).compareTo(_severityChipRank(b));
      if (rankDelta != 0) return rankDelta;
      return a.compareTo(b);
    });

    final issues = filterAndSortIssues(
      issues: widget.report.issues,
      query: _searchController.text,
      severityFilter: _severityFilter,
      fileFilter: _fileFilter,
      ruleFilter: _ruleFilter,
      sortField: _sortField,
      ascending: _sortAscending,
    );

    final fileValues = widget.report.issues
        .map((issue) => issue.file)
        .toSet()
        .toList()
      ..sort();
    final ruleValues = widget.report.issues
        .map((issue) => issue.rule)
        .toSet()
        .toList()
      ..sort();

    final hasActiveFilter =
        _severityFilter != 'all' ||
        _fileFilter != 'all' ||
        _ruleFilter != 'all' ||
        _searchController.text.isNotEmpty;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Wrap(
          spacing: 12,
          runSpacing: 8,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            LayoutBuilder(
              builder: (context, constraints) => ConstrainedBox(
                constraints: BoxConstraints(
                  maxWidth: constraints.maxWidth > 420 ? 420 : constraints.maxWidth,
                ),
                child: TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    isDense: true,
                    prefixIcon: const Icon(Icons.search),
                    labelText: 'Search message/file/rule/field/value/row/severity',
                    suffixIcon: _searchController.text.isEmpty
                        ? null
                        : IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: () {
                              _searchController.clear();
                              setState(() {});
                            },
                          ),
                  ),
                  onChanged: (_) => setState(() {}),
                ),
              ),
            ),
            DropdownButton<String>(
              value: _fileFilter,
              items: <DropdownMenuItem<String>>[
                const DropdownMenuItem(value: 'all', child: Text('All files')),
                ...fileValues.map(
                  (value) => DropdownMenuItem(value: value, child: Text(value)),
                ),
              ],
              onChanged: (value) {
                if (value == null) return;
                setState(() => _fileFilter = value);
              },
            ),
            DropdownButton<String>(
              value: _ruleFilter,
              items: <DropdownMenuItem<String>>[
                const DropdownMenuItem(value: 'all', child: Text('All rules')),
                ...ruleValues.map(
                  (value) => DropdownMenuItem(value: value, child: Text(value)),
                ),
              ],
              onChanged: (value) {
                if (value == null) return;
                setState(() => _ruleFilter = value);
              },
            ),
            TextButton(
              onPressed: hasActiveFilter
                  ? () {
                      _searchController.clear();
                      setState(() {
                        _severityFilter = 'all';
                        _fileFilter = 'all';
                        _ruleFilter = 'all';
                      });
                    }
                  : null,
              child: const Text('Reset filters'),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: [
            FilterChip(
              label: Text('All (${scopedIssues.length})'),
              selected: _severityFilter == 'all',
              onSelected: (_) => setState(() => _severityFilter = 'all'),
            ),
            ...severityOptions.map(
              (severity) => FilterChip(
                label:
                    Text('${_formatSeverityLabel(severity)} (${severityCounts[severity]})'),
                selected: _severityFilter == severity,
                onSelected: (_) => setState(() => _severityFilter = severity),
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        Row(
          children: [
            Text('Showing ${issues.length} of ${widget.report.issues.length} issue(s)'),
            if (hasActiveFilter) ...[
              const SizedBox(width: 10),
              const Text('(filtered)'),
            ],
          ],
        ),
        if (hasActiveFilter) ...[
          const SizedBox(height: 4),
          Row(
            children: [
              Expanded(
                child: Text(
                  _activeFilterSummary(visibleCount: issues.length),
                  style: Theme.of(context).textTheme.bodySmall,
                ),
              ),
              IconButton(
                tooltip: 'Copy active filter summary',
                icon: const Icon(Icons.content_copy, size: 18),
                onPressed: () async {
                  await Clipboard.setData(
                    ClipboardData(
                      text: _activeFilterSummary(visibleCount: issues.length),
                    ),
                  );
                  if (!mounted) return;
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                      content: Text('Active filter summary copied'),
                      duration: Duration(milliseconds: 1200),
                    ),
                  );
                },
              ),
            ],
          ),
        ],
        const SizedBox(height: 8),
        if (issues.isEmpty)
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 8),
            child: Row(
              children: [
                const Text('No issues match the current filters.'),
                if (hasActiveFilter) ...[
                  const SizedBox(width: 8),
                  TextButton(
                    onPressed: () {
                      _searchController.clear();
                      setState(() {
                        _severityFilter = 'all';
                        _fileFilter = 'all';
                        _ruleFilter = 'all';
                      });
                    },
                    child: const Text('Reset all filters'),
                  ),
                ],
              ],
            ),
          ),
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          child: DataTable(
            sortColumnIndex: _sortColumnIndex,
            sortAscending: _sortAscending,
            columns: [
              _column('Severity', 0, IssueSortField.severity),
              _column('File', 1, IssueSortField.file),
              _column('Rule', 2, IssueSortField.rule),
              _column('Field', 3, IssueSortField.field),
              _column('Row', 4, IssueSortField.row),
              _column('Value', 5, IssueSortField.value),
              _column('Message', 6, IssueSortField.message),
            ],
            rows: issues
                .map(
                  (i) => DataRow(
                    cells: [
                      DataCell(Text(i.severity)),
                      DataCell(Text(i.file)),
                      DataCell(Text(i.rule)),
                      DataCell(Text(i.field)),
                      DataCell(Text(i.row?.toString() ?? '-')),
                      DataCell(Text(i.value ?? '-')),
                      DataCell(SizedBox(width: 320, child: Text(i.message))),
                    ],
                  ),
                )
                .toList(growable: false),
          ),
        ),
      ],
    );
  }

  DataColumn _column(
    String label,
    int index,
    IssueSortField field,
  ) {
    return DataColumn(
      label: Text(label),
      onSort: (_, ascending) {
        setState(() {
          _sortColumnIndex = index;
          _sortAscending = ascending;
          _sortField = field;
        });
      },
    );
  }

  String _formatSeverityLabel(String value) {
    final parts = value
        .trim()
        .split(RegExp(r'[^A-Za-z0-9]+'))
        .where((part) => part.isNotEmpty);
    if (parts.isEmpty) return value;

    return parts
        .map((part) => '${part[0].toUpperCase()}${part.substring(1).toLowerCase()}')
        .join(' ');
  }

  String _activeFilterSummary({required int visibleCount}) {
    final parts = <String>[];
    parts.add('visible=$visibleCount/${widget.report.issues.length}');
    if (_severityFilter != 'all') {
      parts.add('severity=${_formatSeverityLabel(_severityFilter)}');
    }
    if (_fileFilter != 'all') {
      parts.add('file=$_fileFilter');
    }
    if (_ruleFilter != 'all') {
      parts.add('rule=$_ruleFilter');
    }
    final query = _searchController.text.trim();
    if (query.isNotEmpty) {
      parts.add('search="$query"');
    }
    return 'Active filters: ${parts.join(', ')}';
  }

  int _severityChipRank(String severity) {
    switch (severity.toLowerCase()) {
      case 'critical':
      case 'fatal':
        return 0;
      case 'error':
        return 1;
      case 'warning':
        return 2;
      case 'info':
        return 3;
      default:
        return 4;
    }
  }
}
