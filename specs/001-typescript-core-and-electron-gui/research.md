# Research: TypeScript Core + Electron GUI

## 1. CSV Parsing Library

- **Decision**: papaparse
- **Rationale**: Most mature RFC 4180 parser in JS ecosystem. Zero dependencies. Handles quoted fields, multiline values, and streaming. Used by 4M+ weekly downloads. Correctness-first aligns with constitution.
- **Alternatives considered**: csv-parse (heavier, stream-oriented API adds complexity for sync use), custom parser (error-prone, violates "correctness first")

## 2. Go Channel Pattern Translation

- **Decision**: TypeScript generator functions (`function*`)
- **Rationale**: Generators provide lazy evaluation semantics identical to Go channels. `for...of` stops early (matching Go's channel close behavior). The exists validator's early-break pattern works naturally. No need for async generators since file I/O is synchronous in the main process.
- **Alternatives considered**: Async iterators (unnecessary complexity since validation runs in Node.js main process), eager arrays (would break the lazy caching optimization in exists validator)

## 3. Error Flow Pattern

- **Decision**: Collect all errors into `ValidationError[]` instead of panic on first failure
- **Rationale**: GUI needs all issues at once. Go's panic/recover is an anti-pattern in JS. Error collection is strictly more useful for the Electron GUI use case. A `failFast` option can restore Go's single-error behavior.
- **Alternatives considered**: Throw on first error (matches Go behavior but poor GUI experience), callback pattern (unnecessarily complex)

## 4. Renderer Framework

- **Decision**: React 19
- **Rationale**: Best ecosystem for data tables (core UI component). Widest contributor pool. Mature Electron integration patterns. Flutter's stateful widget architecture maps naturally to React components with hooks.
- **Alternatives considered**: Svelte (smaller ecosystem for data tables), Vue (viable but less Electron ecosystem support), vanilla HTML (would need to build table sort/filter from scratch)

## 5. Build Tooling

- **Decision**: Vite for renderer, tsc for core library
- **Rationale**: Vite provides fast HMR, native TS/JSX support, integrates with vitest. tsc for the core keeps it simple with no bundler dependency. electron-vite wraps both.
- **Alternatives considered**: webpack (unnecessarily complex), esbuild alone (less mature HMR story)

## 6. Package Manager

- **Decision**: bun
- **Rationale**: User preference. Existing package.json already has `@types/bun`. Fast installation and execution. Compatible with Node.js APIs used by Electron.
- **Alternatives considered**: npm (viable but slower), yarn (no advantage over bun)

## 7. Test Framework

- **Decision**: vitest
- **Rationale**: Native TypeScript support, ESM-first, fast execution, built-in coverage via istanbul, integrates with Vite. Avoids Jest's CommonJS transform overhead.
- **Alternatives considered**: jest (legacy CJS baggage), node:test (no coverage integration, no watch mode)

## 8. Integer Precision

- **Decision**: Use `Number.parseInt` with documented 2^53 precision limit
- **Rationale**: CSV data rarely contains integers exceeding 2^53. The validation use case only checks "is this parseable as an integer?" — precision beyond Number.MAX_SAFE_INTEGER is not practically needed.
- **Alternatives considered**: BigInt (adds complexity for parsing and range comparison, minimal practical benefit)

## 9. TS Core Placement

- **Decision**: Inside `gui/electron/src/core/` as part of the standalone Electron project
- **Rationale**: User wants `gui/electron/` to be fully standalone with no external dependencies on other project directories. Simpler project structure with no workspace configuration needed.
- **Alternatives considered**: `ts-core/` at repo root (adds workspace complexity, user prefers standalone)
