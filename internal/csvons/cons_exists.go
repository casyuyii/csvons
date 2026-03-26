package csvons

import (
	"log"
)

// ExistsTest validates that values in specified columns of a source CSV file
// also exist in corresponding columns of destination CSV files.
//
// For each rule in the ruler slice, this function:
//  1. Reads the source CSV file using the stem parameter
//  2. Reads the destination CSV file specified by each rule's DstFileStem
//  3. For each field pair, extracts values using field expressions
//  4. Verifies every source value exists in the destination values
//
// The function uses a cache (cacheDstFieldVals) to avoid redundant lookups
// and a searchedFields map to skip already-verified source values.
//
// Calls log.Fatalf if any source value is not found in the destination,
// or if required parameters are invalid.
func ExistsTest(stem string, ruler []Exists, metadata *Metadata) {
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

	// Check each existence rule against its destination file.
	for _, exist := range ruler {
		// Read the destination CSV file.
		dstRecords := ReadCsvFile(exist.DstFileStem, metadata)
		if dstLen := len(dstRecords); dstLen <= dataIndex {
			failf("dst_records length [%d] <= data_index [%d]", dstLen, dataIndex)
			return
		}
		log.Printf("checking dst file %s ...", exist.DstFileStem)

		dstFields := dstRecords[nameIndex]
		log.Printf("dst_fields: %q", dstFields)

		// Validate each pair of source and destination fields.
		for _, field := range exist.Fields {
			// Create field expression for the source column.
			srcFieldExpr := GenerateFieldExpr(metadata, field.Src)
			srcFieldVals := requiredFieldValues(srcFieldExpr, field.Src, srcFields, srcRecords)

			// Create field expression for the destination column.
			dstFieldExpr := GenerateFieldExpr(metadata, field.Dst)
			dstFieldVals := requiredFieldValues(dstFieldExpr, field.Dst, dstFields, dstRecords)

			// Track already-searched source values and cache destination values.
			searchedFields := make(map[string]int)
			cacheDstFieldVals := make(map[string]bool)

			for fieldVal := range srcFieldVals {
				// Skip source values we've already verified.
				if _, ok := searchedFields[fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] already searched at row [%d]", field.Src, fieldVal, searchedFields[fieldVal])
					continue
				}

				// Check cache first before iterating destination values.
				if _, ok := cacheDstFieldVals[fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] hit cache", field.Src, fieldVal)
					continue
				}

				// Iterate through destination values until we find a match.
				// Each consumed destination value is cached for future lookups.
				for dstFieldVal := range dstFieldVals {
					cacheDstFieldVals[dstFieldVal] = true
					if dstFieldVal == fieldVal {
						log.Printf("found src_field [%s] value [%s] in dst_records", field.Src, fieldVal)
						break
					}
				}

				// If the value was not found after exhausting destination values, fail.
				if _, ok := cacheDstFieldVals[fieldVal]; !ok {
					failf("src_field [%s] value [%s] not found in dst_records", field.Src, fieldVal)
				}
			}
		}
	}
}
