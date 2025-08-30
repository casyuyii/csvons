package csvons

import (
	"encoding/json"
	"log"
	"testing"
)

func TestExists(t *testing.T) {
	configFileName := "./ruler.json"
	m := readConfigFile(configFileName)
	if m == nil {
		log.Fatal("read config file error", "file_name", configFileName)
		return
	}

	metadata := getMetadata(m)
	if metadata == nil {
		log.Fatal("get metadata error", "file_name", configFileName)
		return
	}

	for stem, v := range m {
		rulers := map[string]json.RawMessage{}
		err := json.Unmarshal(v, &rulers)
		if err != nil {
			log.Fatal("error unmarshalling rulers", "error", err)
			return
		}

		for k, v := range rulers {
			if k == "exists" {
				var exists []Exists
				err := json.Unmarshal(v, &exists)
				if err != nil {
					log.Fatal("error unmarshalling exists", "error", err)
					return
				}
				ExistsTest(stem, exists, metadata)
			}
		}
	}
}
