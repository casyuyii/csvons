package csvons

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"
)

// projectRoot returns the absolute path to the project root directory.
// It uses runtime.Caller to resolve the path relative to this test file,
// navigating up from internal/csvons/ to the project root.
func projectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..")
}

// TestExists validates the "exists" constraint using the original ruler.json.
// It reads the config file, extracts exists rules, and runs ExistsTest
// against the CSV test data in testdata/.
func TestExists(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "ruler.json")

	rules, metadata := ReadConfigFile(configFileName)
	if rules == nil || metadata == nil {
		t.Fatalf("read config file error: file_name=%s", configFileName)
	}

	// Update the CSV file folder to use the absolute path.
	metadata.CSVFileFolder = filepath.Join(root, metadata.CSVFileFolder)

	for stem, v := range rules {
		rulers := map[string]json.RawMessage{}
		if err := json.Unmarshal(v, &rulers); err != nil {
			t.Fatalf("error unmarshalling rulers for %s: %v", stem, err)
		}

		for k, v := range rulers {
			if k == "exists" {
				var exists []Exists
				if err := json.Unmarshal(v, &exists); err != nil {
					t.Fatalf("error unmarshalling exists: %v", err)
				}
				ExistsTest(stem, exists, metadata)
			}
		}
	}
}

// TestExistsProducts validates exists constraints using the products test data.
func TestExistsProducts(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "testdata", "ruler_products.json")

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
			if k == "exists" {
				var exists []Exists
				if err := json.Unmarshal(v, &exists); err != nil {
					t.Fatalf("error unmarshalling exists: %v", err)
				}
				ExistsTest(stem, exists, metadata)
			}
		}
	}
}

// TestExistsOrders validates exists constraints using the orders test data.
func TestExistsOrders(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "testdata", "ruler_orders.json")

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
			if k == "exists" {
				var exists []Exists
				if err := json.Unmarshal(v, &exists); err != nil {
					t.Fatalf("error unmarshalling exists: %v", err)
				}
				ExistsTest(stem, exists, metadata)
			}
		}
	}
}

// TestExistsEmployees validates exists constraints using the employees test data.
func TestExistsEmployees(t *testing.T) {
	root := projectRoot()
	configFileName := filepath.Join(root, "testdata", "ruler_employees.json")

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
			if k == "exists" {
				var exists []Exists
				if err := json.Unmarshal(v, &exists); err != nil {
					t.Fatalf("error unmarshalling exists: %v", err)
				}
				ExistsTest(stem, exists, metadata)
			}
		}
	}
}
