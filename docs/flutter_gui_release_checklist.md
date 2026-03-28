# Flutter GUI release checklist

Use this checklist for desktop-first release candidates.

## Before build

1. Start from the merge target commit you intend to ship.
2. Run backend verification: `go test ./...`.
3. In a Flutter-enabled environment, run `cd gui/csvons_gui && flutter create . && make check`.

## Build artifacts

1. Build the Go validator for each target OS and place it under `gui/csvons_gui/bin/<platform>/`.
2. Build the desktop Flutter artifact for each target OS.
3. Stage the validator sidecar into the packaged desktop bundle (`bin/<platform>/csvons` next to the app executable).

## Smoke tests

1. Launch the packaged app on a fresh machine or clean VM.
2. Open the sample workspace and run validation against at least one passing ruler and one failing ruler.
3. Confirm the packaged app resolves the bundled validator without manually editing the binary path.
4. Export JSON and Markdown reports and confirm both files are written successfully.

## Ship notes

1. Record the exact git commit used for both the GUI and bundled validator binaries.
2. Note any signing/notarization steps performed for the target platform.
3. Capture any platform-specific caveats discovered during the smoke test.
