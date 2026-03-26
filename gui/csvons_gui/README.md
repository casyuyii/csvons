# csvons_gui (starter)

Minimal Flutter starter wiring for running the Go `csvons` validator binary.

## What is included

- `lib/core/validation_runner.dart`: process runner for `csvons`.
- `lib/core/local_state_store.dart`: local persistence for recent binary/ruler paths.
- `lib/models/validation_report.dart`: basic JSON models for report parsing.
- `lib/screens/home_page.dart`: starter screen to run validation, reuse recent paths, filter issues, and sort issues in a table.
- `lib/main.dart`: app entry.

## How to use

1. Create/prepare a Flutter project in this folder (`flutter create .` if needed).
2. Build `csvons` binary and place it at a path you can reference from the UI.
3. Start the app and enter:
   - binary path (e.g., `bin/csvons_linux`)
   - absolute `ruler.json` path
4. Click **Run Validation**.

> Note: Go CLI now supports `--format json`, but runner still keeps fallback to raw stdout/stderr for resilience.
> Tip: use `--quiet` with JSON mode to suppress validator logs during GUI integration.

## What’s next (recommended order)

1. **Stabilize Go output contract**
   - Keep `--format json|text` and `--output <path>` backward compatible.
   - Define/lock `ValidationReport` and `ValidationIssue` fields.
   - Keep stable exit code behavior (`0` pass, `1` validation failures, `2` runtime errors).

2. **Turn starter into proper Flutter app**
   - Run `flutter create .` in this folder (if not already initialized).
   - Add `pubspec.yaml` dependencies and desktop platform folders.
   - Add linting/formatting configs.

3. **Improve UX**
   - Add file/folder pickers for binary and ruler paths.
   - Add workspace screen and persisted recent paths.
   - Add sortable/filterable issues table.

4. **Add tests**
   - Unit tests for `ValidationRunner` (mock process execution).
   - Model parsing tests for valid/invalid JSON payloads.
   - Widget tests for run button states + result rendering.

5. **Packaging**
   - Add per-platform binary bundling strategy.
   - Validate one desktop target end-to-end in CI.
   - Add release checklist (signing/notarization as needed).
