package csvons

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

// TestUnique validates the "unique" constraint using the original ruler.json.
func TestUnique(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "ruler", "ruler.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}

		for k, v := range rulers {
			if k == "unique" {
				var unique Unique
				if err := json.Unmarshal(v, &unique); err != nil {
					t.Fatalf("error unmarshalling unique: %v", err)
				}
				UniqueTest(stem, &unique, metadata)
			}
		}
	}
}

// TestUniqueProducts validates unique constraints using the products test data.
func TestUniqueProducts(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "ruler", "ruler_products.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}
		for k, v := range rulers {
			if k == "unique" {
				var unique Unique
				if err := json.Unmarshal(v, &unique); err != nil {
					t.Fatalf("error unmarshalling unique: %v", err)
				}
				UniqueTest(stem, &unique, metadata)
			}
		}
	}
}

// TestUniqueOrders validates unique constraints using the orders test data.
func TestUniqueOrders(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "ruler", "ruler_orders.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}
		for k, v := range rulers {
			if k == "unique" {
				var unique Unique
				if err := json.Unmarshal(v, &unique); err != nil {
					t.Fatalf("error unmarshalling unique: %v", err)
				}
				UniqueTest(stem, &unique, metadata)
			}
		}
	}
}

// TestUniqueEmployees validates unique constraints using the employees test data.
func TestUniqueEmployees(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "ruler", "ruler_employees.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}
		for k, v := range rulers {
			if k == "unique" {
				var unique Unique
				if err := json.Unmarshal(v, &unique); err != nil {
					t.Fatalf("error unmarshalling unique: %v", err)
				}
				UniqueTest(stem, &unique, metadata)
			}
		}
	}
}
