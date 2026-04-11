# Quickstart: TypeScript Core + Electron GUI

## Prerequisites

- bun >= 1.0
- Node.js >= 18

## Setup

```bash
cd gui/electron
bun install
```

## Run Tests

```bash
cd gui/electron

# All tests
bun test

# Tests with coverage
bun test --coverage
```

## Development

```bash
cd gui/electron
bun run dev      # Starts Electron with HMR for renderer
```

## Usage (Core Library)

```typescript
import { validate } from './src/core';

const report = validate('path/to/ruler.json');

if (report.summary.failed > 0) {
  for (const issue of report.issues) {
    console.log(`${issue.file}:${issue.row} field "${issue.field}": ${issue.message}`);
  }
}
```

## Build

```bash
cd gui/electron
bun run build
```
