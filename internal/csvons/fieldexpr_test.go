package csvons

import (
	"testing"
)

// TestPlainField_FieldValue verifies that PlainField correctly extracts
// values from a single column across all data rows.
func TestPlainField_FieldValue(t *testing.T) {
	metadata := &Metadata{
		DataIndex: 1,
	}

	field := &PlainField{
		metadata:  metadata,
		fieldName: "field1",
	}

	fields := []string{"field1", "field2"}
	records := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
		{"value3", "value4"},
	}

	ch := field.FieldValue(fields, records)
	var results []string
	for val := range ch {
		results = append(results, val)
	}

	expected := []string{"value1", "value3"}
	if len(results) != len(expected) {
		t.Errorf("PlainField.FieldValue() returned %d values, expected %d", len(results), len(expected))
	}
	for i, val := range results {
		if val != expected[i] {
			t.Errorf("PlainField.FieldValue() result[%d] = %v, expected %v", i, val, expected[i])
		}
	}
}

// TestPlainField_typeString verifies the type identifier string for PlainField.
func TestPlainField_typeString(t *testing.T) {
	field := &PlainField{}
	result := field.typeString()
	expected := "plain"
	if result != expected {
		t.Errorf("PlainField.typeString() = %v, expected %v", result, expected)
	}
}

// TestPlainField_Init verifies that Init correctly sets metadata and field name.
func TestPlainField_Init(t *testing.T) {
	metadata := &Metadata{DataIndex: 1}
	field := &PlainField{}
	field.Init(metadata, "testField")

	if field.metadata != metadata {
		t.Errorf("PlainField.Init() metadata not set correctly")
	}
	if field.fieldName != "testField" {
		t.Errorf("PlainField.Init() fieldName = %v, expected testField", field.fieldName)
	}
}

// TestRepeatField_FieldValue verifies that RepeatField correctly splits
// cell values by the level-1 separator and yields each element.
func TestRepeatField_FieldValue(t *testing.T) {
	metadata := &Metadata{
		DataIndex:     1,
		Lev1Separator: ",",
	}

	field := &RepeatField{
		metadata:  metadata,
		fieldName: "field1",
	}

	fields := []string{"field1", "field2"}
	records := [][]string{
		{"header1", "header2"},
		{"a,b,c", "value2"},
		{"x,y", "value4"},
	}

	ch := field.FieldValue(fields, records)
	var results []string
	for val := range ch {
		results = append(results, val)
	}

	expected := []string{"a", "b", "c", "x", "y"}
	if len(results) != len(expected) {
		t.Errorf("RepeatField.FieldValue() returned %d values, expected %d", len(results), len(expected))
	}
	for i, val := range results {
		if val != expected[i] {
			t.Errorf("RepeatField.FieldValue() result[%d] = %v, expected %v", i, val, expected[i])
		}
	}
}

// TestRepeatField_typeString verifies the type identifier string for RepeatField.
func TestRepeatField_typeString(t *testing.T) {
	field := &RepeatField{}
	result := field.typeString()
	expected := "repeat"
	if result != expected {
		t.Errorf("RepeatField.typeString() = %v, expected %v", result, expected)
	}
}

// TestRepeatField_Init verifies that Init correctly strips the "[]" suffix.
func TestRepeatField_Init(t *testing.T) {
	metadata := &Metadata{DataIndex: 1}
	field := &RepeatField{}
	field.Init(metadata, "testField[]")

	if field.metadata != metadata {
		t.Errorf("RepeatField.Init() metadata not set correctly")
	}
	if field.fieldName != "testField" {
		t.Errorf("RepeatField.Init() fieldName = %v, expected testField", field.fieldName)
	}
}

