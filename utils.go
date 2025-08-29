package csvons

import (
	"encoding/csv"
	"encoding/json"
	"log/slog"
	"os"
)

func readCsvFile(csvFileName string) [][]string {
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		slog.Error("error opening file", "error", err)
		return nil
	}

	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		slog.Error("error reading file", "error", err)
		return nil
	}

	return records
}

func readConfigFile(configFileName string) map[string]json.RawMessage {
	data, err := os.ReadFile(configFileName)
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

func getExistsRuler(m map[string]json.RawMessage) []Exists {
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

func getMetadata(m map[string]json.RawMessage) *Metadata {
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

func getFieldPos(typeField []string, fieldName string) int {
	if typeField == nil {
		return -1
	}

	for i, record := range typeField {
		if record == fieldName {
			return i
		}
	}

	return -1
}
