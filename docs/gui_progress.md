# GUI Build Progress (csvons + Flutter)

## Current status

Estimated completion: **~75%** of the planned desktop-first V1 scope.

## What is already done

### Go CLI integration contract
- Added CLI flags: `--format`, `--output`.
- Added structured report model (`summary` + `issues`) and JSON/text output emission.
- Hardened recovered failure output so JSON issues now carry `file`, `rule`, `field`, `row`, `value`, `message`, and `severity`.
- Added `runWithArgs(...)` so CLI behavior can be tested without mutating global flag state.
- Validation failures now preserve exit code `1`, while runtime/config failures preserve exit code `2`.
- Added tests for output emission plus `runWithArgs(...)` validation/runtime report behavior.

### Validation error plumbing
- Replaced validator hard exits with recoverable `failf(...)` panics.
- Added shared helper `requiredFieldValues(...)` to fail early on invalid field expressions/nil channels.
- Added tests for helper panic behavior.

### Flutter starter
- Added starter app shell, process runner (`ValidationRunner`), and report models.
- Added local persistence for recent binary/ruler paths.
- Improved home screen UX with:
  - recent-path quick-select chips,
  - run status banner (pass/issues/runtime error),
  - searchable/filterable/sortable issues table for JSON report output,
  - search now includes issue `value`, `row`, and `severity` in addition to message/file/rule/field,
  - issues table now includes a sortable `value` column for direct payload ordering,
  - value sorting now places null/empty values last for cleaner scans,
  - semantic severity sorting order (`critical/fatal` -> `error` -> `warning` -> `info` -> others),
  - deterministic tie-break sorting (row/message) when primary sort values are equal,
  - row sorting places unknown/null row indices last for clearer ordering,
  - file/rule filter controls with reset for faster triage,
  - quick, dynamic severity chips with scope-aware counts for one-click issue slicing,
  - severity chips follow semantic order (`critical/fatal`, `error`, `warning`, `info`, others),
  - reset control now disables when no active filters are applied.
  - explicit empty-report message when a validation run returns zero issues.
  - active-filter summary line explains exactly which filters/search are applied.
  - active-filter summary now includes visible/total issue counts for quick context sharing.
  - copy-to-clipboard action for active filter summary text.
  - copy action now confirms success with a brief snackbar message.
  - issue count label now shows both visible and total issues (`showing X of Y`).
  - empty-filter state now includes a one-click `Reset all filters` action.
  - pre-run path existence validation and clearer empty state before first run,
  - empty-filter feedback when no issues match current table filters.
  - report export controls (JSON + Markdown) from parsed validation results.
  - recent export path history for quicker repeated exports.
  - Clear Recents action for resetting local path history from the validate screen.
  - quote-aware CSV preview parsing with multiline-field support.
  - bundled-validator resolution that prefers packaged `bin/<platform>/csvons` sidecars before legacy dev fallback paths.

## Remaining steps to finish V1 (desktop-first)

1. **Go report schema hardening** *(in progress)*
   - ✅ Added richer issue fields in report model (`file`, `rule`, `field`, `row`, `value`).
   - ✅ Added structured validator/runtime failures so every recovered issue emits consistent metadata.
   - ✅ Added explicit JSON `schema_version` to the report contract.
   - ✅ Split exit-code handling so validation issues remain `1` and runtime/config failures remain `2`.

2. **Flutter project productionization**
   - ✅ Added starter `pubspec.yaml` + `analysis_options.yaml` with lint baseline and dependency declarations.
   - ✅ Added GitHub Actions workflow (`.github/workflows/gui_checks.yml`) for GUI checks (`make check`) and root Go tests.
   - ✅ Enabled Flutter dependency caching in CI workflow for faster repeated runs.
   - ✅ CI Flutter job now runs `make check` to mirror local developer checks.
   - ✅ Added CI formatting check (`dart format --set-exit-if-changed lib test`).
   - ✅ Added local bootstrap helper (`gui/csvons_gui/tool/bootstrap.sh`) to run create/pub/analyze/format/test in one command.
   - ✅ Added `gui/csvons_gui/Makefile` convenience targets (`make check`, `make analyze`, `make format`, `make test`).
   - ✅ Added GUI cleanup command (`make clean`) for generated tool/platform artifacts.
   - ✅ Added Makefile preflight checks with explicit missing `flutter`/`dart` error messages.
   - ✅ Added GUI module `.gitignore` for Flutter tool outputs and generated platform directories.
   - ✅ Committed the scaffolding strategy: generated Flutter platform folders stay untracked and are regenerated via `flutter create .` in bootstrap/CI.
   - ⏳ Execute a full local `flutter create .` + `make check` pass in an environment with Flutter installed and record results.

3. **UX completeness**
   - ✅ Added file/folder picker buttons for ruler, binary, and workspace paths.
   - ✅ Added picker error handling with inline UI feedback.
   - ✅ Added a dedicated workspace screen with CSV discovery, empty/error states, and recent-workspace quick select.
   - ✅ Added CSV header + sample-row preview panel when selecting a workspace file.
   - ✅ Added stale preview guard for fast file switching in workspace list.

4. **Test coverage for GUI layer**
   - ✅ Added first Dart unit tests for issues filtering/sorting logic.
   - ✅ Added Dart unit tests for runner/model parsing.
   - ✅ Added Dart unit tests for CSV preview parsing/loading.
   - ✅ Added first widget tests for issues-table results filtering and empty-filter state messaging.
   - ✅ Added widget test coverage for DataTable value-column sorting interaction.

5. **Packaging and release pipeline**
   - ✅ Added initial packaging location note (`docs/flutter_gui_packaging_note.md`) for bundled per-OS `csvons` binaries.
   - ✅ Updated Flutter runner resolution to prefer bundled sidecar binaries before legacy dev defaults.
   - ✅ Added Linux sidecar staging helper (`gui/csvons_gui/tool/stage_bundle_binary.sh`).
   - ✅ Added one Linux desktop bundle build/smoke-test path to CI.
   - ✅ Added release checklist doc (`docs/flutter_gui_release_checklist.md`).
   - ⏳ Validate the new CI desktop build on one real PR run and capture any follow-up fixes.

## Finish estimate

If worked sequentially, this is about **5 major steps** left to reach a practical V1 desktop release candidate.

## Collaboration workflow (recommended)

- **Yes** — create a PR from your branch and merge it (or rebase onto the latest merged branch) before asking for the next increment.
- I work **incrementally from the current checked-out git state** in this environment.
- If you ask me to continue before merging, I will keep building on the current branch/commit chain here.
- To avoid drift from your GitHub `master` branch:
  1. Merge the current PR.
  2. Update local branch from `master` (`git fetch && git checkout master && git pull`).
  3. Start the next task/PR from that updated head.

See also: `docs/flutter_gui_finish_process.md` for the current step-by-step **dynamic** finish process/checklist.
