// Package csvons provides CSV constraint validation based on JSON configuration rules.
//
// It supports validating CSV files against three types of constraints:
//   - exists: values in a column must exist in another CSV file's column
//   - unique: values in a column must be unique across all rows
//   - vtype: values must conform to a specified type (int, float64, bool) and optional range
//
// Field expressions allow validation of values within nested data structures,
// not just simple column values. See FieldExpr for supported expression types.
package csvons

// ConstrainsConfig represents the complete configuration for CSV constraint validation.
// It combines all constraint types (exists, unique, vtype) with the CSV metadata
// that describes how to read and interpret the CSV files.
type ConstrainsConfig struct {
	Exists   []Exists `json:"exists"`          // Rules for cross-file value existence validation.
	Unique   Unique   `json:"unique"`          // Rules for column uniqueness validation.
	VType    []VType  `json:"vtype"`           // Rules for value type and range validation.
	Metadata Metadata `json:"csvons_metadata"` // Metadata describing CSV file structure.
}

// Metadata describes the structure and location of CSV files being validated.
// It specifies where to find the files, which row contains column headers,
// where actual data begins, and the separators used for nested data structures.
//
// Example:
//
//	Metadata{
//	    CSVFileFolder: "testdata",
//	    NameIndex:     0,
//	    DataIndex:     1,
//	    Extension:     ".csv",
//	    Lev1Separator: ";",
//	    Lev2Separator: ":",
//	}
type Metadata struct {
	CSVFileFolder  string `json:"csv_file_folder"` // Directory containing the CSV files.
	NameIndex      int    `json:"name_index"`      // Row index where column names are defined (0-based).
	DataIndex      int    `json:"data_index"`      // Row index where actual data starts (0-based, must be > NameIndex).
	Extension      string `json:"extension"`       // File extension for CSV files (typically ".csv").
	Lev1Separator  string `json:"lev1_separator"`  // Separator for first-level array values (e.g., ";").
	Lev2Separator  string `json:"lev2_separator"`  // Separator for second-level nested values (e.g., ":").
	FieldConnector string `json:"field_connector"` // Connector string for combining complex field values (e.g., "|").
}

// Exists defines a cross-file existence constraint.
// It specifies that values in source columns must also exist in corresponding
// columns of a destination CSV file.
//
// Example JSON:
//
//	{
//	    "dst_file_stem": "username-d1",
//	    "fields": [{"src": "Username", "dst": "Username"}]
//	}
type Exists struct {
	DstFileStem string `json:"dst_file_stem"` // Base name (stem) of the target CSV file.
	Fields      []struct {
		Src string `json:"src"` // Field expression in the source file.
		Dst string `json:"dst"` // Field expression in the destination file.
	} `json:"fields"` // Pairs of source-destination field expressions to compare.
}

// Unique defines a column uniqueness constraint.
// It specifies that all values within each listed field must be unique across all rows.
//
// Example JSON:
//
//	{"fields": ["Username", "marks{0}"]}
type Unique struct {
	Fields []string `json:"fields"` // Field expressions whose values must be unique.
}

// VType defines a value type and optional range constraint.
// It validates that values in the specified field can be parsed as the given type,
// and optionally fall within a numeric range.
//
// Supported types: "int", "float64", "bool".
// Range is only applicable to "int" and "float64" types.
//
// Example JSON:
//
//	{"field": "Age", "type": "int", "range": {"min": 1, "max": 100}}
type VType struct {
	Field string `json:"field"` // Field expression to validate.
	Type  string `json:"type"`  // Expected value type: "int", "float64", or "bool".
	Range *struct {
		Min float64 `json:"min"` // Minimum allowed value (inclusive).
		Max float64 `json:"max"` // Maximum allowed value (inclusive).
	} `json:"range,omitempty"` // Optional numeric range constraint.
}
