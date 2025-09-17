package csvons

import (
	"log"
)

// UniqueTest tests if the values in a column of a CSV file are unique.
// @param stem the stem (base name) of the CSV file
// @param ruler the rules to be tested
// @param metadata the metadata of the CSV file
func UniqueTest(stem string, ruler *Unique, metadata *Metadata) {
	if ruler == nil || metadata == nil {
		log.Fatalf("ruler [%v] or metadata [%v] is nil", ruler, metadata)
		return
	}
	log.Printf("checking src file %s ...", stem)

	nameIndex := metadata.NameIndex
	if nameIndex < 0 {
		log.Fatalf("name_index [%d] is less than 0", nameIndex)
		return
	}
	log.Printf("name_index: %d", nameIndex)

	dataIndex := metadata.DataIndex
	if dataIndex <= nameIndex {
		log.Fatalf("data_index [%d] is less than or equal to name_index [%d]", dataIndex, nameIndex)
		return
	}
	log.Printf("data_index: %d", dataIndex)

	srcRecords := readCsvFile(stem, metadata)
	if srcLen := len(srcRecords); srcLen <= dataIndex {
		log.Fatalf("src_records length [%d] <= data_index [%d]", srcLen, dataIndex)
		return
	}
	srcFields := srcRecords[nameIndex]
	log.Printf("src_fields: %q", srcFields)

	for _, fieldName := range ruler.Fields {
		fieldExpr := GenerateFieldExpr(metadata, fieldName)
		if fieldExpr == nil {
			log.Fatalf("field expression [%s] is nil", fieldName)
			return
		}
		fieldVals := fieldExpr.FieldValue(srcFields, srcRecords)

		existingFields := make(map[string]int)
		for fieldVal := range fieldVals {
			existingFields[fieldVal] += 1
			if existingFields[fieldVal] > 1 {
				log.Fatalf("src_field [%s] value [%s] already exists", fieldName, fieldVal)
			}
		}

		log.Printf("src_field [%s] values are unique", fieldName)
	}
}
