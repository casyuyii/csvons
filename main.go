package main

import (
	"encoding/json"
	"log"

	csvons "csvons/src"
)

func main() {
	configFileName := "./ruler.json"
	rules, metadata := csvons.ReadConfigFile(configFileName)
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
			switch k {
			case "exists":
				var exists []csvons.Exists
				err := json.Unmarshal(v, &exists)
				if err != nil {
					log.Fatal("error unmarshalling exists", "error", err)
					return
				}
				csvons.ExistsTest(stem, exists, metadata)
			case "unique":
				var unique csvons.Unique
				err := json.Unmarshal(v, &unique)
				if err != nil {
					log.Fatal("error unmarshalling unique", "error", err)
					return
				}
				csvons.UniqueTest(stem, &unique, metadata)
			case "vtype":
				var vtype []csvons.VType
				err := json.Unmarshal(v, &vtype)
				if err != nil {
					log.Fatal("error unmarshalling vtype", "error", err)
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
