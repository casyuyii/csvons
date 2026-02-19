package csvons

import (
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// FieldExpr is the interface for all field expression types.
// A field expression defines how to extract values from CSV records
// based on a column name and optional nested data access pattern.
//
// There are four implementations:
//   - PlainField: direct column reference (e.g., "Username")
//   - RepeatField: array expansion (e.g., "Tags[]")
//   - NestedField: second-level array index (e.g., "marks{1}")
//   - ComplexField: multi-field concatenation (e.g., "{data}{key}")
type FieldExpr interface {
	// FieldValue returns a channel that yields extracted values from the given records.
	// The fields parameter contains column names (header row), and records is the full CSV data.
	// Returns nil if the target field is not found in the column names.
	FieldValue(fields []string, records [][]string) <-chan string

	// typeString returns a human-readable identifier for this field expression type.
	typeString() string

	// Init initializes the field expression with metadata and the raw expression string.
	Init(metadata *Metadata, expr string)
}

// -------------------------------------------------------
// PlainField: direct column name reference.
// Example expressions: "Username", "Age", "field1"
// -------------------------------------------------------

// PlainField extracts values directly from a single column by its name.
// Each row yields exactly one value from the specified column.
type PlainField struct {
	metadata  *Metadata // CSV metadata for data index and separators.
	fieldName string    // Column name to extract values from.
}

// FieldValue yields one value per data row from the column matching fieldName.
// Returns nil if the field name does not exist in the column headers.
func (p *PlainField) FieldValue(fields []string, records [][]string) <-chan string {
	// Find the index of the target column in the header row.
	fieldIndex := slices.Index(fields, p.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan string, 128)
	go func() {
		defer close(output)
		// Iterate from the data start index through all records.
		for i := p.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				output <- record[fieldIndex]
			} else {
				log.Printf("plain field [%s] not found in record [%d]", p.fieldName, i)
			}
		}
	}()

	return output
}

// typeString returns "plain" to identify this as a plain field expression.
func (p *PlainField) typeString() string {
	return "plain"
}

// Init sets the metadata and field name from the raw expression string.
// For plain fields, the expression is used directly as the column name.
func (p *PlainField) Init(metadata *Metadata, expr string) {
	p.metadata = metadata
	p.fieldName = expr
}

// -------------------------------------------------------
// RepeatField: array expansion field.
// Example expressions: "Tags[]", "Items[]"
// Splits cell values by lev1_separator and yields each element.
// -------------------------------------------------------

// RepeatField extracts values from a column and expands array-like values.
// Cell values are split by the Lev1Separator, and each element is yielded separately.
// For example, if a cell contains "tag1;tag2;tag3" and separator is ";",
// three separate values are yielded: "tag1", "tag2", "tag3".
type RepeatField struct {
	metadata  *Metadata // CSV metadata for data index and separators.
	fieldName string    // Column name (without the "[]" suffix).
}

// FieldValue splits each cell value by Lev1Separator and yields individual elements.
// Returns nil if the field name does not exist in the column headers.
func (r *RepeatField) FieldValue(fields []string, records [][]string) <-chan string {
	fieldIndex := slices.Index(fields, r.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan string, 128)
	go func() {
		defer close(output)
		for i := r.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				// Split the cell value by the first-level separator.
				lev1Vals := strings.Split(record[fieldIndex], r.metadata.Lev1Separator)
				for _, lev1Val := range lev1Vals {
					output <- lev1Val
				}
			} else {
				log.Printf("repeat field [%s] not found in record [%d]", r.fieldName, i)
			}
		}
	}()

	return output
}

// typeString returns "repeat" to identify this as a repeat field expression.
func (r *RepeatField) typeString() string {
	return "repeat"
}

// Init sets the metadata and extracts the field name by stripping the "[]" suffix.
func (r *RepeatField) Init(metadata *Metadata, expr string) {
	r.metadata = metadata
	r.fieldName = expr[:len(expr)-2] // Remove trailing "[]".
}

// -------------------------------------------------------
// NestedField: second-level array index access.
// Example expressions: "marks{0}", "scores{1}"
// Splits by lev1_separator, then by lev2_separator, and extracts value at index.
// -------------------------------------------------------

// NestedField extracts values from a two-dimensional nested structure within a cell.
// Cell values are first split by Lev1Separator into groups, then each group is
// split by Lev2Separator, and the value at the specified index is yielded.
//
// For example, with data "9012:30;90:14", Lev1Sep=";", Lev2Sep=":", index=1:
// Split by ";" → ["9012:30", "90:14"]
// Split each by ":" and take index 1 → yields "30", "14"
type NestedField struct {
	metadata  *Metadata // CSV metadata for data index and separators.
	fieldName string    // Column name (without the "{N}" suffix).
	index     int       // Zero-based index into the second-level array.
}

// FieldValue yields values at the nested index from each split cell value.
// Returns nil if the field name does not exist in the column headers.
func (n *NestedField) FieldValue(fields []string, records [][]string) <-chan string {
	fieldIndex := slices.Index(fields, n.fieldName)
	if fieldIndex == -1 {
		return nil
	}

	output := make(chan string, 128)
	go func() {
		defer close(output)
		for i := n.metadata.DataIndex; i < len(records); i++ {
			record := records[i]
			if fieldIndex < len(record) {
				// Split by first-level separator into groups.
				lev1Vals := strings.Split(record[fieldIndex], n.metadata.Lev1Separator)
				for _, lev1Val := range lev1Vals {
					// Split each group by second-level separator and extract by index.
					lev2Vals := strings.Split(lev1Val, n.metadata.Lev2Separator)
					if n.index < len(lev2Vals) {
						output <- lev2Vals[n.index]
					} else {
						log.Printf("nested field [%s] level 2 value [%s] length [%d] not found in record [%d]", n.fieldName, lev1Val, len(lev2Vals), i)
					}
				}
			} else {
				log.Printf("nested field [%s] not found in record [%d]", n.fieldName, i)
			}
		}
	}()

	return output
}

