package csvons

import (
	"log"
	"slices"
)

func ExistsTest(stem string, ruler []Exists, metadata *Metadata) {
	if len(ruler) == 0 || metadata == nil {
		log.Fatal("ruler or metadata is nil", "ruler", ruler, "metadata", metadata)
		return
	}
	log.Printf("checking src file %s ...", stem)

	fieldNameIndex := metadata.FieldNameIndex
	if fieldNameIndex < 0 {
		log.Fatal("field_name_index is less than 0")
		return
	}
	log.Printf("field_name_index: %d", fieldNameIndex)

	srcRecords := readCsvFile(stem, metadata)
	srcLen := len(srcRecords)
	log.Printf("src_records length: %d", srcLen)
	if srcLen <= fieldNameIndex {
		log.Fatal("src_records <= field_name_index", "field_name_index", fieldNameIndex, "src_records length", srcLen)
		return
	}
	srcFileds := srcRecords[fieldNameIndex]
	log.Printf("src_fields: %s", srcFileds)

	for _, exist := range ruler {
		dstRecords := readCsvFile(exist.DstFileStem, metadata)
		if dstLen := len(dstRecords); dstLen <= fieldNameIndex || dstLen < srcLen {
			log.Fatal("dst_records length not enough", "dst_records length", dstLen)
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

			for i := fieldNameIndex + 1; i < len(srcRecords); i++ {
				srcField := srcRecords[i][srcFieldPos]

				found := -1
				for j := fieldNameIndex + 1; j < len(dstRecords); j++ {
					if slices.Contains(dstRecords[j], srcField) {
						found = j
						break
					}
				}

				if found < 0 {
					log.Fatalf("can't find src_field [%s] value [%s] in dst_records", field.Src, srcField)
					return
				}
				log.Printf("found src_field [%s] value [%s] in dst_records at pos [%d]", field.Src, srcField, found)
			}
		}
	}
}
