package csvons

import (
	"fmt"
	"log/slog"
)

func exists_test(ruler []Exists, metadata *Metadata) {
	if len(ruler) == 0 || metadata == nil {
		slog.Error("ruler or metadata is nil")
		return
	}

	field_name_index := metadata.FieldNameIndex

	for _, exist := range ruler {
		src_records := read_csv_file(exist.Src.FileName)
		dst_records := read_csv_file(exist.Dst.FileName)

		if field_name_index < 0 {
			slog.Error("field_name_index is less than 0")
			return
		}

		if len(src_records) <= field_name_index {
			slog.Error("src_records is less than field_name_index", "field_name_index", field_name_index, "src_records length", len(src_records))
			return
		}

		if len(dst_records) <= field_name_index {
			slog.Error("dst_records is less than field_name_index", "field_name_index", field_name_index, "dst_records length", len(dst_records))
			return
		}

		src_field_pos := get_field_pos(src_records[field_name_index], exist.Src.FieldName)
		dst_field_pos := get_field_pos(dst_records[field_name_index], exist.Dst.FieldName)
		if src_field_pos < 0 || dst_field_pos < 0 {
			slog.Error("src_field_pos or dst_field_pos is less than 0", "src_field_pos", src_field_pos, "dst_field_pos", dst_field_pos)
			return
		}

		for i := field_name_index + 1; i < len(src_records); i++ {
			if i >= len(dst_records) {
				slog.Error(fmt.Sprintf("field %s is not found in dst_records(dst_records length not enough)", exist.Src.FieldName))
				return
			}

			src_field := src_records[i][src_field_pos]
			dst_field := dst_records[i][dst_field_pos]

			slog.Info("src_field", "src_field", src_field)
			slog.Info("dst_field", "dst_field", dst_field)

			if src_field != dst_field {
				slog.Error("src_field and dst_field are not equal", "src_field", src_field, "dst_field", dst_field)
				return
			}
		}
	}
}
