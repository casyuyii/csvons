package csvons

import (
	"log"
)

// UniqueTest validates that all values in specified columns of a CSV file are unique.
//
// For each field name in the ruler's Fields list, it:
//  1. Creates a field expression from the field name
//  2. Extracts all values from the corresponding column
//  3. Counts occurrences and fails if any value appears more than once
//
// Calls log.Fatalf if any duplicate is found or if parameters are invalid.
func UniqueTest(stem string, ruler *Unique, metadata *Metadata) {
	// Validate input parameters.
	if ruler == nil || metadata == nil {
		failf("ruler [%v] or metadata [%v] is nil", ruler, metadata)
		return
	}
	log.Printf("checking src file %s ...", stem)

	// Validate metadata indices.
	nameIndex := metadata.NameIndex
	if nameIndex < 0 {
		failf("name_index [%d] is less than 0", nameIndex)
		return
	}
	log.Printf("name_index: %d", nameIndex)

	dataIndex := metadata.DataIndex
	if dataIndex <= nameIndex {
		failf("data_index [%d] is less than or equal to name_index [%d]", dataIndex, nameIndex)
		return
	}
	log.Printf("data_index: %d", dataIndex)

	// Read the source CSV file and validate it has enough rows.
	srcRecords := ReadCsvFile(stem, metadata)
	if srcLen := len(srcRecords); srcLen <= dataIndex {
		failf("src_records length [%d] <= data_index [%d]", srcLen, dataIndex)
		return
	}
	srcFields := srcRecords[nameIndex]
	log.Printf("src_fields: %q", srcFields)

	// Check uniqueness for each specified field.
	for _, fieldName := range ruler.Fields {
		// Create field expression to extract values from the column.
		fieldExpr := GenerateFieldExpr(metadata, fieldName)
		fieldVals := requiredFieldValues(fieldExpr, fieldName, srcFields, srcRecords)

		// Track value occurrences; fail on any duplicate.
		existingFields := make(map[string]int)
		for fieldVal := range fieldVals {
			existingFields[fieldVal] += 1
			if existingFields[fieldVal] > 1 {
				failf("src_field [%s] value [%s] already exists", fieldName, fieldVal)
			}
		}

		log.Printf("src_field [%s] values are unique", fieldName)
	}
}
