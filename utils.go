package csvons

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"os"
)

func read_csv_file(csv_file_name string) [][]string {
	csv_file, err := os.Open(csv_file_name)
	if err != nil {
		slog.Error("error opening file", "error", err)
		return nil
	}

	csv_reader := csv.NewReader(csv_file)
	records, err := csv_reader.ReadAll()
	if err != nil {
		slog.Error("error reading file", "error", err)
		return nil
	}

	return records
}

func read_config_file(config_file_name string) map[string]json.RawMessage {
	data, err := os.ReadFile(config_file_name)
	if err != nil {
		slog.Error("error opening file", "error", err)
		return nil
	}

	var m map[string]json.RawMessage
	err = json.Unmarshal(data, &m)
	if err != nil {
		slog.Error("error unmarshalling file", "error", err)
		return nil
	}

	return m
}

func get_exists_ruler(m map[string]json.RawMessage) []Exists {
	if m == nil {
		return nil
	}

	for k, v := range m {
		if k == "exists" {
			var exists []Exists
			err := json.Unmarshal(v, &exists)
			if err != nil {
				slog.Error("error unmarshalling exists", "error", err)
				return nil
			}
			return exists
		}
	}

	slog.Error("exists ruler not found")
	return nil
}

func get_metadata(m map[string]json.RawMessage) *Metadata {
	if m == nil {
		return nil
	}

	for k, v := range m {
		if k == "metadata" {
			var metadata Metadata
			err := json.Unmarshal(v, &metadata)
			if err != nil {
				slog.Error("error unmarshalling metadata", "error", err)
				return nil
			}
			return &metadata
		}
	}

	return nil
}

func get_field_pos(type_field []string, field_name string) int {
	if type_field == nil {
		return -1
	}

	for i, record := range type_field {
		if record == field_name {
			return i
		}
	}

	return -1
}
