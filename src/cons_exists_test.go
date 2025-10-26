package csvons

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestExists(t *testing.T) {
	oldDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting working directory: error=%v", err)
		return
	}
	defer func() {
		err := os.Chdir(oldDir)
		if err != nil {
			log.Fatalf("error changing directory back to old directory: error=%v", err)
			return
		}
	}()
	err = os.Chdir("..")
	if err != nil {
		log.Fatalf("error changing directory: error=%v", err)
		return
	}

	configFileName := "./ruler.json"
	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		log.Fatalf("read config file error: file_name=%s", configFileName)
		return
	}

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		err := json.Unmarshal(v, &rulers)
		if err != nil {
			log.Fatalf("error unmarshalling rulers: error=%v", err)
			return
		}

		for k, v := range rulers {
			if k == "exists" {
				var exists []Exists
				err := json.Unmarshal(v, &exists)
				if err != nil {
					log.Fatalf("error unmarshalling exists: error=%v", err)
					return
				}
				ExistsTest(stem, exists, metadata)
			}
		}
	}
}
