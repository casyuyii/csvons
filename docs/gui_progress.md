# GUI Build Progress (csvons + Flutter)

## Current status

Estimated completion: **~55%** of the planned desktop-first V1 scope.

## What is already done

### Go CLI integration contract
- Added CLI flags: `--format`, `--output`, `--quiet`.
- Added structured report model (`summary` + `issues`) and JSON/text output emission.
- Added stable run wrappers (`run`, `runWithArgs`) and recovery path for structured failure output.
- Added tests for output and `runWithArgs` behavior (success, invalid format, missing args/config, validation failure).

### Validation error plumbing
- Replaced validator hard exits with recoverable `failf(...)` panics.
- Added shared helper `requiredFieldValues(...)` to fail early on invalid field expressions/nil channels.
- Added tests for helper panic behavior.

### Flutter starter
- Added starter app shell, process runner (`ValidationRunner`), and report models.
- Added local persistence for recent binary/ruler paths.
- Added results filtering and sortable issues table.

## Remaining steps to finish V1 (desktop-first)

1. **Go report schema hardening** *(in progress)*
   - ✅ Added richer issue fields in report model (`file`, `rule`, `field`, `row`, `value`).
   - ✅ Added recovered-failure context population for `file`/`rule`, plus best-effort `field` extraction from validator messages.
   - ⏳ Ensure these fields are populated consistently for every validator failure path.
   - ⏳ Add schema/version note to avoid future breaking changes.

2. **Flutter project productionization**
   - Convert starter folder into fully initialized Flutter project structure and dependency config.
   - Add lint/format/analyzer CI checks.

3. **UX completeness**
   - Add file/folder pickers for ruler and binary/workspace.
   - Add a dedicated workspace screen and better empty/error states.

4. **Test coverage for GUI layer**
   - Add Dart unit tests for runner/model parsing and widget tests for results views/filtering.

5. **Packaging and release pipeline**
   - Bundle per-OS Go binaries with Flutter desktop artifacts.
   - Add one end-to-end desktop build in CI and a release checklist.

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
