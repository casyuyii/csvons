# Data Model: TypeScript Core

## Entities

### Metadata
Configuration describing CSV file structure and location.

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| csvFileFolder | string | csv_file_folder | Directory containing CSV files |
| nameIndex | number | name_index | Row index for column headers (0-based) |
| dataIndex | number | data_index | Row index where data starts (0-based, > nameIndex) |
| extension | string | extension | File extension (e.g., ".csv") |
| lev1Separator | string | lev1_separator | First-level array separator |
| lev2Separator | string | lev2_separator | Second-level nested separator |
| fieldConnector | string | field_connector | Multi-field concatenation connector |

**Validation**: dataIndex > nameIndex, extension non-empty, separators non-empty

### ExistsRule
Cross-file value existence constraint.

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| dstFileStem | string | dst_file_stem | Target CSV file base name |
| fields | {src: string, dst: string}[] | fields | Source-destination field pairs |

### UniqueRule
Column uniqueness constraint.

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| fields | string[] | fields | Field expressions that must be unique |

### VTypeRule
Value type and range constraint.

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| field | string | field | Field expression to validate |
| type | "int" \| "float64" \| "bool" | type | Expected value type |
| range | {min: number, max: number} \| undefined | range | Optional numeric range (inclusive) |

### ValidationError
Structured validation failure.

| Field | Type | Description |
|-------|------|-------------|
| file | string | CSV filename |
| rule | string | Rule type: "exists" \| "unique" \| "vtype" |
| field | string | Field name/expression |
| row | number \| undefined | 1-based row number (undefined if not row-specific) |
| value | string | Value that failed validation |
| message | string | Human-readable error message |
| severity | string | "error" |
| code | number | Exit code: 1 (validation) or 2 (runtime) |

### ValidationReport
Complete validation result.

| Field | Type | Description |
|-------|------|-------------|
| schemaVersion | string | "csvons.validation_report.v1" |
| summary | ValidationSummary | Aggregate results |
| issues | ValidationIssue[] | All validation failures |

### ValidationSummary

| Field | Type | Description |
|-------|------|-------------|
| filesChecked | number | Total files validated |
| passed | number | Files with no errors |
| failed | number | Files with errors |
| durationMs | number | Validation duration in milliseconds |

### FieldExpr (interface)
Abstract field expression for extracting values from CSV cells.

| Method | Signature | Description |
|--------|-----------|-------------|
| fieldValues | (fields: string[], records: string[][]) => Iterable\<string\> | Yields extracted values |
| fieldOccurrences | (fields: string[], records: string[][]) => Iterable\<FieldOccurrence\> | Yields values with row numbers |
| typeString | () => string | Returns expression type name |

**Implementations**: PlainField, RepeatField, NestedField, ComplexField

### FieldOccurrence
Value with row tracking.

| Field | Type | Description |
|-------|------|-------------|
| row | number | 1-based CSV row number |
| value | string | Extracted field value |

## Relationships

```
ruler.json
  └── csvons_metadata → Metadata
  └── [file_stem] → ConstraintsConfig
        ├── exists[] → ExistsRule → references another [file_stem]
        ├── unique → UniqueRule
        └── vtype[] → VTypeRule

validate() → ValidationReport
  ├── summary → ValidationSummary
  └── issues[] → ValidationIssue (derived from ValidationError)

FieldExpr ← GenerateFieldExpr(metadata, expr)
  ├── PlainField   (regex: /^[a-zA-Z0-9 ]+$/)
  ├── RepeatField  (regex: /^[a-zA-Z0-9 ]+\[\]$/)
  ├── NestedField  (regex: /^[a-zA-Z0-9 ]+\{\d+\}$/)
  └── ComplexField (regex: /^\{[a-zA-Z0-9 ]+\}+$/)
```
