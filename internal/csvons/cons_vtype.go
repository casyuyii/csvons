package csvons

import (
	"log"
	"strconv"
)

// VTypeTest validates that values in specified columns conform to expected types
// and optionally fall within specified numeric ranges.
//
// Supported types:
//   - "int": values must be parseable as 64-bit integers
//   - "float64": values must be parseable as 64-bit floats
//   - "bool": values must be parseable as booleans (true/false, 1/0, etc.)
//
// For "int" and "float64" types, an optional Range constraint can specify
// minimum and maximum allowed values (inclusive).
//
// A per-field cache (typedSearchedFieldCache) skips re-checking values that
// have already been validated, improving performance for repeated values.
func VTypeTest(stem string, ruler []VType, metadata *Metadata) {
	// Validate input parameters.
	if len(ruler) == 0 || metadata == nil {
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

	// Validate each vtype rule against the CSV data.
	for _, vtype := range ruler {
		// Create field expression to extract values from the specified column.
		fieldExpr := GenerateFieldExpr(metadata, vtype.Field)
		fieldVals := requiredFieldValues(fieldExpr, vtype.Field, srcFields, srcRecords)

		// Cache already-checked values to avoid redundant type parsing.
		// Map structure: field_name → { value → already_checked }
		typedSearchedFieldCache := make(map[string]map[string]bool)
		for fieldVal := range fieldVals {
			log.Printf("checking src_field [%s] value [%s] of type [%s]", vtype.Field, fieldVal, vtype.Type)

			// Initialize the cache entry for this field if not present.
			if _, ok := typedSearchedFieldCache[vtype.Field]; !ok {
				typedSearchedFieldCache[vtype.Field] = make(map[string]bool)
			}

			switch vtype.Type {
			case "int":
				// Skip if this value was already validated for this field.
				if _, ok := typedSearchedFieldCache[vtype.Field][fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] already checked", vtype.Field, fieldVal)
					continue
				}

				// Parse the value as a 64-bit integer.
				v, ok := strconv.ParseInt(fieldVal, 10, 64)
				if ok != nil {
					failf("src_field [%s] value [%s] is not an int", vtype.Field, fieldVal)
					return
				}
				// Check range constraint if specified.
				if vtype.Range != nil {
					if v > int64(vtype.Range.Max) || v < int64(vtype.Range.Min) {
						failf("src_field [%s] value [%s] is not in the range [%v, %v]", vtype.Field, fieldVal, vtype.Range.Min, vtype.Range.Max)
						return
					}
				}

			case "float64":
				// Skip if this value was already validated for this field.
				if _, ok := typedSearchedFieldCache[vtype.Field][fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] already checked", vtype.Field, fieldVal)
					continue
				}

				// Parse the value as a 64-bit float.
				v, ok := strconv.ParseFloat(fieldVal, 64)
				if ok != nil {
					failf("src_field [%s] value [%s] is not a float64", vtype.Field, fieldVal)
					return
				}
				// Check range constraint if specified.
				if vtype.Range != nil {
					if v > vtype.Range.Max || v < vtype.Range.Min {
						failf("src_field [%s] value [%s] is not in the range [%v, %v]", vtype.Field, fieldVal, vtype.Range.Min, vtype.Range.Max)
						return
					}
				}

			case "bool":
				// Skip if this value was already validated for this field.
				if _, ok := typedSearchedFieldCache[vtype.Field][fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] already checked", vtype.Field, fieldVal)
					continue
				}

				// Validate boolean parsing (accepts: true/false, 1/0, t/f, yes/no, etc.)
				if _, ok := strconv.ParseBool(fieldVal); ok != nil {
					failf("src_field [%s] value [%s] is not a bool", vtype.Field, fieldVal)
					return
				}

			default:
				failf("src_field [%s] value [%s] is not a valid type", vtype.Field, fieldVal)
				return
			}
			// Mark this value as checked in the cache.
			typedSearchedFieldCache[vtype.Field][fieldVal] = true
		}
	}
}
