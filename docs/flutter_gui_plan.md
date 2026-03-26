# Flutter GUI Plan for csvons

## 1) What this project currently is

`csvons` is a Go CLI tool that validates CSV files against rule definitions in a JSON config ("ruler").

- Entry point: `cmd/csvons/main.go` accepts a single argument (`<ruler.json>`) and runs validation.  
- Validation rules:
  - `exists`: cross-file foreign-key-like checks.
  - `unique`: uniqueness checks in one CSV.
  - `vtype`: type/range checks (int/float64/bool).
- Config parsing and CSV reading are in `internal/csvons/utils.go`.
- Rule and metadata data models are in `internal/csvons/types.go`.

This means a GUI should primarily solve:
1. Config creation/editing (`ruler.json`).
2. CSV file selection and preview.
3. Running validation and presenting failures clearly.

## 2) Recommended GUI strategy (Flutter desktop + mobile)

Build a Flutter app as the primary GUI and keep the Go validator as the source of truth.

### Should we translate core validation code to Dart?

Short answer: **No for V1**. Keep Go core code and call it from Flutter.

- **Keep Go core (recommended initially):**
  - Preserves existing, tested validation behavior.
  - Avoids dual-maintenance and logic drift across languages.
  - Delivers GUI faster with lower regression risk.
- **Translate to Dart (optional future):**
  - Consider only if native mobile/offline execution of validator becomes a hard requirement.
  - If pursued, do it incrementally with shared test fixtures and parity tests against Go outputs.

### Integration approach

Use **Process-based integration** first:
- Bundle `csvons` binaries for each desktop platform (Windows/macOS/Linux).
- Flutter launches the binary with a selected `ruler.json` path.
- Parse stdout/stderr logs and render structured results in the UI.

Why first:
- Fastest path to production.
- Reuses current CLI with minimal risk.
- Avoids early complexity of `cgo`/FFI.

Later (optional): add a JSON-output mode in Go for stable machine-readable results.

## 3) Product scope for V1

### User workflows
1. **Open workspace** (folder with CSVs).
2. **Create/edit ruler** with form-based editor + raw JSON mode.
3. **Map fields** via headers discovered from CSV previews.
4. **Run validation** and view grouped errors by file/rule/field/value.
5. **Export report** (JSON and Markdown).

### Non-goals (V1)
- Cloud sync and collaboration.
- Real-time streaming for extremely large datasets.
- Full mobile execution of Go binary (desktop-first recommended).

## 4) Proposed architecture

## 4.1 Flutter layers

- **Presentation**: screens/widgets (workspace, ruler editor, run results).
- **Application**: state + use-cases (`riverpod` or `bloc`).
- **Domain**: entities (`Workspace`, `RulerConfig`, `ValidationIssue`).
- **Infrastructure**:
  - Local file access (file_picker/path_provider).
  - CSV preview parser (Dart csv package).
  - Runner service for executing Go validator process.

## 4.2 Go interaction contract

### Short-term contract
- Input: path to `ruler.json`.
- Output: process exit code + logs.

### Recommended small Go enhancement
Add optional flags:
- `--format json|text` (default text).
- `--output <path>` optional output report path.

JSON output schema suggestion:
- `summary`: total files checked, passed/failed counts, duration.
- `issues[]`: `{file, rule, field, row, value, message, severity}`.

### 4.2.1 How core logic should be implemented (recommended)

Yes — for V1, implement core validation by building Go binaries and calling them from Flutter.

#### Go side (keep validation in Go)

1. Keep existing validators (`exists`, `unique`, `vtype`) in `internal/csvons`.
2. Add a result model in Go:
   - `ValidationReport { summary, issues }`
   - `ValidationIssue { file, rule, field, row, value, message, severity }`
3. Add CLI flags:
   - `--format text|json`
   - `--output <path>` (optional; default stdout)
4. Ensure exit codes are stable:
   - `0`: validation success (no issues)
   - `1`: validation issues found
   - `2`: runtime/config/system error

#### Flutter side (call Go process)

