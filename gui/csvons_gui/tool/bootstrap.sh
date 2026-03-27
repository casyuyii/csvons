#!/usr/bin/env bash
set -euo pipefail

if ! command -v flutter >/dev/null 2>&1; then
  echo "error: flutter is not installed or not on PATH" >&2
  exit 1
fi

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "==> Generating Flutter platform scaffolding (if missing)"
flutter create .

echo "==> Installing dependencies"
flutter pub get

echo "==> Running static checks"
flutter analyze
dart format --set-exit-if-changed lib test

echo "==> Running tests"
flutter test

echo "Bootstrap + checks completed."
