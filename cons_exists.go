package csvons

import (
	"fmt"
	"log/slog"
)

func existsTest(ruler []Exists, metadata *Metadata) {
	if len(ruler) == 0 || metadata == nil {
		slog.Error("ruler or metadata is nil", "ruler", ruler, "metadata", metadata)
		return
	}

	fieldNameIndex := metadata.FieldNameIndex

	for _, exist := range ruler {
		srcRecords := readCsvFile(exist.Src.FileName)
		dstRecords := readCsvFile(exist.Dst.FileName)

		if fieldNameIndex < 0 {
			slog.Error("field_name_index is less than 0")
			return
		}

		if len(srcRecords) <= fieldNameIndex {
			slog.Error("src_records is less than field_name_index", "field_name_index", fieldNameIndex, "src_records length", len(srcRecords))
			return
		}

		if len(dstRecords) <= fieldNameIndex {
			slog.Error("dst_records is less than field_name_index", "field_name_index", fieldNameIndex, "dst_records length", len(dstRecords))
			return
		}

		srcFieldPos := getFieldPos(srcRecords[fieldNameIndex], exist.Src.FieldName)
		dstFieldPos := getFieldPos(dstRecords[fieldNameIndex], exist.Dst.FieldName)
		if srcFieldPos < 0 || dstFieldPos < 0 {
			slog.Error("src_field_pos or dst_field_pos is less than 0", "src_field_pos", srcFieldPos, "dst_field_pos", dstFieldPos)
			return
		}

		for i := fieldNameIndex + 1; i < len(srcRecords); i++ {
			if i >= len(dstRecords) {
				slog.Error(fmt.Sprintf("field %s is not found in dst_records(dst_records length not enough)", exist.Src.FieldName))
				return
			}

			srcField := srcRecords[i][srcFieldPos]
			dstField := dstRecords[i][dstFieldPos]

			slog.Info("src_field", "src_field", srcField)
			slog.Info("dst_field", "dst_field", dstField)

			if srcField != dstField {
				slog.Error("src_field and dst_field are not equal", "src_field", srcField, "dst_field", dstField)
				return
			}
		}
	}
}
