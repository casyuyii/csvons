# Contract: @csvons/core Public API

## validate(configPath: string, options?: ValidateOptions): ValidationReport

Main entry point. Reads a ruler.json config, validates all referenced CSV files.

**Parameters:**
- `configPath: string` — Absolute or relative path to ruler.json
- `options.failFast?: boolean` — Stop on first error (default: false, collect all)
- `options.basePath?: string` — Base directory for resolving relative paths in metadata (default: dirname of configPath)

**Returns:** `ValidationReport` matching JSON schema `csvons.validation_report.v1`

**Errors:** Throws on file I/O failure or invalid config. Validation failures are returned in `report.issues`, not thrown.

## readCsvFile(stem: string, metadata: Metadata): string[][]

Reads and parses a CSV file.

**Parameters:**
- `stem: string` — File stem (base name without extension)
- `metadata: Metadata` — CSV structure descriptor

**Returns:** 2D array of strings. `records[metadata.nameIndex]` = headers, `records[metadata.dataIndex:]` = data rows.

## readConfigFile(configPath: string): { rules: Record<string, RawRules>, metadata: Metadata }

Parses a ruler.json config file.

**Parameters:**
- `configPath: string` — Path to ruler.json

**Returns:** Parsed rules keyed by file stem, plus extracted metadata.

## generateFieldExpr(metadata: Metadata, expr: string): FieldExpr | null

Factory for field expressions.

**Parameters:**
- `metadata: Metadata` — CSV metadata with separators
- `expr: string` — Field expression string

**Returns:** Appropriate FieldExpr implementation, or null if pattern doesn't match.
