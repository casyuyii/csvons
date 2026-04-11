# Feature Specification: TypeScript Core + Electron GUI

## Summary

Port the csvons Go core validation library to TypeScript and build an Electron desktop application that uses the TypeScript core directly for CSV constraint validation.

## Goals

1. Create a TypeScript library (`ts-core/`) that replicates all validation logic from `internal/csvons/`
2. Build an Electron GUI (`gui/electron/`) with feature parity to the existing Flutter GUI
3. Comprehensive unit tests for the TS core with >=80% coverage

## Requirements

### Functional
- Support all 4 field expression types: PlainField, RepeatField, NestedField, ComplexField
- Support all 3 constraint validators: exists, unique, vtype
- Parse ruler.json config files with identical semantics to Go implementation
- Parse CSV files following RFC 4180 (via papaparse)
- Produce ValidationReport with same JSON schema as Go CLI output
- Electron GUI with Validate and Workspace tabs
- Issue filtering/sorting, CSV preview, report export

### Non-Functional
- Error messages follow format: `[file]:[row] field "<field>": <reason>`
- Boolean parsing accepts same values as Go's strconv.ParseBool
- Unit test coverage >=80%
- Positive AND negative tests for each validator

## Success Criteria

- All Go test scenarios pass in TypeScript equivalent tests
- Electron app can validate any ruler.json that the Go CLI can
- Validation results match Go CLI JSON output