// typeString returns "nested" to identify this as a nested field expression.
func (n *NestedField) typeString() string {
	return "nested"
}

// Init parses the expression to extract the field name and index.
// Expression format: "fieldName{index}" (e.g., "marks{1}").
func (n *NestedField) Init(metadata *Metadata, expr string) {
	n.metadata = metadata
	// Parse "fieldName{index}" using regex to extract name and numeric index.
	matches := regexp.MustCompile(`([a-zA-Z0-9]+)\{(\d+)\}?`).FindStringSubmatch(expr)
	if len(matches) != 3 {
		return
	}

	n.fieldName = matches[1]
	n.index, _ = strconv.Atoi(matches[2])
}

// -------------------------------------------------------
// ComplexField: multi-field concatenation.
// Example expressions: "{data}", "{field1}{field2}"
// Concatenates values from multiple columns using FieldConnector.
// -------------------------------------------------------

// ComplexField concatenates values from multiple columns into a single string.
// Each field's value is joined using the FieldConnector from metadata.
//
// For example, with fields ["field1", "field2"], connector="|":
// Row ["val1", "val2"] → yields "val1|val2|"
type ComplexField struct {
	metadata   *Metadata // CSV metadata for data index and separators.
	fieldNames []string  // List of column names to concatenate.
}

// FieldValue concatenates values from multiple columns for each row.
// Returns nil if any of the required field names are not found in the column headers.
func (c *ComplexField) FieldValue(fields []string, records [][]string) <-chan string {
	// Resolve indices for all required field names.
	fieldIndexes := make([]int, len(c.fieldNames))
	for i, fieldName := range c.fieldNames {
		fieldIndexes[i] = slices.Index(fields, fieldName)
		if fieldIndexes[i] == -1 {
			log.Printf("complex field [%s] not found in fields", fieldName)
			return nil
		}
	}

	output := make(chan string, 128)
	go func() {
		defer close(output)
		for i := c.metadata.DataIndex; i < len(records); i++ {
			cpxStr := ""
			for _, fieldIndex := range fieldIndexes {
				record := records[i]
				if fieldIndex < len(record) {
					// Append each field value with the connector separator.
					cpxStr += record[fieldIndex] + c.metadata.FieldConnector
				} else {
					log.Printf("complex field [%s] not found in record [%d]", c.fieldNames[fieldIndex], i)
					return
				}
			}
			output <- cpxStr
		}
	}()

	return output
}

// typeString returns "complex" to identify this as a complex field expression.
func (c *ComplexField) typeString() string {
	return "complex"
}

// Init parses the expression to extract all field names enclosed in curly braces.
// Expression format: "{field1}{field2}..." — each {name} becomes a column reference.
func (c *ComplexField) Init(metadata *Metadata, expr string) {
	c.metadata = metadata
	// Extract all alphanumeric tokens from the expression (field names inside braces).
	c.fieldNames = regexp.MustCompile(`([a-zA-Z0-9]+)`).FindAllString(expr, -1)
}

// -------------------------------------------------------
// Field Expression Factory
// -------------------------------------------------------

// fieldExprMap maps regex patterns to constructor functions for each field expression type.
// The patterns are matched against the raw field expression string to determine its type:
//   - `^[a-zA-Z0-9]+$`          → PlainField  (e.g., "Username")
//   - `^[a-zA-Z0-9]+\[\]$`      → RepeatField (e.g., "Tags[]")
//   - `^[a-zA-Z0-9]+\{\d+\}$`   → NestedField (e.g., "marks{1}")
//   - `^\{[a-zA-Z0-9]+\}+$`     → ComplexField (e.g., "{data}")
var fieldExprMap = map[string]func(string) FieldExpr{
	`^[a-zA-Z0-9]+$`:        func(string) FieldExpr { return &PlainField{} },
	`^[a-zA-Z0-9]+\[\]$`:    func(string) FieldExpr { return &RepeatField{} },
	`^[a-zA-Z0-9]+\{\d+\}$`: func(string) FieldExpr { return &NestedField{} },
	`^\{[a-zA-Z0-9]+\}+$`:   func(string) FieldExpr { return &ComplexField{} },
}

// GenerateFieldExpr creates and initializes a FieldExpr from a raw expression string.
// It matches the expression against known patterns and returns the appropriate implementation.
//
// Returns nil if no pattern matches the expression. Panics (log.Fatal) if metadata is nil.
//
// Example:
//
//	expr := GenerateFieldExpr(metadata, "Tags[]")
//	// Returns a *RepeatField initialized with fieldName="Tags"
func GenerateFieldExpr(metadata *Metadata, fieldExpr string) FieldExpr {
	if metadata == nil {
		log.Fatal("metadata is nil")
		return nil
	}

	// Try each registered pattern to find a matching field expression type.
	for pattern, makeFunc := range fieldExprMap {
		if match, _ := regexp.MatchString(pattern, fieldExpr); match {
			t := makeFunc(fieldExpr)
			t.Init(metadata, fieldExpr)
			return t
		}
	}

	return nil
}