// TestNestedField_FieldValue verifies that NestedField correctly splits
// by both separators and extracts the value at the specified index.
func TestNestedField_FieldValue(t *testing.T) {
	metadata := &Metadata{
		DataIndex:     1,
		Lev1Separator: ",",
		Lev2Separator: ":",
	}

	field := &NestedField{
		metadata:  metadata,
		fieldName: "field1",
		index:     1,
	}

	fields := []string{"field1", "field2"}
	records := [][]string{
		{"header1", "header2"},
		{"a:b:c,d:e:f", "value2"},
		{"x:y:z", "value4"},
	}

	ch := field.FieldValue(fields, records)
	var results []string
	for val := range ch {
		results = append(results, val)
	}

	expected := []string{"b", "e", "y"}
	if len(results) != len(expected) {
		t.Errorf("NestedField.FieldValue() returned %d values, expected %d", len(results), len(expected))
	}
	for i, val := range results {
		if val != expected[i] {
			t.Errorf("NestedField.FieldValue() result[%d] = %v, expected %v", i, val, expected[i])
		}
	}
}

// TestNestedField_typeString verifies the type identifier string for NestedField.
func TestNestedField_typeString(t *testing.T) {
	field := &NestedField{}
	result := field.typeString()
	expected := "nested"
	if result != expected {
		t.Errorf("NestedField.typeString() = %v, expected %v", result, expected)
	}
}

// TestNestedField_Init verifies that Init correctly parses the field name and index.
func TestNestedField_Init(t *testing.T) {
	metadata := &Metadata{DataIndex: 1}
	field := &NestedField{}
	field.Init(metadata, "testField{2}")

	if field.metadata != metadata {
		t.Errorf("NestedField.Init() metadata not set correctly")
	}
	if field.fieldName != "testField" {
		t.Errorf("NestedField.Init() fieldName = %v, expected testField", field.fieldName)
	}
	if field.index != 2 {
		t.Errorf("NestedField.Init() index = %v, expected 2", field.index)
	}
}

// TestComplexField_FieldValue verifies that ComplexField correctly concatenates
// values from multiple columns using the field connector.
func TestComplexField_FieldValue(t *testing.T) {
	metadata := &Metadata{
		DataIndex:      1,
		FieldConnector: "-",
	}

	field := &ComplexField{
		metadata:   metadata,
		fieldNames: []string{"field1", "field2"},
	}

	fields := []string{"field1", "field2", "field3"}
	records := [][]string{
		{"header1", "header2", "header3"},
		{"value1", "value2", "value3"},
		{"value4", "value5", "value6"},
	}

	ch := field.FieldValue(fields, records)
	var results []string
	for val := range ch {
		results = append(results, val)
	}

	expected := []string{"value1-value2-", "value4-value5-"}
	if len(results) != len(expected) {
		t.Errorf("ComplexField.FieldValue() returned %d values, expected %d", len(results), len(expected))
	}
	for i, val := range results {
		if val != expected[i] {
			t.Errorf("ComplexField.FieldValue() result[%d] = %v, expected %v", i, val, expected[i])
		}
	}
}

// TestComplexField_typeString verifies the type identifier string for ComplexField.
func TestComplexField_typeString(t *testing.T) {
	field := &ComplexField{}
	result := field.typeString()
	expected := "complex"
	if result != expected {
		t.Errorf("ComplexField.typeString() = %v, expected %v", result, expected)
	}
}

// TestComplexField_Init verifies that Init correctly parses field names from curly braces.
func TestComplexField_Init(t *testing.T) {
	metadata := &Metadata{DataIndex: 1}
	field := &ComplexField{}
	field.Init(metadata, "{field1}{field2}")

	if field.metadata != metadata {
		t.Errorf("ComplexField.Init() metadata not set correctly")
	}
	expected := []string{"field1", "field2"}
	if len(field.fieldNames) != len(expected) {
		t.Errorf("ComplexField.Init() fieldNames length = %d, expected %d", len(field.fieldNames), len(expected))
	}
	for i, name := range field.fieldNames {
		if name != expected[i] {
			t.Errorf("ComplexField.Init() fieldNames[%d] = %v, expected %v", i, name, expected[i])
		}
	}
}

