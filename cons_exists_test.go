package csvons

import (
	"encoding/json"
	"log"
	"testing"
)

func TestExists(t *testing.T) {
	configFileName := "./ruler.json"
	rules, metadata := readConfigFile(configFileName)
	if rules == nil || metadata == nil {
		log.Fatal("read config file error", "file_name", configFileName)
		return
	}

	for stem, v := range rules {
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
