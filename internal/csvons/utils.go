package csvons

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// METADATA_KEY is the reserved key in the ruler JSON configuration file
// that holds the CSV metadata (file paths, separators, indices).
// This key is extracted from the config and not treated as a CSV file stem.
var METADATA_KEY = "csvons_metadata"

// ReadConfigFile reads and parses a ruler JSON configuration file.
// It extracts the metadata section (keyed by METADATA_KEY) and returns
// the remaining keys as a map of CSV file stems to their raw JSON rule definitions.
//
// Returns (nil, nil) if the file cannot be read, parsed, or lacks valid metadata.
//
// Example:
//
//	rules, metadata := ReadConfigFile("ruler.json")
//	// rules["username"] → raw JSON containing exists/unique/vtype rules
//	// metadata → parsed Metadata struct
func ReadConfigFile(configFileName string) (map[string]json.RawMessage, *Metadata) {
	// Read the entire config file into memory.
	data, err := os.ReadFile(configFileName)
	if err != nil {
		log.Printf("error opening file %s: %v", configFileName, err)
		return nil, nil
	}

	// Parse the top-level JSON object into a map of raw messages.
	// Each key is either a CSV file stem or the metadata key.
	var cfg map[string]json.RawMessage
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Printf("error unmarshalling file %s: %v", configFileName, err)
		return nil, nil
	}

	// Extract and parse the metadata section from the configuration.
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

	// Remove the metadata key so only CSV file stem rules remain.
	delete(cfg, METADATA_KEY)
	return cfg, metadata
}

// ReadCsvFile reads a CSV file identified by its stem (base name) and metadata.
// It constructs the full file path from the metadata's CSVFileFolder and Extension fields,
// then reads all records from the CSV file.
//
// Note: Files are not cached; each call reads the file from disk.
// Callers should ensure each file is read only once for performance.
//
// Returns nil if metadata is nil or the file cannot be opened/parsed.
//
// Example:
//
//	records := ReadCsvFile("username", metadata)
//	// records[0] → header row
//	// records[1:] → data rows
func ReadCsvFile(stem string, metadata *Metadata) [][]string {
	if metadata == nil {
		log.Println("metadata is nil")
		return nil
	}

	// Build the full file path: <folder>/<stem><extension>
	fullPath := filepath.Join(metadata.CSVFileFolder, stem+metadata.Extension)
	csvFile, err := os.Open(fullPath)
	if err != nil {
		log.Printf("error opening file %s: %v", fullPath, err)
		return nil
	}
	defer csvFile.Close() // Ensure the file handle is released after reading.

	// Parse the entire CSV file into a 2D string slice.
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Printf("error reading file %s: %v", fullPath, err)
		return nil
	}

	return records
}
