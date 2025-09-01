package csvons

import (
	"log"
)

func ExistsTest(stem string, ruler []Exists, metadata *Metadata) {
	if len(ruler) == 0 || metadata == nil {
		log.Fatalf("ruler [%v] or metadata [%v] is nil", ruler, metadata)
		return
	}
	log.Printf("checking src file %s ...", stem)

	fieldNameIndex := metadata.FieldNameIndex
	if fieldNameIndex < 0 {
		log.Fatalf("field_name_index [%d] is less than 0", fieldNameIndex)
		return
	}
	log.Printf("field_name_index: %d", fieldNameIndex)

	dataIndex := metadata.DataIndex
	if dataIndex <= fieldNameIndex {
		log.Fatalf("data_index [%d] is less than or equal to field_name_index [%d]", dataIndex, fieldNameIndex)
		return
	}
	log.Printf("data_index: %d", dataIndex)

	srcRecords := readCsvFile(stem, metadata)
	if srcLen := len(srcRecords); srcLen <= dataIndex {
		log.Fatalf("src_records length [%d] <= data_index [%d]", srcLen, dataIndex)
		return
	}
	srcFileds := srcRecords[fieldNameIndex]
	log.Printf("src_fields: %s", srcFileds)

	for _, exist := range ruler {
		dstRecords := readCsvFile(exist.DstFileStem, metadata)
		if dstLen := len(dstRecords); dstLen <= dataIndex {
			log.Fatalf("dst_records length not enough %d", dstLen)
			return
		}
		log.Printf("checking dst file %s ...", exist.DstFileStem)

		dstFileds := dstRecords[fieldNameIndex]
		log.Printf("dst_fields: %s", dstFileds)

		for _, field := range exist.Fields {
			srcFieldPos := getFieldPos(srcFileds, field.Src)
			if srcFieldPos < 0 {
				log.Fatalf("src_field not found: %s", field.Src)
				return
			}

			dstFieldPos := getFieldPos(dstFileds, field.Dst)
			if dstFieldPos < 0 {
				log.Fatalf("dst_field not found: %s", field.Dst)
				return
			}

			searchedField := make(map[string]int)
			for i := dataIndex; i < len(srcRecords); i++ {
				srcField := srcRecords[i][srcFieldPos]
				if _, ok := searchedField[srcField]; ok {
					log.Printf("src_field [%s] value [%s] already searched at row [%d]", field.Src, srcField, searchedField[srcField])
					continue
				}

				for j := dataIndex; j < len(dstRecords); j++ {
					dstField := dstRecords[j][dstFieldPos]
					searchedField[dstField] = j
					if dstField == srcField {
						break
					}
				}

				rowIndex, ok := searchedField[srcField]
				if !ok {
					log.Fatalf("can't find src_field [%s] value [%s] in dst_records", field.Src, srcField)
					return
				}
				log.Printf("found src_field [%s] value [%s] in dst_records at row [%d]", field.Src, srcField, rowIndex)
			}
		}
	}
}
