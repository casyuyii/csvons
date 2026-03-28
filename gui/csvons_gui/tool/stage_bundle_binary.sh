#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PLATFORM="${1:-linux}"

case "$PLATFORM" in
  linux)
    SOURCE_PATH="$ROOT_DIR/bin/linux/csvons"
    TARGET_DIR="$ROOT_DIR/build/linux/x64/release/bundle/bin/linux"
    TARGET_PATH="$TARGET_DIR/csvons"
    ;;
  *)
    echo "error: unsupported platform '$PLATFORM'" >&2
    exit 1
    ;;
esac

if [[ ! -f "$SOURCE_PATH" ]]; then
  echo "error: bundled validator not found at $SOURCE_PATH" >&2
  exit 1
fi

mkdir -p "$TARGET_DIR"
cp "$SOURCE_PATH" "$TARGET_PATH"
chmod +x "$TARGET_PATH"

echo "Staged bundled validator to $TARGET_PATH"
