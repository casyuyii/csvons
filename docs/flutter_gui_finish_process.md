# Process to finish the Flutter GUI (csvons)

This is the **current recommended process** to finish the desktop-first GUI to a V1-ready state.

## Working agreement (dynamic guide)

Use this file as a **living guide**:

1. Before each GUI change:
   - read `README.md`,
   - read `docs/gui_progress.md`,
   - read this file,
   - read `cmd/csvons/main.go`,
   - read `gui/csvons_gui/lib/models/validation_report.dart`,
   - read `gui/csvons_gui/lib/main.dart`,
   - pick one small checklist item from **Current iteration checklist**.
2. After each GUI change:
   - update the checklist status in this file,
   - add a short “What changed” note in **Iteration log**.
3. Keep increments small and verifiable (`make check` + `go test ./...` where available).
4. Keep V1 scoped to the current desktop-first **Validate + Workspace** app. Do not add ruler editing in this phase.
5. Each handoff increment should have a strict done definition: code changed, tests run, and both progress docs updated.

## 1) Bootstrap a real Flutter project shell in `gui/csvons_gui`

Platform-folder strategy: generated Flutter platform folders remain **untracked** in git and should be recreated with `flutter create .` in local Flutter-enabled environments and in CI.

From repo root:

```bash
cd gui/csvons_gui
flutter create .
flutter pub get
```

Then run local checks:

```bash
make check
```

## 2) Stabilize report/runner contract end-to-end

1. Verify Go JSON output schema fields are populated consistently (`file`, `rule`, `field`, `row`, `value`, `severity`, `message`).
2. Preserve exit-code meaning: `1` for validation failures, `2` for runtime/config failures.
3. Validate mapping in Flutter models and runner parsing for runtime/config failures.
4. Add/refresh fixtures so Flutter tests cover all known report variants.

## 3) Complete UX pass for issues triage and workspace flow

1. Keep improving issues table workflows (filters, sort, empty-state recovery, summary clarity).
2. Add any missing “issue to CSV context” affordances needed for analysts.
3. Validate workspace screen behavior with realistic CSV sets.

## 4) Enforce quality gates in CI

CI already runs GUI checks and Go tests. Keep these green on each PR:

- GUI: `make check`
- Go: `go test ./...`

If formatting/lint policies change, update both:

- `gui/csvons_gui/Makefile`
- `.github/workflows/gui_checks.yml`

## 5) Packaging + release prep (desktop-first)

Packaging location note: see `docs/flutter_gui_packaging_note.md`.
Release checklist: see `docs/flutter_gui_release_checklist.md`.

1. Build and bundle per-OS `csvons` binary with Flutter desktop artifacts.
2. Create a release checklist (build, smoke test, signing/notarization where needed).
3. Add one documented “fresh machine” verification pass.

## 6) Merge cadence

1. Implement a small increment.
2. Run checks (`make check`, `go test ./...`).
3. Open/merge PR.
4. Rebase or restart from latest `main`.

Repeat until the remaining checklist in `docs/gui_progress.md` is complete.

---

## Current iteration checklist

- [x] Commit the generated-project scaffolding strategy: platform folders remain untracked and are regenerated via `flutter create .` in bootstrap/CI.
- [ ] Execute a full GUI check run (`make check`) in an environment with Flutter installed and record results.
- [x] Harden the Go JSON report contract so recovered issues include `file`/`rule`/`field`/`row`/`value` and preserve validation-vs-runtime exit codes.
- [x] Add at least one widget test for table column sorting interaction (tap DataTable column headers and assert ordering changes).
- [x] Provide a repeatable cleanup command for generated Flutter artifacts (`make clean` / `tool/clean_generated.sh`).
- [x] Add a short packaging note describing where bundled `csvons` binaries will live for desktop builds.
- [x] Add a CI desktop build path plus a sidecar-staging helper for the bundled validator.
- [ ] Validate CI workflow behavior on one real PR run and capture follow-up fixes.

## Iteration log

- 2026-03-27: Added this dynamic-guide protocol and checklist so each future increment updates the process doc as part of the change.
- 2026-03-27: Added widget test coverage for DataTable value-column sort interaction and marked that checklist item complete.
- 2026-03-27: Added generated-artifact cleanup helper (`make clean`) and marked cleanup checklist item complete.
- 2026-03-27: Added packaging location note (`docs/flutter_gui_packaging_note.md`) and marked packaging-note checklist item complete.
- 2026-03-27: Hardened Go JSON issue output with structured context (`file`, `rule`, `field`, `row`, `value`) and distinct validation/runtime exit codes.
- 2026-03-27: Updated the Flutter runner to prefer bundled `bin/<platform>/csvons` sidecars and added a Linux sidecar-staging helper plus CI desktop bundle build.
