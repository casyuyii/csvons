package csvons

import (
	"testing"
)

func TestPlainField_Check(t *testing.T) {
	field := &PlainField{
		SrcFields: []string{"field1", "field2"},
		DstFields: []string{"field1", "field2"},
	}

	srcRecords := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
	}
	dstRecords := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
	}

	result := field.Check(srcRecords, dstRecords)
	if !result {
		t.Errorf("PlainField.Check() returned false, expected true")
	}
}

func TestPlainField_typeString(t *testing.T) {
	field := &PlainField{}
	result := field.typeString()
	expected := "plain"
	if result != expected {
		t.Errorf("PlainField.typeString() = %v, expected %v", result, expected)
	}
}

func TestNestedField_Check(t *testing.T) {
	field := &NestedField{
		SrcFields: []string{"field1[]", "field2[]"},
		DstFields: []string{"field1[]", "field2[]"},
	}

	srcRecords := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
	}
	dstRecords := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
	}

	result := field.Check(srcRecords, dstRecords)
	if !result {
		t.Errorf("NestedField.Check() returned false, expected true")
	}
}

func TestNestedField_typeString(t *testing.T) {
	field := &NestedField{}
	result := field.typeString()
	expected := "nested"
	if result != expected {
		t.Errorf("NestedField.typeString() = %v, expected %v", result, expected)
	}
}

func TestComplexField_Check(t *testing.T) {
	field := &ComplexField{
		SrcFields: []string{"{field1}", "{field2}"},
		DstFields: []string{"{field1}", "{field2}"},
	}

	srcRecords := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
	}
	dstRecords := [][]string{
		{"header1", "header2"},
		{"value1", "value2"},
	}

	result := field.Check(srcRecords, dstRecords)
	if !result {
		t.Errorf("ComplexField.Check() returned false, expected true")
	}
}

func TestComplexField_typeString(t *testing.T) {
	field := &ComplexField{}
	result := field.typeString()
	expected := "complex"
	if result != expected {
		t.Errorf("ComplexField.typeString() = %v, expected %v", result, expected)
	}
}

func TestGenerateFieldexpr(t *testing.T) {
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
			name:         "Nested field expression",
			fieldExpr:    "field1[]",
			expectedType: "nested",
			shouldBeNil:  false,
		},
		{
			name:         "Nested field with numbers",
			fieldExpr:    "field123[]",
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
			expectedType: "complex",
			shouldBeNil:  false,
		},
		{
			name:         "Complex field with numbers",
			fieldExpr:    "{field123}",
			expectedType: "complex",
			shouldBeNil:  false,
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
			name:         "Mixed case nested field",
			fieldExpr:    "Field1[]",
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
			result := GenerateFieldExpr(tt.fieldExpr)

			if tt.shouldBeNil {
				if result != nil {
					t.Errorf("GenerateFieldexpr(%v) = %v, expected nil", tt.fieldExpr, result)
				}
				return
			}

			if result == nil {
				t.Errorf("GenerateFieldexpr(%v) = nil, expected non-nil", tt.fieldExpr)
				return
			}

			actualType := result.typeString()
			if actualType != tt.expectedType {
				t.Errorf("GenerateFieldexpr(%v).typeString() = %v, expected %v", tt.fieldExpr, actualType, tt.expectedType)
			}
		})
	}
}

func TestFieldexprInterface(t *testing.T) {
	// Test that all field types implement the Fieldexpr interface
	var _ FieldExpr = (*PlainField)(nil)
	var _ FieldExpr = (*NestedField)(nil)
	var _ FieldExpr = (*ComplexField)(nil)
}

func TestInitMethod(t *testing.T) {
	// Test that Init method can be called on all field types
	plainField := &PlainField{}
	nestedField := &NestedField{}
	complexField := &ComplexField{}

	// These should not panic
	plainField.Init("field1")
	nestedField.Init("field1[]")
	complexField.Init("{field1}")

	// Verify the fields are still functional after Init
	if plainField.typeString() != "plain" {
		t.Errorf("PlainField.typeString() after Init = %v, expected plain", plainField.typeString())
	}
	if nestedField.typeString() != "nested" {
		t.Errorf("NestedField.typeString() after Init = %v, expected nested", nestedField.typeString())
	}
	if complexField.typeString() != "complex" {
		t.Errorf("ComplexField.typeString() after Init = %v, expected complex", complexField.typeString())
	}
}

func TestFieldexprMap(t *testing.T) {
	// Test that the fieldexprMap contains the expected patterns and creates correct types
	expectedPatterns := []string{
		`^([a-zA-Z0-9]+)$`,
		`^([a-zA-Z0-9]+)\[\]$`,
		`^(\{[a-zA-Z0-9]+\})+$`,
	}

	expectedTypes := []string{
		"plain",
		"nested",
		"complex",
	}

	for i, pattern := range expectedPatterns {
		constructor, exists := fieldExprMap[pattern]
		if !exists {
			t.Errorf("fieldexprMap missing pattern: %v", pattern)
			continue
		}

		// Test that the constructor creates the right type
		instance := constructor("test")
		actualType := instance.typeString()
		if actualType != expectedTypes[i] {
			t.Errorf("fieldexprMap[%v] creates type %v, expected %v", pattern, actualType, expectedTypes[i])
		}
	}
}

func TestFieldTypesWithEmptyRecords(t *testing.T) {
	// Test field types with empty records
	emptyRecords := [][]string{}

	plainField := &PlainField{}
	nestedField := &NestedField{}
	complexField := &ComplexField{}

	tests := []struct {
		name   string
		field  FieldExpr
		expect bool
	}{
		{"PlainField with empty records", plainField, true},
		{"NestedField with empty records", nestedField, true},
		{"ComplexField with empty records", complexField, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.field.Check(emptyRecords, emptyRecords)
			if result != tt.expect {
				t.Errorf("%s.Check() = %v, expected %v", tt.name, result, tt.expect)
			}
		})
	}
}

func TestFieldTypesWithNilRecords(t *testing.T) {
	// Test field types with nil records
	plainField := &PlainField{}
	nestedField := &NestedField{}
	complexField := &ComplexField{}

	tests := []struct {
		name   string
		field  FieldExpr
		expect bool
	}{
		{"PlainField with nil records", plainField, true},
		{"NestedField with nil records", nestedField, true},
		{"ComplexField with nil records", complexField, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.field.Check(nil, nil)
			if result != tt.expect {
				t.Errorf("%s.Check() = %v, expected %v", tt.name, result, tt.expect)
			}
		})
	}
}
