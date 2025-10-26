package csvons

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestUnique(t *testing.T) {
	oldDir, err := os.Getwd()
	if err != nil {
		log.Fatal("error getting working directory", "error", err)
		return
	}
	defer func() {
		err := os.Chdir(oldDir)
		if err != nil {
			log.Fatal("error changing directory back to old directory", "error", err)
			return
		}
	}()
	err = os.Chdir("..")
	if err != nil {
		log.Fatal("error changing directory", "error", err)
		return
	}

	configFileName := "./ruler.json"
	rules, metadata := ReadConfigFile(configFileName)
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
			if k == "unique" {
				var unique Unique
				err := json.Unmarshal(v, &unique)
				if err != nil {
					log.Fatal("error unmarshalling unique", "error", err)
					return
				}
				UniqueTest(stem, &unique, metadata)
			}
		}
	}
}