Create a `ValidationRunner` service that:
1. Resolves platform binary path (Windows/macOS/Linux).
2. Writes/uses the selected `ruler.json`.
3. Executes the process with `Process.start(...)`.
4. Reads `stdout`/`stderr`.
5. If `--format json`, decode to `ValidationReport`; otherwise show raw logs.
6. Maps exit code to UI status (Passed / Failed / Error).

Minimal flow:
- User clicks **Run Validation**
- Flutter executes: `csvons --format json /path/to/ruler.json`
- Flutter parses report and renders summary/issues table

#### Dart runner skeleton (concept)

```dart
final proc = await Process.start(binaryPath, [
  '--format', 'json',
  rulerPath,
]);
final out = await proc.stdout.transform(utf8.decoder).join();
final err = await proc.stderr.transform(utf8.decoder).join();
final code = await proc.exitCode;
```

#### Why this is the right implementation for now

- No duplication of validation logic.
- Reuses existing tested Go behavior.
- Faster delivery for desktop GUI.
- Easier to add Dart unit tests around orchestration while Go keeps rule correctness.

### 4.2.2 When to consider moving logic out of process execution

Only consider full Dart port or FFI/plugin migration if one of these becomes mandatory:
- Native on-device mobile validation without server/desktop helper.
- Very low-latency, high-frequency validation where process startup overhead is unacceptable.
- Tight platform integration constraints that process execution cannot satisfy.

### 4.2.3 Native method options (Flutter)

If by “native method” you mean not spawning an external CLI process, there are two alternatives:

1. **Flutter plugin + FFI/native bridge (keep Go logic)**
   - Build a platform plugin that calls a native library.
   - Practical note: Go does not plug into Flutter FFI as smoothly as C/C++/Rust; per-platform wrappers are usually needed.
   - Pros:
     - Can reduce process-start overhead.
     - Better embedding into app lifecycle.
   - Cons:
     - Higher build complexity (especially Windows/macOS/Linux differences).
     - Harder debugging/release signing/tooling.
     - More maintenance than process runner.

2. **Rewrite validator in Dart (pure Flutter-native logic)**
   - Re-implement `exists`, `unique`, `vtype`, and field-expression behavior in Dart.
   - Pros:
     - Single language stack in app layer.
     - Easiest path to true mobile-native execution.
   - Cons:
     - Highest risk of rule-parity bugs.
     - Requires strong regression suite against Go outputs.
     - Duplicates maintenance across implementations unless Go is retired.

#### Recommendation

- **V1:** Process runner (Flutter -> Go binary) for fastest, safest delivery.
- **V2+ only if needed:** evaluate plugin/FFI or Dart rewrite based on measured constraints (mobile requirement, latency, ops complexity).

## 4.3 Packaging

- Build per-OS binaries during CI.
- Place binaries in Flutter desktop assets or sidecar install dir.
- Runtime binary selection by `Platform.operatingSystem`.

## 5) UX design outline

### Main navigation
- **Workspace**
- **Rules (Ruler Editor)**
- **Validate**
- **Reports**
- **Settings**

### Key screens

1. **Workspace Screen**
   - Select CSV folder.
   - Auto-detect CSV files.
   - Display columns from header row based on metadata (`name_index`, `data_index`).

2. **Ruler Editor Screen**
   - Metadata section form.
   - Rule builders:
     - Exists rule builder with source/destination dropdowns.
     - Unique fields selector.
     - VType builder with type and range constraints.
   - Advanced field-expression helper (simple/array/nested/complex).

3. **Validation Results Screen**
   - Run button and progress indicator.
   - Summary cards.
   - Issues table with filtering (file/rule/severity/field).
   - Click-to-jump to CSV row preview.

4. **Reports Screen**
   - Export JSON/Markdown.
   - Save last N runs.

## 6) Delivery plan (phased)

## Phase 0: Foundation (1 week)
- Create Flutter project (`csvons_gui`) with desktop targets.
- App shell, routing, theme, error handling.
- Local workspace selection and persistence.

