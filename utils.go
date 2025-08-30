package csvons

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func readConfigFile(configFileName string) map[string]json.RawMessage {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		log.Fatal("error opening file", "error", err)
		return nil
	}

	var m map[string]json.RawMessage
	err = json.Unmarshal(data, &m)
	if err != nil {
		log.Fatal("error unmarshalling file", "error", err)
		return nil
	}

	return m
}

// read csv file
// @note file not cached, each file should only be read once
// @param stem the stem of the csv file
// @example readCsvFile("username", &Metadata{CSVFileFolder: "testdata", Extension: ".csv"})
func readCsvFile(stem string, metadata *Metadata) [][]string {
	if metadata == nil {
		log.Fatal("metadata is nil")
		return nil
	}

	fullPath := filepath.Join(metadata.CSVFileFolder, stem+metadata.Extension)
	csvFile, err := os.Open(fullPath)
	if err != nil {
		log.Fatal("error opening file", "full_path", fullPath, "error", err)
		return nil
	}

	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("error reading file", "error", err)
		return nil
	}

	return records
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
				log.Fatal("error unmarshalling metadata", "error", err)
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
