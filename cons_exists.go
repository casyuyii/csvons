package csvons

import (
	"log"
)

// ExistsTest tests if the values in a column of a CSV file exist in a specified column of another file.
// @param stem the stem (base name) of the CSV file
// @param ruler the rules to be tested
// @param metadata the metadata of the CSV file
func ExistsTest(stem string, ruler []Exists, metadata *Metadata) {
	if len(ruler) == 0 || metadata == nil {
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

	for _, exist := range ruler {
		dstRecords := readCsvFile(exist.DstFileStem, metadata)
		if dstLen := len(dstRecords); dstLen <= dataIndex {
			log.Fatalf("dst_records length [%d] <= data_index [%d]", dstLen, dataIndex)
			return
		}
		log.Printf("checking dst file %s ...", exist.DstFileStem)

		dstFields := dstRecords[nameIndex]
		log.Printf("dst_fields: %q", dstFields)

		for _, field := range exist.Fields {
			srcFieldExpr := GenerateFieldExpr(metadata, field.Src)
			if srcFieldExpr == nil {
				log.Fatalf("field expression [%s] is nil", field.Src)
				return
			}
			srcFieldVals := srcFieldExpr.FieldValue(srcFields, srcRecords)

			dstFieldExpr := GenerateFieldExpr(metadata, field.Dst)
			if dstFieldExpr == nil {
				log.Fatalf("field expression [%s] is nil", field.Dst)
				return
			}
			dstFieldVals := dstFieldExpr.FieldValue(dstFields, dstRecords)

			searchedFields := make(map[string]int)
			cacheDstFieldVals := make(map[string]bool)
			for fieldVal := range srcFieldVals {
				if _, ok := searchedFields[fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] already searched at row [%d]", field.Src, fieldVal, searchedFields[fieldVal])
					continue
				}

				if _, ok := cacheDstFieldVals[fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] hit cache", field.Src, fieldVal)
					continue
				}

				for dstFieldVal := range dstFieldVals {
					cacheDstFieldVals[dstFieldVal] = true
					if dstFieldVal == fieldVal {
						log.Printf("found src_field [%s] value [%s] in dst_records", field.Src, fieldVal)
						break
					}
				}

				if _, ok := cacheDstFieldVals[fieldVal]; !ok {
					log.Fatalf("src_field [%s] value [%s] not found in dst_records", field.Src, fieldVal)
				}
			}
		}
	}
}
