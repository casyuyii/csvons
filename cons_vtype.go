package csvons

import (
	"log"
	"slices"
	"strconv"
)

// VTypeTest tests if the values in a column of a CSV file are of a specified type.
// @param stem the stem (base name) of the CSV file
// @param ruler the rules to be tested
// @param metadata the metadata of the CSV file
func VTypeTest(stem string, ruler []VType, metadata *Metadata) {
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

	for _, vtype := range ruler {
		fieldName, factory := fieldsFactory(vtype.Field, metadata)
		if fieldName == "" || factory == nil {
			log.Fatalf("get field name or factory not found: %s", vtype.Field)
			return
		}

		srcFieldPos := slices.Index(srcFields, fieldName)
		if srcFieldPos < 0 {
			log.Fatalf("src_field not found: %s", fieldName)
			return
		}

		for i := dataIndex; i < len(srcRecords); i++ {
			srcFieldVals := factory(srcRecords[i][srcFieldPos])
			if len(srcFieldVals) == 0 {
				log.Fatalf("src_field [%s] value [%s] is empty", vtype.Field, srcRecords[i][srcFieldPos])
				return
			}

			for _, srcField := range srcFieldVals {
				log.Printf("checking src_field [%s] value [%s] of type [%s]", vtype.Field, srcField, vtype.Type)

				switch vtype.Type {
				case "int":
					v, ok := strconv.ParseInt(srcField, 10, 64)
					if ok != nil {
						log.Fatalf("src_field [%s] value [%s] is not an int", vtype.Field, srcField)
						return
					}
					if vtype.Range != nil {
						if v > int64(vtype.Range.Max) || v < int64(vtype.Range.Min) {
							log.Fatalf("src_field [%s] value [%s] is not in the range [%v, %v]", vtype.Field, srcField, vtype.Range.Min, vtype.Range.Max)
							return
						}
					}
				case "float64":
					v, ok := strconv.ParseFloat(srcField, 64)
					if ok != nil {
						log.Fatalf("src_field [%s] value [%s] is not a float64", vtype.Field, srcField)
						return
					}
					if vtype.Range != nil {
						if v > vtype.Range.Max || v < vtype.Range.Min {
							log.Fatalf("src_field [%s] value [%s] is not in the range [%v, %v]", vtype.Field, srcField, vtype.Range.Min, vtype.Range.Max)
							return
						}
					}
				case "bool":
					if _, ok := strconv.ParseBool(srcField); ok != nil {
						log.Fatalf("src_field [%s] value [%s] is not a bool", vtype.Field, srcField)
						return
					}
				default:
					log.Fatalf("src_field [%s] value [%s] is not a valid type", vtype.Field, srcField)
					return
				}
			}
		}
	}
}
