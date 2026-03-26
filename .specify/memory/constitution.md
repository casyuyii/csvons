<!--
SYNC IMPACT REPORT
==================
Version change: [TEMPLATE] → 1.0.0
Modified principles: N/A (initial population from template)

Added sections:
  - Core Principles (I–IV)
  - Development Workflow
  - Quality Gates
  - Governance

Removed sections: N/A (template placeholders replaced)

Templates reviewed:
  ✅ .specify/templates/plan-template.md — Constitution Check section present; aligns with principles
  ✅ .specify/templates/spec-template.md — Success Criteria / measurable outcomes align with perf/UX principles
  ✅ .specify/templates/tasks-template.md — Polish phase covers testing, perf, and UX cross-cutting tasks

Deferred TODOs:
  - RATIFICATION_DATE set to today (2026-03-26); update to original adoption date if this predates today.
-->

# csvons Constitution

## Core Principles

### I. Code Quality

All production code MUST be correct, readable, and maintainable in that priority order.
Correctness takes precedence over cleverness or performance. Code MUST:

- Pass all linting and formatting checks (`go vet`, `gofmt`) with zero warnings before merge.
- Use clear, self-documenting identifiers; avoid abbreviations unless domain-standard (e.g., `csv`, `dst`).
- Limit function scope: a function MUST do one thing; functions exceeding ~50 lines SHOULD be decomposed.
- Avoid premature abstractions — introduce shared helpers only when the same logic appears three or more times.
- Treat `internal/csvons` as the canonical library boundary; cross-package dependencies MUST flow inward only.

**Rationale**: csvons is a correctness-critical validation tool. Subtle bugs in constraint evaluation
have downstream data-quality consequences. Readable code ensures contributors can audit logic confidently.

### II. Testing Standards

Automated testing is non-negotiable. The following rules MUST be observed on every feature branch:

- Unit tests MUST cover every public function in `internal/csvons`; coverage MUST NOT drop below 80 %.
- Tests MUST be written before or alongside implementation — no untested code reaches `main`.
- Each validator rule (`exists`, `unique`, `vtype`) MUST have at least one positive test (valid data passes)
  and one negative test (invalid data is correctly rejected).
- Integration tests MUST exercise end-to-end scenarios using real CSV + `ruler.json` fixtures in `testdata/`.
- Tests MUST be deterministic and isolated — no shared mutable state, no reliance on execution order.
- GUI widget tests MUST verify that validation results are rendered correctly (error list, success state).

**Rationale**: The test suite is the primary safety net against regressions in constraint logic.
A test that passes with invalid data is worse than no test; both positive and negative cases are required.

### III. User Experience Consistency

All user-facing surfaces (CLI output and GUI) MUST present information consistently:

- Error messages MUST follow the format: `[file]:[row] field "<field>": <reason>` for validation failures.
- The CLI MUST exit with code `0` on success and non-zero on any validation failure or configuration error.
- The GUI MUST mirror CLI error semantics — the same constraint violation MUST produce equivalent
  human-readable text in both interfaces.
- Terminology MUST be consistent across docs, CLI flags, GUI labels, and `ruler.json` keys
  (e.g., "constraint", "field expression", "ruler" are the canonical terms).
- Breaking changes to `ruler.json` schema MUST be accompanied by a migration guide and a deprecation
  warning in the CLI output for the previous schema version.

**Rationale**: Users interact with csvons through multiple surfaces. Inconsistent language or error formats
create cognitive overhead and erode trust in the tool's output.

### IV. Performance Requirements

Performance is a third-tier priority (after correctness and features) but MUST not regress without justification:

- Validation of a single CSV file up to 100 000 rows MUST complete in under 2 seconds on reference hardware
  (single core, 2 GHz, 512 MB available memory).
- Memory usage MUST scale linearly with input size; unbounded in-memory accumulation of rows is prohibited.
- Performance benchmarks (`go test -bench`) MUST be included for any new validator that processes per-row logic.
- Any change that regresses a benchmark by more than 20 % MUST include a documented justification in the PR.

**Rationale**: csvons targets data-engineering workflows where files can be large. Predictable, linear
performance is more valuable than micro-optimised worst-case speed.

## Development Workflow

All feature work MUST follow this sequence:

1. Create a feature branch from `main` (`###-feature-name` convention).
2. Populate `specs/###-feature-name/spec.md` via `/speckit.specify` before writing any code.
3. Generate `plan.md` via `/speckit.plan`; the Constitution Check section MUST be completed.
4. Generate `tasks.md` via `/speckit.tasks`; tasks MUST be organised by user story.
5. Implement in task order; each task MUST be committed individually with a descriptive message.
6. Open a PR to `main`; PRs MUST reference the feature spec and pass all quality gates below.

Hotfixes (critical correctness bugs) MAY skip spec/plan steps but MUST include a regression test.

## Quality Gates

A PR MUST satisfy all of the following before merge:

- `go vet ./...` — zero issues.
- `gofmt -l .` — zero unformatted files.
- `go test ./... -count=1` — all tests pass.
- Coverage check: `go test ./internal/csvons/... -coverprofile=cov.out && go tool cover -func=cov.out`
  reports ≥ 80 % statement coverage.
- No new `ruler.json` schema keys introduced without a corresponding spec entry and migration note.
- Performance benchmarks pass (no regression > 20 % vs. `main`).

## Governance

This constitution supersedes all implicit conventions and undocumented team norms. Amendments require:

1. A PR that edits this file with a version bump per the semantic versioning policy below.
2. The PR description MUST explain the rationale and list any affected templates or workflows.
3. Approval from at least one other contributor before merge.

**Versioning Policy**:
- MAJOR: Removal or redefinition of a principle in a backward-incompatible way.
- MINOR: New principle or section added; material expansion of existing guidance.
- PATCH: Clarifications, wording fixes, non-semantic refinements.

All PRs and code reviews MUST verify compliance with the principles above. Violations noted during review
are blockers unless a justified exception is recorded in the PR description and the Complexity Tracking
table of `plan.md`.

**Version**: 1.0.0 | **Ratified**: 2026-03-26 | **Last Amended**: 2026-03-26
