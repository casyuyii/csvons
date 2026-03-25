# csvons_gui (starter)

Minimal Flutter starter wiring for running the Go `csvons` validator binary.

## What is included

- `lib/core/validation_runner.dart`: process runner for `csvons`.
- `lib/models/validation_report.dart`: basic JSON models for report parsing.
- `lib/screens/home_page.dart`: starter screen to run validation and view output.
- `lib/main.dart`: app entry.

## How to use

1. Create/prepare a Flutter project in this folder (`flutter create .` if needed).
2. Build `csvons` binary and place it at a path you can reference from the UI.
3. Start the app and enter:
   - binary path (e.g., `bin/csvons_linux`)
   - absolute `ruler.json` path
4. Click **Run Validation**.

> Note: current Go CLI may not emit JSON yet. If JSON parse fails, raw stdout/stderr is shown.
