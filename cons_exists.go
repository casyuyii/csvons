package csvons

import (
	"log"
	"slices"
)

// ExistsTest tests if the values in a column of a CSV file exist in a specified column of another file.
// @param stem the stem (base name) of the CSV file
// @param ruler the rules to be tested
// @param metadata the metadata of the CSV file
// @example ExistsTest("username",
//
//	[
//	 {
//	  "dst_file_stem": "username-d1",
//	  "fields": [
//	   {
//	    "src": "Username",
//	    "dst": "Username"
//	   }
//	  ]
//	 }
//	],
//	&Metadata{NameIndex: 0, DataIndex: 1, Extension: ".csv"})
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
			srcFieldPos := slices.Index(srcFields, field.Src)
			if srcFieldPos < 0 {
				log.Fatalf("src_field not found: %s", field.Src)
				return
			}

			dstFieldPos := slices.Index(dstFields, field.Dst)
			if dstFieldPos < 0 {
				log.Fatalf("dst_field not found: %s", field.Dst)
				return
			}

			searchedFields := make(map[string]int)
			for i := dataIndex; i < len(srcRecords); i++ {
				srcField := srcRecords[i][srcFieldPos]
				if _, ok := searchedFields[srcField]; ok {
					log.Printf("src_field [%s] value [%s] already searched at row [%d]", field.Src, srcField, searchedFields[srcField])
					continue
				}

				for j := dataIndex; j < len(dstRecords); j++ {
					dstField := dstRecords[j][dstFieldPos]
					searchedFields[dstField] = j
					if dstField == srcField {
						break
					}
				}

				rowIndex, ok := searchedFields[srcField]
				if !ok {
					log.Fatalf("can't find src_field [%s] value [%s] in dst_records", field.Src, srcField)
					return
				}
				log.Printf("found src_field [%s] value [%s] in dst_records at row [%d]", field.Src, srcField, rowIndex)
			}
		}
	}
}
