# csvons_gui

Desktop-first Flutter shell for running the Go `csvons` validator binary.

Current V1 scope is intentionally limited to the **Validate** and **Workspace** flows. Ruler editing is deferred until after the desktop validator release candidate is stable.

## What is included

- `lib/core/validation_runner.dart`: process runner for `csvons`.
- `lib/core/local_state_store.dart`: local persistence for recent binary/ruler/workspace/export paths.
- `lib/core/csv_preview.dart`: CSV preview loader/parser for workspace sample inspection (including quoted/multiline cell handling).
- `lib/core/report_exporter.dart`: JSON/Markdown report export helpers for parsed validation results.
- `lib/core/issue_filters.dart`: issue filtering/sorting logic extracted for testability.
- `lib/models/validation_report.dart`: JSON models for report parsing (including `schema_version`).
- `lib/main.dart`: app entry with the current `Validate` + `Workspace` navigation shell.
- `lib/screens/home_page.dart`: starter screen to run validation, browse binary/ruler paths with native picker buttons, reuse/clear recent paths, validate file existence before launch, and inspect issues in a searchable/sortable table.
- `lib/screens/workspace_page.dart`: workspace scanner that lists CSV files, supports workspace directory picking, and includes recent-workspace quick select + preview (with stale-preview guard when switching files quickly).
- `test/core/issue_filters_test.dart`: unit tests for issue filtering and sorting behavior.
- `test/core/csv_preview_test.dart`: unit tests for CSV preview parsing/loading behavior.
- `test/core/report_exporter_test.dart`: unit tests for report export formatting and file writing.
- `test/core/local_state_store_test.dart`: unit tests for local recents persistence (including export paths).
## How to use

1. Install dependencies: `flutter pub get`.
2. If platform folders are missing, generate them locally with `flutter create .`.
3. Run checks:
   - `flutter analyze`
   - `flutter test`
   - or one-shot bootstrap + checks: `./tool/bootstrap.sh`
   - or use make targets: `make check` (or `make analyze`, `make format`, `make test`)
   - remove generated/local artifacts with: `make clean`
   - make targets include a preflight check that prints a clear error if `flutter`/`dart` are missing
4. Build `csvons` and either:
   - point the UI at an explicit binary path, or
   - place a bundled binary under `bin/<platform>/csvons` (`bin/windows/csvons.exe` on Windows).
5. Start the app and enter:
   - binary path (optional if the bundled sidecar path exists)
   - absolute `ruler.json` path
6. Click **Run Validation**.

> Note: current Go CLI may not emit JSON yet. If JSON parse fails, raw stdout/stderr is shown.


## Notes

- Path picker buttons use Flutter's `file_selector` package; ensure it is included in `pubspec.yaml` when initializing this starter.
- `ValidationRunner` resolves binaries in this order: explicit path, bundled sidecar next to the app executable, `bin/<platform>/csvons`, then the legacy dev fallback path.

- Picker failures are surfaced as inline UI errors instead of silently failing.

- Home screen includes a **Clear Recents** action to reset locally stored path history.

- CI workflow for GUI checks lives at `.github/workflows/gui_checks.yml` and runs `make check` (deps/analyze/format/test) plus root-level `go test ./...`.
- CI workflow enables Flutter dependency caching to reduce rerun time.
- Local bootstrap/check helper lives at `tool/bootstrap.sh` and runs `flutter create .`, dependency install, analyze, format check, and tests.
- Local cleanup helper lives at `tool/clean_generated.sh` (also `make clean`) and removes generated platform and tool artifacts.
- Generated Flutter platform folders are intentionally ignored via `.gitignore`; recreate them with `flutter create .` in local/CI environments.
- Bundled validator binaries under `bin/` are intentionally ignored; CI/release builds generate and stage them per platform.

- Project-level finish workflow is documented in `../../docs/flutter_gui_finish_process.md`.
- Packaging location note is documented in `../../docs/flutter_gui_packaging_note.md`.
- Release checklist is documented in `../../docs/flutter_gui_release_checklist.md`.
