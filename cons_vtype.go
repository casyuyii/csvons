package csvons

import (
	"log"
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
		fieldExpr := GenerateFieldExpr(metadata, vtype.Field)
		if fieldExpr == nil {
			log.Fatalf("field expression [%s] is nil", vtype.Field)
			return
		}
		fieldVals := fieldExpr.FieldValue(srcFields, srcRecords)

		typedSearchedFieldCache := make(map[string]map[string]bool)
		for fieldVal := range fieldVals {
			log.Printf("checking src_field [%s] value [%s] of type [%s]", vtype.Field, fieldVal, vtype.Type)

			if _, ok := typedSearchedFieldCache[vtype.Field]; !ok {
				typedSearchedFieldCache[vtype.Field] = make(map[string]bool)
			}

			switch vtype.Type {
			case "int":
				if _, ok := typedSearchedFieldCache[vtype.Field][fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] already checked", vtype.Field, fieldVal)
					continue
				}

				v, ok := strconv.ParseInt(fieldVal, 10, 64)
				if ok != nil {
					log.Fatalf("src_field [%s] value [%s] is not an int", vtype.Field, fieldVal)
					return
				}
				if vtype.Range != nil {
					if v > int64(vtype.Range.Max) || v < int64(vtype.Range.Min) {
						log.Fatalf("src_field [%s] value [%s] is not in the range [%v, %v]", vtype.Field, fieldVal, vtype.Range.Min, vtype.Range.Max)
						return
					}
				}
			case "float64":
				if _, ok := typedSearchedFieldCache[vtype.Field][fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] already checked", vtype.Field, fieldVal)
					continue
				}

				v, ok := strconv.ParseFloat(fieldVal, 64)
				if ok != nil {
					log.Fatalf("src_field [%s] value [%s] is not a float64", vtype.Field, fieldVal)
					return
				}
				if vtype.Range != nil {
					if v > vtype.Range.Max || v < vtype.Range.Min {
						log.Fatalf("src_field [%s] value [%s] is not in the range [%v, %v]", vtype.Field, fieldVal, vtype.Range.Min, vtype.Range.Max)
						return
					}
				}
			case "bool":
				if _, ok := typedSearchedFieldCache[vtype.Field][fieldVal]; ok {
					log.Printf("src_field [%s] value [%s] already checked", vtype.Field, fieldVal)
					continue
				}

				if _, ok := strconv.ParseBool(fieldVal); ok != nil {
					log.Fatalf("src_field [%s] value [%s] is not a bool", vtype.Field, fieldVal)
					return
				}
			default:
				log.Fatalf("src_field [%s] value [%s] is not a valid type", vtype.Field, fieldVal)
				return
			}
			typedSearchedFieldCache[vtype.Field][fieldVal] = true
		}
	}
}
