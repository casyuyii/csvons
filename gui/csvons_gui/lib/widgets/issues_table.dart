import 'package:flutter/material.dart';

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
    final issues = filterAndSortIssues(
      issues: widget.report.issues,
      query: _searchController.text,
      severityFilter: _severityFilter,
      sortField: _sortField,
      ascending: _sortAscending,
    );

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Expanded(
              child: TextField(
                controller: _searchController,
                decoration: InputDecoration(
                  isDense: true,
                  prefixIcon: const Icon(Icons.search),
                  labelText: 'Search message/file/rule/field',
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
            const SizedBox(width: 12),
            DropdownButton<String>(
              value: _severityFilter,
              items: const [
                DropdownMenuItem(value: 'all', child: Text('All severities')),
                DropdownMenuItem(value: 'error', child: Text('Error only')),
                DropdownMenuItem(value: 'warning', child: Text('Warning only')),
              ],
              onChanged: (value) {
                if (value == null) return;
                setState(() => _severityFilter = value);
              },
            ),
          ],
        ),
        const SizedBox(height: 8),
        Text('Showing ${issues.length} issue(s)'),
        const SizedBox(height: 8),
        if (issues.isEmpty)
          const Padding(
            padding: EdgeInsets.symmetric(vertical: 8),
            child: Text('No issues match the current filters.'),
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
              _column('Message', 5, IssueSortField.message),
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
}