## Phase 1: Ruler authoring (1–2 weeks)
- Data models + JSON serialization.
- Metadata form + validations.
- Rule editors for exists/unique/vtype.
- Raw JSON view with round-trip validation.

## Phase 2: Validation execution (1 week)
- Process runner service for bundled Go binary.
- Parse logs into structured UI messages.
- Summary + issues list.

## Phase 3: CSV preview + mapping assist (1–2 weeks)
- Header detection and first-N-row preview.
- Field-expression assistant and validation hints.
- Better diagnostics linking issue -> CSV cell context.

## Phase 4: Packaging & release (1 week)
- CI pipelines for Flutter desktop artifacts.
- Go binary build matrix and bundling.
- Installer/notarization/signing tasks per OS.

## 7) Suggested repo layout

```text
csvons/
  cmd/csvons/
  internal/csvons/
  gui/
    csvons_gui/
      lib/
      test/
      pubspec.yaml
  docs/
    flutter_gui_plan.md
```

## 8) Risks and mitigations

- **Risk:** Parsing plain text logs is brittle.  
  **Mitigation:** Add JSON output mode in Go early.

- **Risk:** Mobile platforms cannot easily execute bundled Go binary.  
  **Mitigation:** Desktop-first V1; mobile becomes config editor + remote runner later.

- **Risk:** Field expressions are powerful but hard for users.  
  **Mitigation:** Guided rule-builder UI + expression examples + inline validation.

## 9) Immediate next tasks (actionable)

1. Add Go flag support: `--format json` and machine-readable issue output.
2. Bootstrap Flutter desktop app with workspace + ruler JSON open/save.
3. Implement `ValidationRunner` service to invoke binary and parse output.
4. Build minimal results page (summary + issue list).
5. Add sample project import using existing `testdata/` and `ruler/` fixtures.
6. Add a tech-decision record: “Go core retained for V1; no full Dart port yet”.

## 10) Kickoff plan (start now)

Yes — we can start immediately while keeping Go source code as-is for now.

### Day 1–2
- Create `gui/csvons_gui` Flutter desktop project scaffold.
- Add routing + app shell (Workspace / Rules / Validate / Reports).
- Add local settings persistence (last workspace path, last ruler path).

### Day 3–4
- Implement `ValidationRunner`:
  - Resolve OS-specific binary path.
  - Execute `csvons --format json <ruler.json>`.
  - Parse stdout JSON and map exit code to UI status.
- Add a minimal validation results screen (summary + issues list).

### Day 5
- Integrate sample data flow using existing `testdata/` and `ruler/` fixtures.
- Add packaging smoke check for one target OS.
- Write short developer runbook (`gui/README.md`) for local run/build.

### Definition of done for kickoff
- User can open a workspace, select ruler, run validation, and view issue list in GUI.
- Same input files produce same pass/fail behavior as Go CLI.
- One desktop artifact is produced successfully in CI/local pipeline.

## 11) Practical quick start (commands)

Use these commands to start implementation immediately.

### 11.1 Verify current Go validator

```bash
go test ./...
go build -o bin/csvons ./cmd/csvons
./bin/csvons ruler/ruler_employees.json
```

### 11.2 Create Flutter desktop app

```bash
mkdir -p gui
cd gui
flutter create csvons_gui
cd csvons_gui
flutter config --enable-linux-desktop --enable-macos-desktop --enable-windows-desktop
flutter run -d linux
```

### 11.3 Add runner integration (first pass)

1. Copy/package `csvons` binary into an app-accessible folder per OS.
2. Add a Dart `ValidationRunner` service using `Process.start`.
3. Run command from Flutter:

```bash
csvons --format json /absolute/path/to/ruler.json
```

4. Parse stdout JSON into `ValidationReport`.
5. Show pass/fail/error based on exit code.

### 11.4 Build first vertical slice

- Screen 1: choose workspace folder + ruler path.
- Screen 2: run validation button.
- Screen 3: summary + issue table.

When these 3 screens work end-to-end, kickoff is complete.
