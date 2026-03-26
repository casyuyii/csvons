# csvons_gui (starter)

Minimal Flutter starter wiring for running the Go `csvons` validator binary.

## What is included

- `lib/core/validation_runner.dart`: process runner for `csvons`.
- `lib/core/local_state_store.dart`: local persistence for recent binary/ruler/workspace paths.
- `lib/core/issue_filters.dart`: issue filtering/sorting logic extracted for testability.
- `lib/models/validation_report.dart`: JSON models for report parsing (including `schema_version`).
- `lib/screens/home_page.dart`: starter screen to run validation, reuse recent paths, validate file existence before launch, and inspect issues in a searchable/sortable table.
- `lib/screens/workspace_page.dart`: workspace scanner that lists CSV files in a selected directory and supports recent-workspace quick select.
- `test/core/issue_filters_test.dart`: unit tests for issue filtering and sorting behavior.
- `lib/main.dart`: app entry.

## How to use

1. Create/prepare a Flutter project in this folder (`flutter create .` if needed).
2. Build `csvons` binary and place it at a path you can reference from the UI.
3. Start the app and enter:
   - binary path (e.g., `bin/csvons_linux`)
   - absolute `ruler.json` path
4. Click **Run Validation**.

> Note: current Go CLI may not emit JSON yet. If JSON parse fails, raw stdout/stderr is shown.
