# Flutter GUI packaging note (desktop-first)

This note defines where bundled `csvons` binaries should live for desktop packaging.

## Proposed binary layout

Under `gui/csvons_gui/bin/`:

- `gui/csvons_gui/bin/linux/csvons`
- `gui/csvons_gui/bin/macos/csvons`
- `gui/csvons_gui/bin/windows/csvons.exe`

## Runtime selection

`ValidationRunner` should resolve binary path in this order:

1. Explicit user-entered binary path (already supported).
2. Platform default in `gui/csvons_gui/bin/<platform>/`.
3. Fallback to existing dev default.

## Release checklist hook

For each release candidate:

1. Build per-OS `csvons` binaries from the same commit as the GUI artifact.
2. Place binaries into the paths above.
3. Smoke test one validation run per platform artifact.
