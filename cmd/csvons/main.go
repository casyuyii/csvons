// Command csvons validates CSV files against constraint rules defined in a JSON configuration file.
//
// Usage:
//
//	csvons <ruler.json>
//
// The program reads the specified ruler JSON file, parses the metadata
// and constraint rules, then validates each referenced CSV file against its rules.
//
// Supported constraints:
//   - exists: values in a column must exist in another CSV file's column
//   - unique: values in a column must be unique across all rows
//   - vtype: values must conform to a specified type and optional range
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	csvons "csvons/internal/csvons"
)

func main() {
	// Require the ruler JSON file path as a command-line argument.
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <ruler.json>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nValidate CSV files against constraint rules defined in a JSON configuration file.\n")
		os.Exit(1)
	}
	configFileName := os.Args[1]

	// Read and parse the configuration file into rules (per-file constraints)
	// and metadata (CSV file structure information).
	rules, metadata := csvons.ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		log.Fatalf("read config file error: file_name=%s", configFileName)
		return
	}

	// Iterate over each CSV file stem and its associated constraint rules.
	for stem, v := range rules {
		// Parse the raw JSON into a map of constraint type → raw JSON data.
		rulers := map[string]json.RawMessage{}
		err := json.Unmarshal(v, &rulers)
		if err != nil {
			log.Fatalf("error unmarshalling rulers: error=%v", err)
			return
		}

		// Process each constraint type for the current CSV file.
		for k, v := range rulers {
			switch k {
			case "exists":
				// Validate cross-file value existence.
				var exists []csvons.Exists
				err := json.Unmarshal(v, &exists)
				if err != nil {
					log.Fatalf("error unmarshalling exists: error=%v", err)
					return
				}
				csvons.ExistsTest(stem, exists, metadata)

			case "unique":
				// Validate column value uniqueness.
				var unique csvons.Unique
				err := json.Unmarshal(v, &unique)
				if err != nil {
					log.Fatalf("error unmarshalling unique: error=%v", err)
					return
				}
				csvons.UniqueTest(stem, &unique, metadata)

			case "vtype":
				// Validate value types and ranges.
				var vtype []csvons.VType
				err := json.Unmarshal(v, &vtype)
				if err != nil {
					log.Fatalf("error unmarshalling vtype: error=%v", err)
					return
				}
				csvons.VTypeTest(stem, vtype, metadata)

			default:
				log.Fatalf("unknown key %s", k)
				return
			}
		}
	}
}
