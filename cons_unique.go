package csvons

import (
	"log"
	"slices"
)

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
		srcFieldPos := slices.Index(srcFields, fieldName)
		if srcFieldPos < 0 {
			log.Fatalf("src_field not found: %s", fieldName)
			return
		}

		existingFields := make(map[string]int)
		for i := dataIndex; i < len(srcRecords); i++ {
			srcField := srcRecords[i][srcFieldPos]
			if rowIndex, ok := existingFields[srcField]; ok {
				log.Fatalf("src_field [%s] value [%s] already exists at row [%d]", fieldName, srcField, rowIndex)
				return
			}
			existingFields[srcField] = i
		}
		log.Printf("src_field [%s] values are unique", fieldName)
	}
}