// TestGenerateFieldExpr uses table-driven tests to verify that GenerateFieldExpr
// correctly identifies and creates the appropriate FieldExpr type for each expression pattern.
func TestGenerateFieldExpr(t *testing.T) {
	metadata := &Metadata{
		DataIndex:      1,
		Lev1Separator:  ",",
		Lev2Separator:  ":",
		FieldConnector: "-",
	}

	tests := []struct {
		name         string
		fieldExpr    string
		expectedType string
		shouldBeNil  bool
	}{
		{
			name:         "Plain field expression",
			fieldExpr:    "field1",
			expectedType: "plain",
			shouldBeNil:  false,
		},
		{
			name:         "Plain field with numbers",
			fieldExpr:    "field123",
			expectedType: "plain",
			shouldBeNil:  false,
		},
		{
			name:         "Repeat field expression",
			fieldExpr:    "field1[]",
			expectedType: "repeat",
			shouldBeNil:  false,
		},
		{
			name:         "Repeat field with numbers",
			fieldExpr:    "field123[]",
			expectedType: "repeat",
			shouldBeNil:  false,
		},
		{
			name:         "Nested field expression",
			fieldExpr:    "field1{0}",
			expectedType: "nested",
			shouldBeNil:  false,
		},
		{
			name:         "Nested field with numbers",
			fieldExpr:    "field123{2}",
			expectedType: "nested",
			shouldBeNil:  false,
		},
		{
			name:         "Complex field expression single",
			fieldExpr:    "{field1}",
			expectedType: "complex",
			shouldBeNil:  false,
		},
		{
			name:         "Complex field expression multiple",
			fieldExpr:    "{field1}{field2}",
			expectedType: "",
			shouldBeNil:  true,
		},
		{
			name:         "Invalid field expression",
			fieldExpr:    "field-1",
			expectedType: "",
			shouldBeNil:  true,
		},
		{
			name:         "Empty field expression",
			fieldExpr:    "",
			expectedType: "",
			shouldBeNil:  true,
		},
		{
			name:         "Field with special characters",
			fieldExpr:    "field@1",
			expectedType: "",
			shouldBeNil:  true,
		},
		{
			name:         "Mixed case field",
			fieldExpr:    "Field1",
			expectedType: "plain",
			shouldBeNil:  false,
		},
		{
			name:         "Mixed case repeat field",
			fieldExpr:    "Field1[]",
			expectedType: "repeat",
			shouldBeNil:  false,
		},
		{
			name:         "Mixed case nested field",
			fieldExpr:    "Field1{1}",
			expectedType: "nested",
			shouldBeNil:  false,
		},
		{
			name:         "Mixed case complex field",
			fieldExpr:    "{Field1}",
			expectedType: "complex",
			shouldBeNil:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateFieldExpr(metadata, tt.fieldExpr)

			if tt.shouldBeNil {
				if result != nil {
					t.Errorf("GenerateFieldExpr(%v) = %v, expected nil", tt.fieldExpr, result)
				}
				return
			}

			if result == nil {
				t.Errorf("GenerateFieldExpr(%v) = nil, expected non-nil", tt.fieldExpr)
				return
			}

			actualType := result.typeString()
			if actualType != tt.expectedType {
				t.Errorf("GenerateFieldExpr(%v).typeString() = %v, expected %v", tt.fieldExpr, actualType, tt.expectedType)
			}
		})
	}
}

// TestGenerateFieldExpr_NilMetadata tests the nil metadata case.
// Skipped because log.Fatal will exit the program in a unit test environment.
func TestGenerateFieldExpr_NilMetadata(t *testing.T) {
	t.Skip("Skipping nil metadata test - log.Fatal exits program")
}

// TestFieldExprInterface verifies that all field types correctly implement
// the FieldExpr interface at compile time.
func TestFieldExprInterface(t *testing.T) {
	var _ FieldExpr = (*PlainField)(nil)
	var _ FieldExpr = (*RepeatField)(nil)
	var _ FieldExpr = (*NestedField)(nil)
	var _ FieldExpr = (*ComplexField)(nil)
}

