#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "Removing generated Flutter/Dart artifacts..."
rm -rf .dart_tool build .flutter-plugins .flutter-plugins-dependencies .packages

echo "Removing generated platform folders (if present)..."
rm -rf android ios linux macos windows web

echo "Cleanup complete."
