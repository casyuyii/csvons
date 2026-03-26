# csvons_gui (starter)

Minimal Flutter starter wiring for running the Go `csvons` validator binary.

## What is included

- `lib/core/validation_runner.dart`: process runner for `csvons`.
- `lib/core/local_state_store.dart`: local persistence for recent binary/ruler/workspace/export paths.
- `lib/core/csv_preview.dart`: CSV preview loader/parser for workspace sample inspection (including quoted/multiline cell handling).
- `lib/core/report_exporter.dart`: JSON/Markdown report export helpers for parsed validation results.
- `lib/core/issue_filters.dart`: issue filtering/sorting logic extracted for testability.
- `lib/models/validation_report.dart`: JSON models for report parsing (including `schema_version`).
- `lib/screens/home_page.dart`: starter screen to run validation, browse binary/ruler paths with native picker buttons, reuse/clear recent paths, validate file existence before launch, and inspect issues in a searchable/sortable table.
- `lib/screens/workspace_page.dart`: workspace scanner that lists CSV files, supports workspace directory picking, and includes recent-workspace quick select + preview (with stale-preview guard when switching files quickly).
- `test/core/issue_filters_test.dart`: unit tests for issue filtering and sorting behavior.
- `test/core/csv_preview_test.dart`: unit tests for CSV preview parsing/loading behavior.
- `test/core/report_exporter_test.dart`: unit tests for report export formatting and file writing.
- `test/core/local_state_store_test.dart`: unit tests for local recents persistence (including export paths).
- `lib/main.dart`: app entry.

## How to use

1. Create/prepare a Flutter project in this folder (`flutter create .` if needed).
2. Build `csvons` binary and place it at a path you can reference from the UI.
3. Start the app and enter:
   - binary path (e.g., `bin/csvons_linux`)
   - absolute `ruler.json` path
4. Click **Run Validation**.

> Note: current Go CLI may not emit JSON yet. If JSON parse fails, raw stdout/stderr is shown.


## Notes

- Path picker buttons use Flutter's `file_selector` package; ensure it is included in `pubspec.yaml` when initializing this starter.

- Picker failures are surfaced as inline UI errors instead of silently failing.

- Home screen includes a **Clear Recents** action to reset locally stored path history.