// TestFieldExprMap verifies that fieldExprMap contains the expected regex patterns
// and that each constructor creates an instance of the correct type.
func TestFieldExprMap(t *testing.T) {
	expectedPatterns := []string{
		`^[a-zA-Z0-9]+$`,
		`^[a-zA-Z0-9]+\[\]$`,
		`^[a-zA-Z0-9]+\{\d+\}$`,
		`^\{[a-zA-Z0-9]+\}+$`,
	}

	expectedTypes := []string{
		"plain",
		"repeat",
		"nested",
		"complex",
	}

	for i, pattern := range expectedPatterns {
		constructor, exists := fieldExprMap[pattern]
		if !exists {
			t.Errorf("fieldExprMap missing pattern: %v", pattern)
			continue
		}

		// Verify the constructor creates the right type.
		instance := constructor("test")
		actualType := instance.typeString()
		if actualType != expectedTypes[i] {
			t.Errorf("fieldExprMap[%v] creates type %v, expected %v", pattern, actualType, expectedTypes[i])
		}
	}
}

// TestFieldTypesWithEmptyRecords verifies that all field types handle
// empty records gracefully by returning zero values through the channel.
func TestFieldTypesWithEmptyRecords(t *testing.T) {
	metadata := &Metadata{DataIndex: 1}
	emptyRecords := [][]string{}

	plainField := &PlainField{metadata: metadata, fieldName: "field1"}
	repeatField := &RepeatField{metadata: metadata, fieldName: "field1"}
	nestedField := &NestedField{metadata: metadata, fieldName: "field1", index: 0}
	complexField := &ComplexField{metadata: metadata, fieldNames: []string{"field1"}}

	fields := []string{"field1"}

	tests := []struct {
		name   string
		field  FieldExpr
		expect int // Expected number of output values.
	}{
		{"PlainField with empty records", plainField, 0},
		{"RepeatField with empty records", repeatField, 0},
		{"NestedField with empty records", nestedField, 0},
		{"ComplexField with empty records", complexField, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.field.FieldValue(fields, emptyRecords)
			count := 0
			for range ch {
				count++
			}
			if count != tt.expect {
				t.Errorf("%s.FieldValue() returned %d values, expected %d", tt.name, count, tt.expect)
			}
		})
	}
}

// TestFieldTypesWithNilRecords verifies that all field types handle
// nil records gracefully by returning zero values through the channel.
func TestFieldTypesWithNilRecords(t *testing.T) {
	metadata := &Metadata{DataIndex: 1}
	fields := []string{"field1"}

	plainField := &PlainField{metadata: metadata, fieldName: "field1"}
	repeatField := &RepeatField{metadata: metadata, fieldName: "field1"}
	nestedField := &NestedField{metadata: metadata, fieldName: "field1", index: 0}
	complexField := &ComplexField{metadata: metadata, fieldNames: []string{"field1"}}

	tests := []struct {
		name   string
		field  FieldExpr
		expect int // Expected number of output values.
	}{
		{"PlainField with nil records", plainField, 0},
		{"RepeatField with nil records", repeatField, 0},
		{"NestedField with nil records", nestedField, 0},
		{"ComplexField with nil records", complexField, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.field.FieldValue(fields, nil)
			count := 0
			for range ch {
				count++
			}
			if count != tt.expect {
				t.Errorf("%s.FieldValue() returned %d values, expected %d", tt.name, count, tt.expect)
			}
		})
	}
}

// TestFieldNotFound verifies that all field types return nil when
// the specified field is not present in the column headers.
func TestFieldNotFound(t *testing.T) {
	metadata := &Metadata{DataIndex: 1}
	fields := []string{"otherField"}
	records := [][]string{
		{"header1"},
		{"value1"},
	}

	plainField := &PlainField{metadata: metadata, fieldName: "field1"}
	repeatField := &RepeatField{metadata: metadata, fieldName: "field1"}
	nestedField := &NestedField{metadata: metadata, fieldName: "field1", index: 0}
	complexField := &ComplexField{metadata: metadata, fieldNames: []string{"field1"}}

	tests := []struct {
		name  string
		field FieldExpr
	}{
		{"PlainField field not found", plainField},
		{"RepeatField field not found", repeatField},
		{"NestedField field not found", nestedField},
		{"ComplexField field not found", complexField},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := tt.field.FieldValue(fields, records)
			if ch != nil {
				t.Errorf("%s.FieldValue() should return nil when field not found", tt.name)
			}
		})
	}
}
