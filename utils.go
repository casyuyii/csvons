package csvons

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var (
	METADATA_KEY = "csvons_metadata"
)

// read config file
// @param configFileName the name of the config file
// @return a map of the config file and the metadata
// @example readConfigFile("ruler.json")
func readConfigFile(configFileName string) (map[string]json.RawMessage, *Metadata) {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		log.Printf("error opening file %s: %v", configFileName, err)
		return nil, nil
	}

	var cfg map[string]json.RawMessage
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Printf("error unmarshalling file %s: %v", configFileName, err)
		return nil, nil
	}

	var metadata *Metadata
	if v, ok := cfg[METADATA_KEY]; ok {
		var m Metadata
		err := json.Unmarshal(v, &m)
		if err != nil {
			log.Printf("error unmarshalling metadata: %v", err)
			return nil, nil
		}
		metadata = &m
	}

	if metadata == nil {
		log.Printf("metadata is nil")
		return nil, nil
	}

	delete(cfg, METADATA_KEY)
	return cfg, metadata
}

// read csv file
// @note file not cached, each file should only be read once
// @param stem the stem of the csv file
// @example readCsvFile("username", &Metadata{CSVFileFolder: "testdata", Extension: ".csv"})
func readCsvFile(stem string, metadata *Metadata) [][]string {
	if metadata == nil {
		log.Println("metadata is nil")
		return nil
	}

	fullPath := filepath.Join(metadata.CSVFileFolder, stem+metadata.Extension)
	csvFile, err := os.Open(fullPath)
	if err != nil {
		log.Printf("error opening file %s: %v", fullPath, err)
		return nil
	}

	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("error reading file %s: %v", fullPath, err)
		return nil
	}

	return records
}
